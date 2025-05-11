package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spaghetti-people/khuthon-ai/internal/model"
)

const maxConversationHistory = 15

// PlantChatService 식물 채팅 서비스 인터페이스
type PlantChatService interface {
	Chat(ctx context.Context, req *model.PlantChatRequest) (*model.PlantChatResponse, error)
	Close() error
}

// plantChatService 식물 채팅 서비스 구현체
type plantChatService struct {
	geminiClient      GeminiClient
	conversationStore ConversationStore
}

// NewPlantChatService 새로운 식물 채팅 서비스 인스턴스를 생성합니다
func NewPlantChatService(geminiClient GeminiClient) (PlantChatService, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "conversations.db"
	}
	store, err := NewSQLiteConversationStore(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation store: %w", err)
	}

	return &plantChatService{
		geminiClient:      geminiClient,
		conversationStore: store,
	}, nil
}

// Close 서비스를 종료합니다
func (s *plantChatService) Close() error {
	if store, ok := s.conversationStore.(*SQLiteConversationStore); ok {
		return store.Close()
	}
	return nil
}

// Chat 식물과의 대화를 처리합니다
func (s *plantChatService) Chat(ctx context.Context, req *model.PlantChatRequest) (*model.PlantChatResponse, error) {
	// 저장된 대화 기록 조회
	plantID := req.PlantInfo.Name // 식물 이름을 ID로 사용
	savedMessages, err := s.conversationStore.GetConversation(plantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// 새로운 대화 메시지 추가
	updatedMessages := append(savedMessages, model.Message{
		UserID:  req.UserID,
		Role:    "user",
		Content: req.UserMessage,
	})

	// 최근 15개의 메시지만 사용
	if len(updatedMessages) > maxConversationHistory {
		updatedMessages = updatedMessages[len(updatedMessages)-maxConversationHistory:]
	}

	// 프롬프트 템플릿 생성
	prompt := s.buildPrompt(req.PlantInfo, req.ChildProfile, req.DailyData, updatedMessages)
	// 프롬프트 로깅
	fmt.Printf("생성된 프롬프트:\n%s\n", prompt)
	// Gemini API 호출
	response, err := s.geminiClient.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// AI 응답 메시지 추가
	updatedMessages = append(updatedMessages, model.Message{
		UserID:  req.UserID,
		Role:    "assistant",
		Content: response,
	})

	// 대화 기록 저장
	if err := s.conversationStore.SaveConversation(plantID, updatedMessages); err != nil {
		return nil, fmt.Errorf("failed to save conversation: %w", err)
	}

	return &model.PlantChatResponse{
		Message: response,
	}, nil
}

// buildPrompt 요청 정보를 바탕으로 프롬프트를 생성합니다
func (s *plantChatService) buildPrompt(plantInfo model.PlantInfo, childProfile model.ChildProfile, dailyData model.DailyPlantData, messages []model.Message) string {
	var promptBuilder strings.Builder

	// 기본 설정
	promptBuilder.WriteString(fmt.Sprintf("당신은 %s라는 %s의 종류인 식물 케릭터입니다. 다음 규칙을 따라 대화해주세요:\n\n", plantInfo.Name, plantInfo.Type))

	// 대화 규칙
	promptBuilder.WriteString("1. 귀엽고 친근한 말투를 사용하세요. 대화는 간결하게 한문장으로 하세요.\n")
	promptBuilder.WriteString(fmt.Sprintf("2. %s의 이름을 자주 불러주세요\n", childProfile.Name))
	promptBuilder.WriteString(fmt.Sprintf("3. %s의 관심사(%s)를 대화에 자연스럽게 포함하세요\n", childProfile.Name, strings.Join(childProfile.FavoriteTopics, ", ")))
	promptBuilder.WriteString(fmt.Sprintf("4. %s의 나이(%d세)에 맞는 단어와 문장을 사용하세요\n", childProfile.Name, childProfile.Age))
	promptBuilder.WriteString("5. 현재 식물의 상태를 자연스럽게 언급하세요\n")
	promptBuilder.WriteString("6. 반드시 한글로만 응답해주세요\n\n")

	// 현재 상태 정보
	promptBuilder.WriteString("현재 상태:\n")
	promptBuilder.WriteString(fmt.Sprintf("- 온도: %.1f°C\n", dailyData.Conditions.Temperature))
	promptBuilder.WriteString(fmt.Sprintf("- 습도: %.1f%%\n", dailyData.Conditions.Humidity))
	promptBuilder.WriteString(fmt.Sprintf("- 토양 수분: %.1f%%\n", dailyData.Conditions.SoilMoisture))
	promptBuilder.WriteString(fmt.Sprintf("- 일조량: %.1f시간\n", dailyData.Conditions.SunlightExposure))
	promptBuilder.WriteString(fmt.Sprintf("- 해충 압력: %.1f%%\n", dailyData.Conditions.PestPressure))
	promptBuilder.WriteString(fmt.Sprintf("- 서리 위험: %.1f%%\n\n", dailyData.Conditions.FrostRisk))

	// 마지막 사용자 메시지
	if len(messages) > 0 {
		lastMessage := messages[len(messages)-1]
		promptBuilder.WriteString(fmt.Sprintf("사용자: %s\n\n", lastMessage.Content))
	}

	// 응답 형식 지정
	promptBuilder.WriteString("다음과 같은 JSON 형식으로만 응답해주세요. 다른 텍스트는 포함하지 마세요:\n")
	promptBuilder.WriteString("{\"message\": \"대화 내용 (인사말, 현재 상태 언급, 대화, 조언, 마무리 인사를 자연스럽게 포함)\"}")

	return promptBuilder.String()
}
