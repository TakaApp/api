package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type getTripRequest struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Time     string `json:"time"`
	Date     string `json:"date"`
	ArriveBy string `json:"arriveBy"`
}

// GetTrip returns possible itineraries between two points
func GetTrip(c echo.Context) error {
	// build a GetTripRequest
	request := new(getTripRequest)
	err := c.Bind(request)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return c.String(http.StatusBadRequest, "Sorry (1)")
	}

	// create a request to our OTP server
	req, err := http.NewRequest("GET", "http://gtfs.aksels.io/otp/routers/default/plan", nil)
	if err != nil {
		log.Print(err)
		return c.String(http.StatusInternalServerError, "Sorry (2)")
	}

	// add OTP required params
	q := req.URL.Query()
	q.Add("fromPlace", request.From)
	q.Add("toPlace", request.To)
	q.Add("time", request.Time)
	q.Add("date", request.Date)
	q.Add("mode", "TRANSIT,WALK")
	q.Add("maxWalkDistance", "500000")
	q.Add("arriveBy", request.ArriveBy)
	q.Add("locale", "fr")

	req.URL.RawQuery = q.Encode()

	// request our server
	response, err := http.Get(req.URL.String())
	if err != nil {
		log.Print(err)
		return c.String(http.StatusInternalServerError, "Sorry (3)")
	}

	rawResponseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return c.String(http.StatusInternalServerError, "Sorry (4)")
	}

	var result GTFSResult
	if err := json.Unmarshal(rawResponseBody, &result); err != nil {
		log.Print(err)
	}

	for iIDX, itinerary := range result.Plan.Itineraries {
		fmt.Println("----------------")
		for lIDX, leg := range itinerary.Legs {
			fromStopID := leg.From.StopID
			toStopID := leg.To.StopID

			routeID := leg.RouteID
			if fromStopID == "" || toStopID == "" || leg.RouteID == "" {
				continue
			}

			// fetch patterns
			patternReq, _ := http.NewRequest("GET", fmt.Sprintf("http://gtfs.aksels.io/otp/routers/default/index/routes/%s/patterns", routeID), nil)
			patternReqResult, _ := http.Get(patternReq.URL.String())

			var patterns []Pattern
			patternReqResponseBody, _ := ioutil.ReadAll(patternReqResult.Body)
			if err := json.Unmarshal(patternReqResponseBody, &patterns); err != nil {
				log.Print(err)
			}

			fmt.Printf("(!) %s >> %s\n", fromStopID, toStopID)
			// fmt.Printf("patterns %v\n", patterns)

			var stopDetails []Stop
			for _, pattern := range patterns {
				stopDetails = nil

				// fetch data about a single pattern
				singlePatternReq, _ := http.NewRequest("GET", fmt.Sprintf("http://gtfs.aksels.io/otp/routers/default/index/patterns/%s", pattern.ID), nil)
				singlePatternReqResult, _ := http.Get(singlePatternReq.URL.String())

				var p Pattern
				singlePatternReqResponseBody, _ := ioutil.ReadAll(singlePatternReqResult.Body)
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
					if stopDetails[0].ID == leg.To.StopID {
						// reverse method someone ?
						for i := len(stopDetails)/2 - 1; i >= 0; i-- {
							opp := len(stopDetails) - 1 - i
							stopDetails[i], stopDetails[opp] = stopDetails[opp], stopDetails[i]
						}
					}
					break
				}
			}

			result.Plan.Itineraries[iIDX].Legs[lIDX].Stops = stopDetails
		}
	}

	response.Body.Close()

	simplifiedResult, _ := json.Marshal(result)

	return c.String(http.StatusOK, string(simplifiedResult))
}
