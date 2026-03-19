# 🧩 Customer Service

This project is a **Go** microservice for managing **companies** and **users**.  
It uses **PostgreSQL** as the database, **RabbitMQ** for asynchronous communication between microservices, and runs entirely via **Docker Compose**.

## 🚀 Features

- Management of **companies** (`companies`)
- Management of **users** (`users-company`)
- Lookup of company users by **Keycloak ID**
- Global **CORS** support with `OPTIONS` preflight handling
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
│   ├── company/        # "Company" business logic
│   ├── user/           # "User" business logic
│   └── event/          # RabbitMQ publishing
├── docs/
│   ├── API.md          # Route-by-route API reference
│   └── available-keycloak-users-contract.md
│                      # Proposed dropdown contract for frontend integration
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

The service only needs database and RabbitMQ configuration from `.env`.

## 🐳 Running the project

From the root of the project:

```bash
docker compose up --build
```

## 🧠 Usage

📋 API (via Postman, curl, etc.)

| Method | Endpoint | Description |
| ------ | -------- | ----------- |
| GET | `/health` | Health check |
| POST | `/companies` | Create a company |
| GET | `/companies` | List companies |
| GET | `/companies/:id` | Retrieve a company |
| PUT | `/companies/:id` | Update a company |
| DELETE | `/companies/:id` | Delete a company |
| POST | `/companies/:id/users-company` | Create a user linked to a company |
| GET | `/companies/:id/users-company` | List users linked to a company |
| GET | `/users-company/keycloak/:keycloak_id` | Retrieve a user by Keycloak ID |
| GET | `/users-company/keycloak/:keycloak_id/role` | Retrieve only the user role by Keycloak ID |
| GET | `/users-company/keycloak/:keycloak_id/company-id` | Retrieve only the company ID by Keycloak ID |
| GET | `/users-company/:id` | Retrieve a user |
| PUT | `/users-company/:id` | Update a user |
| DELETE | `/users-company/:id` | Delete a user |

Full route-by-route documentation with direct URLs and `curl` examples is available in [`docs/API.md`](docs/API.md).

The proposed backend contract for the frontend "available Keycloak users" dropdown is documented in [`docs/available-keycloak-users-contract.md`](docs/available-keycloak-users-contract.md).

## 🌐 CORS

The API exposes global CORS headers for browser clients:

- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Origin, Content-Type, Accept, Authorization`
- `OPTIONS` preflight requests return `204 No Content`

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

## Create User
Please run this create for adding your keycloak account to an admin account.
```bash
curl -X POST http://localhost:8080/companies/1/users-company \
  -H "Content-Type: application/json" \
  -d '{
    "id_auth_kc": "08d697a6-8b89-40fc-aaf3-bdaaa65e0ae4",
    "role": "admin"
  }'
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
