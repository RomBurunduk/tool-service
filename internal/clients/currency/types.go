package currency

// LatestAPIResponse mirrors currencyapi.com v3 latest JSON (subset).
type LatestAPIResponse struct {
	Meta struct {
		LastUpdatedAt string `json:"last_updated_at"`
	} `json:"meta"`
	Data map[string]struct {
		Code  string  `json:"code"`
		Value float64 `json:"value"`
	} `json:"data"`
}
