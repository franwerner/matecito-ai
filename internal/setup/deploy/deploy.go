package deploy

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"

	matecitoai "github.com/franwerner/matecito-ai"
)

type Mode int

const (
	ModeFile Mode = iota
	ModeDir
	ModeGrouped
)

type Mapping struct {
	SourceRel string
	TargetRel string
	Mode      Mode
}

var Mappings = []Mapping{
	{SourceRel: "CLAUDE.md", TargetRel: "matecito-ai.md", Mode: ModeFile},
	{SourceRel: "agents", TargetRel: "agents", Mode: ModeDir},
	{SourceRel: "skills", TargetRel: "skills", Mode: ModeGrouped},
	{SourceRel: "references", TargetRel: "references", Mode: ModeDir},
}

type FileStatus int

const (
	StatusNew FileStatus = iota
	StatusChanged
	StatusSame
)

// FileOp describe la copia de un archivo del payload (Source, ruta interna al
// fs.FS) hacia el filesystem real del usuario (Target).
type FileOp struct {
	Source string
	Target string
	Status FileStatus
}

// FindPayloadDir busca payload/ en el filesystem real desde start hacia
// arriba. Solo se usa cuando el binario se ejecuta desde un clone del repo
// (modo dev). En binarios distribuidos se prefiere el embedded FS.
func FindPayloadDir(start string) (string, error) {
	dir := start
	for {
		p := filepath.Join(dir, "payload")
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return p, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no se encontró payload/ en cwd ni en ningún directorio padre")
		}
		dir = parent
	}
}

func ClaudeHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude"), nil
}

func BackupRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".matecito-ai", "backups"), nil
}

// BackupDir devuelve la carpeta de backup para la corrida actual:
// ~/.matecito-ai/backups/<timestamp>/. NO crea el directorio; el caller
// lo crea solo si llega a haber algo que respaldar (lazy).
func BackupDir() (string, error) {
	root, err := BackupRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, time.Now().Format("20060102-150405")), nil
}

// Plan arma la lista de operaciones de copia para todos los Mappings.
// payloadFS es la raíz del payload (sea os.DirFS sobre un payload/ local, o
// un sub-FS de PayloadFS embebido). claudeHome es la carpeta real de destino.
func Plan(payloadFS fs.FS, claudeHome string) ([]FileOp, error) {
	var ops []FileOp
	for _, m := range Mappings {
		more, err := expandMapping(payloadFS, claudeHome, m)
		if err != nil {
			return nil, err
		}
		ops = append(ops, more...)
	}

	seen := map[string]string{}
	for _, op := range ops {
		if prev, dup := seen[op.Target]; dup {
			return nil, fmt.Errorf("clash de nombres en deploy: %q y %q apuntan al mismo destino %q", prev, op.Source, op.Target)
		}
		seen[op.Target] = op.Source
	}

	for i := range ops {
		ops[i].Status = computeStatus(payloadFS, ops[i].Source, ops[i].Target)
	}
	return ops, nil
}

func expandMapping(payloadFS fs.FS, claudeHome string, m Mapping) ([]FileOp, error) {
	source := m.SourceRel
	info, err := fs.Stat(payloadFS, source)
	if err != nil {
		return nil, fmt.Errorf("payload no tiene %q: %w", m.SourceRel, err)
	}

	switch m.Mode {
	case ModeFile:
		if info.IsDir() {
			return nil, fmt.Errorf("se esperaba archivo en %q, es directorio", source)
		}
		return []FileOp{{Source: source, Target: filepath.Join(claudeHome, m.TargetRel)}}, nil

	case ModeDir:
		return walkDir(payloadFS, source, filepath.Join(claudeHome, m.TargetRel))

	case ModeGrouped:
		groups, err := fs.ReadDir(payloadFS, source)
		if err != nil {
			return nil, err
		}
		var ops []FileOp
		for _, g := range groups {
			if !g.IsDir() {
				continue
			}
			groupPath := path.Join(source, g.Name())
			entries, err := fs.ReadDir(payloadFS, groupPath)
			if err != nil {
				return nil, err
			}
			targetBase := filepath.Join(claudeHome, m.TargetRel)
			for _, e := range entries {
				childSrc := path.Join(groupPath, e.Name())
				childDst := filepath.Join(targetBase, e.Name())
				if e.IsDir() {
					more, err := walkDir(payloadFS, childSrc, childDst)
					if err != nil {
						return nil, err
					}
					ops = append(ops, more...)
				} else {
					ops = append(ops, FileOp{Source: childSrc, Target: childDst})
				}
			}
		}
		return ops, nil
	}
	return nil, fmt.Errorf("modo desconocido: %d", m.Mode)
}

func walkDir(payloadFS fs.FS, sourceRoot, targetRoot string) ([]FileOp, error) {
	var ops []FileOp
	err := fs.WalkDir(payloadFS, sourceRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(sourceRoot, p)
		if err != nil {
			return err
		}
		ops = append(ops, FileOp{Source: p, Target: filepath.Join(targetRoot, rel)})
		return nil
	})
	return ops, err
}

func computeStatus(payloadFS fs.FS, source, target string) FileStatus {
	dst, err := os.ReadFile(target)
	if err != nil {
		return StatusNew
	}
	src, err := fs.ReadFile(payloadFS, source)
	if err != nil {
		return StatusNew
	}
	if bytes.Equal(src, dst) {
		return StatusSame
	}
	return StatusChanged
}

type Summary struct {
	New     int
	Changed int
	Same    int
}

func Summarize(ops []FileOp) Summary {
	var s Summary
	for _, op := range ops {
		switch op.Status {
		case StatusNew:
			s.New++
		case StatusChanged:
			s.Changed++
		case StatusSame:
			s.Same++
		}
	}
	return s
}

// ResolvePayloadFS prefiere un payload/ local (modo dev/source) si existe en
// el cwd o algún directorio padre. Si no, cae al payload embebido en el
// binario via go:embed.
func ResolvePayloadFS() (payloadFS fs.FS, source string, err error) {
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		if local, findErr := FindPayloadDir(cwd); findErr == nil {
			return os.DirFS(local), local, nil
		}
	}
	sub, err := fs.Sub(matecitoai.PayloadFS, "payload")
	if err != nil {
		return nil, "", err
	}
	return sub, "embedded", nil
}

// Apply ejecuta las copias en disco. Lee de payloadFS y escribe en
// claudeHome. Si hay cambios sobre archivos existentes, los respalda en
// backupDir/<rel-path>. backupDir es la carpeta concreta de esta corrida
// (típicamente de deploy.BackupDir()) — Apply lo crea on-demand solo si
// llega a haber al menos un archivo cambiado.
// Los archivos de payload se copian byte a byte sin ninguna transformación.
func Apply(payloadFS fs.FS, ops []FileOp, claudeHome, backupDir string) (bool, error) {
	backupCreated := false
	for _, op := range ops {
		if op.Status == StatusSame {
			continue
		}

		if err := os.MkdirAll(filepath.Dir(op.Target), 0o755); err != nil {
			return backupCreated, err
		}

		// backup solo para archivos que ya existían y cambiaron
		if op.Status == StatusChanged {
			rel, err := filepath.Rel(claudeHome, op.Target)
			if err != nil {
				return backupCreated, err
			}
			bk := filepath.Join(backupDir, rel)
			if err := os.MkdirAll(filepath.Dir(bk), 0o755); err != nil {
				return backupCreated, err
			}
			if err := copyDiskFile(op.Target, bk); err != nil {
				return backupCreated, err
			}
			backupCreated = true
		}

		if err := copyFromFS(payloadFS, op.Source, op.Target); err != nil {
			return backupCreated, err
		}
	}
	return backupCreated, nil
}

func copyFromFS(payloadFS fs.FS, src, dst string) error {
	in, err := payloadFS.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func copyDiskFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
