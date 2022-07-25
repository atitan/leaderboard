package controller

import (
	"github.com/go-redis/redis/v8"
)

type ControllerBase struct {
	redisClient redis.UniversalClient
}

func NewController(redisClient redis.UniversalClient) *ControllerBase {
	return &ControllerBase{
		redisClient: redisClient,
	}
}
