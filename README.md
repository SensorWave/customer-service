# 🧩 Customer Service

Ce projet est un microservice en **Go** permettant la gestion des **entreprises (companies)** et des **utilisateurs (users)**.  
Il utilise **PostgreSQL** comme base de données, **RabbitMQ** pour la communication asynchrone entre microservices, et tourne entièrement via **Docker Compose**.

## 🚀 Fonctionnalités

- Gestion des **entreprises** (`companies`)
- Gestion des **utilisateurs** (`users`)
- Publication d’événements RabbitMQ (`CompanyCreated`, `UserCreated`)
- Architecture micro-service modulaire et extensible
- Persistance via PostgreSQL

## 🧱 Architecture

```bash
├── cmd/
│   └── main.go         # Point d’entrée de l’application
├── config/
│   └── config.go       # Chargement de la configuration & env
├── internal/
│   ├── customer/
│   │ ├── company/      # Logique métier "Company"
│   │ └── user/         # Logique métier "User"
│   └── event/          # Publication RabbitMQ
├── docker-compose.yml  # Stack complète (Go + Postgres + RabbitMQ)
├── Dockerfile          # Build du microservice Go
├── go.mod / go.sum     # Dépendances Go
└── .env.example        # Variables d’environnement
```

## ⚙️ Prérequis

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- (Optionnel) [psql](https://www.postgresql.org/download/) pour se connecter à la DB

## 🔧 Configuration

Crée un fichier `.env` à la racine du projet et remplir les variables :

```bash
cp .env.example .env
```

## 🐳 Lancer le projet

Depuis la racine du projet :

```bash
docker-compose up --build
```

## 🧠 Utilisation

📋 API (via Postman, curl, etc.)

| Méthode | Endpoint         | Description                              |
|----------|------------------|------------------------------------------|
| POST     | `/companies`       | Créer une entreprise                     |
| GET      | `/companies/:id`   | Récupérer une entreprise                 |
| POST     | `/users`           | Créer un utilisateur lié à une entreprise |
| GET      | `/users/:id`       | Récupérer un utilisateur                 |

## 🗄️ Accès à PostgreSQL

Depuis ton terminal :

```bash
psql -h localhost -p 5433 -U user -d customers
```

Puis :

```sql
\dt -- Liste des tables
SELECT * FROM companies;
SELECT * FROM users;
```

## 🧹 Nettoyage

Pour arrêter et supprimer les containers :

```bash
docker-compose down
```

Pour tout supprimer (containers + volumes + images) :

```bash
docker-compose down -v --rmi all
```

## 📘 Notes

- Le projet est extensible : tu peux ajouter d’autres microservices qui consomment les événements RabbitMQ.
- Le code est organisé selon une structure modulaire, facilitant la maintenance et les tests unitaires.