package endpoints

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDisksRoute(t *testing.T) {
	r := gin.Default()
	DisksRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/disks", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := []any{}
	json.Unmarshal([]byte(rawBody), &body)

	if assert.Equal(t, 200, w.Code) {
		if assert.IsType(t, map[string]any{}, body[0]) {
			entry := body[0].(map[string]any)
			assert.IsType(t, "", entry["device"])
			assert.IsType(t, "", entry["fs"])
			assert.IsType(t, "", entry["mountpoint"])
			assert.IsType(t, "", entry["opts"])
			if assert.IsType(t, map[string]any{}, entry["usage"]) {
				usage := entry["usage"].(map[string]any)
				assert.IsType(t, 0.0, usage["free"])
				assert.IsType(t, 0.0, usage["percent"])
				assert.IsType(t, 0.0, usage["total"])
				assert.IsType(t, 0.0, usage["used"])
			}
		}
	}
}

func TestDisksHumanRoute(t *testing.T) {
	r := gin.Default()
	DisksRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/disks?human=true", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := []any{}
	json.Unmarshal([]byte(rawBody), &body)

	if assert.Equal(t, 200, w.Code) {
		if assert.IsType(t, map[string]any{}, body[0]) {
			entry := body[0].(map[string]any)
			assert.IsType(t, "", entry["device"])
			assert.IsType(t, "", entry["fs"])
			assert.IsType(t, "", entry["mountpoint"])
			assert.IsType(t, "", entry["opts"])
			if assert.IsType(t, map[string]any{}, entry["usage"]) {
				usage := entry["usage"].(map[string]any)
				assert.IsType(t, "", usage["free"])
				assert.IsType(t, 0.0, usage["percent"])
				assert.IsType(t, "", usage["total"])
				assert.IsType(t, "", usage["used"])
			}
		}
	}
}

func TestDisksBrokenRoute(t *testing.T) {
	r := gin.Default()
	DisksRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/disks?human=xd", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := make(map[string]any)
	json.Unmarshal([]byte(rawBody), &body)

	if assert.Equal(t, 400, w.Code) {
		assert.Equal(t, "error", body["status"])
	}
}
