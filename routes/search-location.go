package routes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/labstack/echo"
)

// Result is one object contained in the array returned by the API
type Result struct {
	// either FREE or STOP
	Type string `json:"type"`
	Name string `json:"name"`

	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type algoliaGeoLoc struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type algoliaPlacesHit struct {
	City []string `json:"city"`

	LatLng     algoliaGeoLoc `json:"_geoloc"`
	LocalNames []string      `json:"locale_names"`
}

// AlgoliaPlacesSuggestion describes the structure returned
// by the Algolia Places API
type AlgoliaPlacesSuggestion struct {
	Hits []algoliaPlacesHit `json:"hits"`
}

var (
	algoliaClient         algoliasearch.Client
	algoliaIndex          algoliasearch.Index
	algoliaSearchSettings algoliasearch.Map
)

// AlgoliaPlacesURL is the url where we can post our query
// TODO we should implement the 2 others fall back urls
const AlgoliaPlacesURL = "https://places-dsn.algolia.net/1/places/query"

func init() {
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
	var algoliaResult AlgoliaPlacesSuggestion
	var hits []Result

	values := map[string]string{
		// search query
		"query": query,
		// we only have a french audience
		"language": "fr",
		// 5 is way enough
		"hitsPerPage": "5",
		// and we limit to
		//  - France
		"countries": "fr",
		//  - around Nantes
		"aroundLatLng": "47.215033,-1.553952",
		//  - in a 15 km radius
		"aroundRadius": "15000",
	}

	q, _ := json.Marshal(values)
	response, err := http.Post(AlgoliaPlacesURL, "application/json", bytes.NewBuffer(q))

	if err != nil {
		log.Printf("err :%v\n", err)
		return nil, err
	}

	defer response.Body.Close()

	jsonData, err := ioutil.ReadAll(response.Body)

	json.Unmarshal(jsonData, &algoliaResult)

	for _, suggestion := range algoliaResult.Hits {
		hits = append(hits, Result{
			Type:      "FREE",
			Name:      suggestion.LocalNames[0],
			Latitude:  suggestion.LatLng.Lat,
			Longitude: suggestion.LatLng.Lng,
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
		existingStops := map[string]bool{}
		for _, result := range algoliaResults {
			if !existingStops[result.Name] {
				results = append(results, result)
				existingStops[result.Name] = true
			}
		}
		m.Unlock()
		wg.Done()
	}()

	wg.Wait()

	return c.JSON(http.StatusOK, results)
}
