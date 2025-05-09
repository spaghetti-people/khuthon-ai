package model

// ChildProfile 아이 프로필 정보
type ChildProfile struct {
	Name           string   `json:"name"`            // 아이 이름
	Age            int      `json:"age"`             // 나이
	FavoriteTopics []string `json:"favorite_topics"` // 관심사
	FavoriteEmojis []string `json:"favorite_emojis"` // 좋아하는 이모티콘
	LanguageLevel  string   `json:"language_level"`  // 언어 수준 (초급/중급/고급)
	Personality    string   `json:"personality"`     // 성격
	LearningStyle  string   `json:"learning_style"`  // 학습 스타일
	SpecialNeeds   []string `json:"special_needs"`   // 특별한 요구사항
}
