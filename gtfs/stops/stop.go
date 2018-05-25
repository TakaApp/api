package gtfs

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

// Stop a stop from the gtfs data list
type Stop struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

var (
	stops []Stop
)

func init() {
	file, err := os.Open("./gtfs/stops/stops.txt")
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(file))

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		latitude, _ := strconv.ParseFloat(line[3], 64)
		longitude, _ := strconv.ParseFloat(line[4], 64)

		//
		stops = append(stops, Stop{
			Name:      line[1],
			Latitude:  latitude,
			Longitude: longitude,
		})
	}
}

func SearchStop(text string) []string {
	var result []string

	return result
}
