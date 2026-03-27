# Context and Architecture Decision (Keycloak + Customer Service)

## Context
We currently have:
- **Keycloak** already running on the frontend side (Vue) for authentication, including username/password login and JWT issuance.
- A **Customer Service (Go + Postgres)** that manages the business model: **companies** and related logic.
- An initial **duplication** issue: Keycloak stores identity users, and the Customer Service also stores users.

Goal: **a single login service** (Keycloak), while still keeping the business relationship `user ↔ company` inside the Customer Service, without duplicating identity data.

## Decision
1. **Keycloak is the source of truth for identity and authentication**
   - login / password / MFA / password reset / sessions
   - Keycloak user identifiers (UUID) exposed in tokens through the `sub` claim

2. **Customer Service is the source of truth for the business model**
   - membership in a company
   - business role inside the company, for example `company_owner`, `company_admin`, `member`
   - business permissions, such as who can manage members

3. **No authentication is performed against the Customer Service database**
   - Keycloak does not query the Customer Service database for login or registration.
   - The backend validates the Keycloak JWT, reads `sub`, then resolves the business context through its own database.

## Selected Data Model
Inside Customer Service, we do not store identity data such as email, first name, last name, or password as required business records.
We only store a linking table, for example `company_members` or `user_company`:

- `keycloak_user_id` (string/UUID) = the `sub` claim from the Keycloak JWT
- `company_id`
- `role` (business role)

> Option: enrich the UI with display data such as email or name on demand through the Keycloak Admin API, without persisting those fields.

## Runtime Flow
### Authentication and API calls
1. The user logs in through **Keycloak** from the Vue frontend via `keycloak-js`.
2. The frontend calls the API with `Authorization: Bearer <access_token>`.
3. The backend:
   - validates the JWT (issuer + JWKS + exp)
   - extracts `sub`
   - resolves `company_id` and `role` through the `company_members` table

### Minimum required endpoint
- `GET /me` (protected):
  - reads `sub`
  - returns `{ company_id, role }`
  - returns `404 { code: "NOT_LINKED" }` if no link exists

## Member Management (company owner)
Long-term goal: a `company_owner` can add or remove members from their company.

Decision:
- The "CRUD users" frontend for the owner calls **Customer Service**
- Customer Service applies business rules, for example owner/admin only
- Customer Service updates the relationship table (`company_members`)
- Keycloak account creation and deletion happens **through the backend** and never directly from the frontend

### Phasing (current state)
- Right now, Keycloak accounts are created manually by a Keycloak admin.
- The backend must therefore expose an internal bootstrap endpoint to create the database link:
  - `POST /admin/memberships/link { keycloak_user_id, company_id, role }`

### Phasing (future evolution)
- Later, the backend will create Keycloak accounts from the owner UI through the Keycloak Admin API
  - service account (`backend-admin`) + `client_credentials`
  - user creation + required actions such as verify email / set password
  - then insertion of the `company_members` link

## Non-goals (to avoid unnecessary complexity)
- Do not implement a Keycloak User Federation/SPI to authenticate against the Customer Service database.
- Do not expose Keycloak admin credentials or admin access to the frontend.
- Do not systematically duplicate identity data such as email or name in the business database.

## Consequences
- There are two "representations" of a user:
  - **Keycloak** for identity and security
  - **Customer Service** for business context: company and role via `keycloak_user_id`
- This is not duplicated authentication: there is still **only one login flow**, through Keycloak.
- Permissions such as "company owner manages members" are enforced by the backend through the business database.
