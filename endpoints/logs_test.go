package endpoints

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogsRoute(t *testing.T) {
	r := gin.Default()
	LogsRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs?level=warn&boot=0", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := []any{}
	json.Unmarshal([]byte(rawBody), &body)

	if assert.Equal(t, 200, w.Code) {
		if assert.IsType(t, map[string]any{}, body[0]) {
			entry := body[0].(map[string]any)
			assert.IsType(t, "", entry["id"])
			assert.IsType(t, "", entry["level"])
			assert.IsType(t, 0.0, entry["time"])
		}
	}
}

func TestLogsRouteInvalidLevel(t *testing.T) {
	r := gin.Default()
	LogsRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs?level=123&boot=0", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestLogsRouteInvalidBoot(t *testing.T) {
	r := gin.Default()
	LogsRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs?level=error&boot=alo", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestLogsRouteEmptyOut(t *testing.T) {
	r := gin.Default()
	LogsRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs?level=emerg&boot=0", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := []any{}
	json.Unmarshal([]byte(rawBody), &body)

	if assert.Equal(t, 200, w.Code) {
		assert.Empty(t, body)
	}
}

func TestLogsRouteJournalctlFail(t *testing.T) {
	r := gin.Default()
	LogsRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs?level=emerg&since=123-123-123&until=222-22-22", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := make(map[string]string)
	json.Unmarshal([]byte(rawBody), &body)

	if assert.Equal(t, 500, w.Code) {
		assert.Equal(t, "error", body["status"])
	}
}

func TestLogsRouteDayParam(t *testing.T) {
	r := gin.Default()
	LogsRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs?level=error&day=2023-12-14", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestLogsRouteIdParam(t *testing.T) {
	r := gin.Default()
	LogsRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs?level=error&id=kernel", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
