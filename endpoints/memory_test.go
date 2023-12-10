package endpoints

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRoute(t *testing.T) {
	r := gin.Default()
	MemRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/memory", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := make(map[string]any)
	json.Unmarshal([]byte(rawBody), &body)

	if assert.Equal(t, 200, w.Code) {
		if assert.IsType(t, map[string]any{}, body["ram"]) {
			assert.IsType(t, 0.0, body["ram"].(map[string]any)["percent"])
			assert.IsType(t, 0.0, body["ram"].(map[string]any)["total"])
			assert.IsType(t, 0.0, body["ram"].(map[string]any)["available"])
			assert.IsType(t, 0.0, body["ram"].(map[string]any)["used"])
			assert.IsType(t, 0.0, body["ram"].(map[string]any)["free"])
		}
		if assert.IsType(t, map[string]any{}, body["swap"]) {
			assert.IsType(t, 0.0, body["ram"].(map[string]any)["percent"])
			assert.IsType(t, 0.0, body["ram"].(map[string]any)["total"])
			assert.IsType(t, 0.0, body["ram"].(map[string]any)["used"])
			assert.IsType(t, 0.0, body["ram"].(map[string]any)["free"])
		}
	}
}

func TestMemoryHumanRoute(t *testing.T) {
	r := gin.Default()
	MemRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/memory?human=true", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := make(map[string]any)
	json.Unmarshal([]byte(rawBody), &body)

	if assert.Equal(t, 200, w.Code) {
		if assert.IsType(t, map[string]any{}, body["ram"]) {
			assert.IsType(t, 0.0, body["ram"].(map[string]any)["percent"])
			assert.IsType(t, "", body["ram"].(map[string]any)["total"])
			assert.IsType(t, "", body["ram"].(map[string]any)["available"])
			assert.IsType(t, "", body["ram"].(map[string]any)["used"])
			assert.IsType(t, "", body["ram"].(map[string]any)["free"])
		}
		if assert.IsType(t, map[string]any{}, body["swap"]) {
			assert.IsType(t, 0.0, body["ram"].(map[string]any)["percent"])
			assert.IsType(t, "", body["ram"].(map[string]any)["total"])
			assert.IsType(t, "", body["ram"].(map[string]any)["used"])
			assert.IsType(t, "", body["ram"].(map[string]any)["free"])
		}
	}
}
