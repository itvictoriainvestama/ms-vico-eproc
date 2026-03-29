Unit Testing Workflow (Per Function / Per Module) — Ultra Modern (Rating 1–10)
Instruction

Create unit tests that validate correctness of logic at the function/class/module level. Unit tests must be:

repo-aware (match existing test stack and conventions),

idiomatic to the detected language/framework,

deterministic (no flaky time/random/network),

fast (no real DB/network/services unless explicitly classified as integration tests),

focused on business logic and edge cases.

Unit tests are NOT E2E.
They should isolate the unit under test by controlling dependencies via mocks/fakes/stubs only at the boundary (I/O, DB, network, filesystem, time).

0) Mandatory Recon (Always First)
0.1 Detect Stack & Test Tooling

The agent MUST inspect the repo to identify:

language(s), framework(s)

test runner and assertion library (e.g., Jest/Vitest/Mocha, Go test, JUnit, PHPUnit, etc.)

mocking strategy used in repo (mocks vs fakes vs DI)

file naming conventions (*.spec.*, *_test.go, etc.)

test folder conventions (test/, __tests__/, etc.)

existing test utilities and helpers

Rule: Reuse existing tools and patterns. Do not introduce a new test framework unless explicitly required.

0.2 Identify Unit Boundaries

The agent MUST classify components into:

pure logic (ideal unit test targets)

orchestrators (services/use-cases needing dependency mocks)

I/O boundaries (DB, HTTP, queue, fs) → unit tests should mock these

0.3 Map What “Module” Means in This Repo

Define module boundaries based on codebase:

package/module folder boundaries

domain module boundaries

public API surface of the module

1) Scope of Unit Testing
1.1 What to Unit Test (Priority Order)

Pure functions / domain logic

Business rules (validation beyond schema, invariants, calculations)

Services/use-cases (orchestration, branching, error mapping)

Adapters/wrappers (only logic inside them; avoid real network/DB)

1.2 What NOT to Unit Test

UI rendering details best covered by component tests (unless “unit test” in that ecosystem includes it)

database correctness (integration tests)

full HTTP request lifecycle (integration/e2e)

third-party behavior (mock at boundary)

2) Testing Strategy (Modern, Deterministic)
2.1 Test Design Rules

Each unit test must follow:

Arrange: setup inputs and controlled dependencies

Act: call unit

Assert: verify output + side effects (only within controlled boundaries)

Assertions must verify:

returned values (or thrown errors)

state changes (only within unit)

calls to dependencies (only for orchestration units, and only meaningful ones)

2.2 Determinism Rules (Mandatory)

No real time dependencies:

inject a clock/time provider or mock time

No randomness:

inject RNG or seed it deterministically

No network/DB:

use mocks/fakes

No filesystem (unless unit is fs-related; then use temp dirs or mocks)

2.3 Coverage Approach (Risk-Based)

Prioritize unit tests for:

high branching complexity

money/PII logic

permission/business rule checks

state transitions and edge-case prone logic

historically bug-prone modules (if commit history indicates)

3) Scenarios to Cover (Per Function/Class/Module)
3.1 Happy Path

valid input produces expected output

expected interactions happen (if orchestrator)

3.2 Edge Cases

empty/null/undefined cases (as applicable)

min/max boundaries

unusual but valid data

rounding/precision (numeric logic)

encoding/timezone (if relevant)

3.3 Error Paths

invalid input rejected with correct error type/code

dependency failure mapped correctly (e.g., repo throws → service returns domain error)

forbidden operations return proper error (if authz logic is inside unit)

retries/backoff logic behaves deterministically (if present)

3.4 State Transitions / Invariants (If Applicable)

allowed transitions succeed

invalid transitions fail predictably

invariants are enforced consistently

4) Module-Level Unit Test Requirements
4.1 Public API Coverage

For each module:

list public functions/classes

ensure each has at least:

1 happy path test

1 edge case test

1 error path test (where applicable)

4.2 Contract Tests (Unit-Level)

If module exposes DTO-like contracts internally:

verify schema validation behavior (if kept within module boundary)

verify error mapping consistency

5) Mocks / Fakes / Stubs Policy (Modern Best Practice)
5.1 Prefer Fakes over Excessive Mocks

Use fakes for repositories/services when behavior matters.

Use mocks mainly to verify orchestration and interactions.

Avoid mocking internal implementation details.

5.2 Mock at Boundaries Only

DB, HTTP clients, queues, filesystem, time, randomness

Do not mock core business logic collaborators unless necessary

5.3 Avoid Over-Specification

Don’t assert exact call ordering unless required.

Don’t assert internal calls that are likely to change with refactor.

Focus on outputs and meaningful side effects.

6) Test Code Quality Standards
6.1 Readability & Naming

tests should read like documentation:

should_return_X_when_Y

returns_error_when_invalid_input

avoid massive setup blocks; use helpers/builders

6.2 Test Data Builders & Fixtures

create small deterministic builders:

makeUser(), makeOrder()

avoid huge JSON blobs unless necessary

keep fixtures close to tests unless shared widely

6.3 Maintainability

one reason to change per test file

minimize duplication via builders/helpers

keep tests independent

7) Output Requirements (What the Agent Must Produce)

The answer MUST include:

Recon summary (stack, test runner, conventions)

Target list of functions/classes/modules to unit test (prioritized)

For each target:

test cases list (happy/edge/error)

key assertions

dependency control approach (mock/fake/DI)

Test structure plan (folders, naming, helpers)

Pass/Fail reporting structure (how results will be presented)

Issues found (if any design makes unit testing hard) + recommended refactor seams

Scoring table (rubric)

Final rating: Rating: X/10

Structured Unit Test Quality Rubric (Score 1–10)
Aspect	Score (1–10)
Recon completeness (tooling + conventions understood)	
Coverage of core logic (happy paths)	
Coverage of edge cases	
Coverage of error paths	
Determinism (no flaky time/random/network)	
Isolation quality (boundaries mocked correctly)	
Quality of assertions (meaningful, not brittle)	
Test readability and maintainability	
Module public API coverage completeness	
Use of builders/fixtures/helpers appropriately	
Idiomatic use of language/framework testing style	
Speed/efficiency of the unit test suite	
Testability feedback (seams/refactor suggestions)	
Overall coherence of the unit testing plan	

Final Rating = average of all aspect scores
Output format: Rating: X/10

Rating Scale Meaning

9–10 = Excellent: comprehensive, deterministic, idiomatic, maintainable, fast

8 = Very good: minor gaps in edges/errors or structure

7 = Good: some gaps in coverage or brittle tests exist

6 = Acceptable: lacks depth in edges/errors or isolation

5 = Borderline: needs significant redesign of tests or code seams

3–4 = Weak: misses core logic or tests are flaky/brittle

1–2 = Very poor: unreliable, unmaintainable

0 = Broken / Unusable