package main

import (
	"customer-service/config"
	company "customer-service/internal/company"
	"customer-service/internal/event"
	user "customer-service/internal/user"
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func main() {
	cfg := config.LoadConfig()

	// Connexion à PostgreSQL
	log.Println("Connecting to Postgres with DSN:", cfg.PostgresDSN())
	db, err := sql.Open("postgres", cfg.PostgresDSN())
	if err != nil {
		log.Fatal("Cannot open PostgreSQL:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to PostgreSQL:", err)
	}
	log.Println("✅ Connected to PostgreSQL")

	// Connexion à RabbitMQ avec retry
	var conn *amqp.Connection
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(cfg.RabbitMQURL())
		if err == nil {
			break
		}
		log.Println("RabbitMQ not ready, retrying in 3s...")
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ after retries:", err)
	}
	defer conn.Close()
	log.Println("✅ Connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open RabbitMQ channel:", err)
	}
	defer ch.Close()

	// Déclarer l'échange pour les événements
	err = ch.ExchangeDeclare("customer_events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare exchange:", err)
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

	// Healthcheck pour Docker
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	companyHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)

	log.Println("Customer Service running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
