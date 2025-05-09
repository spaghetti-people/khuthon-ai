package model

import "encoding/json"

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text json.RawMessage `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}
