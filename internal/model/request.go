package model

import (
	"encoding/json"
)

type Part struct {
	Text string `json:"text,omitempty"`
}
type ContentWrapper struct {
	Parts []Part `json:"parts"`
}
type GenerationConfig struct {
	ResponseMimeType string      `json:"responseMimeType"`
	ResponseSchema   interface{} `json:"responseSchema"`
	Temperature      float64     `json:"temperature"`
}

type GeminiRequest struct {
	Contents         []ContentWrapper `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig"`
}

func BuildGeminiRequest(prompt string, schema map[string]interface{}, temperature float64) []byte {
	req := GeminiRequest{
		Contents: []ContentWrapper{{Parts: []Part{{Text: prompt}}}},
		GenerationConfig: GenerationConfig{
			Temperature: temperature,
		},
	}
	b, _ := json.Marshal(req)
	return b
}
