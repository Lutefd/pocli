package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"poke-repl/internal/config"
	"reflect"
	"testing"
)

func TestLocationResult_GetLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		responseBody := `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
				{
					"name": "Location 1",
					"url": "https://pokeapi.co/api/v2/location/1/"
				},
				{
					"name": "Location 2",
					"url": "https://pokeapi.co/api/v2/location/2/"
				}
			]
		}`

		_, _ = w.Write([]byte(responseBody))
	}))
	defer server.Close()

	mockConfig := &config.Config{
		NextUrl:     "",
		PreviousUrl: "",
	}

	locationResult := &LocationResult{}

	locations, err := locationResult.GetLocation(server.URL, mockConfig)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedLocations := LocationList{
		{Name: "Location 1", URL: "https://pokeapi.co/api/v2/location/1/"},
		{Name: "Location 2", URL: "https://pokeapi.co/api/v2/location/2/"},
	}
	if !reflect.DeepEqual(locations, expectedLocations) {
		t.Errorf("Expected locations %v, but got %v", expectedLocations, locations)
	}

	if mockConfig.NextUrl != "" {
		t.Errorf("Expected NextUrl to be empty, but got %s", mockConfig.NextUrl)
	}
	if mockConfig.PreviousUrl != "" {
		t.Errorf("Expected PreviousUrl to be empty, but got %s", mockConfig.PreviousUrl)
	}

}

func TestLocationResult_GetLocation_WithPreviousUrl(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		responseBody := `{
            "count": 2,
            "next": "https://pokeapi.co/api/v2/location/?offset=20&limit=2",
            "previous": "https://pokeapi.co/api/v2/location/?offset=0&limit=2",
            "results": [
                {
                    "name": "Location 3",
                    "url": "https://pokeapi.co/api/v2/location/3/"
                },
                {
                    "name": "Location 4",
                    "url": "https://pokeapi.co/api/v2/location/4/"
                }
            ]
        }`

		_, _ = w.Write([]byte(responseBody))
	}))
	defer server.Close()

	mockConfig := &config.Config{
		NextUrl:     "",
		PreviousUrl: "https://pokeapi.co/api/v2/location/?offset=0&limit=2",
	}

	locationResult := &LocationResult{}

	locations, err := locationResult.GetLocation(server.URL, mockConfig)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedLocations := LocationList{
		{Name: "Location 3", URL: "https://pokeapi.co/api/v2/location/3/"},
		{Name: "Location 4", URL: "https://pokeapi.co/api/v2/location/4/"},
	}
	if !reflect.DeepEqual(locations, expectedLocations) {
		t.Errorf("Expected locations %v, but got %v", expectedLocations, locations)
	}

	if mockConfig.NextUrl != "https://pokeapi.co/api/v2/location/?offset=20&limit=2" {
		t.Errorf("Expected NextUrl to be updated, but got %s", mockConfig.NextUrl)
	}
}

func TestLocationResult_GetLocation_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	mockConfig := &config.Config{
		NextUrl:     "",
		PreviousUrl: "",
	}

	locationResult := &LocationResult{}

	_, err := locationResult.GetLocation(server.URL, mockConfig)
	if err == nil {
		t.Error("Expected an unauthorized error, got nil")
	}

	if mockConfig.NextUrl != "" {
		t.Errorf("Expected NextUrl to remain empty, but got %s", mockConfig.NextUrl)
	}
	if mockConfig.PreviousUrl != "" {
		t.Errorf("Expected PreviousUrl to remain empty, but got %s", mockConfig.PreviousUrl)
	}
}

func TestLocationResult_GetLocation_DecoderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{invalid JSON"))
	}))
	defer server.Close()

	mockConfig := &config.Config{}

	locationResult := &LocationResult{}

	_, err := locationResult.GetLocation(server.URL, mockConfig)
	if err == nil {
		t.Error("Expected a JSON decoder error, got nil")
	}
}
func TestLocationResult_GetLocation_CachedData(t *testing.T) {
	mockConfig := &config.Config{
		NextUrl:     "",
		PreviousUrl: "",
	}
	locationResult := &LocationResult{}
	cachedLocations := LocationList{
		{Name: "Location 1", URL: "https://pokeapi.co/api/v2/location/1/"},
		{Name: "Location 2", URL: "https://pokeapi.co/api/v2/location/2/"},
	}
	cachedData, _ := json.Marshal(cachedLocations)
	pokeCache.Set("https://pokeapi.co/api/v2/location/", cachedData)
	locations, err := locationResult.GetLocation("https://pokeapi.co/api/v2/location/", mockConfig)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(locations, cachedLocations) {
		t.Errorf("Expected locations %v, but got %v", cachedLocations, locations)
	}
	if mockConfig.NextUrl != "" {
		t.Errorf("Expected NextUrl to be empty, but got %s", mockConfig.NextUrl)
	}
	if mockConfig.PreviousUrl != "" {
		t.Errorf("Expected PreviousUrl to be empty, but got %s", mockConfig.PreviousUrl)
	}
}

func TestLocationResult_GetLocation_ErrorFetching(t *testing.T) {
	mockConfig := &config.Config{
		NextUrl:     "",
		PreviousUrl: "",
	}
	locationResult := &LocationResult{}
	_, err := locationResult.GetLocation("https://pokeapi.co/api/v2/location/", mockConfig)
	if err == nil {
		t.Error("Expected an error, got nil")
	}
	if mockConfig.NextUrl != "" {
		t.Errorf("Expected NextUrl to be empty, but got %s", mockConfig.NextUrl)
	}
	if mockConfig.PreviousUrl != "" {
		t.Errorf("Expected PreviousUrl to be empty, but got %s", mockConfig.PreviousUrl)
	}
}

func TestLocationResult_GetLocation_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{invalid JSON"))
	}))
	defer server.Close()
	mockConfig := &config.Config{
		NextUrl:     "",
		PreviousUrl: "",
	}
	locationResult := &LocationResult{}
	_, err := locationResult.GetLocation(server.URL, mockConfig)
	if err == nil {
		t.Error("Expected a JSON decoder error, got nil")
	}
	if mockConfig.NextUrl != "" {
		t.Errorf("Expected NextUrl to be empty, but got %s", mockConfig.NextUrl)
	}
	if mockConfig.PreviousUrl != "" {
		t.Errorf("Expected PreviousUrl to be empty, but got %s", mockConfig.PreviousUrl)
	}
}
