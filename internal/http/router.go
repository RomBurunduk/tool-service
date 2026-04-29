package httpapi

import (
	"net/http"

	"tool-service/internal/middleware/requestid"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handlers) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(requestid.Middleware)

	r.Route("/tools/api", func(r chi.Router) {
		r.Get("/region", h.Region)
		r.Get("/wordstat", h.WordStat)
		r.Get("/phone-prices", h.PhonePrices)
		r.Get("/currency", h.Currency)
	})
	return r
}
