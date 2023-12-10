package main

import (
	"pilot-sysmon-backend/endpoints"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	endpoints.CPURoutes(r)
	endpoints.MemRoutes(r)
	endpoints.DisksRoutes(r)
	endpoints.OSRoutes(r)
	endpoints.TimeRoutes(r)
	r.Run()
}
