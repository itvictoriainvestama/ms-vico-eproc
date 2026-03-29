Code Audit — Ultra Comprehensive (Code Quality) + Rating 1–10
Instruction

Perform a code audit focusing strictly on code quality (not DevOps, infrastructure, or process).
The audit must be language-agnostic and framework-agnostic, adapting to whatever the codebase uses (backend/frontend, SPA/SSR, Go/Java/JS/PHP/etc.).

Mandatory sequence (do not skip):

Study the codebase structure first (Recon).

Identify language/framework/tooling and existing conventions.

Evaluate quality using the rubric.

Produce prioritized findings and a minimal-diff improvement plan.

Provide scoring and final rating.

Constraints:

Prefer additive changes.

Keep diffs minimal.

Don’t remove working behavior unless explicitly requested.

0) Mandatory Pre-Work: Codebase Recon (Do This First)
0.1 Repository Topology & Module Map

Identify:

mono-repo vs multi-repo

apps/packages/libs layout

key modules and how they depend on each other

runtime entry points and boundaries (API layer, UI layer, domain layer, data layer)

0.2 Language/Framework/Tooling Detection

Determine:

programming language(s)

frameworks (backend/frontend/SSR)

build tooling

formatting/linting tools and code style conventions

test tooling and current test coverage patterns

typing model (strict types vs dynamic)

0.3 Architectural Pattern Recognition (Code-Level Only)

From code (not diagrams), infer:

layering pattern (controller/service/repo, clean architecture, MVC, etc.)

dependency direction (does domain depend on infra?)

shared libraries usage and boundaries

data flow patterns (events, queues, synchronous calls)

0.4 “Quality Boundary Map”

Document where the codebase enforces:

input validation boundaries

error handling boundaries

business rules location (service/domain layer)

infrastructure boundaries (db/http/fs)

shared utilities boundaries

Rule: Do not assume. Trace actual imports and call paths.

Scope of Evaluation (Code Quality Dimensions)
1) Core Fundamentals

DRY (Don’t Repeat Yourself) – no duplicated logic/knowledge

SOLID – SRP/OCP/LSP/ISP/DIP applied consistently

Big-O / Algorithmic Complexity – time & space efficiency is reasonable for expected scale

2) Additional Quality Aspects

Readability & clarity

naming clarity, intention-revealing code

reasonable function sizes and modularity

Simplicity (KISS, YAGNI)

avoid over-engineering

minimal abstractions that still solve the problem cleanly

Function/class cohesion

single purpose per module

avoid god objects and “utility dumping grounds”

Immutability & side effects

reduce shared mutable state

pure functions where feasible

Error handling discipline

consistent error model and propagation

no swallowed exceptions

Null-safety & edge cases

safe handling of null/undefined/empty/0 boundaries

API/function contracts & parameter sanity

clear preconditions/postconditions

avoid long parameter lists and flag arguments

Consistency of patterns & naming

uniform conventions across modules

Dependency hygiene & circular dependencies

clean dependency graph

minimal coupling, avoid cycles

Concurrency safety (where relevant)

race condition avoidance, safe async flows

Memory/resource safety

close resources, avoid leaks, bounded collections

Idiomatic usage

align with language/framework best practices

Testability (from design)

seams for DI/mocking

deterministic business logic

Code smells

long method, god class, long params, deep if/else, duplicated logic, primitive obsession, feature envy

3) Findings Requirements (Traceable + Actionable)

For each finding, include:

ID (stable)

Category (DRY/SOLID/Perf/Readability/etc.)

Severity (High/Medium/Low) — severity = impact on maintainability/bug risk

Location (file/module/function/class)

Evidence (what in code indicates the issue)

Why it matters (risk: bugs, cost, scaling, readability)

Recommendation (minimal-diff fix)

Effort (S/M/L)

Regression risk (Low/Med/High)

Verification (how to confirm improvement: tests/static checks)

Prioritization rule:
Sort findings by:

Severity

Breadth (how many areas affected)

Effort (quick wins first when impact is high)

4) Minimal-Diff Improvement Plan (Phased)

Provide a structured plan:

4.1 Quick Wins (Low Effort / High Impact)

naming improvements

extracting small pure helpers

removing obvious duplication

adding guards for edge cases

4.2 Medium Refactors (Controlled Scope)

reduce function size

improve layering boundaries

introduce interfaces/abstractions only where justified

replace fragile patterns with consistent ones

4.3 Strategic Refactors (Only If Necessary)

larger reorganizations must be proposed as optional

include migration strategy and risk containment

do not proceed unless requested

5) Optional Static Checks (Code-Only)

If repo supports it, recommend:

lint rules to prevent regressions

formatting consistency tools

dependency cycle detection

type checks

complexity thresholds (cyclomatic complexity, max params, etc.)

Rule: Recommend only tools that align with repo ecosystem.

Structured Scoring Rubric (Score Each Aspect 1–10)
Aspect	Score (1–10)
DRY	
SOLID	
Big-O & algorithmic complexity	
Readability & clarity	
Simplicity (KISS/YAGNI)	
Function/class cohesion	
Immutability & side effects	
Error handling	
Null-safety & edge cases	
API/function contracts	
Consistency	
Dependency hygiene	
Concurrency safety	
Memory/resource safety	
Idiomatic usage	
Testability	
Code smells overall	
Recon quality (stack + conventions understood)	
Findings quality (evidence + actionable fixes)	
Improvement plan quality (minimal diff + phased)	

Final score calculation:

Final Rating = average of all aspect scores (1–10)
Output format: Rating: X/10

Rating Scale Meaning

9–10 = excellent, production-grade, strong best practices

8 = very good, only minor improvements needed

7 = good, but some structural issues exist

6 = acceptable, but many important areas need improvement

5 = borderline, major refactors required

3–4 = low quality, many critical code smells

1–2 = very poor, not maintainable

0 = broken / unreadable / heavy anti-patterns

Expected Output from the AI Agent

The answer must include:

Recon summary (stack/frameworks/conventions + module map)

Short summary of overall code quality

Main strengths

Main issues found (traceable findings list)

Concrete recommendations + phased minimal-diff plan

Scoring table per category

Final rating:

Rating: X/10