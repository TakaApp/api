package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/labstack/echo"
)

type Stop struct {
	Name      string  `json:"stop_name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

// GetSearchStop returns a list of stops
func GetSearchStop(c echo.Context) error {
	text := c.Param("text")

	client := algoliasearch.NewClient(os.Getenv("ALGOLIA_APP_ID"), os.Getenv("ALGOLIA_SECRET"))
	index := client.InitIndex("stops")

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

	var hits []Stop
	for _, result := range res.Hits {
		fmt.Printf("result: %v\n", result)

		stopName := result["stop_name"].(string)

		stopLatitude := result["stop_lat"].(float64)
		stopLongitude := result["stop_lon"].(float64)

		hits = append(hits, Stop{
			Name:      stopName,
			Latitude:  stopLatitude,
			Longitude: stopLongitude,
		})
	}

	return c.JSON(http.StatusOK, hits)
}
