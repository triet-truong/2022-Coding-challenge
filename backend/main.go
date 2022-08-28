package main

import (
	"backend-axon-challenge-2022/handlers"
	"backend-axon-challenge-2022/models"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var mapIncidentById models.ObjectMapByID[int, models.Incident]
var mapOfficerById models.ObjectMapByID[int, models.Officer]

func main() {
	go handlers.ReadEvents()
	StartServer()
}

func StartServer() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/api/v1/state", ServeState)
	r.Run(":8080")
}

func ServeState(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"data": &models.Data{
			Incidents: handlers.MapIncidentById.Values(),
			Officers:  handlers.MapOfficerById.Values(),
		},
	})
}
