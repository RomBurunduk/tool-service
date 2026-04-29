package toollog

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"tool-service/internal/middleware/requestid"
	"tool-service/internal/model"
	toollogrepo "tool-service/internal/repository/toollog"
)

type Service struct {
	repo *toollogrepo.Repository
	log  *slog.Logger
}

func New(repo *toollogrepo.Repository, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	return &Service{repo: repo, log: log}
}

// Record schedules audit for a tool invocation (success or error). Non-blocking.
func (s *Service) Record(parent context.Context, toolName string, input json.RawMessage, output any, callErr error, start time.Time) {
	dur := int(time.Since(start).Milliseconds())
	tc := model.ToolCall{
		ToolName:   toolName,
		Input:      input,
		DurationMs: dur,
	}
	if callErr != nil {
		tc.Status = model.ToolCallError
		es := callErr.Error()
		tc.Error = &es
	} else {
		tc.Status = model.ToolCallSuccess
		if output == nil {
			tc.Output = []byte("null")
		} else {
			raw, err := json.Marshal(output)
			if err != nil {
				tc.Output = []byte(`{"error":"marshal_output"}`)
			} else {
				tc.Output = raw
			}
		}
	}
	s.LogAsync(parent, tc)
}

func (s *Service) LogAsync(parent context.Context, tc model.ToolCall) {
	rid, ok := requestid.FromContext(parent)
	if !ok {
		return
	}
	tc.RequestID = rid

	go func(tc model.ToolCall) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := s.repo.Insert(ctx, tc); err != nil {
			s.log.Error("toollog insert failed", "err", err, "tool", tc.ToolName, "request_id", tc.RequestID.String())
		}
	}(tc)
}
