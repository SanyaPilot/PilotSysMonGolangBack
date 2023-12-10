package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestAnswerError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	AnswerGopsutilError(errors.New("TEST"), ctx)
	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	rawBody := w.Body.String()
	body := make(map[string]any)
	json.Unmarshal([]byte(rawBody), &body)
	assert.Equal(t, "error", body["status"])
	assert.Equal(t, "gopsutil error: TEST", body["msg"])
}

func TestAnswerErrorNil(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	AnswerError("testing test", nil, ctx)
	res := w.Result()
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)

	rawBody := w.Body.String()
	body := make(map[string]any)
	json.Unmarshal([]byte(rawBody), &body)
	assert.Equal(t, "error", body["status"])
	assert.Equal(t, "testing test", body["msg"])
}
