
## Weeks 1–2
Before writing any production code, I want to properly 
understand what Phase 1 actually built. Read through the 
BYOA doc, ADR-002, sit with the codebase until it makes 
sense. Then write a design doc and get mentor sign-off.

## Weeks 3–5
Wire the skills engine into the agent. Discovery, 
validation, and dynamic context injection so loaded 
skills actually change how the agent behaves at runtime, 
not just sit in a registry doing nothing.

## Weeks 6–7
Natural language search and trace explanation that pulls 
from loaded skills instead of hardcoded logic. This is 
where the framework becomes actually useful.

## Weeks 8–9
Local model support. Ollama and Llama.cpp compatibility 
so skills work fully offline without sending trace data 
to public clouds.

## Weeks 10–11
React UI to make the reasoning visible — which skill 
fired, which tools the agent called, what it concluded. 
Right now that's all hidden. It shouldn't be.

## Week 12
Documentation guide for writing custom skills, with real 
debugging workflow examples. And whatever the mentors 
think needs more work.
