package handlers

import (
	"backend-axon-challenge-2022/models"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/streadway/amqp"
)

var MapIncidentById models.ObjectMapByID[int, models.Incident] = models.NewObjectMapByID[int, models.Incident]()
var MapOfficerById models.ObjectMapByID[int, models.Officer] = models.NewObjectMapByID[int, models.Officer]()
var taskQueue models.TaskQueue = make(models.TaskQueue, 10000)
var messageQueue chan models.Message = make(chan models.Message, 10000)

func ReadEvents() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"events", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func(taskQueue models.TaskQueue) {
		for d := range msgs {
			var m *models.Message
			if err := json.Unmarshal(d.Body, &m); err != nil {
				log.Println(err)
				continue
			}

			// b, _ := json.Marshal(m)
			// log.Printf("Received a message: %v\n", string(b))
			messageQueue <- *m
		}
	}(taskQueue)
	go consumeMessage()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

}

func consumeMessage() {
	for m := range messageQueue {
		// time.Sleep(1000 * time.Millisecond)

		// log.Printf("MapOfficerById %v\n", MapOfficerById)
		// log.Printf("MapIncidentById %v\n", MapIncidentById)
		b, _ := json.Marshal(m)
		log.Printf("Received a message: %s\n", b)

		if m.Type == "IncidentOccurred" {
			handleIncidentOccurred(models.Incident{
				Id:       *m.IncidentId,
				CodeName: *m.Codename,
				Loc:      *m.Loc,
			})
		} else if m.Type == "IncidentResolved" {
			handleIncidentResolved(models.Incident{
				Id: *m.IncidentId,
			})
		} else if m.Type == "OfficerGoesOnline" {
			handleOfficerGoesOnline(models.Officer{
				Id:        *m.OfficerId,
				BadgeName: *m.BadgeName,
			})
		} else if m.Type == "OfficerGoesOffline" {
			handleOfficerGoesOffline(*m.OfficerId)
		} else if m.Type == "OfficerLocationUpdated" {
			handleOfficerLocationUpdated(models.Officer{
				Id:  *m.OfficerId,
				Loc: *m.Loc,
			})
		}
	}
}

func handleIncidentOccurred(incident models.Incident) {
	if chosenOfficer := findNearestAvailableOfficer(incident); chosenOfficer != nil {
		fmt.Println(chosenOfficer)
		incident.OfficerId = chosenOfficer.Id

		chosenOfficer.IncidentId = incident.Id
		MapOfficerById.Set(chosenOfficer.Id, *chosenOfficer)
	}
	MapIncidentById.Set(incident.Id, incident)
}

func handleIncidentResolved(in models.Incident) {
	incident := MapIncidentById.Get(in.Id)
	if incident.Id != 0 && incident.OfficerId != 0 {
		officer := MapOfficerById.Get(incident.OfficerId)
		if officer.Id == 0 {
			return
		}

		if nextIncident := findNearestUnassignedIncident(officer); nextIncident != nil {
			nextIncident.OfficerId = officer.Id
			MapIncidentById.Set(nextIncident.Id, *nextIncident)

			officer.IncidentId = nextIncident.Id
			MapOfficerById.Set(officer.Id, officer)
		} else {
			MapOfficerById.Set(incident.OfficerId, officer)
		}
	}
	MapIncidentById.Delete(incident.Id)

}

func handleOfficerGoesOnline(officer models.Officer) {
	MapOfficerById.Set(officer.Id, officer)
}

func handleOfficerGoesOffline(officerId int) {
	officer := MapOfficerById.Get(officerId)
	if officer.Id == 0 {
		return
	}
	incidentId := officer.IncidentId
	MapOfficerById.Delete(officerId)

	if incidentId != 0 {
		assignedIncident := MapIncidentById.Get(incidentId)
		if assignedIncident.Id == 0 {
			return
		}

		if nextOfficer := findNearestAvailableOfficer(assignedIncident); nextOfficer != nil {
			nextOfficer.IncidentId = assignedIncident.Id
			MapOfficerById.Set(nextOfficer.Id, *nextOfficer)

			assignedIncident.OfficerId = nextOfficer.Id
		} else {
			assignedIncident.OfficerId = 0
		}
		MapIncidentById.Set(incidentId, assignedIncident)

	}
}

func handleOfficerLocationUpdated(officer models.Officer) {
	foundOfficer := MapOfficerById.Get(officer.Id)
	if foundOfficer.Id == 0 {
		return
	}
	foundOfficer.Loc = officer.Loc
	MapOfficerById.Set(officer.Id, foundOfficer)
	if foundOfficer.IncidentId == 0 {
		if incident := findNearestUnassignedIncident(foundOfficer); incident != nil {
			foundOfficer.IncidentId = incident.Id
			MapOfficerById.Set(foundOfficer.Id, foundOfficer)

			incident.OfficerId = foundOfficer.Id
			MapIncidentById.Set(incident.Id, *incident)
		}
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func distance(first *models.Location, last *models.Location) float64 {
	return math.Sqrt(float64(first.X-last.X)*float64(first.X-last.X) + float64(first.Y-last.Y)*float64(first.Y-last.Y))
}

func findNearestAvailableOfficer(incident models.Incident) *models.Officer {
	shortest := 10000000.0
	var chosenOfficer *models.Officer
	for _, officer := range MapOfficerById.Copy() {
		if officer.Id == 0 || officer.IncidentId != 0 || (officer.Loc.X == 0 && officer.Loc.Y == 0) {
			continue
		}
		d := distance(&officer.Loc, &incident.Loc)
		if d < shortest {
			shortest = d
			chosenOfficer = &officer
		}
	}
	return chosenOfficer
}

func findNearestUnassignedIncident(officer models.Officer) *models.Incident {
	shortest := 10000000.0
	var pickedTask *models.Incident
	for _, incident := range MapIncidentById.Copy() {
		if incident.Id == 0 || incident.OfficerId != 0 {
			continue
		}
		d := distance(&officer.Loc, &incident.Loc)
		if d < shortest {
			shortest = d
			pickedTask = &incident
		}
	}
	return pickedTask
}
