package importer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"tool-service/internal/model"

	"github.com/samber/lo"
)

// ImportMobileCSV reads mobile-style CSV (header row with Brand, Model, Price) and returns rows for DB insert.
func ImportMobileCSV(path string) ([]model.Phone, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseMobileCSV(f)
}

func ParseMobileCSV(r io.Reader) ([]model.Phone, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = false
	cr.LazyQuotes = true
	cr.FieldsPerRecord = -1

	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("csv header: %w", err)
	}

	var out []model.Phone
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read: %w", err)
		}
		row := recordToMap(header, rec)
		brand := strings.TrimSpace(row["Brand"])
		modelName := strings.TrimSpace(row["Model"])
		if brand == "" && modelName == "" {
			continue
		}
		var pricePtr *float64
		if pStr := strings.TrimSpace(row["Price"]); pStr != "" {
			if v, err := strconv.ParseFloat(strings.ReplaceAll(pStr, ",", "."), 64); err == nil {
				pricePtr = &v
			}
		}
		payload := lo.OmitByKeys(row, []string{"Brand", "Model", "Price"})
		raw, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		out = append(out, model.Phone{
			Brand:   brand,
			Model:   modelName,
			Price:   pricePtr,
			Payload: raw,
		})
	}
	return out, nil
}

func recordToMap(header, rec []string) map[string]string {
	out := make(map[string]string, len(header))
	for i, h := range header {
		key := strings.TrimSpace(h)
		val := ""
		if i < len(rec) {
			val = rec[i]
		}
		out[key] = val
	}
	return out
}
