package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
)

// ProjectContext carries the result of project detection from the cwd.
type ProjectContext struct {
	Name      string
	RepoRoot  string
	InProject bool
}

// DetectProject walks up from the current working directory looking for a
// .git/ directory or a .matecito-ai/ directory. When found it shells out to
// git to obtain the remote URL and delegates name derivation to the pure
// agentmodel.DeriveProjectName so that parsing stays tested and side-effect-free.
func DetectProject() ProjectContext {
	cwd, err := os.Getwd()
	if err != nil {
		return ProjectContext{}
	}

	root := findRepoRoot(cwd)
	if root == "" {
		return ProjectContext{}
	}

	remoteURL := gitRemoteURL(root)
	name := agentmodel.DeriveProjectName(remoteURL, root)

	return ProjectContext{
		Name:      name,
		RepoRoot:  root,
		InProject: true,
	}
}

// findRepoRoot walks ancestor directories from dir until it finds one that
// contains .git/ or .matecito-ai/, returning that directory path.
// Returns "" when neither marker is found up to the filesystem root.
func findRepoRoot(dir string) string {
	current := filepath.Clean(dir)
	for {
		if dirExists(filepath.Join(current, ".git")) || dirExists(filepath.Join(current, ".matecito-ai")) {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
}

// gitRemoteURL runs git -C root remote get-url origin and returns the trimmed
// output. Returns "" on any error (no remote, git not found, etc.).
func gitRemoteURL(root string) string {
	out, err := exec.Command("git", "-C", root, "remote", "get-url", "origin").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
