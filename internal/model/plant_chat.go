package model

// PlantInfo 식물의 기본 정보를 담는 구조체
type PlantInfo struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	Environment      string `json:"environment"`
	Characteristics  string `json:"characteristics"`
	CareInstructions string `json:"careInstructions"`
	Health           string `json:"health"`
	Temperature      string `json:"temperature"`
	Humidity         string `json:"humidity"`
}

// Message 대화 메시지를 담는 구조체
type Message struct {
	UserID  string `json:"user_id"`
	Role    string `json:"role"`
	Content string `json:"content"`
	Time    string `json:"time,omitempty"` // 메시지 생성 시간
}

// PlantChatRequest 식물 채팅 요청을 담는 구조체
type PlantChatRequest struct {
	UserID              string       `json:"user_id"`
	PlantInfo           PlantInfo    `json:"plant_info"`
	UserMessage         string       `json:"user_message"`
	ChildProfile        ChildProfile `json:"child_profile"`
	ConversationHistory []Message    `json:"conversationHistory,omitempty"`
}

// PlantChatResponse 식물 채팅 응답을 담는 구조체
type PlantChatResponse struct {
	Message string `json:"message"`
}
