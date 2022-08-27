package models

type Incident struct {
	Location Location `json:"location`
	Codename string   `json:"codename`
	OfficerID int
}

type Officer struct {
	Id       int      `json:"id`
	Officername string   `json:"codename`
	Location Location `json:"location`
}

type Location struct {
	X int `json:"x"`
	Y int `json:"y"`
}

OfficersMap := make([int]Officer, 10)
OfficersIdleList := make([]Officer, 10)
IncidentsMap := make([int]Incident, 10)
IncidentsWaitingList := make([]Incident, 10)


func (of *Officer) handleOfficerOnline(obj map[string]string)
{
	// check incidentQueues
	
}

func distance(first *Location, last *Location) float
{
	return math.Sqrt((first.X - last.X)*(first.X - last.X) + (first.Y - last.Y)*(first.Y - last.Y))
}

func (incident *Incident) handleIncidentOccurred(incident map[string]string)
{
	// Looking for idle and closest officer
	shortest := -1
	var chosenOfficer Officer;
	for _, officer := range OfficersIdleList {
		d = distance(officer.Location, incident.Location)
		if shortest > 0 && d < shortest {
			shortest = d
			chosenOfficer = officer
		}
	}
	incident.OfficerID = chosenOfficer.ID
	//remove officer out idle list
	OfficersIdleList.remove(chosenOfficer)
	IncidentsMap[incident.ID] = incident
	
}
