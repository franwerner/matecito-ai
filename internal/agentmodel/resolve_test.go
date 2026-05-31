package agentmodel_test

import (
	"testing"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
)

// --- ResolveTdd ---

func TestResolveTdd(t *testing.T) {
	boolPtr := func(v bool) *bool { return &v }

	tests := []struct {
		name    string
		global  agentmodel.Config
		project *agentmodel.Config
		want    bool
	}{
		{
			// S5.9: per-project key present wins over global
			name:    "perproject_wins",
			global:  agentmodel.Config{StrictTdd: boolPtr(false)},
			project: &agentmodel.Config{StrictTdd: boolPtr(true)},
			want:    true,
		},
		{
			// S5.10: no per-project file → global fallback
			name:    "global_fallback",
			global:  agentmodel.Config{StrictTdd: boolPtr(true)},
			project: nil,
			want:    true,
		},
		{
			// S5.11: neither file → default false
			name:    "default_false",
			global:  agentmodel.Config{},
			project: nil,
			want:    false,
		},
		{
			// S8.5: per-project file present but key nil → falls back to global
			name:    "perproject_key_absent_uses_global",
			global:  agentmodel.Config{StrictTdd: boolPtr(true)},
			project: &agentmodel.Config{StrictTdd: nil},
			want:    true,
		},
		{
			// global also nil → false
			name:    "both_keys_absent",
			global:  agentmodel.Config{StrictTdd: nil},
			project: &agentmodel.Config{StrictTdd: nil},
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := agentmodel.ResolveTdd(tc.global, tc.project)
			if got != tc.want {
				t.Errorf("ResolveTdd() = %v, want %v", got, tc.want)
			}
		})
	}
}

// --- DeriveProjectName ---

func TestDeriveProjectName(t *testing.T) {
	tests := []struct {
		name      string
		remoteURL string
		dir       string
		want      string
	}{
		{
			// S6.1: HTTPS remote
			name:      "https_remote",
			remoteURL: "https://github.com/franwerner/matecito-ai.git",
			dir:       "/home/user/matecito-ai",
			want:      "matecito-ai",
		},
		{
			// S6.2: SSH remote
			name:      "ssh_remote",
			remoteURL: "git@github.com:franwerner/matecito-ai.git",
			dir:       "/home/user/matecito-ai",
			want:      "matecito-ai",
		},
		{
			// S6.3: no remote → fallback to dir name
			name:      "no_remote_dir_fallback",
			remoteURL: "",
			dir:       "/home/user/my-project",
			want:      "my-project",
		},
		{
			// HTTPS without .git suffix
			name:      "https_no_git_suffix",
			remoteURL: "https://github.com/franwerner/my-repo",
			dir:       "/home/user/x",
			want:      "my-repo",
		},
		{
			// trailing slash in dir
			name:      "dir_trailing_slash",
			remoteURL: "",
			dir:       "/home/user/my-project/",
			want:      "my-project",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := agentmodel.DeriveProjectName(tc.remoteURL, tc.dir)
			if got != tc.want {
				t.Errorf("DeriveProjectName(%q, %q) = %q, want %q", tc.remoteURL, tc.dir, got, tc.want)
			}
		})
	}
}

// --- IsDevBuild ---

func TestIsDevBuild(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"0.1.0-dev", true}, // S7.1
		{"0.2.0", false},    // S7.1
		{"1.0.0-dev", true},
		{"0.1.0-devXYZ", true}, // any -dev suffix
		{"", false},
		{"dev", false}, // no hyphen
		{"0.1.0-beta", false},
	}

	for _, tc := range tests {
		got := agentmodel.IsDevBuild(tc.version)
		if got != tc.want {
			t.Errorf("IsDevBuild(%q) = %v, want %v", tc.version, got, tc.want)
		}
	}
}

// --- NormalizeVersion ---

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"v0.2.0", "0.2.0"},
		{"0.2.0", "0.2.0"},
		{"v1.0.0-dev", "1.0.0-dev"},
		{"", ""},
	}
	for _, tc := range tests {
		got := agentmodel.NormalizeVersion(tc.input)
		if got != tc.want {
			t.Errorf("NormalizeVersion(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// --- ShouldShowBadge ---

func TestShouldShowBadge(t *testing.T) {
	tests := []struct {
		name    string
		current string
		latest  string
		want    bool
	}{
		{
			// S7.2: newer tag shows badge
			name:    "newer_tag",
			current: "0.2.0",
			latest:  "v0.3.0",
			want:    true,
		},
		{
			// S7.3: same version no badge
			name:    "same_version",
			current: "0.2.0",
			latest:  "0.2.0",
			want:    false,
		},
		{
			// S7.6: v-prefix normalization → same → no badge
			name:    "v_prefix_equal",
			current: "0.2.0",
			latest:  "v0.2.0",
			want:    false,
		},
		{
			// S7.1: dev build never badges
			name:    "dev_build",
			current: "0.1.0-dev",
			latest:  "v0.3.0",
			want:    false,
		},
		{
			// empty latest → no badge
			name:    "empty_latest",
			current: "0.2.0",
			latest:  "",
			want:    false,
		},
		{
			// both empty → no badge
			name:    "both_empty",
			current: "",
			latest:  "",
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := agentmodel.ShouldShowBadge(tc.current, tc.latest)
			if got != tc.want {
				t.Errorf("ShouldShowBadge(%q, %q) = %v, want %v", tc.current, tc.latest, got, tc.want)
			}
		})
	}
}
