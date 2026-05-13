package skills

import (
	"fmt"
	"sort"
	"strings"
)

// Registry holds all loaded and validated skills.
// In the full BYOA implementation, this would be
// injected into the agent's tool dispatcher.
type Registry struct {
	skills map[string]*Skill
}

func NewRegistry() *Registry {
	return &Registry{skills: make(map[string]*Skill)}
}

// Register adds a validated skill to the registry.
// Returns an error if a skill with the same name already exists.
func (r *Registry) Register(s *Skill) error {
	if _, exists := r.skills[s.Name]; exists {
		return fmt.Errorf("conflict: skill %q already registered - names must be unique", s.Name)
	}
	r.skills[s.Name] = s
	return nil
}

// Get retrieves a skill by name.
func (r *Registry) Get(name string) (*Skill, bool) {
	s, ok := r.skills[name]
	return s, ok
}

// All returns all registered skills.
func (r *Registry) All() []*Skill {
	names := make([]string, 0, len(r.skills))
	for name := range r.skills {
		names = append(names, name)
	}
	sort.Strings(names)

	result := make([]*Skill, 0, len(names))
	for _, name := range names {
		result = append(result, r.skills[name])
	}
	return result
}

// Summary prints a human-readable overview.
func (r *Registry) Summary() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Registry: %d skill(s) loaded\n", len(r.skills))
	for _, s := range r.All() {
		fmt.Fprintf(&b, "  - %-30s tools: %v\n", s.Name, s.Tools)
	}
	return b.String()
}
