package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr string

	DatabaseURL string

	WordStatAPIKey      string
	WordStatFolderID    string
	WordStatBaseURL     string
	WordStatPath        string
	WordStatHTTPTimeout time.Duration

	CurrencyAPIKey      string
	CurrencyAPIURL      string
	CurrencyHTTPTimeout time.Duration

	ImporterCSVPath string
}

func Load() (Config, error) {
	var c Config
	c.HTTPAddr = getenv("HTTP_ADDR", ":8081")
	c.DatabaseURL = strings.TrimSpace(os.Getenv("DATABASE_URL"))

	c.WordStatAPIKey = strings.TrimSpace(os.Getenv("WORDSTAT_API_KEY"))
	c.WordStatFolderID = strings.TrimSpace(os.Getenv("WORDSTAT_FOLDER_ID"))
	c.WordStatBaseURL = strings.TrimSpace(getenv("WORDSTAT_BASE_URL", "https://searchapi.api.cloud.yandex.net"))
	c.WordStatPath = strings.TrimSpace(getenv("WORDSTAT_PATH", "/v2/wordstat/dynamics"))
	c.WordStatHTTPTimeout = durationEnv("WORDSTAT_HTTP_TIMEOUT", 30*time.Second)

	c.CurrencyAPIKey = strings.TrimSpace(os.Getenv("CURRENCY_API_KEY"))
	c.CurrencyAPIURL = strings.TrimSpace(getenv(
		"CURRENCY_API_URL",
		"https://api.currencyapi.com/v3/latest?base_currency=RUB&currencies=EUR,USD,CNY",
	))
	c.CurrencyHTTPTimeout = durationEnv("CURRENCY_HTTP_TIMEOUT", 15*time.Second)

	c.ImporterCSVPath = strings.TrimSpace(getenv("CSV_PATH", "./mobile.csv"))

	if c.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if c.WordStatFolderID == "" {
		return Config{}, fmt.Errorf("WORDSTAT_FOLDER_ID is required")
	}
	if c.WordStatAPIKey == "" {
		return Config{}, fmt.Errorf("WORDSTAT_API_KEY is required")
	}
	if c.CurrencyAPIKey == "" {
		return Config{}, fmt.Errorf("CURRENCY_API_KEY is required")
	}

	return c, nil
}

func LoadImporter() (Config, error) {
	var c Config
	c.DatabaseURL = strings.TrimSpace(os.Getenv("DATABASE_URL"))
	c.ImporterCSVPath = strings.TrimSpace(getenv("CSV_PATH", "./mobile.csv"))
	if c.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if c.ImporterCSVPath == "" {
		return Config{}, fmt.Errorf("CSV_PATH is required")
	}
	return c, nil
}

func getenv(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func durationEnv(key string, def time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

func MustParseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
