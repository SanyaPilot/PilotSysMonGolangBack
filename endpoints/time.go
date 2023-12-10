package endpoints

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/host"
)

func TimeRoutes(router *gin.Engine) {
	router.GET("/time", func(ctx *gin.Context) {
		t := time.Now()
		zone, offset := t.Zone()
		uptime, _ := host.Uptime()
		fmt.Println(zone, offset)
		ctx.JSON(200, gin.H{
			"time":     t.Unix() + int64(offset),
			"timezone": zone,
			"uptime":   uptime,
		})
	})
}
