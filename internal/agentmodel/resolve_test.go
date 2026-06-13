package agentmodel_test

import (
	"testing"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
)

// --- ResolveModel ---

// devModels builds a config with per-agent model overrides under the default domain.
func devModels(m map[string]string) *agentmodel.Config {
	return &agentmodel.Config{DomainConfig: map[string]*agentmodel.DomainConfig{
		agentmodel.DefaultDomain: {Models: m},
	}}
}

func TestResolveModel(t *testing.T) {
	tests := []struct {
		name       string
		global     *agentmodel.Config
		project    *agentmodel.Config
		agent      string
		wantModel  string
		wantSource string
	}{
		{
			name:       "project_wins",
			global:     devModels(map[string]string{"sdd-design": "sonnet"}),
			project:    devModels(map[string]string{"sdd-design": "haiku"}),
			agent:      "sdd-design",
			wantModel:  "haiku",
			wantSource: "project",
		},
		{
			name:       "global_fallback",
			global:     devModels(map[string]string{"sdd-design": "sonnet"}),
			project:    nil,
			agent:      "sdd-design",
			wantModel:  "sonnet",
			wantSource: "global",
		},
		{
			name:       "default_when_neither",
			global:     devModels(map[string]string{}),
			project:    nil,
			agent:      "sdd-spec",
			wantModel:  "",
			wantSource: "default",
		},
		{
			name:       "unknown_agent",
			global:     devModels(map[string]string{"sdd-design": "sonnet"}),
			project:    devModels(map[string]string{"sdd-design": "haiku"}),
			agent:      "nonexistent-agent",
			wantModel:  "",
			wantSource: "default",
		},
		{
			name:       "both_nil",
			global:     nil,
			project:    nil,
			agent:      "sdd-apply",
			wantModel:  "",
			wantSource: "default",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotModel, gotSource := agentmodel.ResolveModel(tc.global, tc.project, agentmodel.DefaultDomain, tc.agent)
			if gotModel != tc.wantModel {
				t.Errorf("ResolveModel() model = %q, want %q", gotModel, tc.wantModel)
			}
			if gotSource != tc.wantSource {
				t.Errorf("ResolveModel() source = %q, want %q", gotSource, tc.wantSource)
			}
		})
	}
}

// TestResolveModel_MigratedLegacy verifies a flat (pre-M7) config still resolves
// for development after Load normalizes it into the default domain.
func TestResolveModel_MigratedLegacy(t *testing.T) {
	legacy := &agentmodel.Config{Models: map[string]string{"sdd-design": "opus"}}
	// normalize is invoked by Load; simulate by round-tripping the accessor path
	// through ResolveModel after a normalize via Validate (which calls normalize).
	_ = agentmodel.Validate(legacy)
	got, src := agentmodel.ResolveModel(legacy, nil, agentmodel.DefaultDomain, "sdd-design")
	if got != "opus" || src != "global" {
		t.Errorf("migrated legacy resolve = (%q,%q), want (opus,global)", got, src)
	}
}

// --- ResolveTdd ---

// devTdd builds a value config with strictTdd set under the default domain.
func devTdd(v *bool) agentmodel.Config {
	return agentmodel.Config{DomainConfig: map[string]*agentmodel.DomainConfig{
		agentmodel.DefaultDomain: {StrictTdd: v},
	}}
}

// devTddP is the pointer variant for per-project configs.
func devTddP(v *bool) *agentmodel.Config {
	c := devTdd(v)
	return &c
}

func TestResolveTdd(t *testing.T) {
	boolPtr := func(v bool) *bool { return &v }

	tests := []struct {
		name    string
		global  agentmodel.Config
		project *agentmodel.Config
		want    bool
	}{
		{
			name:    "perproject_wins",
			global:  devTdd(boolPtr(false)),
			project: devTddP(boolPtr(true)),
			want:    true,
		},
		{
			name:    "global_fallback",
			global:  devTdd(boolPtr(true)),
			project: nil,
			want:    true,
		},
		{
			name:    "default_false",
			global:  agentmodel.Config{},
			project: nil,
			want:    false,
		},
		{
			name:    "perproject_key_absent_uses_global",
			global:  devTdd(boolPtr(true)),
			project: devTddP(nil),
			want:    true,
		},
		{
			name:    "both_keys_absent",
			global:  devTdd(nil),
			project: devTddP(nil),
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := agentmodel.ResolveTdd(tc.global, tc.project, agentmodel.DefaultDomain)
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
