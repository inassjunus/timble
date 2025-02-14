package entity_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"timble/module/users/entity"
)

func TestReaction_NewReactionPayload(t *testing.T) {
	type args struct {
		configName string
		configData string
	}
	tests := []struct {
		name           string
		body           string
		expectedResult entity.ReactionParams
		expectedErr    error
	}{
		{
			name: "normal case",
			body: `
		    {
		      "target_id":  2,
		      "type": 1
		    }
		  `,
			expectedResult: entity.ReactionParams{
				UserID:   1,
				TargetID: 2,
				Type:     1,
			},
		},
		{
			name: "error case with invalid target ID",
			body: `
		    {
		      "type": 1
		    }
		  `,
			expectedErr: errors.New("Invalid target user"),
		},
		{
			name: "error case with same target ID as user ID",
			body: `
		    {
		      "target_id":  1,
		      "type": 1
		    }
		  `,
			expectedErr: errors.New("Invalid target user"),
		},
		{
			name: "error case with invalid type",
			body: `
		    {
		      "target_id":  2,
		      "type": 3
		    }
		  `,
			expectedErr: errors.New("Invalid reaction type"),
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
