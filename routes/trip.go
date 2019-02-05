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

	response.Body.Close()

	simplifiedResult, _ := json.Marshal(result)

	return c.String(http.StatusOK, string(simplifiedResult))
}
