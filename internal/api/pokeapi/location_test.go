package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"poke-repl/internal/config"
	"reflect"
	"testing"
)

func TestGetLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := LocationResult{
			Count:    1,
			Next:     "https://pokeapi.co/api/v2/location/?offset=20&limit=20",
			Previous: nil,
			Results: []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{
				{
					Name: "location1",
					URL:  "https://pokeapi.co/api/v2/location/1/",
				},
			},
		}

		responseJSON, _ := json.Marshal(response)

		w.Write(responseJSON)
	}))
	defer server.Close()

	cfg := &config.Config{
		NextUrl:     "",
		PreviousUrl: "",
	}

	locationResult := &LocationResult{}

	locations, err := locationResult.GetLocation(server.URL, cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedLocations := LocationList{
		LocationInfo{Name: "location1", URL: "https://pokeapi.co/api/v2/location/1/"},
	}
	if !reflect.DeepEqual(locations, expectedLocations) {
		t.Errorf("expected locations %v, got %v", expectedLocations, locations)
	}

	expectedNextUrl := "https://pokeapi.co/api/v2/location/?offset=20&limit=20"
	if cfg.NextUrl != expectedNextUrl {
		t.Errorf("expected NextUrl %q, got %q", expectedNextUrl, cfg.NextUrl)
	}
	expectedPreviousUrl := ""
	if cfg.PreviousUrl != expectedPreviousUrl {
		t.Errorf("expected PreviousUrl %q, got %q", expectedPreviousUrl, cfg.PreviousUrl)
	}
}
func TestGetLocation_WithCachedData(t *testing.T) {
	cachedData := LocationResult{
		Count:    1,
		Next:     "https://pokeapi.co/api/v2/location/?offset=40&limit=20",
		Previous: "https://pokeapi.co/api/v2/location/?offset=0&limit=20",
		Results: []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			{
				Name: "cached_location",
				URL:  "https://pokeapi.co/api/v2/location/cached/",
			},
		},
	}
	cachedDataJSON, _ := json.Marshal(cachedData)
	pokeCache.Set("https://pokeapi.co/api/v2/location/", cachedDataJSON)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("HTTP request was made to the server, but cached data should have been used")
	}))
	defer server.Close()

	cfg := &config.Config{
		NextUrl:     "",
		PreviousUrl: "",
	}

	locationResult := &LocationResult{}

	locations, err := locationResult.GetLocation("https://pokeapi.co/api/v2/location/", cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedLocations := LocationList{
		LocationInfo{Name: "cached_location", URL: "https://pokeapi.co/api/v2/location/cached/"},
	}
	if !reflect.DeepEqual(locations, expectedLocations) {
		t.Errorf("expected locations %v, got %v", expectedLocations, locations)
	}
	expectedNextUrl := "https://pokeapi.co/api/v2/location/?offset=40&limit=20"
	if cfg.NextUrl != expectedNextUrl {
		t.Errorf("expected NextUrl %q, got %q", expectedNextUrl, cfg.NextUrl)
	}
	expectedPreviousUrl := "https://pokeapi.co/api/v2/location/?offset=0&limit=20"
	if cfg.PreviousUrl != expectedPreviousUrl {
		t.Errorf("expected PreviousUrl %q, got %q", expectedPreviousUrl, cfg.PreviousUrl)
	}
}
