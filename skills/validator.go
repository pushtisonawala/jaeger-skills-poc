package skills

import (
	"fmt"
	"strings"
)

// KnownMCPTools are the tools registered in the Jaeger MCP server.
// From ADR-002: search_traces, get_trace, get_trace_topology,
// get_span_details, get_services, get_operations.
var KnownMCPTools = map[string]bool{
	"search_traces":      true,
	"get_trace":          true,
	"get_trace_topology": true,
	"get_span_details":   true,
	"get_services":       true,
	"get_operations":     true,
}

type ValidationResult struct {
	SkillName string
	Valid     bool
	Errors    []string
	Warnings  []string
}

// Validate checks a skill for correctness.
// Returns errors (blocking) and warnings (informational).
func Validate(s *Skill) ValidationResult {
	result := ValidationResult{
		SkillName: s.Name,
		Valid:     true,
	}

	if strings.TrimSpace(s.Name) == "" {
		result.Errors = append(result.Errors, "missing required field: name")
	}
	if strings.TrimSpace(s.SystemPrompt) == "" {
		result.Errors = append(result.Errors, "missing required field: system_prompt")
	}
	if len(s.Tools) == 0 {
		result.Errors = append(result.Errors, "skill must declare at least one tool")
	}

	for _, tool := range s.Tools {
		if !KnownMCPTools[tool] {
			result.Errors = append(result.Errors,
				fmt.Sprintf("unknown tool %q - not registered in Jaeger MCP server", tool))
		}
	}

	if len(s.StepPipeline) > 0 {
		for i, step := range s.StepPipeline {
			if !KnownMCPTools[step.Tool] {
				result.Errors = append(result.Errors,
					fmt.Sprintf("step[%d]: unknown tool %q in pipeline", i, step.Tool))
			}
		}

		first := s.StepPipeline[0].Tool
		if first != "search_traces" && first != "get_trace_topology" {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("step_pipeline starts with %q - consider starting with search_traces or get_trace_topology per ADR-002 progressive disclosure pattern", first))
		}
	}

	if s.MaxTokenBudget == 0 {
		result.Warnings = append(result.Warnings,
			"no max_token_budget set - risk of 53K+ token context from get_trace on large traces")
	} else if s.MaxTokenBudget > 50000 {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("max_token_budget %d is very high - large traces may exceed this in production", s.MaxTokenBudget))
	}

	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result
}
