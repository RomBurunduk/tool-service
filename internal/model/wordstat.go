package model

type WordStatResult struct {
	Results []WordStatPoint `json:"results"`
}

type WordStatPoint struct {
	Date  string  `json:"date"`
	Count string  `json:"count"`
	Share float64 `json:"share"`
}

type WordStatPair struct {
	Used WordStatResult `json:"used"`
	New  WordStatResult `json:"new"`
}
