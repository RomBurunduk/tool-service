package requestid

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey int

const keyRequestID ctxKey = 1

func WithRequestID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

func FromContext(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(keyRequestID)
	if v == nil {
		return uuid.UUID{}, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}
