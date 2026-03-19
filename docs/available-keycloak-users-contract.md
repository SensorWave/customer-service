# Available Keycloak Users Contract Proposal

Status: proposal only, not implemented in `customer-service` yet.

## Goal

Support the company-user creation flow with a backend endpoint that returns Keycloak users who can still be linked through:

`POST /companies/:id/users-company`

Because `user_companies.id_auth_kc` is unique, this endpoint should only return users that are not already linked in `user_companies`.

## Proposed Route

`GET /companies/:id/available-keycloak-users`

## Query Parameters

- `search` optional string
Used to filter by username, email, first name, last name, or display name.

- `limit` optional integer
Maximum number of rows returned. Default `25`, max `100`.

## Response Contract

`200 OK`

```json
{
  "items": [
    {
      "keycloak_id": "08d697a6-8b89-40fc-aaf3-bdaaa65e0ae4",
      "username": "alice",
      "email": "alice@example.com",
      "first_name": "Alice",
      "last_name": "Martin",
      "display_name": "Alice Martin"
    },
    {
      "keycloak_id": "2f7f4ebf-5a0d-4a83-9f17-7cf4d7c9ac25",
      "username": "bob",
      "email": "bob@example.com",
      "first_name": "Bob",
      "last_name": "Nguyen",
      "display_name": "Bob Nguyen"
    }
  ],
  "meta": {
    "company_id": 12,
    "search": "ali",
    "limit": 25,
    "returned": 2,
    "source": "keycloak",
    "excluded_already_linked": true
  }
}
```

## Field Semantics

- `keycloak_id`
Stable Keycloak user identifier. This is the value the frontend sends as `id_auth_kc` when creating a company user.

- `display_name`
Backend-computed label for dropdown usage. Prefer `first_name + last_name`; fall back to `username`; fall back to `email`.

- `meta.excluded_already_linked`
Must be `true` for the default behavior so the frontend can trust this endpoint as a list of selectable candidates.

## Sorting

Return rows sorted by:

1. exact `search` matches first when `search` is present
2. then alphabetical `display_name`
3. then alphabetical `username`

## Error Contract

`400 Bad Request`

```json
{
  "error": "invalid company ID"
}
```

`401 Unauthorized`

```json
{
  "error": "missing or invalid access token"
}
```

`403 Forbidden`

```json
{
  "error": "forbidden"
}
```

`502 Bad Gateway`

```json
{
  "error": "failed to fetch users from Keycloak"
}
```

## Why This Shape

- The route is company-scoped because the dropdown is used from the company user creation screen.
- `items + meta` leaves room for paging, tracing, and source diagnostics without breaking the payload later.
- The backend, not the frontend, should decide which Keycloak users are still available because it owns the uniqueness rule on `id_auth_kc`.
- `display_name` avoids duplicating label-building logic in each frontend client.

## Frontend Usage

The create-user screen should call:

`GET /companies/:companyId/available-keycloak-users`

Then use `items[].display_name` for the dropdown label and `items[].keycloak_id` as the submitted value for `id_auth_kc`.
