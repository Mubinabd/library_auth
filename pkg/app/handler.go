package app

import (
	"log/slog"

	kafka "github.com/Mubinabd/library_auth/pkg/kafka/consumer"
	"github.com/Mubinabd/library_auth/service"
)

func Reader(brokers []string, kcm *kafka.KafkaConsumerManager, authService *service.AuthService, userService *service.UserService) {

	if err := kcm.RegisterConsumer(brokers, "reg-user", "auth", kafka.UserRegisterHandler(authService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			slog.Warn("Consumer for topic 'reg-user' already exists")
		} else {
			slog.Error("Error registering consumer: %v", "err", err)
		}
	}

	if err := kcm.RegisterConsumer(brokers, "upd-user", "auth", kafka.UserEditProfileHandler(userService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			slog.Warn("Consumer for topic 'upd-user' already exists")
		} else {
			slog.Error("Error registering consumer: %v", "err", err)
		}
	}

	if err := kcm.RegisterConsumer(brokers, "upd-pass", "auth", kafka.UserEditPasswordHandler(userService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			slog.Warn("Consumer for topic 'upd-pass' already exists")
		} else {
			slog.Error("Error registering consumer: %v", "err", err)
		}
	}

	if err := kcm.RegisterConsumer(brokers, "upd-setting", "auth", kafka.UserEditSettingHandler(userService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			slog.Warn("Consumer for topic 'upd-setting' already exists")
		} else {
			slog.Error("Error registering consumer: %v", "err", err)
		}
	}
}
