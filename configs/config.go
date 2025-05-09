package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 구조체는 파일 또는 환경 변수에서 설정을 로드합니다.
type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	API struct {
		GeminiKey string `yaml:"geminiKey"`
		Model     string `yaml:"model"`
	} `yaml:"api"`
}

func Load(path string) (*Config, error) {
	var cfg Config

	// 1) 설정 파일 로드 시도
	data, err := os.ReadFile(path)
	if err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("설정 파일 파싱 실패: %w", err)
		}
	}
	// 2) 환경 변수로 오버라이드 (파일 또는 .env 로드)
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		cfg.API.GeminiKey = key
	}
	if model := os.Getenv("GEMINI_MODEL"); model != "" {
		cfg.API.Model = model
	}

	// Server.Address는 필수
	if cfg.Server.Address == "" {
		cfg.Server.Address = ":8081"
	}
	// API 키/모델 유효성 검사
	if cfg.API.GeminiKey == "" {
		return nil, fmt.Errorf("API 키가 설정되지 않았습니다. 환경변수 GEMINI_API_KEY를 확인하세요")
	}
	if cfg.API.Model == "" {
		return nil, fmt.Errorf("모델명이 설정되지 않았습니다. 환경변수 GEMINI_MODEL을 확인하세요")
	}

	return &cfg, nil
}
