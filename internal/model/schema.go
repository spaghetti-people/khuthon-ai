// model/schema.go
package model

var QuizSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"title":       map[string]any{"type": "string"},
		"description": map[string]any{"type": "string"},
		"question_type": map[string]any{
			"type": "string",
			"enum": []string{"multiple_choice", "short_answer", "true_false"},
		},
		"hint":   map[string]any{"type": "string"},
		"answer": map[string]any{"type": "string"},
		"point":  map[string]any{"type": "integer", "default": 10},
	},
	"required": []string{"title", "description", "question_type", "hint", "answer", "point"},
}

var CommentarySchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]any{
		"commentary": map[string]any{"type": "string"},
	},
	"required": []string{"summary"},
}
