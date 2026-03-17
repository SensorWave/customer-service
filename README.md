# 🧩 Customer Service

This project is a **Go** microservice for managing **companies** and **users**.  
It uses **PostgreSQL** as the database, **RabbitMQ** for asynchronous communication between microservices, and runs entirely via **Docker Compose**.

## 🚀 Features

- Management of **companies** (`companies`)
- Management of **users** (`users-company`)
- Publishing RabbitMQ events (`CompanyCreated`, `UserCreated`)
- Modular and extensible microservice architecture
- Persistence via PostgreSQL

## 🧱 Architecture

```bash
├── cmd/
│   └── main.go         # Application entry point
├── config/
│   └── config.go       # Configuration & env loading
├── internal/
│   ├── customer/
│   │ ├── company/      # "Company" business logic
│   │ └── user/         # "User" business logic
│   └── event/          # RabbitMQ publishing
├── routes/
│   └── route.go        # Application routes
├── docker-compose.yml  # Full stack (Go + Postgres + RabbitMQ)
├── Dockerfile          # Go microservice build
├── go.mod / go.sum     # Go dependencies
└── .env.example        # Environment variables
```

## ⚙️ Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- (Optionnel) [psql](https://www.postgresql.org/download/) to connect to the DB

## 🔧 Configuration

Create a `.env` file at the root of the project and fill in the variables:

```bash
cp .env.example .env
```

For local Keycloak, set at least these values:

```bash
KEYCLOAK_ISSUER=http://localhost:8888/realms/dev
KEYCLOAK_JWKS_URL=http://localhost:8888/realms/dev/protocol/openid-connect/certs
KEYCLOAK_AUDIENCE=
```

`KEYCLOAK_AUDIENCE` is optional and can stay empty until audience validation is added.
The service now refuses to start if `KEYCLOAK_ISSUER` or `KEYCLOAK_JWKS_URL` is missing.
If `customer-service` runs in Docker while Keycloak runs on your host machine, replace `localhost` with `host.docker.internal`.

## 🐳 Running the project

From the root of the project:

```bash
docker compose up --build
```

## 🧠 Usage

📋 API (via Postman, curl, etc.)

| Method | Endpoint               | Description                       |
| ------ | ---------------------- | --------------------------------- |
| POST   | `/companies`           | Create a company                  |
| GET    | `/companies/:id`       | Retrieve a company                |
| POST   | `/companies/:id/users-company` | Create a user linked to a company |
| GET    | `/users-company/:id`           | Retrieve a user                   |

> Sample CURLs are available in each resource-specific `repository`.

## 🗄️ Accessing PostgreSQL

From your terminal:

```bash
psql -h localhost -p 5434 -U user -d customer_service_db
```

Then:

```sql
\dt -- List tables
SELECT * FROM companies;
SELECT * FROM user_companies;
```

## 🧹 Cleanup

To stop and remove containers:

```bash
docker compose down
```

To remove everything (containers + volumes + images):

```bash
docker compose down -v --rmi all
```

## 📘 Notes

- The project is extensible: you can add other microservices that consume RabbitMQ events.
- The code is organized in a modular structure, making maintenance and unit testing easier.
