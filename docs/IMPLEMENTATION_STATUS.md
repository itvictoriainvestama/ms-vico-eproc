# Implementation Status

This document tracks the current backend alignment against:

- `BRD_E-Procurement.md`
- `FSD_E-Procurement.md`
- `TSD_E-Procurement.md`

## Current Phase

Phase 1 Foundation

## Completed In This Phase

### Architecture foundation

- backend remains a single Go application for local simplicity
- internal structure has been reshaped toward logical service boundaries
- API namespaces now follow the TSD direction:
  - `/api/v1/auth/*`
  - `/api/v1/internal/*`
  - `/api/v1/vendor/*`
  - `/api/v1/files/*`
  - `/api/v1/reports/*`
  - `/api/v1/admin/*`

### HTTP foundation

- standard success response implemented
- standard error response implemented
- per-request trace ID implemented via middleware
- health endpoint now returns standard response format

### Data foundation

The schema foundation has been aligned toward the TSD by introducing or reshaping these areas:

- `entities`
- `roles`
- `users`
- `user_roles`
- `delegate_approvers`
- `vendors`
- `vendor_users`
- `vendor_blacklists`
- `reference_prices`
- `budgets`
- `procurement_policy_rules`
- `approval_models`
- `approval_matrices`
- `purchase_requests`
- `purchase_request_items`
- `pr_attachments`
- `pr_approvals`
- `rfqs`
- `rfq_vendors`
- `quotations`
- `quotation_items`
- `vendor_evaluations`
- `bafo_rounds`
- `vendor_selections`
- `direct_appointments`
- `purchase_orders`
- `purchase_order_items`
- `po_approvals`
- `vendor_confirmations`
- `notifications`
- `audit_logs`
- `app_logs`

### Authentication foundation

- JWT claims expanded with:
  - `entity_id`
  - `role_code`
  - `role_name`
  - `scope_type`
  - `subject_type`
- auth flow now reads primary role context from user-role relations

### Access control foundation

- role-based route guards now protect internal and admin namespaces
- entity-scoped access checks now apply in service-level detail, status update, and approval actions
- admin management endpoints added for `entities` and `users`

### Verified

- `go test ./...` passes
- migration completes successfully on local MySQL after the migration strategy adjustment
- API boot verified successfully on alternate port `8081`
- local bootstrap flow now supports fresh DB reset + migrate + baseline seed
- baseline seed now creates `entities`, `roles`, `departments`, admin `users`, and `user_roles`
- admin user/entity APIs verified with seeded login and RBAC smoke checks

## Important Notes

### Existing database caution

This refactor moves the schema significantly closer to the TSD. If a developer already has an old local database from the earlier version of the project:

- tables may be altered in place
- some legacy tables may still remain beside the new TSD-aligned tables
- a fresh local database is recommended for clean development going forward

### Migration strategy note

Migration uses GORM auto-migrate with foreign key constraint creation disabled during migration to avoid collisions with legacy constraint names in older local databases.

This was necessary because the pre-refactor schema and the TSD-aligned schema overlap semantically but not structurally.

## Still Pending For Next Phases

### Immediate next implementation slice

- change-password flow for seeded and managed users
- role assignment update endpoints
- department management endpoints
- approval/policy master data APIs

### Master and governance modules

- entity management endpoints
- user management endpoints
- role assignment endpoints
- delegate approver management
- budget management endpoints
- procurement policy configuration endpoints
- approval workflow configuration endpoints

### Procurement flow depth

- PR attachment upload flow
- PR revision and resubmit flow
- method selection flow
- RFQ publish flow
- RFQ reopen/close flow
- vendor quotation submission flow
- direct appointment full workflow

### Evaluation domain

- technical evaluation
- commercial evaluation
- weighted scoring
- BAFO workflow
- vendor selection workflow

### PO domain

- PO approval chain
- send-to-vendor flow
- vendor confirmation flow
- void flow

### Cross-cutting modules

- audit service enrichment
- notification dispatch
- report/export generation
- file upload signed flow
- queue/worker integration
- dashboard/report aggregation

### Security and control

- full RBAC enforcement by role matrix
- data isolation by entity scope in all queries
- SoD enforcement
- force-change-password flow
- reset password flow
- failed login lock policy
- vendor portal-specific authorization policy

## Recommended Next Implementation Order

1. seed and master data bootstrap for `entities`, `roles`, and admin user
2. internal auth and RBAC hardening
3. entity and user management APIs
4. budget and procurement policy modules
5. PR workflow with approval model resolution
6. RFQ and vendor participation flow
7. evaluation and BAFO
8. PO and vendor confirmation
9. reporting, notifications, and background jobs
