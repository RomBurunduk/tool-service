package model

type Phone struct {
	ID      int64
	Brand   string
	Model   string
	Price   *float64
	Payload []byte // JSON
}

type PhonePriceResponse struct {
	Price float64 `json:"price"`
}
