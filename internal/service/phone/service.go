package phone

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"tool-service/internal/model"
	phonerepo "tool-service/internal/repository/phone"
	toollogsvc "tool-service/internal/service/toollog"
)

var (
	ErrMissingParams = errors.New("brand and model are required")
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

func (s *Service) PriceByBrandModel(ctx context.Context, brand, modelName string) (out model.PhonePriceResponse, err error) {
	start := time.Now()
	input, _ := json.Marshal(map[string]string{"brand": brand, "model": modelName})
	defer func() {
		if s.tl == nil {
			return
		}
		s.tl.Record(ctx, "phone-prices", input, out, err, start)
	}()

	if brand == "" || modelName == "" {
		return model.PhonePriceResponse{}, ErrMissingParams
	}
	p, err := s.repo.GetByBrandModel(ctx, brand, modelName)
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
