package entity

import (
	"encoding/json"
	"io"
	"timble/internal/utils"
	"time"
)

type UserReaction struct {
	UserID    uint      `json:"user_id"`
	TargetID  uint      `json:"target_id"`
	Type      int       `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReactionParams struct {
	UserID   uint `json:"user_id"`
	TargetID uint `json:"target_id"`
	Type     int  `json:"type"`
}

var (
	ReactionTypes = map[int]bool{
		0: true, // undecided
		1: true, // not interested (pass)
		2: true, // like
	}
)

func NewReactionPayload(body io.Reader, userID uint) (ReactionParams, error) {
	params := ReactionParams{
		UserID: userID,
	}
	err := json.NewDecoder(body).Decode(&params)
	if err != nil {
		return params, err
	}

	if params.TargetID <= 0 || params.TargetID == userID {
		return params, utils.BadRequestParamError("Invalid target user", "target_id")
	}

	if !ReactionTypes[params.Type] {
		return params, utils.BadRequestParamError("Invalid reaction type", "type")
	}
	return params, nil
}
