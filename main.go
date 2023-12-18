package main

import (
	"pilot-sysmon-backend/endpoints"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	endpoints.CPURoutes(r)
	endpoints.MemRoutes(r)
	endpoints.DisksRoutes(r)
	endpoints.OSRoutes(r)
	endpoints.TimeRoutes(r)
	endpoints.LogsRoutes(r)
	endpoints.BackInfoEndpoints(r)
	r.Run()
}
