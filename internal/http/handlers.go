package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"tool-service/internal/service/currency"
	"tool-service/internal/service/phone"
	"tool-service/internal/service/region"
	"tool-service/internal/service/wordstat"
)

type Handlers struct {
	WordStatSvc *wordstat.Service
	PhoneSvc    *phone.Service
	CurrencySvc *currency.Service
	RegionSvc   *region.Service
}

type errJSON struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handlers) Region(w http.ResponseWriter, r *http.Request) {
	out, err := h.RegionSvc.Region(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handlers) WordStat(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	out, err := h.WordStatSvc.Get(r.Context(), q)
	if err != nil {
		if errors.Is(err, wordstat.ErrQueryRequired) {
			writeJSON(w, http.StatusBadRequest, errJSON{Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusBadGateway, errJSON{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handlers) PhonePrices(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	out, err := h.PhoneSvc.PriceByQuery(r.Context(), q)
	if err != nil {
		switch {
		case errors.Is(err, phone.ErrQueryRequired):
			writeJSON(w, http.StatusBadRequest, errJSON{Error: err.Error()})
		case errors.Is(err, phone.ErrNotFound), errors.Is(err, phone.ErrNoPrice):
			writeJSON(w, http.StatusNotFound, errJSON{Error: err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, errJSON{Error: err.Error()})
		}
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handlers) Currency(w http.ResponseWriter, r *http.Request) {
	out, err := h.CurrencySvc.Rates(r.Context())
	if err != nil {
		writeJSON(w, http.StatusBadGateway, errJSON{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, out)
}
