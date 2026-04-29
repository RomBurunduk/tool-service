package currency

import (
	"context"
	"encoding/json"
	"time"

	curclient "tool-service/internal/clients/currency"
	"tool-service/internal/config"
	"tool-service/internal/model"
	toollogsvc "tool-service/internal/service/toollog"
)

type Service struct {
	client *curclient.Client
	tl     *toollogsvc.Service
}

func New(cfg config.Config, tl *toollogsvc.Service) *Service {
	return &Service{
		client: curclient.NewClient(cfg),
		tl:     tl,
	}
}

func (s *Service) Rates(ctx context.Context) (rates model.CurrencyRates, err error) {
	start := time.Now()
	input := json.RawMessage([]byte("{}"))
	defer func() {
		if s.tl == nil {
			return
		}
		s.tl.Record(ctx, "currency", input, rates, err, start)
	}()

	raw, err := s.client.FetchLatest(ctx)
	if err != nil {
		return model.CurrencyRates{}, err
	}
	return curclient.ToOpenAPIRates(raw)
}
