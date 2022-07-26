package controller

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"homework/service/leaderboard"
	"homework/utility/mytype"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/uptrace/bunrouter"
)

func TestUploadWithoutClientId(t *testing.T) {
	redisClient, _ := redismock.NewClientMock()

	r := httptest.NewRequest(http.MethodPost, "/api/v1/score", nil)
	req := bunrouter.NewRequest(r)
	w := httptest.NewRecorder()

	controllerBase := NewController(redisClient)
	err := controllerBase.Upload(w, req)

	apiError, ok := err.(mytype.ApiError)
	if !ok {
		t.Error("Controller should return ApiError")
	}

	if apiError.Status != http.StatusBadRequest {
		t.Error("ApiError.Status should be http.StatusBadRequest")
	}

	if apiError.Hint != "Invalid ClientId" {
		t.Error("Unexpected ApiError.Hint")
	}
}

func TestUploadWithInvalidJSON(t *testing.T) {
	redisClient, _ := redismock.NewClientMock()

	r := httptest.NewRequest(http.MethodPost, "/api/v1/score", nil)
	r.Header.Add("ClientId", "Rick")
	req := bunrouter.NewRequest(r)
	w := httptest.NewRecorder()

	controllerBase := NewController(redisClient)
	err := controllerBase.Upload(w, req)

	apiError, ok := err.(mytype.ApiError)
	if !ok {
		t.Error("Controller should return ApiError")
	}

	if apiError.Status != http.StatusBadRequest {
		t.Error("ApiError.Status should be http.StatusBadRequest")
	}

	if apiError.Hint != "Malformed json" {
		t.Error("Unexpected ApiError.Hint")
	}
}

func TestUploadSuccess(t *testing.T) {
	redisClient, redisMock := redismock.NewClientMock()

	zMember := redis.Z{Score: 33.89, Member: "Rick"}
	redisMock.ExpectZAdd(leaderboard.RedisStorageKey, &zMember).SetVal(1)

	reqBody := strings.NewReader(`{"score": 33.89}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/score", reqBody)
	r.Header.Add("ClientId", "Rick")
	req := bunrouter.NewRequest(r)
	w := httptest.NewRecorder()

	controllerBase := NewController(redisClient)
	err := controllerBase.Upload(w, req)
	if err != nil {
		t.Error(err)
	}

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Error("Status code should be http.StatusOK")
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, res.Body)
	if err != nil {
		t.Error(err)
	}

	if buf.String() != "{\"status\":\"ok\"}\n" {
		t.Error("Incorrect response")
	}
}
