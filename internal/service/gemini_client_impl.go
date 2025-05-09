package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spaghetti-people/khuthon-ai/configs"
	"github.com/spaghetti-people/khuthon-ai/internal/model"
)

type geminiClientImpl struct {
	cfg        *configs.Config
	httpClient *http.Client
}

// NewGeminiClient 생성자
func NewGeminiClient(cfg *configs.Config, client *http.Client) GeminiClient {
	return &geminiClientImpl{
		cfg:        cfg,
		httpClient: client,
	}
}

// GenerateContent Gemini API 호출 후 반드시 {"message": "..."} 형태로 리턴
func (c *geminiClientImpl) GenerateContent(ctx context.Context, prompt string) (string, error) {
	// 1) 요청 본문 생성
	reqBody := model.BuildGeminiRequest(prompt, nil, 0.7)

	// 2) 엔드포인트 URL
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		c.cfg.API.Model,
		c.cfg.API.GeminiKey,
	)

	// 3) HTTP 요청
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 4) 요청 실행
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 5) 응답 바디 읽기
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// 6) GeminiResponse 언마샬
	var geminiResp model.GeminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// 7) 텍스트 추출 및 전처리
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini API")
	}
	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	responseText = strings.TrimSpace(responseText)
	responseText = strings.ReplaceAll(responseText, "`", "")

	// 만약 "json\n" 프리픽스가 붙어 있다면 제거
	responseText = strings.TrimPrefix(responseText, "json\n")

	// 8) 응답 그대로 반환
	return responseText, nil
}
