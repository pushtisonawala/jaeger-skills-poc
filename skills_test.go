package main

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/pushtisonawala/jaeger-skills-poc/skills"
)

func TestLoadFromDir(t *testing.T) {
	loaded, errs := skills.LoadFromDir("examples")
	if len(errs) > 0 {
		t.Fatalf("unexpected load errors: %v", errs)
	}
	if len(loaded) == 0 {
		t.Fatal("expected skills to be loaded from examples/")
	}
}

func TestValidSkill(t *testing.T) {
	skill := &skills.Skill{
		Name:           "test-skill",
		SystemPrompt:   "you are helpful",
		Tools:          []string{"search_traces", "get_trace"},
		MaxTokenBudget: 8000,
	}
	result := skills.Validate(skill)
	if !result.Valid {
		t.Errorf("expected valid skill, got errors: %v", result.Errors)
	}
}

func TestInvalidSkill_MissingName(t *testing.T) {
	skill := &skills.Skill{
		SystemPrompt: "you are helpful",
		Tools:        []string{"search_traces"},
	}
	result := skills.Validate(skill)
	if result.Valid {
		t.Error("expected invalid skill with missing name")
	}
}

func TestInvalidSkill_UnknownTool(t *testing.T) {
	skill := &skills.Skill{
		Name:         "bad-skill",
		SystemPrompt: "you are helpful",
		Tools:        []string{"fake_tool_that_doesnt_exist"},
	}
	result := skills.Validate(skill)
	if result.Valid {
		t.Error("expected invalid skill with unknown tool")
	}
}

func TestRegistry_ConflictDetection(t *testing.T) {
	r := skills.NewRegistry()
	s := &skills.Skill{Name: "my-skill", SystemPrompt: "x", Tools: []string{"get_trace"}}

	if err := r.Register(s); err != nil {
		t.Fatalf("first register should succeed: %v", err)
	}
	if err := r.Register(s); err == nil {
		t.Error("second register should fail with conflict error")
	}
}

func TestDryRun_NoRegistration(t *testing.T) {
	r := skills.NewRegistry()
	if len(r.All()) != 0 {
		t.Error("registry should be empty in dry-run mode")
	}
}

func writeTempSkill(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "skill.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp skill: %v", err)
	}
	return dir
}

func TestLoadFromDir_InvalidYAML(t *testing.T) {
	dir := writeTempSkill(t, "name: bad\nstep_pipeline: [")
	loaded, errs := skills.LoadFromDir(dir)
	if len(loaded) != 0 {
		t.Fatalf("expected no loaded skills, got %d", len(loaded))
	}
	if len(errs) == 0 {
		t.Fatal("expected YAML parsing error")
	}
}

func TestTokenBudgetWarning(t *testing.T) {
	skill := &skills.Skill{
		Name:         "no-budget",
		SystemPrompt: "test",
		Tools:        []string{"get_trace"},
	}
	result := skills.Validate(skill)
	if len(result.Warnings) == 0 {
		t.Error("expected warning for missing token budget")
	}
}

func TestFixtures(t *testing.T) {
	validData, err := os.ReadFile(filepath.Join("testdata", "valid-skill.yaml"))
	if err != nil {
		t.Fatalf("failed to read valid fixture: %v", err)
	}

	var validSkill skills.Skill
	if err := yaml.Unmarshal(validData, &validSkill); err != nil {
		t.Fatalf("failed to parse valid fixture: %v", err)
	}
	if result := skills.Validate(&validSkill); !result.Valid {
		t.Fatalf("expected valid fixture to pass validation, got: %v", result.Errors)
	}

	invalidData, err := os.ReadFile(filepath.Join("testdata", "invalid-skill.yaml"))
	if err != nil {
		t.Fatalf("failed to read invalid fixture: %v", err)
	}

	var invalidSkill skills.Skill
	if err := yaml.Unmarshal(invalidData, &invalidSkill); err != nil {
		t.Fatalf("failed to parse invalid fixture: %v", err)
	}
	if result := skills.Validate(&invalidSkill); result.Valid {
		t.Fatal("expected invalid fixture to fail validation")
	}
}
