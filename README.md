# 🧩 Customer Service

This project is a **Go** microservice for managing **companies** and **users**.  
It uses **PostgreSQL** as the database, **RabbitMQ** for asynchronous communication between microservices, and runs entirely via **Docker Compose**.

## 🚀 Features

- Management of **companies** (`companies`)
- Management of **users** (`users`)
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
| POST   | `/companies/:id/users` | Create a user linked to a company |
| GET    | `/users/:id`           | Retrieve a user                   |

> Sample CURLs are available in each resource-specific `repository`.

## 🗄️ Accessing PostgreSQL

From your terminal:

```bash
psql -h localhost -p 5433 -U user -d customers
```

Then:

```sql
\dt -- List tables
SELECT * FROM companies;
SELECT * FROM users;
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
