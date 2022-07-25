package controller

import (
	"encoding/json"
	"net/http"

	"homework/service/leaderboard"
	"homework/utility/mytype"

	"github.com/uptrace/bunrouter"
)

func (c *ControllerBase) Upload(w http.ResponseWriter, req bunrouter.Request) error {
	clientId := req.Header.Get("ClientId")
	if clientId == "" {
		return mytype.ApiError{
			Status: http.StatusBadRequest,
			Hint:   "Invalid ClientId",
		}
	}

	uploadPayload := leaderboard.UploadScoreFormat{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&uploadPayload)
	if err != nil {
		return err
	}

	err = leaderboard.UploadScore(req.Context(), c.redisClient, clientId, uploadPayload.Score)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(w)
	return encoder.Encode(map[string]any{"status": "ok"})
}
