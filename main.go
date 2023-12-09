package main

import (
	cpu "pilot-sysmon-backend/endpoints"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	cpu.Routes(r)
	r.Run()
}
