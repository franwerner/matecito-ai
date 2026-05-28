package deploy

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
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
}

type FileStatus int

const (
	StatusNew FileStatus = iota
	StatusChanged
	StatusSame
)

type FileOp struct {
	Source string
	Target string
	Status FileStatus
}

func FindPayload(start string) (string, error) {
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

func Plan(payloadDir, claudeHome string) ([]FileOp, error) {
	var ops []FileOp
	for _, m := range Mappings {
		more, err := expandMapping(payloadDir, claudeHome, m)
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
		ops[i].Status = computeStatus(ops[i].Source, ops[i].Target)
	}
	return ops, nil
}

func expandMapping(payloadDir, claudeHome string, m Mapping) ([]FileOp, error) {
	source := filepath.Join(payloadDir, m.SourceRel)
	info, err := os.Stat(source)
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
		return walkDir(source, filepath.Join(claudeHome, m.TargetRel))

	case ModeGrouped:
		groups, err := os.ReadDir(source)
		if err != nil {
			return nil, err
		}
		var ops []FileOp
		for _, g := range groups {
			if !g.IsDir() {
				continue
			}
			groupPath := filepath.Join(source, g.Name())
			entries, err := os.ReadDir(groupPath)
			if err != nil {
				return nil, err
			}
			targetBase := filepath.Join(claudeHome, m.TargetRel)
			for _, e := range entries {
				childSrc := filepath.Join(groupPath, e.Name())
				childDst := filepath.Join(targetBase, e.Name())
				if e.IsDir() {
					more, err := walkDir(childSrc, childDst)
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

func walkDir(sourceRoot, targetRoot string) ([]FileOp, error) {
	var ops []FileOp
	err := filepath.WalkDir(sourceRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(sourceRoot, path)
		if err != nil {
			return err
		}
		ops = append(ops, FileOp{Source: path, Target: filepath.Join(targetRoot, rel)})
		return nil
	})
	return ops, err
}

func computeStatus(source, target string) FileStatus {
	dst, err := os.ReadFile(target)
	if err != nil {
		return StatusNew
	}
	src, err := os.ReadFile(source)
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

func Apply(ops []FileOp, claudeHome, backupRoot string) (string, error) {
	backupDir := ""
	for _, op := range ops {
		if op.Status == StatusSame {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(op.Target), 0o755); err != nil {
			return backupDir, err
		}
		if op.Status == StatusChanged {
			if backupDir == "" {
				backupDir = filepath.Join(backupRoot, time.Now().Format("20060102-150405"))
			}
			rel, err := filepath.Rel(claudeHome, op.Target)
			if err != nil {
				return backupDir, err
			}
			bk := filepath.Join(backupDir, rel)
			if err := os.MkdirAll(filepath.Dir(bk), 0o755); err != nil {
				return backupDir, err
			}
			if err := copyFile(op.Target, bk); err != nil {
				return backupDir, err
			}
		}
		if err := copyFile(op.Source, op.Target); err != nil {
			return backupDir, err
		}
	}
	return backupDir, nil
}

func copyFile(src, dst string) error {
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
