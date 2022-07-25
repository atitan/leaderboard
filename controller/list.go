package controller

import (
	"encoding/json"
	"net/http"

	"homework/service/leaderboard"

	"github.com/uptrace/bunrouter"
)

func (c *ControllerBase) List(w http.ResponseWriter, req bunrouter.Request) error {
	topPlayersList, err := leaderboard.GetTopTenPlayers(req.Context(), c.redisClient)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(w)
	return encoder.Encode(map[string]any{"topPlayers": topPlayersList})
}
