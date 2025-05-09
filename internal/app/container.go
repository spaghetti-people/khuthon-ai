package app

import (
	"net/http"

	"github.com/spaghetti-people/khuthon-ai/configs"
	"github.com/spaghetti-people/khuthon-ai/internal/service"
)

type Container struct {
	PlantChat service.PlantChatService
}

func NewContainer(cfg *configs.Config, client *http.Client) (*Container, error) {
	geminiClient := service.NewGeminiClient(cfg, client)
	plantChat, err := service.NewPlantChatService(geminiClient)
	if err != nil {
		return nil, err
	}

	return &Container{
		PlantChat: plantChat,
	}, nil
}

// Close 컨테이너의 리소스를 정리합니다
func (c *Container) Close() error {
	return c.PlantChat.Close()
}
