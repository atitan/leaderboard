package main

import (
	"context"
	"net/http"
	"time"

	"homework/controller"
	"homework/middleware/rescuer"
	"homework/service/leaderboard"

	"github.com/go-redis/redis/v8"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Clear redis storage on startup and every 10 minutes
	go func() {
		var err error

		for {
			err = leaderboard.ClearAllScores(context.Background(), redisClient)
			if err != nil {
				// Currently not handled
			}

			<-time.After(10 * time.Minute)
		}
	}()

	handler := controller.NewController(redisClient)

	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
		bunrouter.Use(rescuer.NewMiddleware()),
	)

	router.WithGroup("/api/v1", func(group *bunrouter.Group) {
		group.GET("/leaderboard", handler.List)
		group.POST("/score", handler.Upload)
	})

	http.ListenAndServe(":8080", router)
}
