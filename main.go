package main

import (
	"fmt"

	"github.com/carmandomx/tza/config"
	"github.com/carmandomx/tza/meetings"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

var DB *pg.DB

func main() {
	r := gin.Default()
	DB, err := config.SetupDb()
	if err != nil {
		panic(err)
	}
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080", "https://editor.swagger.io/"}
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}

	r.Use(cors.Default()) // Corregir en deployment
	defer DB.Close()
	err = meetings.CreateSchema(DB)
	if err != nil {
		fmt.Println(err)
	}
	// meetings.CronJob(DB)
	v1 := r.Group("/v1")
	mtngs := v1.Group("/meetings")
	p := v1.Group("/participants")
	mtngs.GET("/:id", meetings.GetMeeting(DB))
	mtngs.GET("/:id/instances", meetings.FetchAllMeetingInstances(DB))
	p.GET("/:participantId", meetings.FetchParticipant(DB))
	v1.POST("/token", meetings.SetToken)
	r.Static("/swaggerui", "./swaggerui")
	r.Run()
}
