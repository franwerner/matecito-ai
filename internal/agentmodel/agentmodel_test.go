package agentmodel_test

import (
	"testing"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
)

// --- ReadModel ---

func TestReadModel(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			// S4.1: model line present in frontmatter
			name:  "present",
			input: "---\nname: sdd-spec\nmodel: sonnet\n---\n",
			want:  "sonnet",
		},
		{
			// S4.2: model line absent from frontmatter
			name:  "absent",
			input: "---\nname: sdd-spec\n---\n",
			want:  "",
		},
		{
			// S4.3: indented model line must be ignored
			name:  "indented_decoy",
			input: "---\nname: sdd-spec\n  model: sonnet\n---\n",
			want:  "",
		},
		{
			// S4.4: model line outside the first frontmatter block must be ignored
			name:  "outside_frontmatter",
			input: "---\nname: sdd-spec\n---\nbody\nmodel: opus\n",
			want:  "",
		},
		{
			// no frontmatter at all
			name:  "no_frontmatter",
			input: "just text\nmodel: opus\n",
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := agentmodel.ReadModel([]byte(tc.input))
			if (err != nil) != tc.wantErr {
				t.Fatalf("ReadModel() error = %v, wantErr %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("ReadModel() = %q, want %q", got, tc.want)
			}
		})
	}
}

// --- IsValidModel / IsValidAgent ---

func TestIsValidModel(t *testing.T) {
	for _, m := range []string{"opus", "sonnet", "haiku"} {
		if !agentmodel.IsValidModel(m) {
			t.Errorf("IsValidModel(%q) = false, want true", m)
		}
	}
	for _, m := range []string{"gpt-4", "", "OPUS", "claude"} {
		if agentmodel.IsValidModel(m) {
			t.Errorf("IsValidModel(%q) = true, want false", m)
		}
	}
}

func TestIsValidAgent(t *testing.T) {
	for _, a := range agentmodel.Agents {
		if !agentmodel.IsValidAgent(a) {
			t.Errorf("IsValidAgent(%q) = false, want true", a)
		}
	}
	for _, a := range []string{"sdd-unknown", "", "intake", "sdd-INTAKE"} {
		if agentmodel.IsValidAgent(a) {
			t.Errorf("IsValidAgent(%q) = true, want false", a)
		}
	}
}

func TestAgentsSlice(t *testing.T) {
	want := []string{
		"sdd-intake", "sdd-explore", "sdd-propose", "sdd-spec",
		"sdd-design", "sdd-tasks", "sdd-apply", "sdd-verify", "sdd-archive",
	}
	if len(agentmodel.Agents) != len(want) {
		t.Fatalf("Agents length = %d, want %d", len(agentmodel.Agents), len(want))
	}
	for i, a := range want {
		if agentmodel.Agents[i] != a {
			t.Errorf("Agents[%d] = %q, want %q", i, agentmodel.Agents[i], a)
		}
	}
}
