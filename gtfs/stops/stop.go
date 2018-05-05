package gtfs

// Stop a stop from the gtfs data list
type Stop struct {
	Name      string  `json:"name"`
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}
