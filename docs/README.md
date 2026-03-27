# Project Diagrams

This folder contains Mermaid diagrams that describe the `customer-service` runtime architecture, request flow, and core data model.

## Available diagrams

- [Architecture diagram](./architecture-diagram.md): infrastructure, application layers, and event publication flow.
- [Request sequence diagram](./request-sequence-diagram.md): how a typical request moves from Gin routing to services, repositories, PostgreSQL, and RabbitMQ.
- [Domain data model diagram](./domain-data-model-diagram.md): the relationship between companies, user memberships, and emitted domain events.

## Notes

- The diagrams reflect the current implementation in `cmd/main.go`, `routes/`, `internal/company/`, `internal/user/`, `internal/event/`, `config/config.go`, and `scripts/create_database.sql`.
- Existing API reference files in this folder were left unchanged.
