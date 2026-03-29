Antigravity Workflow — Backend → Frontend API Contract + Integration Brief (Rating 1–10)
Mission

You are the Backend Agent. Your job is to:

Study the backend modules (actual implemented code) to understand real behavior.

Produce a complete API contract that reflects reality (or explicitly notes gaps/TODO).

Produce a Frontend Integration Brief so a Frontend Agent can integrate reliably with minimal back-and-forth.

Key output: a brief that answers FE questions before they’re asked.

0) Mandatory Recon (Backend First, Always)
0.1 Map Backend Modules & Entry Points

MUST identify:

service/module boundaries (domain modules)

route registrations (REST/GraphQL/WebSocket)

controllers/handlers/resolvers

services/use-cases orchestration

repositories/data access

authn/authz middleware/guards

response envelope + error conventions

common interceptors/middlewares (logging, validation, transform)

0.2 Identify “Integration Surface”

MUST list:

all endpoints the frontend is expected to call for this feature

whether the endpoint is:

stable/public vs internal

v1/v2/versioned/unversioned

base URL(s) per environment (dev/staging/prod) if present in code/config

required headers and tokens

0.3 Confirm What Is REAL vs PLANNED

MUST separate:

implemented behavior (confirmed in code)

partially implemented

not implemented (proposal only)

Rule: Never present a proposed field/endpoint as implemented unless confirmed.

1) Mandatory Best-Practice Alignment (Stack-Aware)

After recon, adapt the brief to the repo’s conventions:

DTO validation pattern used in backend

naming and response formats already used

error_code taxonomy style

pagination/filter/sort pattern used in other modules

auth/role/tenant patterns used by the system

Prefer existing patterns over inventing new ones.

2) Deliverable A — API Contract (Frontend-Facing)
2.1 Endpoint Catalog (Required)

For each endpoint, provide a contract block:

Endpoint Block Template (Mandatory Fields)

Name / Purpose

Method + Path

Auth

unauthenticated/authenticated

roles/permissions if applicable

Headers

Authorization format (Bearer/cookie)

request_id/correlation_id rules

idempotency key if mutation

Request

path params

query params (with defaults)

JSON body schema (field-by-field)

validation rules (required/optional, type, min/max, enum, format)

Response (Success)

HTTP status

envelope format (actual)

result schema (field-by-field)

Response (Errors)

status codes used

error_code list for FE handling

error_details shape

Behavior Notes

side effects

ordering guarantees

pagination semantics

sorting/filtering semantics

Examples

example request

example success response

example error response

Strict requirement: examples must match the response envelope used by backend.

2.2 Field-Level Data Dictionary (Required)

Provide a dictionary for all relevant DTO fields:

field name

type

nullable?

source of truth (DB field / computed / external)

format (date, currency, uuid)

FE display hints (if obvious)

validation constraints

2.3 Error Handling Contract (FE Must Know)

Define:

global error envelope

standard error codes FE must handle:

validation_error

unauthorized

forbidden

not_found

conflict

rate_limited (if present)

server_error

which error codes appear on which endpoints

FE recommended behavior per error code (toast, redirect login, show inline field errors, etc.)

3) Deliverable B — End-to-End Flow (User + System)
3.1 User Journey Flow (Frontend Perspective)

Describe the flow FE should implement:

screens/pages involved

user actions

when to call which endpoint

loading states

success states

failure states & user messaging

navigation transitions

3.2 System Sequence Flow (Backend Perspective)

Provide a sequence-like description:
FE → Endpoint A → Service → Repo → DB → (Queue/External) → Response
Include:

validation steps

authz steps

transaction boundaries (conceptual)

async completion if any (polling/webhook/event)

4) Deliverable C — Frontend Integration Brief (Hand-off Document)

This is the document the Frontend Agent will consume.

4.1 Integration Checklist (Required)

base URL and environment variables FE must set

auth prerequisites (how FE gets token/session)

required headers

how to construct requests safely

pagination/filtering rules

idempotency rules (if needed)

retry guidance (only for safe requests)

4.2 API Call Matrix (Required)

Create a matrix that maps:

UI action → endpoint → request payload → expected response → UI update

Example columns:

UI Action

Endpoint

Request (fields)

Success Result fields used by FE

Error codes to handle

UI behavior on success/error

4.3 Data Mapping Guide (Required)

For each screen/component, specify:

which API fields map to which UI fields

formatting notes (dates, currency, status labels)

derived states (e.g., status mapping)

4.4 State Machine / Status Mapping (If Applicable)

If the module has statuses (draft/approved/etc.):

list possible states

allowed transitions

FE rules:

which buttons show in which state

which transitions call which endpoint

what errors occur if invalid transition

4.5 Edge Cases & Gotchas (Required)

List “things that will break integration if missed”:

nullability quirks

optional fields not always present

eventual consistency delay

async job completion requirements

permission-based hiding vs server enforcement

known backend constraints

4.6 Test Data & Manual Testing Steps (Required)

Provide:

test accounts/roles (if in seed)

sample IDs or how to create them

manual test script FE can run

expected outputs

5) Non-Breaking & Additive Evolution Rules (Mandatory)

The brief must:

preserve existing endpoints/fields unless explicitly changed

add new fields as optional by default

add new endpoints instead of changing behavior broadly

document any breaking change explicitly with migration steps

6) Scoring Rubric (Rate the Brief + Contract Quality)
Aspect	Score (1–10)
Backend recon completeness (modules/routes/DTOs found)	
Contract accuracy (matches actual backend behavior)	
Endpoint completeness (all FE-needed endpoints covered)	
Request schema clarity + validation detail	
Response schema clarity + field dictionary	
Error model clarity (codes + FE handling guidance)	
Flow clarity (FE journey + backend sequence)	
Integration brief completeness (checklist + call matrix + mapping)	
Edge cases & gotchas coverage	
Non-breaking evolution plan	
Consistency with repo conventions	
Testability (manual test script + deterministic examples)	
Overall usefulness for FE integration	

Final Rating = average of all scores
Output format: Rating: X/10

7) Output Requirements (Strict)

The backend agent must output:

Recon summary (what was inspected, what patterns found)

API Contract (endpoint blocks + examples)

End-to-end flow (user journey + system sequence)

Frontend Integration Brief (checklist + call matrix + mapping + edge cases + test steps)

Scoring table

Rating: X/10