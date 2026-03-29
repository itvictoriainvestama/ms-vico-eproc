Code Security Audit Workflow (Code-Level Only) — Ultra Modern (Rating 1–10)
Instruction

Perform a code security audit focused strictly on source code risks (not pentesting, not infrastructure, not DevOps).
The agent must:

Study the codebase first (repo structure, stack, conventions, trust boundaries).

Apply idiomatic security best practices for the detected language/framework (FE/BE, SPA/SSR, any language).

Identify security issues with evidence, impact, exploitability, and concrete minimal-diff fixes.

Produce a structured risk rating and a final score from 1–10.

Constraints

Focus on code-level vulnerabilities, unsafe patterns, and security design flaws.

Prefer minimal diff, additive fixes.

For frontend: do not change styling unless directly required to fix a security issue.

0) Mandatory Recon (Always First)
0.1 Stack & Architecture Discovery

The agent MUST map:

languages, frameworks, runtimes

app type: SPA/SSR/hybrid

architecture: monolith/microservices

entry points:

backend controllers/handlers/resolvers

frontend routing + rendering + data fetching

background workers/jobs

shared libraries and cross-cutting middleware

0.2 Security Boundary & Data Classification Map

The agent MUST identify:

trust boundaries (what inputs are untrusted)

authentication mechanism(s) (JWT/cookie/session/OAuth)

authorization model(s) (RBAC/ABAC/tenant scoping)

sensitive data types:

credentials/tokens

PII

payment or financial data

secrets/keys

critical operations:

admin actions

writes/approvals

uploads/imports

exports/reports

0.3 Identify Existing Security Controls (Code)

Find and document:

input validation layer (schema validators, DTO validation)

output encoding/escaping mechanisms (XSS/SSR safety)

CSRF/CORS configuration usage in code

rate limiting usage (if any in code)

logging and redaction behavior

error handling patterns (stack traces exposure risk)

secret loading patterns (env/config, vault clients)

Rule: Do not assume. Point to code evidence and paths.

1) Best-Practice Alignment (Mandatory, Stack-Specific)

After recon, the agent MUST align to:

OWASP Top 10 principles (conceptually) but applied through code review

language/framework-specific secure coding guidance, including:

secure dependency usage patterns

safe templating/rendering

request validation idioms

auth middleware/guards conventions

safe serialization/deserialization

Rule: Reuse existing security utilities in the repo when available.

2) Audit Scope (Code-Level Only)
2.1 Backend Security Areas (If Backend Exists)

Audit for:

Authentication flaws

Authorization bugs (including IDOR / object-level access control)

Input validation and injection risks:

SQL/NoSQL injection

command injection

SSRF

path traversal

Unsafe file upload handling

Cryptography misuse:

weak hashing

insecure random

token signing issues

Session/JWT handling flaws:

missing expiration checks

insecure refresh flow

token leakage

Error handling and information disclosure

Logging of secrets/PII

Rate limiting / brute-force protection (code-level presence)

Deserialization risks

Insecure redirects

Concurrency-related security issues (race conditions enabling privilege bugs)

2.2 Frontend Security Areas (If Frontend Exists)

Audit for:

XSS risks:

unsafe HTML injection

dangerous rendering APIs

SSR hydration/serialization issues

CSRF issues (cookie-based auth flows)

Insecure token storage:

localStorage/sessionStorage risks

leaking tokens to third-party scripts

Sensitive data exposure in UI or logs

Open redirect and unsafe navigation

Insecure CORS usage in client fetch patterns (where relevant)

SSR-specific issues (if SSR):

reflected content injection

unsafe server-side templating

leaking secrets into client bundles

Dependency usage patterns (supply chain at code-level)

Constraint reminder: Do not change UI styling unless required for security.

2.3 Shared / Cross-Cutting Security Areas

Audit for:

secret management and config loading

dependency hygiene:

dangerous packages usage patterns

outdated crypto libs usage patterns (code-level evidence)

insecure defaults and feature flags

unsafe debug endpoints or dev-only code left enabled

permission checks consistency across modules

insecure inter-service calls (missing auth between services, if in code)

3) Vulnerability Reporting Requirements (Traceable + Actionable)

For EACH finding, the agent MUST provide:

ID (stable)

Category (Auth/AuthZ/XSS/Injection/etc.)

Severity (Critical/High/Medium/Low)

Exploitability (High/Medium/Low) — practical likelihood

Impact (what attacker gains)

Evidence

exact location (file/module/function)

what pattern is unsafe

Attack Scenario (Code-Level)

concise “how it could be abused” without providing weaponized steps

Fix Recommendation

minimal-diff patch strategy

framework-idiomatic safe alternative

Regression Risk (Low/Med/High)

Verification

how to confirm it’s fixed (unit/integration tests, static checks)

Rule: Provide remediation that is realistic for the repo’s stack.

4) Prioritization Model (Most Modern)

Sort findings by:

Severity

Exploitability

Blast radius (how many endpoints/users affected)

Ease of fix (quick wins first where impact is high)

Also provide:

“Top 5 Immediate Fixes” list

“Structural Fixes” list (requires coordinated changes)

5) Security Fix Constraints (Minimal Diff + Safe Evolution)
MUST

preserve existing behavior unless change is security-critical

introduce safe wrappers/utilities rather than rewriting entire modules

prefer centralized fixes (middleware/guards) when appropriate

add tests to prevent regression (where feasible)

NEVER

add new security library without checking existing repo equivalents

suppress errors by hiding them (must handle safely)

log secrets/PII during debugging

weaken auth requirements “to make it work”

6) Recommended Security Regression Tests (Code-Level)

Design test recommendations (not pentest):

Auth: login failure throttling logic (if exists), token validation

AuthZ: object-level access checks (IDOR prevention)

Validation: rejects malicious/suspicious inputs safely

XSS: unsafe HTML is not rendered/executed

Upload: rejects unsafe file types / paths

Error leakage: stack traces not exposed

7) Output Requirements (Strict)

The answer MUST include:

Recon summary

stack/frameworks

trust boundaries

sensitive data map

Security posture summary

overall risk level and patterns observed

Findings

prioritized list with full fields (ID, severity, evidence, fix, verification)

Top 5 immediate fixes

Structural improvements

Security regression test recommendations

Scoring table

Final rating: Rating: X/10

🧮 Structured Security Audit Quality Rubric (Score 1–10)
Aspect	Score (1–10)
Recon completeness (stack + boundaries + data map)	
Coverage of auth/authz risks	
Coverage of injection risks (SQL/NoSQL/SSRF/path traversal)	
Coverage of XSS/SSR risks (if frontend/SSR exists)	
Secret management & sensitive data handling review	
Error handling & information disclosure review	
Logging safety & redaction review	
Dependency usage risks (code-level patterns)	
Quality of evidence (traceable locations)	
Quality of remediation (minimal diff, idiomatic)	
Prioritization quality (severity/exploitability/blast radius)	
Regression test recommendations quality	
Overall clarity & usefulness of the audit report	

Final Rating = average of all aspect scores
Output format: Rating: X/10

Rating Scale Meaning

9–10 = Excellent: comprehensive, practical, actionable, strongly aligned with best practices

8 = Very good: minor gaps in coverage or remediation detail

7 = Good: some important areas need deeper review

6 = Acceptable: multiple missing areas or weak prioritization

5 = Borderline: significant redesign of audit approach needed

3–4 = Weak: misses critical classes of vulnerabilities

1–2 = Very poor: unclear, untraceable, not actionable

0 = Broken / Unusable / Not code-level security audit