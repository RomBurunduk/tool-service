package wordstat

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	wsclient "tool-service/internal/clients/wordstat"
	"tool-service/internal/config"
	"tool-service/internal/model"
	toollogsvc "tool-service/internal/service/toollog"
)

var ErrQueryRequired = errors.New("query is required")

type Service struct {
	client   *wsclient.Client
	folderID string
	tl       *toollogsvc.Service
}

func New(cfg config.Config, tl *toollogsvc.Service) *Service {
	return &Service{
		client:   wsclient.NewClient(cfg),
		folderID: cfg.WordStatFolderID,
		tl:       tl,
	}
}

func (s *Service) Get(ctx context.Context, query string) (pair model.WordStatPair, err error) {
	start := time.Now()
	q := strings.TrimSpace(query)
	input, _ := json.Marshal(map[string]string{"query": q})
	defer func() {
		if s.tl == nil {
			return
		}
		s.tl.Record(ctx, "wordstat", input, pair, err, start)
	}()

	if q == "" {
		return model.WordStatPair{}, ErrQueryRequired
	}

	from := FromDateLastThreeMonths(time.Now())
	usedPhrase := q + " БУ"

	used, err := s.client.Dynamics(ctx, usedPhrase, s.folderID, from)
	if err != nil {
		return model.WordStatPair{}, err
	}
	newRes, err := s.client.Dynamics(ctx, q, s.folderID, from)
	if err != nil {
		return model.WordStatPair{}, err
	}
	return model.WordStatPair{Used: used, New: newRes}, nil
}
