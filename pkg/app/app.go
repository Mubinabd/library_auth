package app

import (
	"context"
	"log"

	"github.com/Mubinabd/library_auth/api"
	"github.com/Mubinabd/library_auth/api/handlers"
	"github.com/Mubinabd/library_auth/config"
	kafka "github.com/Mubinabd/library_auth/pkg/kafka/consumer"
	prd "github.com/Mubinabd/library_auth/pkg/kafka/producer"
	"github.com/Mubinabd/library_auth/pkg/storage/postgres"
	"github.com/Mubinabd/library_auth/service"
	"github.com/go-redis/redis/v8"
)

func Run(cfg *config.Config) {

    db, err := postgres.NewPostgresStorage(cfg)
    if err != nil {
        log.Printf("can't connect to db: %v", err)
    }
    log.Println("Connected to Postgres")

    authService := service.NewAuthService(db)
    userService := service.NewUserService(db)


	// Redis Connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Panicf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")


	// Kafka
	brokers := []string{"kafka:9092"}
	cm := kafka.NewKafkaConsumerManager()
	pr, err := prd.NewKafkaProducer(brokers)
	if err != nil {
		log.Println("Failed to create Kafka producer:", err)
		return
	}

	Reader(brokers, cm, authService, userService)

	// HTTP Server
	h := handlers.NewHandler(authService, userService, rdb, &pr)

	router := api.Engine(h)
	router.SetTrustedProxies(nil)

	if err := router.Run(cfg.AUTH_PORT); err != nil {
		log.Panicf("can't start server: %v", err)
	}

	log.Printf("REST server started on port %s", cfg.AUTH_PORT)
}
