package main

import (
	"log"
	"sycretru"
	"sycretru/internal/config"
	"sycretru/internal/handler"
	"sycretru/internal/service"
	"sycretru/pkg/logger"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	logger := logger.New()

	srv := &sycretru.Server{}

	svc := service.NewService()

	cfg := config.Get()

	handlers := handler.NewHandler(svc, logger)
	srv.Run(cfg.Host, cfg.Port, handlers.InitRoutes())

}
