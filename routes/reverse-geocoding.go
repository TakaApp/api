package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo"
)

// AlgoliaPlacesReverseGeocodingURL is the url where we can do our reverse
// geocoding
// NOTE API keys ???
const AlgoliaPlacesReverseGeocodingURL = "https://places-dsn.algolia.net/1/places/reverse"

func fetchAlgoliaReverseGeocoding(lat float64, lng float64) ([]Result, error) {
	var algoliaResult AlgoliaPlacesSuggestion
	var hits []Result

	client := &http.Client{}

	url := fmt.Sprintf("%s?aroundLatLng=%f, %f&hitsPerPage=1&language=fr", AlgoliaPlacesReverseGeocodingURL, lat, lng)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Algolia-Application-Id", os.Getenv("ALGOLIA_PLACES_APP_ID"))
	req.Header.Add("X-Algolia-API-Key", os.Getenv("ALGOLIA_PLACES_API_KEY"))

	response, err := client.Do(req)

	if err != nil {
		log.Printf("err :%v\n", err)
		return nil, err
	}

	defer response.Body.Close()

	jsonData, err := ioutil.ReadAll(response.Body)

	json.Unmarshal(jsonData, &algoliaResult)

	for _, suggestion := range algoliaResult.Hits {

		if len(suggestion.LocalNames) == 0 {
			continue
		}

		name := suggestion.LocalNames[0]
		if len(suggestion.City) > 0 {
			name = name + ", " + suggestion.City[0]
		}

		hits = append(hits, Result{
			Type:      "FREE",
			Name:      name,
			Latitude:  suggestion.LatLng.Lat,
			Longitude: suggestion.LatLng.Lng,
		})
	}

	return hits, nil
}

type reverseGeocodingRequest struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// ReverseGeocoding returns the information about a
// place at lat lng position
func ReverseGeocoding(c echo.Context) error {
	// build a reverseGeocodingRequest
	request := new(reverseGeocodingRequest)
	err := c.Bind(request)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return c.String(http.StatusBadRequest, "Sorry (1)")
	}

	var results []Result

	var wg sync.WaitGroup
	var m sync.Mutex

	wg.Add(1)

	go func() {
		algoliaResults, _ := fetchAlgoliaReverseGeocoding(request.Lat, request.Lng)
		m.Lock()
		results = append(results, algoliaResults...)
		m.Unlock()
		wg.Done()
	}()

	// we could find stops there based on the location

	wg.Wait()

	return c.JSON(http.StatusOK, results)
}
