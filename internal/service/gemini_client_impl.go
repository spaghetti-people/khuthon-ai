package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spaghetti-people/khuthon-ai/configs"
	"github.com/spaghetti-people/khuthon-ai/internal/model"
)

type geminiClientImpl struct {
	cfg        *configs.Config
	httpClient *http.Client
}

func NewGeminiClient(cfg *configs.Config, client *http.Client) GeminiClient {
	return &geminiClientImpl{
		cfg:        cfg,
		httpClient: client,
	}
}

func (c *geminiClientImpl) GenerateContent(ctx context.Context, prompt string) (string, error) {
	// 프롬프트를 Gemini API 요청 형식으로 변환
	reqBody := model.BuildGeminiRequest(prompt, nil, 0.7)

	// API 엔드포인트 URL 구성
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		c.cfg.API.Model, c.cfg.API.GeminiKey)

	// HTTP 요청 생성
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 요청 실행
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 응답 읽기
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// 응답 파싱
	var geminiResp model.GeminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// 응답에서 텍스트 추출
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini API")
	}

	return string(geminiResp.Candidates[0].Content.Parts[0].Text), nil
}
