package leaderboard

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/google/go-cmp/cmp"
)

func TestGetTopTenPlayers(t *testing.T) {
	redisClient, redisMock := redismock.NewClientMock()

	zMember := redis.Z{Score: 765, Member: "Konami"}
	zRangeBy := redis.ZRangeBy{Min: "-inf", Max: "+inf", Offset: 0, Count: 10}
	redisMock.ExpectZRevRangeByScoreWithScores(RedisStorageKey, &zRangeBy).SetVal([]redis.Z{zMember})

	result, err := GetTopTenPlayers(context.Background(), redisClient)
	if err != nil {
		t.Error(err)
	}

	err = redisMock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}

	if len(result) != 1 {
		t.Error("Incorrect size of redis result")
	}

	expectedResult := TopPlayerFormat{ClientId: "Konami", Score: 765}
	if !cmp.Equal(result[0], expectedResult) {
		t.Error("Unexpected redis result")
	}
}

func TestUploadScore(t *testing.T) {
	redisClient, redisMock := redismock.NewClientMock()

	zMember := redis.Z{Score: 95.27, Member: "JohnCena"}
	redisMock.ExpectZAdd(RedisStorageKey, &zMember).SetVal(1)

	err := UploadScore(context.Background(), redisClient, "JohnCena", 95.27)
	if err != nil {
		t.Error(err)
	}

	err = redisMock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func TestClearAllScores(t *testing.T) {
	redisClient, redisMock := redismock.NewClientMock()

	redisMock.ExpectDel(RedisStorageKey).SetVal(1)

	err := ClearAllScores(context.Background(), redisClient)
	if err != nil {
		t.Error(err)
	}

	err = redisMock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
