package entity

import (
	"encoding/json"
	"io"
)

type UserReaction struct {
	UserID    uint   `json:"user_id"`
	TargetID  uint   `json:"target_id"`
	Type      int    `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ReactionParams struct {
	UserID   uint `json:"user_id"`
	TargetID uint `json:"target_id"`
	Type     int  `json:"type"`
}

var (
	ReactionTypes = map[int]bool{
		1: true, // pass
		2: true, // like
		3: true, // block
	}
)

func NewReactionPayload(body io.Reader, userID uint) (ReactionParams, error) {
	params := ReactionParams{
		UserID: userID,
	}
	err := json.NewDecoder(body).Decode(&params)
	return params, err
}
