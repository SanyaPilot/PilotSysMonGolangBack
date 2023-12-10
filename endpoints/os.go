package endpoints

import (
	"bufio"
	"fmt"
	"os"
	"pilot-sysmon-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/host"
)

type FreeDesktopInfo struct {
	Name    string
	Version string
	Url     string
}

func getFreeDesktopInfo(filePath string) (FreeDesktopInfo, error) {
	res := FreeDesktopInfo{}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open FreeDesktop info file. Not running on Linux?")
		return res, err
	}
	s := bufio.NewScanner(file)

	for s.Scan() {
		line := s.Text()
		entry := strings.Split(line, "=")
		switch entry[0] {
		case "NAME":
			res.Name = strings.Trim(entry[1], "\"")
		case "VERSION":
			res.Version = strings.Trim(entry[1], "\"")
		case "HOME_URL":
			res.Url = strings.Trim(entry[1], "\"")
		}
	}
	return res, nil
}

func OSRoutes(router *gin.Engine) {
	router.GET("/os", func(ctx *gin.Context) {
		info, err := host.Info()
		if err != nil {
			utils.AnswerGopsutilError(err, ctx)
			return
		}
		payload := gin.H{
			"family":  strings.ToUpper(info.OS[:1]) + info.OS[1:],
			"version": info.PlatformVersion,
			"release": info.KernelVersion,
		}
		if info.OS == "linux" {
			fdInfo, err := getFreeDesktopInfo("/etc/os-release")
			if err == nil {
				payload["name"] = fdInfo.Name
				payload["version"] = fdInfo.Version
				payload["url"] = fdInfo.Url
			}
		}
		ctx.JSON(200, payload)
	})
}
