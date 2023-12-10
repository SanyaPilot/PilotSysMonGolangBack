package endpoints

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetCPUFreqPaths(t *testing.T) {
	paths := getCPUFreqPaths()
	assert.NotNilf(t, paths, "CPUFreq paths aren't found. Not running on Linux?")
	assert.NotEmptyf(t, paths, "CPUFreq paths len is 0!")
	wantPrefix := "/sys/devices/system/cpu/"
	if !strings.HasPrefix(paths[0], wantPrefix) {
		t.Fatalf("CPUFreq path has a wrong prefix!\nWant: %s\nGot: %s", wantPrefix, paths[0])
	}
}

func TestGetCPUFreqPathsInvalid(t *testing.T) {
	paths := getCPUFreqPaths("/some/invalid/prefix/")
	assert.Nil(t, paths)
}

func TestReadCPUFreqNode(t *testing.T) {
	tryPath := "/sys/devices/system/cpu/cpufreq/policy0"
	var path string
	_, err := os.Stat(tryPath)
	if err == nil {
		path = tryPath
	} else {
		tryPath := "/sys/devices/system/cpu/cpu0/cpufreq"
		_, err := os.Stat(tryPath)
		if err != nil {
			t.Skipf("The system doesn't support CPUFreq scaling")
		}
		path = tryPath
	}
	res := readCPUFreqNode(path, "scaling_cur_freq")
	assert.GreaterOrEqualf(t, res, float32(100), "Invalid scaling_cur_freq value! Got %f", res)
}

func TestReadCPUFreqInvalidNode(t *testing.T) {
	res := readCPUFreqNode("/some/invalid/sysfs/path", "scaling_cur_freq")
	assert.Zerof(t, res, "A non zero value returned while reading invalid sysfs node!")
}

func TestGetCPUFreqs(t *testing.T) {
	freqs := getCPUFreqs(getCPUFreqPaths())
	assert.NotEmpty(t, freqs)
	for _, freq := range freqs {
		assert.NotZero(t, freq.Min)
		assert.NotZero(t, freq.Max)
		assert.NotZero(t, freq.Current)
	}
}

func TestGetInvalidCPUFreqs(t *testing.T) {
	freqs := getCPUFreqs(nil)
	assert.Nil(t, freqs)
}

func TestGetCPUCoreCount(t *testing.T) {
	log, phys := getCPUCoreCount()
	assert.Greater(t, log, 0)
	assert.Greater(t, phys, 0)
}

func TestReadCPUInfo(t *testing.T) {
	res := readCPUInfo("/proc/cpuinfo")
	assert.NotEmptyf(t, res, "CPU info map is empty!")
	for k := range res {
		assert.NotContainsf(t, k, " ", "CPU info map key contains spaces!")
	}
}

func TestReadInvalidCPUInfo(t *testing.T) {
	res := readCPUInfo("/invalid/file")
	assert.Nil(t, res)
}

func TestMainCPURoute(t *testing.T) {
	r := gin.Default()
	CPURoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cpu", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := make(map[string]any)
	json.Unmarshal([]byte(rawBody), &body)

	assert.Equal(t, 200, w.Code)
	assert.IsType(t, "", body["name"])
	if assert.IsType(t, map[string]any{}, body["cpus"]) {
		assert.IsType(t, 0.0, body["cpus"].(map[string]any)["count"])
		assert.IsType(t, 0.0, body["cpus"].(map[string]any)["physical"])
	}
	if assert.IsType(t, map[string]any{}, body["freq"]) {
		assert.IsType(t, 0.0, body["freq"].(map[string]any)["min"])
		assert.IsType(t, 0.0, body["freq"].(map[string]any)["max"])
		assert.IsType(t, 0.0, body["freq"].(map[string]any)["current"])
		assert.IsType(t, []any{}, body["freq"].(map[string]any)["per_cpu"])
	}
	if assert.IsType(t, map[string]any{}, body["load_percent"]) {
		assert.IsType(t, 0.0, body["load_percent"].(map[string]any)["current"])
		assert.IsType(t, []any{}, body["freq"].(map[string]any)["per_cpu"])
	}
	assert.IsType(t, []any{}, body["load"])
}

func TestCPUInfoRoute(t *testing.T) {
	r := gin.Default()
	CPURoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cpu/info", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := make(map[string]any)
	json.Unmarshal([]byte(rawBody), &body)

	assert.Equal(t, 200, w.Code)
	assert.IsType(t, "", body["name"])
	assert.IsType(t, "", body["arch"])
}
