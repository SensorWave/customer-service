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

The service only needs database and RabbitMQ configuration from `.env`.

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

Full route-by-route documentation with direct URLs and `curl` examples is available in [`docs/API.md`](docs/API.md).

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
