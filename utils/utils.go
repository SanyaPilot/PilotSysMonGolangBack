package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AnswerError(reason string, err error, ctx *gin.Context) {
	var msg string
	if err != nil {
		msg = reason + err.Error()
	} else {
		msg = reason
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"status": "error", "msg": msg,
	})
}

func AnswerGopsutilError(err error, ctx *gin.Context) {
	AnswerError("gopsutil error: ", err, ctx)
}
