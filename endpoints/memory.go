package endpoints

import (
	"math"
	"pilot-sysmon-backend/utils"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/mem"
)

type MemParams struct {
	Human bool `form:"human"`
}

func MemRoutes(router *gin.Engine) {
	router.GET("/memory", func(ctx *gin.Context) {
		var params MemParams
		if err := ctx.ShouldBind(&params); err != nil {
			ctx.JSON(400, gin.H{"status": "error", "msg": err})
			return
		}
		memory, err := mem.VirtualMemory()
		if err != nil {
			utils.AnswerGopsutilError(err, ctx)
			return
		}
		swap, err := mem.SwapMemory()
		if err != nil {
			utils.AnswerGopsutilError(err, ctx)
			return
		}
		memPercent := math.Round((float64(memory.Total-memory.Available) / float64(memory.Total) * 100 * 10)) / 10

		var payload gin.H
		if params.Human {
			payload = gin.H{
				"ram": gin.H{
					"percent":   memPercent,
					"total":     humanize.Bytes(memory.Total),
					"available": humanize.Bytes(memory.Available),
					"used":      humanize.Bytes(memory.Used),
					"free":      humanize.Bytes(memory.Free),
					"active":    humanize.Bytes(memory.Active),
					"inactive":  humanize.Bytes(memory.Inactive),
					"buffers":   humanize.Bytes(memory.Buffers),
					"cached":    humanize.Bytes(memory.Cached),
					"shared":    humanize.Bytes(memory.Shared),
					"slab":      humanize.Bytes(memory.Slab),
				},
				"swap": gin.H{
					"percent": swap.UsedPercent,
					"total":   humanize.Bytes(swap.Total),
					"used":    humanize.Bytes(swap.Used),
					"free":    humanize.Bytes(swap.Free),
					"sin":     humanize.Bytes(swap.Sin),
					"sout":    humanize.Bytes(swap.Sout),
				},
			}
		} else {
			payload = gin.H{
				"ram": gin.H{
					"percent":   memPercent,
					"total":     memory.Total,
					"available": memory.Available,
					"used":      memory.Used,
					"free":      memory.Free,
					"active":    memory.Active,
					"inactive":  memory.Inactive,
					"buffers":   memory.Buffers,
					"cached":    memory.Cached,
					"shared":    memory.Shared,
					"slab":      memory.Slab,
				},
				"swap": gin.H{
					"percent": swap.UsedPercent,
					"total":   swap.Total,
					"used":    swap.Used,
					"free":    swap.Free,
					"sin":     swap.Sin,
					"sout":    swap.Sout,
				},
			}
		}
		ctx.JSON(200, payload)
	})
}
