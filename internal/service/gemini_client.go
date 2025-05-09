package service

import "context"

// GeminiClient Gemini API와 통신하기 위한 인터페이스
type GeminiClient interface {
	// GenerateContent Gemini API를 호출하여 응답을 생성합니다
	GenerateContent(ctx context.Context, prompt string) (string, error)
}
