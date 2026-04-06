package usecase_test

import (
	"testing"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_ListByTag_LimitBounds(t *testing.T) {
	userID := uint64(1)
	tagID := uint64(3)
	records := []domain.Record{{ID: 1, UserID: userID}}

	tests := []struct {
		name        string
		limit       int
		expectedLim int
	}{
		{
			name:        "default limit when zero",
			limit:       0,
			expectedLim: 50,
		},
		{
			name:        "default limit when too large",
			limit:       150,
			expectedLim: 50,
		},
		{
			name:        "uses provided limit",
			limit:       15,
			expectedLim: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup.RecordServiceTest(t)
			defer suite.Ctrl.Finish()

			suite.RecordRepository.EXPECT().
				ListByTag(gomock.Any(), tagID, userID, tt.expectedLim, nil, nil).
				Return(records, nil)

			result, err := suite.RecordService.ListByTag(suite.Ctx, tagID, userID, tt.limit, nil, nil)
			require.NoError(t, err)
			assert.Equal(t, records, result)
		})
	}
}
