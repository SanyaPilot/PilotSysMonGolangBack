package endpoints

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTimeRoutes(t *testing.T) {
	r := gin.Default()
	TimeRoutes(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/time", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := make(map[string]any)
	json.Unmarshal([]byte(rawBody), &body)

	if assert.Equal(t, 200, w.Code) {
		assert.IsType(t, 0.0, body["time"])
		assert.IsType(t, "", body["timezone"])
		assert.IsType(t, 0.0, body["uptime"])
	}
}
