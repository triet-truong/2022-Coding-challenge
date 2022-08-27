package main

import (
	"backend-axon-challenge-2022/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func main() {
	ReadEvents()
	StartServer()
}

func StartServer() {
	r := gin.Default()
	r.GET("/api/v1/state", ServeState)
	r.Run(":8080")
}

func ServeState(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": &models.Data{
			Incidents: []models.Incident{},
			Officers:  []models.Officer{},
		},
	})
}

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
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			var m models.Message
			json.Unmarshal(d.Body, &m)
			b, _ := json.Marshal(m)
			log.Printf("Received a message: %s\n", b)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

}
