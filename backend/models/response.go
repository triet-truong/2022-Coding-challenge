package models

type Message struct {
	Type       string    `json:"type,omitempty"`
	IncidentId *int      `json:"incidentId,omitempty"`
	Codename   *string   `json:"codeName"`
	OfficerId  *int      `json:"officerId,omitempty"`
	BadgeName  *string   `json:"badgeName"`
	Loc        *Location `json:"loc"`
}

type Response struct {
	Data Data        `json:"data"`
	Err  interface{} `json:"error"`
}

type Data struct {
	Incidents []Incident `json:"incidents"`
	Officers  []Officer  `json:"officers"`
}
type Incident struct {
	Id        int      `json:"id"`
	OfficerId int      `json:"officerId"`
	CodeName  string   `json:"codename"`
	Loc       Location `json:"loc"`
}

type Officer struct {
	Id         int      `json:"id"`
	BadgeName  string   `json:"badgeName"`
	IncidentId int      `json:"-"`
	Loc        Location `json:"loc"`
}

type Location struct {
	X int `json:"x"`
	Y int `json:"y"`
}
