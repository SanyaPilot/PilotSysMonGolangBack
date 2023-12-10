package endpoints

import (
	"pilot-sysmon-backend/utils"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/disk"
)

type DiskParams struct {
	Human bool `form:"human"`
}

func DisksRoutes(router *gin.Engine) {
	router.GET("/disks", func(ctx *gin.Context) {
		var params MemParams
		if err := ctx.ShouldBind(&params); err != nil {
			ctx.JSON(400, gin.H{"status": "error", "msg": err})
			return
		}
		parts, err := disk.Partitions(false)
		if err != nil {
			utils.AnswerGopsutilError(err, ctx)
			return
		}
		var payload []gin.H
		for _, part := range parts {
			usage, err := disk.Usage(part.Mountpoint)
			if err != nil {
				utils.AnswerGopsutilError(err, ctx)
				return
			}
			var usagePayload gin.H
			if params.Human {
				usagePayload = gin.H{
					"total":   humanize.Bytes(usage.Total),
					"used":    humanize.Bytes(usage.Used),
					"free":    humanize.Bytes(usage.Free),
					"percent": usage.UsedPercent,
				}
			} else {
				usagePayload = gin.H{
					"total":   usage.Total,
					"used":    usage.Used,
					"free":    usage.Free,
					"percent": usage.UsedPercent,
				}
			}
			payload = append(payload, gin.H{
				"device":     part.Device,
				"mountpoint": part.Mountpoint,
				"fs":         part.Fstype,
				"opts":       strings.Join(part.Opts, ","),
				"usage":      usagePayload,
			})
		}
		ctx.JSON(200, payload)
	})
}
