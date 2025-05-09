package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/spaghetti-people/khuthon-ai/configs"
	"github.com/spaghetti-people/khuthon-ai/internal/api"
	"github.com/spaghetti-people/khuthon-ai/internal/app"
)

func main() {
	// 설정 로드
	_ = godotenv.Load()
	cfg, err := configs.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}
	client := &http.Client{Timeout: 5 * time.Second}

	container, err := app.NewContainer(cfg, client)
	if err != nil {
		log.Fatalf("failed to create container: %v", err)
	}
	defer container.Close()

	router := api.NewRouter(container)

	// 서버 종료 시그널 처리
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		if err := container.Close(); err != nil {
			log.Printf("Error closing container: %v", err)
		}
		os.Exit(0)
	}()

	router.Run(cfg.Server.Address)
	log.Printf("server starting on port %s\n", os.Getenv("PORT"))
}
