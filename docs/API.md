# API Reference

Base URL: `http://localhost:8080`

Proposal docs:
- See [`available-keycloak-users-contract.md`](available-keycloak-users-contract.md) for the proposed backend contract used by the frontend dropdown when creating a company user.

CORS:
- `Access-Control-Allow-Origin: *`
- Allowed methods: `GET, POST, PUT, DELETE, OPTIONS`
- Preflight `OPTIONS` requests return `204 No Content`

## Quick Summary

| Method | Route | Purpose |
| --- | --- | --- |
| GET | `/health` | Health check |
| POST | `/companies` | Create a company |
| GET | `/companies` | List all companies |
| GET | `/companies/:id` | Get one company |
| PUT | `/companies/:id` | Update one company |
| DELETE | `/companies/:id` | Delete one company |
| POST | `/companies/:id/users-company` | Create a user for one company |
| GET | `/companies/:id/users-company` | List users for one company |
| GET | `/users-company/keycloak/:keycloak_id` | Get one user by Keycloak ID |
| GET | `/users-company/keycloak/:keycloak_id/role` | Get only the user role by Keycloak ID |
| GET | `/users-company/keycloak/:keycloak_id/company-id` | Get only the user company ID by Keycloak ID |
| GET | `/users-company/:id` | Get one user |
| PUT | `/users-company/:id` | Update one user |
| DELETE | `/users-company/:id` | Delete one user |

## Data Format

### Company JSON

```json
{
  "id": "1",
  "name": "Acme Corp",
  "email": "contact@acme.com",
  "phone": "0601020304",
  "address": "42 rue de Paris",
  "created_at": "2026-03-19T10:00:00Z",
  "updated_at": "2026-03-19T10:00:00Z"
}
```

Write fields for create/update:

```json
{
  "name": "Acme Corp",
  "email": "contact@acme.com",
  "phone": "0601020304",
  "address": "42 rue de Paris"
}
```

### User JSON

```json
{
  "id": 1,
  "company_id": 1,
  "id_auth_kc": "kc-acme-admin-001",
  "role": "admin",
  "created_at": "2026-03-19T10:00:00Z",
  "updated_at": "2026-03-19T10:00:00Z"
}
```

Write fields for create/update:

```json
{
  "id_auth_kc": "kc-acme-admin-001",
  "role": "admin"
}
```

Allowed `role` values: `user`, `manager`, `admin`

## Routes

### 1. Health Check

URL: `http://localhost:8080/health`

Browser / direct URL:

```text
http://localhost:8080/health
```

curl:

```bash
curl http://localhost:8080/health
```

### 2. Create Company

URL: `http://localhost:8080/companies`

curl:

```bash
curl -X POST http://localhost:8080/companies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Acme Corp",
    "email": "contact@acme.com",
    "phone": "0601020304",
    "address": "42 rue de Paris"
  }'
```

### 3. List Companies

URL: `http://localhost:8080/companies`

Browser / direct URL:

```text
http://localhost:8080/companies
```

curl:

```bash
curl http://localhost:8080/companies
```

### 4. Get Company By ID

URL: `http://localhost:8080/companies/{companyId}`

Browser / direct URL example:

```text
http://localhost:8080/companies/1
```

curl:

```bash
curl http://localhost:8080/companies/1
```

### 5. Update Company

URL: `http://localhost:8080/companies/{companyId}`

curl:

```bash
curl -X PUT http://localhost:8080/companies/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Acme Corporation Updated",
    "email": "support@acme.com",
    "phone": "0707070707",
    "address": "100 avenue de Lyon"
  }'
```

### 6. Delete Company

URL: `http://localhost:8080/companies/{companyId}`

curl:

```bash
curl -X DELETE http://localhost:8080/companies/1
```

### 7. Create User For A Company

URL: `http://localhost:8080/companies/{companyId}/users-company`

curl:

```bash
curl -X POST http://localhost:8080/companies/1/users-company \
  -H "Content-Type: application/json" \
  -d '{
    "id_auth_kc": "08d697a6-8b89-40fc-aaf3-bdaaa65e0ae4",
    "role": "admin"
  }'
```

### 8. List Users By Company

URL: `http://localhost:8080/companies/{companyId}/users-company`

Browser / direct URL example:

```text
http://localhost:8080/companies/1/users-company
```

curl:

```bash
curl http://localhost:8080/companies/1/users-company
```

### 9. Get User By Keycloak ID

URL: `http://localhost:8080/users-company/keycloak/{keycloakId}`

Browser / direct URL example:

```text
http://localhost:8080/users-company/keycloak/kc-acme-001
```

curl:

```bash
curl http://localhost:8080/users-company/keycloak/kc-acme-001
```

### 10. Get User Role By Keycloak ID

URL: `http://localhost:8080/users-company/keycloak/{keycloakId}/role`

Browser / direct URL example:

```text
http://localhost:8080/users-company/keycloak/kc-acme-001/role
```

curl:

```bash
curl http://localhost:8080/users-company/keycloak/kc-acme-001/role
```

Example response:

```json
{
  "role": "admin"
}
```

### 11. Get User Company ID By Keycloak ID

URL: `http://localhost:8080/users-company/keycloak/{keycloakId}/company-id`

Browser / direct URL example:

```text
http://localhost:8080/users-company/keycloak/kc-acme-001/company-id
```

curl:

```bash
curl http://localhost:8080/users-company/keycloak/kc-acme-001/company-id
```

Example response:

```json
{
  "company_id": 1
}
```

### 12. Get User By ID

URL: `http://localhost:8080/users-company/{userId}`

Browser / direct URL example:

```text
http://localhost:8080/users-company/1
```

curl:

```bash
curl http://localhost:8080/users-company/1
```

### 13. Update User

URL: `http://localhost:8080/users-company/{userId}`

curl:

```bash
curl -X PUT http://localhost:8080/users-company/1 \
  -H "Content-Type: application/json" \
  -d '{
    "id_auth_kc": "kc-user-002",
    "role": "manager"
  }'
```

### 14. Delete User

URL: `http://localhost:8080/users-company/{userId}`

curl:

```bash
curl -X DELETE http://localhost:8080/users-company/1
```

## Notes

- `companyId` and `userId` must be numeric in the URL.
- In responses, company `id` is serialized as a string, while user `id` is serialized as a number.
- `GET` routes can be called directly from a browser with the URL examples above.
- `POST`, `PUT`, and `DELETE` should be called with `curl`, Postman, or another API client.
