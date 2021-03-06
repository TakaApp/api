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

// LegGeometry : A leg's geometry
type LegGeometry struct {
	// A list of coordinates encoded as a string
	Points string `json:"points"`
	Length int    `json:"length"`
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

	// RouteID self explanatory
	RouteID string `json:"routeID"`

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

	// The leg's geometry
	LegGeometry `json:"legGeometry"`

	Stops []Stop `json:"stops"`
}

// An Itinerary is one complete way of getting from
// the start location to the end location.
type Itinerary struct {
	// Time that the trip departs
	StartTime int `json:"startTime"`
	// Time that the trip arrives
	EndTime int `json:"endTime"`

	// Duration of the trip on this itinerary, in seconds
	Duration int `json:"duration"`
	// How much time is spent walking, in seconds
	WalkTime int `json:"walkTime"`
	// How much time is spent on transit, in seconds
	TransitTime int `json:"transitTime"`
	// How much time is spent waiting for transit to arrive, in seconds
	WaitingTime int `json:"waitingTime"`

	// How far the user has to walk, in meters
	WalkDistance float64 `json:"walkDistance"`
	// The number of transfers this trip has.
	Transfers int `json:"transfers"`

	// A list of Legs. Each Leg is either a walking (cycling, car)
	// portion of the trip, or a transit trip on a particular vehicle.
	// So a trip where the use walks to the Q train, transfers to the 6,
	// then walks to their destination, has four legs.
	Legs []Leg `json:"legs"`
}

// GTFSPlan consists of multiples itineraries
type GTFSPlan struct {
	// The time and date of travel
	Date int `json:"date"`

	// The origin
	From Place `json:"from"`

	// The destination
	To Place `json:"to"`

	// A list of possible itineraries
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

type Stop struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type Pattern struct {
	ID   string `json:"id"`
	Desc string `json:"desc"`

	// Stops are ordered !
	Stops []Stop `json:"stops"`
}
