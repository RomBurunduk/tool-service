package phone

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"tool-service/internal/model"
	phonerepo "tool-service/internal/repository/phone"
	toollogsvc "tool-service/internal/service/toollog"
)

var (
	ErrQueryRequired = errors.New("query is required")
	ErrNotFound      = errors.New("phone not found")
	ErrNoPrice       = errors.New("price not available for this phone")
)

type Service struct {
	repo *phonerepo.Repository
	tl   *toollogsvc.Service
}

func New(repo *phonerepo.Repository, tl *toollogsvc.Service) *Service {
	return &Service{repo: repo, tl: tl}
}

func (s *Service) PriceByQuery(ctx context.Context, query string) (out model.PhonePriceResponse, err error) {
	start := time.Now()
	q := strings.TrimSpace(query)
	input, _ := json.Marshal(map[string]string{"query": q})
	defer func() {
		if s.tl == nil {
			return
		}
		s.tl.Record(ctx, "phone-prices", input, out, err, start)
	}()

	if q == "" {
		return model.PhonePriceResponse{}, ErrQueryRequired
	}
	p, err := s.repo.GetBestMatchByQuery(ctx, q)
	if err != nil {
		return model.PhonePriceResponse{}, err
	}
	if p == nil {
		return model.PhonePriceResponse{}, ErrNotFound
	}
	if p.Price == nil {
		return model.PhonePriceResponse{}, ErrNoPrice
	}
	return model.PhonePriceResponse{Price: *p.Price}, nil
}
