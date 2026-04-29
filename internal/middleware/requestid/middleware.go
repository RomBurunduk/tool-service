package requestid

import (
	"net/http"

	"github.com/google/uuid"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		w.Header().Set("X-Request-ID", id.String())
		next.ServeHTTP(w, r.WithContext(WithRequestID(r.Context(), id)))
	})
}
