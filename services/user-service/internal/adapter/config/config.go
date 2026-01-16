package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App  *App
		HTTP *HTTP
		GRPC *GRPC
		Rabbitmq *Rabbitmq
		DB   *DB
		JWT  *JWT
	}

	App struct {
		Name string
		Env  string
	}

	HTTP struct {
		Host           string
		Port           string
		AllowedOrigins string
	}

	GRPC struct {
		Host string
		Port string
	}

	Rabbitmq struct {
		Host string
		Port string
		User string
		Password string
	}

	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	JWT struct {
		Secret   string
		Duration string
	}
)

func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			errMsg := fmt.Errorf("unable to load .env: %v", err.Error())
			return nil, errMsg
		}
	}

	App := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	HTTP := &HTTP{
		Host:           os.Getenv("HTTP_HOST"),
		Port:           os.Getenv("HTTP_PORT_USER"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}

	GRPC := &GRPC{
		Host: os.Getenv("GRPC_HOST"),
		Port: os.Getenv("GRPC_PORT"),
	}

	Rabbitmq := &Rabbitmq{
		Host: os.Getenv("RABBITMQ_HOST"),
		Port: os.Getenv("RABBITMQ_PORT"),
		User: os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
	}

	DB := &DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	JWT := &JWT{
		Secret:   os.Getenv("JWT_SECRET"),
		Duration: os.Getenv("SESSION_DURATION"),
	}

	return &Container{
		App:  App,
		HTTP: HTTP,
		GRPC: GRPC,
		Rabbitmq: Rabbitmq,
		DB:   DB,
		JWT:  JWT,
	}, nil
}
