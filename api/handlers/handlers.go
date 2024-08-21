package handlers

import (
	kafka "github.com/Mubinabd/library_auth/pkg/kafka/producer"
	"github.com/Mubinabd/library_auth/service"
	"github.com/go-redis/redis/v8"
)

type Handlers struct {
	Auth     *service.AuthService
	User     *service.UserService
	RDB      *redis.Client
	Producer kafka.KafkaProducer
}

func NewHandler(auth *service.AuthService, user *service.UserService, rdb *redis.Client, pr *kafka.KafkaProducer) *Handlers {
	return &Handlers{
		Auth:     auth,
		User:     user,
		RDB:      rdb,
		Producer: *pr,
	}
}
