package endpoints

import (
	"runtime"

	"github.com/gin-gonic/gin"
)

var (
	VERSION     string
	COMMIT_HASH string
	BUILD_TIME  string
)

func BackInfoEndpoints(router *gin.Engine) {
	router.GET("/backend_info", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"lang":         "golang",
			"lang_version": runtime.Version(),
			"version":      VERSION,
			"commit_hash":  COMMIT_HASH,
			"build_time":   BUILD_TIME,
		})
	})
}
