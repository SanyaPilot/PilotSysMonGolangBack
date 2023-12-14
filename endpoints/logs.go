package endpoints

import (
	"encoding/json"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type LogsParams struct {
	Id    string `form:"id"`
	Level string `form:"level"`
	Boot  string `form:"boot"`
	Since string `form:"since"`
	Until string `form:"until"`
	Day   string `form:"day"`
}

func LogsRoutes(router *gin.Engine) {
	router.GET("/logs", func(ctx *gin.Context) {
		logLevels := []string{"emerg", "alert", "crit", "error", "warn", "notice", "info", "debug"}

		var params LogsParams
		ctx.ShouldBind(&params) // All params are strings, assume that no bind errors can happen
		args := []string{"-o", "json"}
		if params.Id != "" {
			args = append(args, []string{"-t", params.Id}...)
		}
		if params.Level != "" {
			logLevel := slices.Index(logLevels, params.Level)
			if logLevel == -1 {
				ctx.JSON(400, gin.H{"status": "error", "msg": "Invalid level param!"})
				return
			}
			args = append(args, []string{"-p", strconv.Itoa(logLevel)}...)
		}
		if params.Boot != "" {
			if _, err := strconv.Atoi(params.Boot); err != nil {
				ctx.JSON(400, gin.H{"status": "error", "msg": "boot param must be an int!"})
				return
			}
			args = append(args, []string{"-b", params.Boot}...)
		}
		if params.Day != "" {
			args = append(args, []string{"-S", params.Day, "-U", params.Day + " 23:59:59"}...)
		} else {
			if params.Since != "" {
				args = append(args, []string{"-S", params.Since}...)
			}
			if params.Until != "" {
				args = append(args, []string{"-U", params.Until}...)
			}
		}

		out, err := exec.Command("journalctl", args...).Output()
		if err != nil {
			ctx.JSON(500, gin.H{"status": "error", "msg": "journalctl failed! " + err.Error()})
			return
		}
		if len(out) == 0 {
			ctx.JSON(200, []any{})
			return
		}
		prepOut := "[" + strings.ReplaceAll(string(out), "\n", ",")[:len(out)-1] + "]"
		var rawLog []gin.H
		json.Unmarshal([]byte(prepOut), &rawLog)
		var res []gin.H
		for _, line := range rawLog {
			rawTS, _ := strconv.ParseFloat(line["__REALTIME_TIMESTAMP"].(string), 64)
			rawPrior, _ := strconv.Atoi(line["PRIORITY"].(string))
			var finalId string
			if id, ok := line["SYSLOG_IDENTIFIER"]; ok {
				finalId = id.(string)
			} else {
				if id, ok := line["_COMM"]; ok {
					finalId = id.(string)
				} else {
					finalId = line["CODE_FUNC"].(string)
				}
			}
			msg := line["MESSAGE"]
			res = append(res, gin.H{
				"time":    rawTS / 1000000,
				"level":   logLevels[rawPrior],
				"id":      finalId,
				"message": msg,
			})
		}
		ctx.JSON(200, res)
	})
}
