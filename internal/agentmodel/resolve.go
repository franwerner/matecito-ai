package agentmodel

import (
	"path/filepath"
	"strings"
)

// ResolveModel returns the effective model for agent given two configs.
// source is one of "project", "global", "default".
// ("", "default") means neither config sets the agent; caller omits the param.
// project == nil means no per-project config file was found.
// Unknown agents (not in Agents) always return ("", "default").
func ResolveModel(global *Config, project *Config, agent string) (model string, source string) {
	if !IsValidAgent(agent) {
		return "", "default"
	}
	if project != nil && project.Models[agent] != "" {
		return project.Models[agent], "project"
	}
	if global != nil && global.Models[agent] != "" {
		return global.Models[agent], "global"
	}
	return "", "default"
}

// ResolveTdd determines the effective strict-TDD mode given the global and per-project configs.
// Precedence: per-project (file present AND strictTdd key set) → global (key set) → false.
// project == nil means no per-project config file was found.
func ResolveTdd(global Config, project *Config) bool {
	if project != nil && project.StrictTdd != nil {
		return *project.StrictTdd
	}
	if global.StrictTdd != nil {
		return *global.StrictTdd
	}
	return false
}

// DeriveProjectName returns the repository name from a git remote URL, or
// filepath.Base(dir) when remoteURL is empty. Pure; never shells out.
// Handles HTTPS (https://github.com/owner/repo.git) and
// SSH (git@github.com:owner/repo.git) formats.
func DeriveProjectName(remoteURL, dir string) string {
	if remoteURL == "" {
		return filepath.Base(filepath.Clean(dir))
	}

	// strip trailing .git
	url := strings.TrimSuffix(remoteURL, ".git")

	// SSH format: git@github.com:owner/repo  → take segment after last ':' or '/'
	if strings.HasPrefix(url, "git@") {
		// "git@github.com:owner/repo" — split on ':'
		if idx := strings.LastIndex(url, ":"); idx >= 0 {
			url = url[idx+1:]
		}
	}

	// HTTPS or post-SSH colon strip: take basename
	if idx := strings.LastIndex(url, "/"); idx >= 0 {
		url = url[idx+1:]
	}

	return url
}

// IsDevBuild reports whether version contains the "-dev" suffix indicator.
func IsDevBuild(version string) bool {
	return strings.Contains(version, "-dev")
}

// NormalizeVersion strips a leading "v" from version strings (e.g. "v0.2.0" → "0.2.0").
func NormalizeVersion(v string) string {
	return strings.TrimPrefix(v, "v")
}

// ShouldShowBadge returns true when a release badge should be displayed.
// Never badges dev builds, empty latest tags, or when versions are equal after normalization.
func ShouldShowBadge(current, latest string) bool {
	if IsDevBuild(current) {
		return false
	}
	if latest == "" {
		return false
	}
	return NormalizeVersion(current) != NormalizeVersion(latest)
}
