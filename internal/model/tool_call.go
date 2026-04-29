package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ToolCallStatus string

const (
	ToolCallSuccess ToolCallStatus = "success"
	ToolCallError   ToolCallStatus = "error"
)

type ToolCall struct {
	RequestID  uuid.UUID
	ToolName   string
	Input      json.RawMessage
	Output     json.RawMessage
	Status     ToolCallStatus
	Error      *string
	DurationMs int
}
