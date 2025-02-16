package entity_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"timble/module/users/entity"
)

func TestReaction_NewReactionPayload(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedResult entity.ReactionParams
		expectedErr    error
	}{
		{
			name: "normal case",
			body: `{
		      "target_id":  2,
		      "type": 1
		    }`,
			expectedResult: entity.ReactionParams{
				UserID:   1,
				TargetID: 2,
				Type:     1,
			},
		},
		{
			name: "error case with invalid body",
			body: `{
		      "type": 1`,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: unexpected EOF; field: payload"),
		},
		{
			name: "error case with invalid target ID",
			body: `{
		      "type": 1
		    }`,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: Invalid target user; field: target_id"),
		},
		{
			name: "error case with same target ID as user ID",
			body: `
		    {
		      "target_id":  1,
		      "type": 1
		    }
		  `,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: Invalid target user; field: target_id"),
		},
		{
			name: "error case with invalid type",
			body: `
		    {
		      "target_id":  2,
		      "type": 3
		    }
		  `,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: Invalid reaction type; field: type"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := entity.NewReactionPayload(strings.NewReader(tc.body), 1)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}
