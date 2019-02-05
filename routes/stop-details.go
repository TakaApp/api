package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func fetchPatterns(routeID string) ([]Pattern, error) {
	var patterns []Pattern

	// build a request
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://gtfs.aksels.io/otp/routers/default/index/routes/%s/patterns", routeID), nil)

	// query our gtfs api
	result, err := http.Get(req.URL.String())
	if err != nil {
		return nil, err
	}

	// read and close at some point
	defer result.Body.Close()
	body, _ := ioutil.ReadAll(result.Body)

	// fit this in our structures
	err = json.Unmarshal(body, &patterns)
	if err != nil {
		return nil, err
	}

	return patterns, nil
}

func fetchPattern(patternID string) (Pattern, error) {
	var p Pattern

	// build a request
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://gtfs.aksels.io/otp/routers/default/index/patterns/%s", patternID), nil)

	// query our gtfs api
	result, err := http.Get(req.URL.String())
	if err != nil {
		return p, err
	}

	// read and close at some point
	defer result.Body.Close()
	singlePatternReqResponseBody, _ := ioutil.ReadAll(result.Body)

	// fit this in our structures
	err = json.Unmarshal(singlePatternReqResponseBody, &p)
	if err != nil {
		log.Print(err)
	}

	return p, nil
}

// GetStopDetails returns stops between two stops
func GetStopDetails(c echo.Context) error {
	// retrieve our parameters
	routeID := c.Param("routeID")
	fromStopID := c.Param("from")
	toStopID := c.Param("to")

	// fetch patterns
	patterns, _ := fetchPatterns(routeID)

	for _, pattern := range patterns {
		var stopDetails []Stop

		// fetch all data from this pattern (note: quite a heavy response)
		p, _ := fetchPattern(pattern.ID)

		haveFoundFrom := false
		haveFoundTo := false

		for _, stop := range p.Stops {
			if stop.ID == fromStopID {
				haveFoundFrom = true
			}
			if stop.ID == toStopID {
				haveFoundTo = true
			}

			if haveFoundFrom || haveFoundTo {
				stopDetails = append(stopDetails, stop)
			}

			if haveFoundFrom && haveFoundTo {
				break
			}
		}

		// we have found a way
		if haveFoundFrom && haveFoundTo {
			// maybe need to reverse?
			if stopDetails[0].ID == toStopID {
				// reverse method someone ?
				for i := len(stopDetails)/2 - 1; i >= 0; i-- {
					opp := len(stopDetails) - 1 - i
					stopDetails[i], stopDetails[opp] = stopDetails[opp], stopDetails[i]
				}
			}

			result, _ := json.Marshal(stopDetails)

			return c.String(http.StatusOK, string(result))
		}
	}

	return c.String(http.StatusNotFound, "404")
}
