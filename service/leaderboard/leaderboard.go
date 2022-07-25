package leaderboard

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const (
	RedisStorageKey = "homework:leaderboard:storage"
)

type UploadScoreFormat struct {
	Score float64 `json:"score"`
}

// Same struct as redis.Z
// But different field name
type TopPlayerFormat struct {
	ClientId any     `json:"clientId"`
	Score    float64 `json:"score"`
}

func GetTopTenPlayers(ctx context.Context, redisClient redis.UniversalClient) ([]TopPlayerFormat, error) {
	// Get top 10 scores descending
	zRangeBy := redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  10,
	}

	result, err := redisClient.ZRevRangeByScoreWithScores(ctx, RedisStorageKey, &zRangeBy).Result()
	if err != nil {
		return []TopPlayerFormat{}, err
	}

	topPlayersList := make([]TopPlayerFormat, len(result))

	for i, item := range result {
		topPlayersList[i] = TopPlayerFormat{
			ClientId: item.Member,
			Score:    item.Score,
		}
	}

	return topPlayersList, nil
}

func UploadScore(ctx context.Context, redisClient redis.UniversalClient, clientId string, score float64) error {
	zMember := redis.Z{Score: score, Member: clientId}

	return redisClient.ZAdd(ctx, RedisStorageKey, &zMember).Err()
}

func ClearAllScores(ctx context.Context, redisClient redis.UniversalClient) error {
	return redisClient.Del(ctx, RedisStorageKey).Err()
}
