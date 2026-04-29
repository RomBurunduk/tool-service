package currency

import (
	"math"
	"testing"

	"tool-service/internal/model"
)

func TestToOpenAPIRates(t *testing.T) {
	resp := LatestAPIResponse{
		Data: map[string]struct {
			Code  string  `json:"code"`
			Value float64 `json:"value"`
		}{},
	}
	resp.Data["CNY"] = struct {
		Code  string  `json:"code"`
		Value float64 `json:"value"`
	}{Code: "CNY", Value: 0.0911230312}
	resp.Data["EUR"] = struct {
		Code  string  `json:"code"`
		Value float64 `json:"value"`
	}{Code: "EUR", Value: 0.0113932361}
	resp.Data["USD"] = struct {
		Code  string  `json:"code"`
		Value float64 `json:"value"`
	}{Code: "USD", Value: 0.0133556624}

	got, err := ToOpenAPIRates(resp)
	if err != nil {
		t.Fatal(err)
	}
	want := model.CurrencyRates{
		CNY: 1 / 0.0911230312,
		EUR: 1 / 0.0113932361,
		USD: 1 / 0.0133556624,
	}
	const eps = 1e-9
	if math.Abs(got.CNY-want.CNY) > eps || math.Abs(got.EUR-want.EUR) > eps || math.Abs(got.USD-want.USD) > eps {
		t.Fatalf("got %+v want %+v", got, want)
	}
}
