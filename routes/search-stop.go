package routes

import (
	"context"
	"log"
	"net/http"
	"os"

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
	googleClient  *maps.Client
	algoliaClient algoliasearch.Client
)

func init() {
	googleClient, _ = maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_PLACES_API_KEY")))
	algoliaClient = algoliasearch.NewClient(os.Getenv("ALGOLIA_APP_ID"), os.Getenv("ALGOLIA_SECRET"))
}

// GetSearchStop returns a list of stops
func GetSearchStop(c echo.Context) error {
	text := c.Param("text")

	// !Algolia
	index := algoliaClient.InitIndex("stops")

	settings := algoliasearch.Map{
		"page":                 0,
		"hitsPerPage":          5,
		"attributesToRetrieve": []string{"stop_name", "stop_lat", "stop_lon"},
	}

	res, err := index.Search(text, settings)

	if err != nil {
		log.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	var hits []Result
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

	// !Google
	request := &maps.PlaceAutocompleteRequest{
		Input:    text,
		Language: "fr",
		Components: map[maps.Component]string{
			maps.ComponentCountry: "fr",
		},
	}
	// TODO can improve results with strictBounds & location

	googleResults, err := googleClient.PlaceAutocomplete(context.Background(), request)
	if err != nil {
		log.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	for _, prediction := range googleResults.Predictions {
		hits = append(hits, Result{
			Type:    "GOOGLE",
			Name:    prediction.Description,
			PlaceID: prediction.PlaceID,
		})
	}

	return c.JSON(http.StatusOK, hits)
}
