package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	mConfig "github.com/RichardKnop/machinery/v1/config"
)

type AppConfig struct {
	Email        EmailConfig
	Security     SecurityConfig
	Database     DatabaseConfig
	Job          JobConfig
	MessageQueue MessageQueueConfig
}

type MessageQueueConfig struct {
	URI string
}

type JobConfig struct {
	Enabled   bool
	Broker    string
	QueueName string
	Store     string
}

func (cfg *JobConfig) GetJobConfig() mConfig.Config {
	return mConfig.Config{
		Broker:        cfg.Broker,
		DefaultQueue:  cfg.QueueName,
		ResultBackend: cfg.Store,
		/*
			AMQP: &mConfig.AMQPConfig{
				Exchange:     "machinery_exchange",
				ExchangeType: "direct",
				BindingKey:   "machinery_task",
			} */
	}
}

type EmailConfig struct {
	Host               string
	Port               int
	From               string
	Username           string
	Password           string
	SiteURL            string
	InsecureSkipVerify bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	SslMode  string
}

func (cfg DatabaseConfig) GetDatabaseStandardUri() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
}
func (cfg DatabaseConfig) GetDatabaseConnectionUri() string {
	connection := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Name,
		cfg.Port,
		cfg.SslMode,
	)
	return connection
}

func GetLogLevel(level string) logrus.Level {
	switch level {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.WarnLevel

	}
}

type SecurityConfig struct {
	AccessTokenExpiration time.Duration
	Secret                []byte
}

func (c MessageQueueConfig) GetMessageQueueClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: c.URI,
	})

	_, err := client.Ping(context.Background()).Result()
	if !errors.Is(err, nil) {
		return nil, err
	}

	return client, nil
}
