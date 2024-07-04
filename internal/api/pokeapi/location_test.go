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
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response status code
		w.WriteHeader(http.StatusOK)

		// Create a mock response body
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

		// Marshal the response body to JSON
		responseJSON, _ := json.Marshal(response)

		// Write the response body
		w.Write(responseJSON)
	}))
	defer server.Close()

	// Create a mock config
	cfg := &config.Config{
		NextUrl:     "",
		PreviousUrl: "",
	}

	// Create a new instance of LocationResult
	locationResult := &LocationResult{}

	// Call the GetLocation method
	locations, err := locationResult.GetLocation(server.URL, cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check the returned locations
	expectedLocations := LocationList{
		LocationInfo{Name: "location1", URL: "https://pokeapi.co/api/v2/location/1/"},
	}
	if !reflect.DeepEqual(locations, expectedLocations) {
		t.Errorf("expected locations %v, got %v", expectedLocations, locations)
	}

	// Check the config values
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
	// Create a mock cache with pre-populated data
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

	// Create a mock HTTP server that should not be called
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("HTTP request was made to the server, but cached data should have been used")
	}))
	defer server.Close()

	// Create a mock config
	cfg := &config.Config{
		NextUrl:     "",
		PreviousUrl: "",
	}

	// Create a new instance of LocationResult
	locationResult := &LocationResult{}

	// Call the GetLocation method with the URL that has cached data
	locations, err := locationResult.GetLocation("https://pokeapi.co/api/v2/location/", cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check the returned locations against the expected cached data
	expectedLocations := LocationList{
		LocationInfo{Name: "cached_location", URL: "https://pokeapi.co/api/v2/location/cached/"},
	}
	if !reflect.DeepEqual(locations, expectedLocations) {
		t.Errorf("expected locations %v, got %v", expectedLocations, locations)
	}

	// Check the config values against the expected cached data
	expectedNextUrl := "https://pokeapi.co/api/v2/location/?offset=40&limit=20"
	if cfg.NextUrl != expectedNextUrl {
		t.Errorf("expected NextUrl %q, got %q", expectedNextUrl, cfg.NextUrl)
	}
	expectedPreviousUrl := "https://pokeapi.co/api/v2/location/?offset=0&limit=20"
	if cfg.PreviousUrl != expectedPreviousUrl {
		t.Errorf("expected PreviousUrl %q, got %q", expectedPreviousUrl, cfg.PreviousUrl)
	}
}
