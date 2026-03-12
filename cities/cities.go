package cities

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func FetchCities() ([]CityData, error) {
	req, err := http.NewRequest(
		"GET",
		"https://alerts-history.oref.org.il/Shared/Ajax/GetDistricts.aspx",
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error(
				"Could not close request body",
				"error", err,
			)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected status code from Pikud HaOref Districts API %d",
			resp.StatusCode,
		)
	}

	var cities []CityData
	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		return nil, err
	}

	return cities, nil
}

func FetchAllCityNames() ([]string, error) {
	allCities, err := FetchCities()
	if err != nil {
		return nil, err
	}

	cityNames := make([]string, len(allCities))
	for _, city := range allCities {
		cityNames = append(cityNames, city.Label)
	}

	return cityNames, nil
}
