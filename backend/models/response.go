package models

type Message struct {
	Type       string    `json:"type,omitempty`
	IncidentId *int      `json:"incidentId,omitempty`
	Codename   *string   `json:"codeName`
	OfficerId  *int      `json:"officerId,omitempty`
	BadgeName  *string   `json:"badgeName`
	Location   *Location `json:"location`
}

type Response struct {
	Data Data        `json:"data`
	Err  interface{} `json:"error`
}

type Data struct {
	Incidents []Incident `json:"incidents`
	Officers  []Officer  `json:"officers`
}
type Incident struct {
	Id       int      `json:"id`
	Codename string   `json:"codename`
	Location Location `json:"location`
}

type Officer struct {
	Id        int      `json:"id`
	BadgeName string   `json:"badgeName`
	Location  Location `json:"location`
}

type Location struct {
	X int `json:"x"`
	Y int `json:"y"`
}
