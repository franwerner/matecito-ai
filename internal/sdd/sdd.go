package sdd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/mcp"
)

type ServerRef struct {
	Prefix string
	Tools  []string
	Agents []string
}

func CrossCheck(dir string) []check.Result {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return []check.Result{{
			Name:     "SDD agents",
			Required: false,
			Status:   check.StatusMissing,
			Detail:   fmt.Sprintf("no se encontró el directorio %s", dir),
			FixHint:  "Indicá otra ubicación con --sdd-dir <path>.",
		}}
	}

	refs, err := loadServerRefs(dir)
	if err != nil {
		return []check.Result{{
			Name:     "SDD agents",
			Required: false,
			Status:   check.StatusMissing,
			Detail:   fmt.Sprintf("error leyendo agentes en %s: %v", dir, err),
		}}
	}
	if len(refs) == 0 {
		return []check.Result{{
			Name:     "SDD agents",
			Required: false,
			Status:   check.StatusMissing,
			Detail:   fmt.Sprintf("ningún tool MCP declarado en %s/sdd-*.md", dir),
		}}
	}

	servers := mcp.ListAll()
	registered := map[string]string{}
	for _, s := range servers {
		registered[strings.ReplaceAll(s, ":", "_")] = s
	}

	results := make([]check.Result, 0, len(refs))
	for _, ref := range refs {
		r := check.Result{
			Name:     ref.Prefix,
			Required: true,
		}
		if original, ok := registered[ref.Prefix]; ok {
			r.Status = check.StatusOK
			r.Detail = fmt.Sprintf("%d tools en SDD → MCP %q", len(ref.Tools), original)
		} else {
			r.Status = check.StatusMissing
			r.Detail = fmt.Sprintf("%d tools en SDD, ningún MCP coincide con %q", len(ref.Tools), ref.Prefix)
			r.FixHint = "Corregí el frontmatter del agente, o registrá el MCP con un nombre que matchee el prefijo."
		}
		results = append(results, r)
	}
	return results
}

func loadServerRefs(dir string) ([]ServerRef, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "sdd-*.md"))
	if err != nil {
		return nil, err
	}

	byPrefix := map[string]*ServerRef{}
	for _, p := range matches {
		tools, err := readFrontmatterTools(p)
		if err != nil {
			continue
		}
		agent := filepath.Base(p)
		for _, t := range tools {
			if !strings.HasPrefix(t, "mcp__") {
				continue
			}
			rest := strings.TrimPrefix(t, "mcp__")
			parts := strings.SplitN(rest, "__", 2)
			if len(parts) != 2 {
				continue
			}
			prefix := parts[0]
			ref, ok := byPrefix[prefix]
			if !ok {
				ref = &ServerRef{Prefix: prefix}
				byPrefix[prefix] = ref
			}
			ref.Tools = append(ref.Tools, t)
			if !contains(ref.Agents, agent) {
				ref.Agents = append(ref.Agents, agent)
			}
		}
	}
	refs := make([]ServerRef, 0, len(byPrefix))
	for _, r := range byPrefix {
		refs = append(refs, *r)
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].Prefix < refs[j].Prefix })
	return refs, nil
}

func readFrontmatterTools(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	inFM := false
	for sc.Scan() {
		line := sc.Text()
		if strings.TrimSpace(line) == "---" {
			if !inFM {
				inFM = true
				continue
			}
			break
		}
		if !inFM {
			continue
		}
		if strings.HasPrefix(line, "tools:") {
			return splitToolList(strings.TrimSpace(strings.TrimPrefix(line, "tools:"))), nil
		}
	}
	return nil, sc.Err()
}

func splitToolList(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
