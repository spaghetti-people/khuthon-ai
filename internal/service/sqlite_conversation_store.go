package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/spaghetti-people/khuthon-ai/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteConversationStore SQLite를 사용하는 대화 저장소 구현체
type SQLiteConversationStore struct {
	db *sql.DB
}

// NewSQLiteConversationStore 새로운 SQLite 대화 저장소를 생성합니다
func NewSQLiteConversationStore(dbPath string) (*SQLiteConversationStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 테이블 생성
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS conversations (
			plant_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			time TEXT NOT NULL,
			PRIMARY KEY (plant_id, user_id, time)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SQLiteConversationStore{db: db}, nil
}

// SaveConversation 대화 기록을 저장합니다
func (s *SQLiteConversationStore) SaveConversation(plantID string, messages []model.Message) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 기존 대화 기록 삭제
	_, err = tx.Exec("DELETE FROM conversations WHERE plant_id = ? AND user_id = ?", plantID, messages[0].UserID)
	if err != nil {
		return fmt.Errorf("failed to delete old conversations: %w", err)
	}

	// 새로운 대화 기록 저장
	stmt, err := tx.Prepare(`
		INSERT INTO conversations (plant_id, user_id, role, content, time)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// 각 메시지마다 약간의 시간 차이를 두어 저장
	baseTime := time.Now()
	for i, msg := range messages {
		// 각 메시지마다 1초씩 차이를 둠
		msgTime := baseTime.Add(time.Duration(i) * time.Second)
		_, err = stmt.Exec(plantID, msg.UserID, msg.Role, msg.Content, msgTime.Format(time.RFC3339))
		if err != nil {
			return fmt.Errorf("failed to insert message: %w", err)
		}
	}

	return tx.Commit()
}

// GetConversation 저장된 대화 기록을 조회합니다
func (s *SQLiteConversationStore) GetConversation(plantID string) ([]model.Message, error) {
	rows, err := s.db.Query(`
		SELECT user_id, role, content, time
		FROM conversations
		WHERE plant_id = ?
		ORDER BY time ASC
	`, plantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations: %w", err)
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message
		err := rows.Scan(&msg.UserID, &msg.Role, &msg.Content, &msg.Time)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// Close 데이터베이스 연결을 종료합니다
func (s *SQLiteConversationStore) Close() error {
	return s.db.Close()
}
