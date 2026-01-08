package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App  *App
		GRPC *GRPC
		DB   *DB
	}

	App struct {
		Name string
		Env  string
	}

	GRPC struct {
		Host           string
		Port           string
		AllowedOrigins string
	}

	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
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

	GRPC := &GRPC{
		Host:           os.Getenv("GRPC_HOST"),
		Port:           os.Getenv("GRPC_PORT"),
		AllowedOrigins: os.Getenv("GRPC_ALLOWED_ORIGINS"),
	}

	DB := &DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	return &Container{
		App:  App,
		GRPC: GRPC,
		DB:   DB,
	}, nil
}
