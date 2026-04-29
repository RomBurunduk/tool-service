package currency

import (
	"fmt"

	"tool-service/internal/model"
)

// ToOpenAPIRates converts CurrencyAPI values (foreign currency per 1 RUB when base=RUB)
// to rubles per one unit of foreign currency, matching OpenAPI examples.
func ToOpenAPIRates(resp LatestAPIResponse) (model.CurrencyRates, error) {
	out := model.CurrencyRates{}
	for code, field := range map[string]*float64{
		"CNY": &out.CNY,
		"EUR": &out.EUR,
		"USD": &out.USD,
	} {
		entry, ok := resp.Data[code]
		if !ok {
			return model.CurrencyRates{}, fmt.Errorf("currency: missing %s", code)
		}
		if entry.Value <= 0 {
			return model.CurrencyRates{}, fmt.Errorf("currency: non-positive value for %s", code)
		}
		*field = 1 / entry.Value
	}
	return out, nil
}
