package main

import (
	"github.com/blazee5/imageChecker/internal/config"
	"github.com/blazee5/imageChecker/internal/handler"
	"github.com/blazee5/imageChecker/internal/repository/redis"
	"github.com/blazee5/imageChecker/internal/service"
	redisLib "github.com/blazee5/imageChecker/lib/db/redis"
	"github.com/blazee5/imageChecker/lib/logger"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()
	logger := logger.NewLogger()

	rdb := redisLib.NewRedisClient(redisLib.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
	})

	r := gin.Default()

	repository := redis.NewRepository(rdb)
	service := service.NewService(logger, repository)
	handlers := handler.NewHandler(logger, service)

	handler.RegisterHandlers(r, handlers)

	slog.Info("Server starting...")

	go func() {
		if err := r.Run(":" + cfg.HTTPServer.Port); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	slog.Info("Server exiting...")

	rdb.Close()
}
