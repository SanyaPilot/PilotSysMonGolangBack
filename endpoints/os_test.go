package endpoints

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetFreeDesktopInfo(t *testing.T) {
	_, err := getFreeDesktopInfo("/etc/os-release")
	assert.Nil(t, err)
}

func TestGetBrokenFreeDesktopInfo(t *testing.T) {
	_, err := getFreeDesktopInfo("/random/path")
	assert.NotNil(t, err)
}

func TestOSRoute(t *testing.T) {
	r := gin.Default()
	OSRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/os", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := make(map[string]any)
	json.Unmarshal([]byte(rawBody), &body)

	if assert.Equal(t, 200, w.Code) {
		assert.IsType(t, "", body["family"])
		assert.IsType(t, "", body["version"])
		assert.IsType(t, "", body["release"])

		if runtime.GOOS == "linux" {
			assert.IsType(t, "", body["name"])
			assert.IsType(t, "", body["url"])
		}
	}
}
