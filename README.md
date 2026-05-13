# jaeger-skills-poc

A working prototype of the **BYOA Skills Framework** proposed for Jaeger Phase 2 LFX Mentorship (Term 2 2026).

Built to demonstrate the skills loading layer before the full implementation, so the design can be validated early and discussed concretely.

## What it does

Users drop `.yaml` skill files into a directory.
The engine discovers, validates, and registers them at startup with no recompile required.

```bash
go run . --skills-dir ./examples
```

Example output:

```text
+------------------------------------------+
|     Jaeger Skills Engine - PoC           |
|     BYOA Skills Framework Prototype      |
+------------------------------------------+

Discovered 3 skill file(s) in ./examples

  OK analyze-critical-path               [VALID]
  OK detect-n-plus-one                   [VALID]
  OK find-slow-services                  [VALID]
      WARNING: step_pipeline starts with "get_services" - consider starting with search_traces or get_trace_topology per ADR-002 progressive disclosure pattern

-- Registry --
Registry: 3 skill(s) loaded
  - analyze-critical-path          tools: [search_traces get_trace_topology get_span_details]
  - detect-n-plus-one              tools: [search_traces get_trace_topology get_span_details]
  - find-slow-services             tools: [search_traces get_trace_topology get_services]

Result: OK - all skills loaded successfully
```

## Design decisions aligned with Jaeger ADR-002

### Progressive disclosure enforcement

The validator warns when a skill's `step_pipeline` does not start with `search_traces` or `get_trace_topology`. This reflects the progressive disclosure pattern from ADR-002: the LLM should avoid reaching for expensive detailed span data before first narrowing scope with search or topology.

### Token budget protection

The `max_token_budget` field addresses the large-context problem that can happen with `get_trace` on large traces. Skills without a budget generate a warning so maintainers can spot risky defaults early.

### MCP tool validation

Skills can only reference tools registered in the Jaeger MCP server:

- `search_traces`
- `get_trace`
- `get_trace_topology`
- `get_span_details`
- `get_services`
- `get_operations`

Unknown tools are rejected.

### Conflict detection

Two skills with the same name cannot coexist in the registry. The engine fails fast on conflicts instead of silently overwriting one definition with another.

## Dry-run mode

```bash
go run . --skills-dir ./examples --dry-run
```

This validates all skills without registering them. It is useful for CI or pre-merge checks.

## How this fits into the full BYOA architecture

```text
skills/ or examples/              <- users drop YAML files here
          |
          v
SkillLoader (this PoC)            <- discovers and parses skills
          |
          v
Validator                         <- enforces schema, tool safety, budget hints
          |
          v
SkillRegistry                     <- stores validated skills, detects conflicts
          |
          v
Agent Tool Dispatcher             <- injects skill context into the BYOA agent
          |
          v
MCP Server                        <- executes approved tool calls
          |
          v
Jaeger Query Service              <- returns trace data
```

## Included example skills

| Skill | Focus | Tool flow |
| --- | --- | --- |
| `analyze-critical-path` | critical path analysis | search -> topology -> span details |
| `detect-n-plus-one` | ORM / query anti-pattern detection | search -> topology -> span details |
| `find-slow-services` | service ranking by latency contribution | services -> search -> topology |

## Running tests

```bash
go test ./...
```

## Why this stands out

This is more than a YAML loader:

- Runtime discovery with no rebuild
- Validation against known MCP tools
- Conflict detection for duplicate skills
- Dry-run mode for CI workflows
- Example skills that show how the mechanism plugs into the broader BYOA architecture

## Built by

Pushti Sonawala  
LFX Term 2 2026 applicant  
Jaeger AI-Powered Trace Analysis Phase 2 - Skills Framework
