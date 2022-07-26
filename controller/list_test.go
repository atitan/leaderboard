package controller

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"homework/service/leaderboard"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/uptrace/bunrouter"
)

func TestList(t *testing.T) {
	redisClient, redisMock := redismock.NewClientMock()

	zMember := redis.Z{Score: 876, Member: "Banamu"}
	zRangeBy := redis.ZRangeBy{Min: "-inf", Max: "+inf", Offset: 0, Count: 10}
	redisMock.ExpectZRevRangeByScoreWithScores(leaderboard.RedisStorageKey, &zRangeBy).SetVal([]redis.Z{zMember})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/leaderboard", nil)
	req := bunrouter.NewRequest(r)
	w := httptest.NewRecorder()

	controllerBase := NewController(redisClient)
	err := controllerBase.List(w, req)
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

	if buf.String() != "{\"topPlayers\":[{\"clientId\":\"Banamu\",\"score\":876}]}\n" {
		t.Error("Incorrect response")
	}
}
