package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Skill represents a user-defined debugging workflow
// that can be loaded into the Jaeger AI agent at runtime.
// This mirrors the BYOA architecture described in ADR-002.
type Skill struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Description string   `yaml:"description"`
	Author      string   `yaml:"author"`
	Tags        []string `yaml:"tags"`

	// Tools this skill is allowed to call.
	// Must be a subset of registered MCP tools.
	Tools []string `yaml:"tools"`

	// SystemPrompt is injected into the agent context
	// when this skill is activated.
	SystemPrompt string `yaml:"system_prompt"`

	// StepPipeline enforces call order.
	// The agent is expected to call tools in this sequence.
	StepPipeline []Step `yaml:"step_pipeline"`

	// MaxTokenBudget caps context window usage.
	// Protects against the large-trace token explosion
	// called out in the proposal.
	MaxTokenBudget int `yaml:"max_token_budget"`
}

type Step struct {
	Tool        string `yaml:"tool"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
}

// LoadFromDir discovers and loads all .yaml skill files
// from the given directory. Non-YAML files are ignored.
func LoadFromDir(dir string) ([]*Skill, []error) {
	var loaded []*Skill
	var errs []error

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, []error{fmt.Errorf("cannot read skills dir %s: %w", dir, err)}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		skill, err := loadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("file %s: %w", entry.Name(), err))
			continue
		}
		loaded = append(loaded, skill)
	}

	return loaded, errs
}

func loadFile(path string) (*Skill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	var skill Skill
	if err := yaml.Unmarshal(data, &skill); err != nil {
		return nil, fmt.Errorf("yaml parse error: %w", err)
	}

	return &skill, nil
}
