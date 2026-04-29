package httpapi

import (
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	regionsvc "tool-service/internal/service/region"
)

func TestRegionSmoke(t *testing.T) {
	rnd := rand.New(rand.NewSource(42))
	h := &Handlers{
		RegionSvc: regionsvc.NewWithRand(nil, rnd),
	}
	srv := httptest.NewServer(NewRouter(h))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/tools/api/region")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	b, _ := io.ReadAll(resp.Body)
	if len(b) == 0 {
		t.Fatal("empty body")
	}
}
