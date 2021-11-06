package Models

/**
* Module have methods for internal routing based on Gin framework
 */

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api")
	{
		v1.GET("laps", GetListLaps) // Full list of laps
		v1.GET("laps/:id", GetLap)  // Get one lap details
		v1.GET("laps/bytagid/:id", GetLapsByTagId)
		v1.GET("laps/byraceid/:id", GetLapsByRaceId)
		v1.GET("laps/results/byraceid/:id", GetResultsByRaceId)
		v1.GET("laps/last", GetLastLapData)

	}

	return r
}
