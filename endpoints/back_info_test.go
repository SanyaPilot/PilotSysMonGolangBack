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

func TestBackInfoBackend(t *testing.T) {
	r := gin.Default()
	BackInfoEndpoints(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/backend_info", nil)
	r.ServeHTTP(w, req)
	rawBody := w.Body.String()
	body := make(map[string]any)
	json.Unmarshal([]byte(rawBody), &body)

	if assert.Equal(t, 200, w.Code) {
		assert.Equal(t, "golang", body["lang"])
		assert.Equal(t, runtime.Version(), body["lang_version"])
		assert.IsType(t, "", body["version"])
		assert.IsType(t, "", body["commit_hash"])
		assert.IsType(t, "", body["build_time"])
	}
}
