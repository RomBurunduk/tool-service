package region

import (
	"math/rand"
)

var cities = []string{
	"г. Москва",
	"г. Санкт-Петербург",
	"г. Новосибирск",
	"г. Екатеринбург",
	"г. Казань",
	"г. Нижний Новгород",
}

// Pick returns a random region name or empty string (user may hide geolocation).
func Pick(r *rand.Rand) string {
	if r.Intn(5) == 0 { // ~20% empty
		return ""
	}
	return cities[r.Intn(len(cities))]
}
