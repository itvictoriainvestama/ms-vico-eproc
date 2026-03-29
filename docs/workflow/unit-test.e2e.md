End-to-End (E2E) Testing Workflow — Ultra Comprehensive (Rating 1–10)
Instruction

Perform End-to-End (E2E) testing by simulating real user behavior across the entire application stack, verifying that Frontend, Backend, Database, and Integrations work together correctly from a real user perspective.

Hard Rules (Non-negotiable)

Tests must run against the actual running application (frontend and/or backend).

Tests must use a real database in a deterministic state:

either seeded test data in a real DB, and/or

data created during tests and verified end-to-end.

No mocks/stubs for internal services (anything inside the system boundary).

External third-party integrations should use:

provider sandbox environments whenever possible, OR

controlled test endpoints; stubs only as a last resort and must be justified.

Assertions must be based on observable behavior and persisted state, not internal implementation details.

The suite must be CI-capable: deterministic, isolated, repeatable, debuggable.

0) Mandatory Pre-Work: System Recon (Do This First)
0.1 Identify the System Under Test (SUT)

Determine and document:

Testing mode(s) required:

Backend E2E: API black-box tests against running backend + real DB

Frontend E2E: UI browser automation against running frontend + real backend + real DB

Full-stack E2E: UI → API → DB verification

App type:

SPA, SSR, or hybrid

Topology:

monolith vs microservices

identify system boundary (what counts as “internal” vs “external”)

Entry points:

API routers/controllers, GraphQL resolvers, WebSockets

SSR request handlers and rendering pipeline

background workers/cron jobs

frontend route entry points and protected areas

0.2 Detect Stack & Test Tooling (Repo-Aware)

Study repo to identify:

languages and frameworks used (per app/package)

build/start commands (dev/test/staging)

existing E2E tools/patterns:

Playwright / Cypress / Selenium / WebdriverIO

API test tooling (Supertest, REST clients, Postman/Newman, etc.)

conventions:

test file naming (*.e2e.*, *.spec.*)

folder location (/test/e2e, /e2e, /cypress, etc.)

reporting format used in CI

environment management:

.env.*, config loader, secrets usage in tests

Rule: Use the repo’s existing E2E tools if present. Don’t introduce a new tool unless necessary.

0.3 Build a Security & Data Boundary Map (Code-Informed)

Map where these controls are enforced:

authentication (session/JWT)

authorization (RBAC/ABAC, object-level checks)

tenant scoping (if multi-tenant)

input validation

SSR escaping & state hydration controls

file upload pipeline controls

integrations: email/webhooks/payments/storage

queue/worker behavior and idempotency boundaries

0.4 Identify High-Risk Modules (Prioritize in E2E)

List flows with high risk / high impact:

login/session refresh/logout

password reset / account recovery

admin actions / approvals

data export / reporting endpoints

payments / PII flows

file import/upload

role changes / permission updates

cross-tenant data access risks

1) E2E Coverage Scope
1.1 Critical User Journeys (CUJ)

Define a set of Critical User Journeys that represent the most valuable and risky flows. For each CUJ specify:

CUJ ID (stable identifier)

title

user persona (Admin/User/Guest/etc.)

entry preconditions (data state)

steps (user actions)

expected outcomes:

UI changes (if UI exists)

API responses and status codes

DB state changes and invariants

negative variants (invalid input, forbidden access, conflicts)

Persona Coverage (Minimum)

Include journeys for:

Admin

Standard user

Guest/unauthed user

Any special roles (approver, auditor, operator, etc.)

1.2 Full-Stack Validation Requirements

Verify end-to-end:

UI action triggers backend logic correctly

backend persists correct data in real DB

API contract behavior is correct

async jobs complete correctly (bounded waiting strategy)

integrations are triggered correctly and handled safely

1.3 UI/UX Consistency (for UI E2E)

Verify:

correct navigation, URL changes, route protection

form behaviors (validation, disabled states, focus)

visible success/error messages

loading indicators appear/disappear as expected

critical accessibility basics:

stable selectors and semantic queries

focus behavior on error

keyboard navigation for critical flows (if feasible)

2) Environment Strategy (Real App + Real DB)
2.1 Environment Requirements

The E2E environment must include:

running app(s) in test/staging config

real database (containerized preferred)

migration strategy:

auto-run migrations before tests

seeding strategy:

baseline entities (roles/permissions/reference data)

test users (admin/user/guest)

deterministic configuration:

fixed time behavior when needed

predictable test user credentials and identifiers

2.2 Deterministic Data Strategy (Mandatory)

Choose ONE primary approach and document it clearly:

Option A — Ephemeral DB Per Test Run (Preferred)

start fresh DB (container)

run migrations

seed baseline

run E2E tests

destroy DB

Option B — Reset Between Tests

truncate/reset tables between tests

reseed baseline

ensure each test independent and order-free

Option C — Transaction Rollback Per Test (Only If Reliable)

wrap test operations in transaction and rollback

ensure background jobs/async don’t break assumptions

if async exists, avoid this option or scope carefully

Universal Rules (All Options)

tests must be order-independent

tests must be parallel-safe if parallel is enabled

use a test-run-id tag for created data (where feasible)

cleanup must be guaranteed even on test failure

2.3 External Integrations Strategy (Sandbox First)

Define per integration:

type (email/webhook/payment/storage/SSO)

strategy:

sandbox (preferred)

controlled endpoint

stub only if unavoidable (must justify)

coverage:

at least one success scenario

at least one failure scenario (timeout/500/invalid payload)

observability:

verify outcomes via provider logs, webhook receivers, or internal persisted records

3) Test Design & Implementation Requirements
3.1 Tooling Selection Rules (Flexible, Repo-Aware)

Choose tools based on detected stack:

UI automation:

Playwright / Cypress / Selenium / WebdriverIO

API verification:

use language-native HTTP client or existing test libs

DB verification:

prefer API read-back for black-box purity

allow direct DB queries for strong invariants when safe and available

Reporting:

integrate with existing CI reporter formats

3.2 Test Organization (Maintainable by Design)

Require:

clear folder structure:

e2e/ with subfolders by module or CUJ

shared utilities:

environment bootstrap

auth helpers

data seeding/reset helpers

API client wrappers (if needed)

page objects / screen objects for UI tests (recommended)

consistent naming:

CUJ-###: <Journey Name> [persona]

3.3 Assertion Standards (Behavior First)

Each CUJ test must include assertions at multiple layers (where applicable):

UI state assertions:

element visible

text/content correctness

disabled/enabled states

API behavior assertions:

status codes and error shapes

response schema/envelope consistency

Persistence assertions:

verify created/updated state (API read-back or DB query)

verify invariants (ownership, tenant scope, uniqueness, transitions)

Rule: Do not assert internal method calls, logs, or implementation detail.

3.4 Async/Eventual Consistency Handling (If Applicable)

If the system uses queues/background jobs:

define strategy:

run real worker in E2E environment, OR

enable deterministic “inline mode” that still executes real code paths

wait strategy must be bounded:

polling with max timeout

explicit event completion checks

verify idempotency for at least one scenario if job is retryable

Prohibition: no unbounded waits; no arbitrary sleep() unless there is no alternative (must justify and minimize).

3.5 Stability & Flake Prevention

control time:

use server-side test clock config if supported

avoid time-based assumptions

control randomness:

deterministic seeds for random generators if used

avoid brittle selectors:

prefer data-testid or accessibility roles

isolate tests:

no shared state across tests

3.6 CI Performance Strategy

Split test suite into layers:

Smoke E2E (fast): critical CUJs only (core login + top 2 flows)

Full E2E (complete): all CUJs and negative paths

optional:

nightly extended suite for heavy flows/integrations

4) Scenario Coverage Requirements (Minimum Set)
4.1 Happy Path Coverage (Per CUJ)

For each CUJ:

execute flow successfully

verify:

correct navigation/UI outcomes

API behavior correctness

DB state changes/invariants

4.2 Negative Path Coverage (Per CUJ)

For each CUJ include at least:

validation error (missing/invalid field)

permission denial (forbidden role)

not found / conflict / duplicate scenarios (as applicable)

4.3 Edge Case & Boundary Coverage

Include:

empty lists / first item / last item

min/max input sizes

pagination boundaries

long text inputs and special characters (security-adjacent)

state transition invalid scenarios

4.4 Resilience & Failure Mode Coverage

Where applicable:

third-party failure (sandbox negative)

network timeout behavior (bounded)

retry logic behavior (bounded)

5) Traceability & Compliance Mapping (Add This)
5.1 Requirement ↔ CUJ ↔ Test Case Traceability

All tests must be traceable to requirements. Create a traceability structure:

Entities

REQ: requirement ID (from spec/ticket/user story)

CUJ: critical user journey ID

TC: test case ID (individual test)

Rules

Every CUJ must map to one or more REQs

Every test case must map to exactly:

1 CUJ

≥1 REQ(s)

Provide a traceability artifact:

a markdown table or structured list mapping:

REQ → CUJ → TC

Include coverage status:

Planned / Implemented / Blocked

5.2 Coverage Gating & Reporting

Define minimum coverage expectations:

all “Must-have” REQs have at least one CUJ and ≥1 TC

smoke suite covers the top-risk REQs

full suite covers all CUJs and their negative variants

6) Security Regression Subset (Add This)
6.1 Purpose

Maintain a security regression E2E subset focused on preventing reintroduction of common vulnerabilities in real flows. This is not a pentest, but a targeted regression suite for critical security controls.

6.2 What to Include (E2E-Level, Realistic)

Include tests that verify:

Authentication/Session

login success and failure

session expiration or invalid token behavior (as applicable)

logout invalidates session (where applicable)

Authorization / IDOR

user cannot access another user’s resource by changing IDs (IDOR test)

role changes affect access as expected

admin routes are forbidden to non-admin

Input Validation / Injection Safety (Black-Box Signals)

invalid payload is rejected consistently

suspicious strings do not result in:

unhandled errors

data corruption

unintended behavior
(Do not claim injection exploit success unless proven; focus on correct rejection/handling.)

XSS / Unsafe Rendering (Frontend/SSR)

user-provided content is displayed safely (no script execution)

SSR-rendered pages do not reflect raw query params unsafely

unsafe HTML insertion is not used in critical views (validate behaviorally)

CSRF/CORS/Cookie Behavior (If Applicable)

cookie-based auth endpoints require CSRF token if system design requires it

cross-origin credentialed requests are rejected unless explicitly allowed

Sensitive Data Exposure

errors do not expose stack traces or secrets

responses do not leak sensitive fields to unauthorized users

6.3 Structure

tag these tests as security-regression

run them:

on every PR (subset)

plus in full regression nightly if needed

ensure these tests are deterministic and fast

7) Output Requirements (Strict)

The answer must include:

Recon summary

SUT type (backend/frontend/full-stack)

stack/tooling detected

entry points and high-risk modules

Environment setup plan

how to start app(s)

real DB provisioning + migrations + seeding
deterministic reset