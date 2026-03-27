package main

import (
	"customer-service/config"
	"customer-service/internal/company"
	"customer-service/internal/event"
	"customer-service/internal/user"
	"customer-service/routes"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Invalid configuration:", err)
	}

	// PostgreSQL connection
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

	// RabbitMQ connection with retry
	var conn *amqp.Connection
	for range 10 {
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

	// Declare exchange
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

	// Router
	r := gin.Default()
	r.Use(routes.WithConfig(cfg), routes.CORSMiddleware())

	// Fetch trusted proxies from .env
	proxies := os.Getenv("TRUSTED_PROXIES")
	if proxies != "" {
		proxyList := strings.Split(proxies, ",")
		if err := r.SetTrustedProxies(proxyList); err != nil {
			panic("invalid trusted proxies: " + err.Error())
		}
	}

	// Register all routes
	routes.RegisterRoutes(r, companyHandler, userHandler)

	log.Println("Customer Service running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
