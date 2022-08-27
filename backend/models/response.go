package models

type Incident struct {
}

type Officer struct {
	Id       int      `json:"id`
	Codename string   `json:"codename`
	Location Location `json:"location`
}

type Location struct {
	X int `json:"x"`
	Y int `json:"y"`
}
