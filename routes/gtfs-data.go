package routes

// A Place is where a journey starts or ends
// or a transit stop along the way.
type Place struct {
	// For transit stops, the name of the stop.
	// For points of interest, the name of the POI
	Name      string  `json:"name"`
	StopID    string  `json:"stopID"`
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`

	// The time the rider will depart the place
	Departure int `json:"departure"`
	// The time the rider will arrive at the place
	Arrival int `json:"arrival"`

	// For transit trips, the stop index
	// (numbered from zero from the start of the trip
	StopIndex int `json:"stopIndex"`

	// The "code" of the stop. Depending on the transit
	// agency, this is often something that users care about
	StopCode string `json:"stopCode"`

	// For transit trips, the sequence number of the stop.
	// Per GTFS, these numbers are increasing
	StopSequence int `json:"stopSequence"`
}

// Leg is a struct defining a main step of
// the whole trip. Example:
// - Walk till the bus      (1st leg)
// - Take the bus           (2nd leg)
// - Take a second bus      (3rd leg)
type Leg struct {
	// timestamps of the start & end time
	StartTime int `json:"startTime"`
	EndTime   int `json:"endTime"`
	// duration in seconds
	Duration float64 `json:"duration"`
	// distance in meters
	Distance float64 `json:"distance"`

	// For transit legs, the type of the route.
	//  - Non transit -1
	//  - When 0-7:
	//    - 0 Tram
	//    - 1 Subway
	//    - 2 Train
	//    - 3 Bus
	//    - 4 Ferry
	//    - 5 Cable Car
	//    - 6 Gondola
	//    - 7 Funicular
	//  - When equal or highter than 100:
	//    it is coded using the Hierarchical Vehicle Type (HVT)
	//    codes from the European TPEG standard
	RouteType int `json:"routeType"`

	// For transit leg:
	//  - the route's (background) color (if one exists)
	// For non-transit legs
	//  - null.
	RouteColor string `json:"routeColor"`

	// For transit leg:
	//  - the route's text color (if one exists)
	// For non-transit legs
	//  - null.
	RouteTextColor string `json:"routeTextColor"`

	// The mode used when traversing this leg.
	// ex: BUS, WALK
	Mode string `json:"mode"`

	// For transit legs:
	//  - the route of the bus or train being used
	// For non-transit legs
	//  - the name of the street being traversed.
	// ex: 4, eq Line 4
	Route string `json:"route"`

	// For transit legs:
	//  - the headsign of the bus or train being used
	// For non-transit legs: null.
	//
	// ex: Foch Cathedrale ~ Direction
	HeadSign string `json:"headSign"`

	From Place `json:"from"`
	To   Place `json:"to"`
}

// Itinerary describes an .. itinerary
type Itinerary struct {
	// timestamps of the start & end time
	StartTime int `json:"startTime"`
	EndTime   int `json:"endTime"`

	// duration in seconds
	Duration int `json:"duration"`
	// duration composition, in seconds
	WalkTime    int `json:"walkTime"`
	TransitTime int `json:"transitTime"`
	WaitingTime int `json:"waitingTime"`

	// walk distance in meters
	WalkDistance float64 `json:"walkDistance"`
	// number of transfers
	Transfers int `json:"transfers"`

	// different main steps of the itinerary
	Legs []Leg `json:"legs"`
}

// GTFSPlan consists of multiples itineraries
type GTFSPlan struct {
	Itineraries []Itinerary `json:"itineraries"`
}

// GTFSResult has the plan which interest us
// but also additional parameters such as:
//  - requestParameters
//  - debugOutput
//  - elevationMetadata
// but we ignore them for the moment
type GTFSResult struct {
	Plan GTFSPlan `json:"plan"`
}
