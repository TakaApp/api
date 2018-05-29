package routes

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/labstack/echo"
	"googlemaps.github.io/maps"
)

type Result struct {
	Type      string  `json:"type"`
	Name      string  `json:"stop_name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	PlaceID   string  `json:"place_id"`
}

var (
	googleClient *maps.Client

	algoliaClient         algoliasearch.Client
	algoliaIndex          algoliasearch.Index
	algoliaSearchSettings algoliasearch.Map
)

func init() {
	googleClient, _ = maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_PLACES_API_KEY")))
	algoliaClient = algoliasearch.NewClient(os.Getenv("ALGOLIA_APP_ID"), os.Getenv("ALGOLIA_SECRET"))
	algoliaIndex = algoliaClient.InitIndex("stops")
	algoliaSearchSettings = algoliasearch.Map{
		"page":                 0,
		"hitsPerPage":          5,
		"attributesToRetrieve": []string{"stop_name", "stop_lat", "stop_lon"},
	}
}

func fetchAlgoliaResults(query string) ([]Result, error) {
	var hits []Result

	res, err := algoliaIndex.Search(query, algoliaSearchSettings)

	if err != nil {
		log.Printf("err: %v\n", err)
		return nil, err
	}

	for _, result := range res.Hits {
		stopName := result["stop_name"].(string)

		stopLatitude := result["stop_lat"].(float64)
		stopLongitude := result["stop_lon"].(float64)

		hits = append(hits, Result{
			Type:      "STOP",
			Name:      stopName,
			Latitude:  stopLatitude,
			Longitude: stopLongitude,
		})
	}

	return hits, nil
}

func fetchGoogleResults(query string) ([]Result, error) {
	var hits []Result

	request := &maps.PlaceAutocompleteRequest{
		Input:    query,
		Language: "fr",
		Components: map[maps.Component]string{
			maps.ComponentCountry: "fr",
		},
	}
	// TODO can improve results with strictBounds & location

	googleResults, err := googleClient.PlaceAutocomplete(context.Background(), request)
	if err != nil {
		log.Printf("err: %v\n", err)
		return nil, err
	}

	for _, prediction := range googleResults.Predictions {
		hits = append(hits, Result{
			Type:    "GOOGLE",
			Name:    prediction.Description,
			PlaceID: prediction.PlaceID,
		})
	}

	return hits, nil
}

// GetSearchLocation returns a list of Results from algolia & google
func GetSearchLocation(c echo.Context) error {
	var results []Result

	text := c.Param("text")

	var wg sync.WaitGroup
	var m sync.Mutex

	wg.Add(2)

	go func() {
		googleResults, _ := fetchGoogleResults(text)
		m.Lock()
		results = append(results, googleResults...)
		m.Unlock()
		wg.Done()
	}()

	go func() {
		algoliaResults, _ := fetchAlgoliaResults(text)
		m.Lock()
		results = append(results, algoliaResults...)
		m.Unlock()
		wg.Done()
	}()

	wg.Wait()

	return c.JSON(http.StatusOK, results)
}
