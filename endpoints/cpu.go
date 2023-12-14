package endpoints

import (
	"bufio"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"pilot-sysmon-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
)

type CPUFreq struct {
	Min     float32
	Max     float32
	Current float32
}

func getCPUFreqPaths(node_prefix ...string) []string {
	prefix := "/sys/devices/system/cpu/"
	if len(node_prefix) > 0 {
		prefix = node_prefix[0]
	}
	var paths []string
	_, err := os.Stat(prefix + "cpufreq/policy0")
	if err == nil {
		paths, _ = filepath.Glob(prefix + "cpufreq/policy[0-9]*")
	} else {
		// Try another path
		_, err := os.Stat(prefix + "cpu0/cpufreq")
		if err != nil {
			fmt.Println("CPUFreq paths are not found. Not running on Linux?")
			return nil
		}
		paths, _ = filepath.Glob(prefix + "cpu[0-9]*/cpufreq")
	}
	slices.SortFunc(paths, func(a string, b string) int {
		pattern := regexp.MustCompile("[0-9]+")
		aNum, _ := strconv.Atoi(pattern.FindString(a))
		bNum, _ := strconv.Atoi(pattern.FindString(b))
		return aNum - bNum
	})
	return paths
}

func readCPUFreqNode(path string, node string) float32 {
	data, err := os.ReadFile(path + "/" + node)
	if err != nil {
		fmt.Println("Opening cpufreq scaling info failed!")
		fmt.Println(err.Error())
		return 0
	}
	resInt, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	res := float32(resInt) / 1000
	return res
}

func getCPUFreqs(paths []string) []CPUFreq {
	if paths == nil {
		return nil
	}
	var freqs []CPUFreq
	for _, path := range paths {
		freqs = append(freqs, CPUFreq{
			readCPUFreqNode(path, "scaling_min_freq"),
			readCPUFreqNode(path, "scaling_max_freq"),
			readCPUFreqNode(path, "scaling_cur_freq"),
		})
	}
	return freqs
}

func getCPUCoreCount() (int, int) {
	// Core count
	logicalCount, _ := cpu.Counts(true)
	physicalCount, _ := cpu.Counts(false)
	return logicalCount, physicalCount
}

func readCPUInfo(cpuinfoFile string) map[string]any {
	data, err := os.ReadFile(cpuinfoFile)
	if err != nil {
		fmt.Println("Error while reading /proc/cpuinfo. Not running on Linux?")
		return nil
	}
	res := make(map[string]any)
	allowed_entries := []string{"vendor_id", "cpu_family", "model", "stepping", "microcode", "flags", "bugs", "features", "cache_size"}
	s := bufio.NewScanner(strings.NewReader(string(data)))
	for s.Scan() {
		entry := strings.Split(s.Text(), ": ")
		key := strings.ToLower(strings.TrimSpace(entry[0]))
		key = strings.ReplaceAll(key, " ", "_")
		if slices.Contains(allowed_entries, key) {
			res[key] = entry[1]
		}
	}
	return res
}

func CPURoutes(router *gin.Engine) {
	router.GET("/cpu", func(ctx *gin.Context) {
		info, _ := cpu.Info()

		// Core count
		logicalCount, physicalCount := getCPUCoreCount()

		// Freq
		freqs := getCPUFreqs(getCPUFreqPaths())

		// Prepare current frequencies
		// Count average
		// Select maximum and minimum
		var curFreqs []float32
		maxFreq := freqs[0].Max
		minFreq := freqs[0].Min
		freqSum := float32(0)
		for _, f := range freqs {
			curFreqs = append(curFreqs, f.Current)
			freqSum += f.Current

			maxFreq = max(maxFreq, f.Max)
			minFreq = min(minFreq, f.Min)
		}
		avgFreq := freqSum / float32(len(freqs))

		// CPU percent load
		rawLoads, _ := cpu.Percent(100000000, true)
		loads := []float64{}
		for _, load := range rawLoads {
			loads = append(loads, math.Round(load*10)/10)
		}
		loadSum := 0.0
		for _, l := range loads {
			loadSum += l
		}
		avgLoad := math.Round(loadSum/float64(len(loads))*10) / 10

		// System overall load
		sysLoad, _ := load.Avg()

		ctx.JSON(http.StatusOK, gin.H{
			"name": info[0].ModelName,
			"cpus": gin.H{
				"count":    logicalCount,
				"physical": physicalCount,
			},
			"freq": gin.H{
				"min":     minFreq,
				"max":     maxFreq,
				"current": avgFreq,
				"per_cpu": curFreqs,
			},
			"load_percent": gin.H{
				"current": avgLoad,
				"per_cpu": loads,
			},
			"load": []float64{sysLoad.Load1, sysLoad.Load5, sysLoad.Load15},
		})
	})

	router.GET("/cpu/info", func(ctx *gin.Context) {
		// Info
		cpuInfo, _ := cpu.Info()

		// Core count
		logicalCount, physicalCount := getCPUCoreCount()

		// Arch
		hostInfo, _ := host.Info()

		payload := readCPUInfo("/proc/cpuinfo")
		if payload == nil {
			utils.AnswerError("Internal server error, sorry", nil, ctx)
			return
		}

		payload["name"] = cpuInfo[0].ModelName
		payload["cpus"] = gin.H{
			"count":    logicalCount,
			"physical": physicalCount,
		}
		payload["arch"] = hostInfo.KernelArch
		ctx.JSON(http.StatusOK, payload)
	})
}
