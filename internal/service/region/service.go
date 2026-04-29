package region

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	regionmock "tool-service/internal/mock/region"
	"tool-service/internal/model"
	toollogsvc "tool-service/internal/service/toollog"
)

type Service struct {
	rnd *rand.Rand
	tl  *toollogsvc.Service
}

func New(tl *toollogsvc.Service) *Service {
	return &Service{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
		tl:  tl,
	}
}

func NewWithRand(tl *toollogsvc.Service, rnd *rand.Rand) *Service {
	return &Service{rnd: rnd, tl: tl}
}

func (s *Service) Region(ctx context.Context) (out model.RegionResponse, err error) {
	start := time.Now()
	input := json.RawMessage([]byte("{}"))
	defer func() {
		if s.tl == nil {
			return
		}
		s.tl.Record(ctx, "region", input, out, err, start)
	}()

	out.Region = regionmock.Pick(s.rnd)
	return out, nil
}
