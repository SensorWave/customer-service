package main

import (
    "customer-service/config"
    company "customer-service/internal/company"
    user "customer-service/internal/user"
    "customer-service/internal/event"
    "database/sql"
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
    "github.com/streadway/amqp"
    "log"
)

func main() {
    cfg := config.LoadConfig()

    // Connexion à PostgreSQL
    db, err := sql.Open("postgres", cfg.PostgresDSN())
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Connexion à RabbitMQ
    conn, err := amqp.Dial(cfg.RabbitMQURL())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatal(err)
    }
    defer ch.Close()

    // Déclarer l'échange pour les événements
    err = ch.ExchangeDeclare("customer_events", "topic", true, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    publisher := event.NewPublisher(ch, "customer_events")

    // Repositories
    companyRepo := company.NewRepository(db)
    userRepo := user.NewRepository(db)

    // Services
    companyService := company.NewService(companyRepo, publisher)
    userService := user.NewService(userRepo, companyRepo, publisher)

    // Handlers
    companyHandler := company.NewHandler(companyService)
    userHandler := user.NewHandler(userService)

    // Router Gin
    r := gin.Default()
    companyHandler.RegisterRoutes(r)
    userHandler.RegisterRoutes(r)

    log.Println("Starting Customer Service on :8080")
    r.Run(":8080")
}
