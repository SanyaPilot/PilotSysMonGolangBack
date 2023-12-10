package main

import (
	"pilot-sysmon-backend/endpoints"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	endpoints.CPURoutes(r)
	endpoints.MemRoutes(r)
	r.Run()
}
