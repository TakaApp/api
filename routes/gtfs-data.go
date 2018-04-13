package routes

// Location describes the main places
// such as the origin, destination or a bus stop
type Location struct {
	Name      string `json:"name"`
	StopID    string `json:"name"`
	Longitude int    `json:"lon"`
	Latitude  int    `json:"lat"`

	Departure int `json:"departure"`
	Arrival   int `json:"arrival"`
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
	Duration int `json:"duration"`

	// ex: BUS, WALK
	Mode string `json:"mode"`
	// ex: 4, eq Line 4
	Route string `json:"route"`
	// ex: Foch Cathedrale ~ Direction
	HeadSign string `json:"headSign"`

	From Location `json:"from"`
	To   Location `json:"to"`
}

// Itinerary describes an .. itinerary
type Itinerary struct {
	// timestamps of the start & end time
	StartTime int `json:"startTime`
	EndTime   int `json:"endTime"`

	// duration in seconds
	Duration int `json:"duration"`
	// duration composition, in seconds
	WalkTime    int `json:"walkTime`
	TransitTime int `json:"transitTime`
	WaitingTime int `json:"waitingTime`

	// walk distance in meters
	WalkDistance int `json:walkDistance`
	// number of transfers
	Transfers int `json:transfers`

	// different main steps of the itinerary
	Legs []Leg `json:steps`
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
