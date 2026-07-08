package sddmodel

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/agentmodel"
)

const dom = "development"

func globalWith(models map[string]string) *agentmodel.Config {
	cfg := &agentmodel.Config{}
	for agent, model := range models {
		cfg.SetDomainModelOverride(dom, agent, model)
	}
	return cfg
}

func defaults() map[string]string {
	return map[string]string{"a1": "opus", "a2": "haiku"}
}

// En scope proyecto sin override propio, la herencia debe mostrar el modelo
// elegido globalmente — no el default del payload.
func TestProject_InheritsGlobalOverride(t *testing.T) {
	m := newModel(globalWith(map[string]string{"a1": "sonnet"}), nil, "", agentmodel.ScopeProject, dom, []string{"a1", "a2"}, defaults())

	if got := m.models["a1"]; got != GlobalSentinel {
		t.Fatalf("a1 debería quedar como sentinel (sin override de proyecto), got %q", got)
	}
	pill := m.renderModelPills("a1", GlobalSentinel)
	if !strings.Contains(pill, "hereda: sonnet") {
		t.Errorf("herencia debería reflejar el global (sonnet), got %q", pill)
	}
	if strings.Contains(pill, "opus") {
		t.Errorf("no debería mostrar el default del payload (opus) cuando el global fija sonnet, got %q", pill)
	}
}

// Sin override global, la herencia cae al default real del payload por agente.
func TestProject_InheritsPayloadDefault(t *testing.T) {
	m := newModel(&agentmodel.Config{}, nil, "", agentmodel.ScopeProject, dom, []string{"a2"}, defaults())

	pill := m.renderModelPills("a2", GlobalSentinel)
	if !strings.Contains(pill, "hereda: haiku") {
		t.Errorf("herencia debería caer al default del payload (haiku), got %q", pill)
	}
}

// En scope global, un agente sin entrada se siembra con su default real del
// payload — no con un sonnet hardcodeado.
func TestGlobal_SeedsPayloadDefaultNotSonnet(t *testing.T) {
	m := newModel(globalWith(map[string]string{"a1": "fable"}), nil, "", agentmodel.ScopeGlobal, dom, []string{"a1", "a2"}, defaults())

	if got := m.models["a1"]; got != "fable" {
		t.Errorf("a1 debería conservar su override global (fable), got %q", got)
	}
	if got := m.models["a2"]; got != "haiku" {
		t.Errorf("a2 debería sembrarse con su default del payload (haiku), got %q", got)
	}
}

// Guardar en scope global no debe persistir overrides iguales al default del
// payload (evita contaminar la herencia de los proyectos), pero sí los que
// difieren.
func TestGlobal_SaveDropsDefaultEqualOverrides(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	m := newModel(&agentmodel.Config{}, nil, path, agentmodel.ScopeGlobal, dom, []string{"a1", "a2"}, defaults())

	m.models["a1"] = "opus"   // igual al default → no debe persistirse
	m.models["a2"] = "sonnet" // distinto del default (haiku) → debe persistirse

	if msg := m.saveAndBack()(); msg == nil {
		t.Fatal("saveAndBack no debería fallar")
	}

	saved, err := agentmodel.Load(path)
	if err != nil {
		t.Fatalf("Load config guardado: %v", err)
	}
	models := saved.DomainModels(dom)
	if _, ok := models["a1"]; ok {
		t.Errorf("a1 (igual al default) no debería persistirse, got %q", models["a1"])
	}
	if models["a2"] != "sonnet" {
		t.Errorf("a2 (distinto del default) debería persistirse como sonnet, got %q", models["a2"])
	}
}
