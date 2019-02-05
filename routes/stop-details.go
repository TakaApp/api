package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// GetStopDetails returns stops between two stops
func GetStopDetails(c echo.Context) error {
	routeID := c.Param("routeID")
	fromStopID := c.Param("from")
	toStopID := c.Param("to")

	// fetch patterns
	patternReq, _ := http.NewRequest("GET", fmt.Sprintf("http://gtfs.aksels.io/otp/routers/default/index/routes/%s/patterns", routeID), nil)
	patternReqResult, _ := http.Get(patternReq.URL.String())

	var patterns []Pattern
	patternReqResponseBody, _ := ioutil.ReadAll(patternReqResult.Body)

	defer patternReqResult.Body.Close()

	if err := json.Unmarshal(patternReqResponseBody, &patterns); err != nil {
		log.Print(err)
	}

	var stopDetails []Stop
	for _, pattern := range patterns {
		stopDetails = nil

		// fetch data about a single pattern
		singlePatternReq, _ := http.NewRequest("GET", fmt.Sprintf("http://gtfs.aksels.io/otp/routers/default/index/patterns/%s", pattern.ID), nil)
		singlePatternReqResult, _ := http.Get(singlePatternReq.URL.String())

		var p Pattern
		singlePatternReqResponseBody, _ := ioutil.ReadAll(singlePatternReqResult.Body)

		defer singlePatternReq.Body.Close()

		if err := json.Unmarshal(singlePatternReqResponseBody, &p); err != nil {
			log.Print(err)
		}

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

		if haveFoundFrom && haveFoundTo {
			// maybe need to reverse?
			if stopDetails[0].ID == toStopID {
				// reverse method someone ?
				for i := len(stopDetails)/2 - 1; i >= 0; i-- {
					opp := len(stopDetails) - 1 - i
					stopDetails[i], stopDetails[opp] = stopDetails[opp], stopDetails[i]
				}
			}
			break
		}
	}

	result, _ := json.Marshal(stopDetails)

	return c.String(http.StatusOK, string(result))
}
