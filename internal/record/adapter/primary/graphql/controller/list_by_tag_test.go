package controller_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/record/adapter/primary/graphql/controller"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListByTag_UserIDMissing(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByTag(t.Context(), 1, 0, 10, nil, nil)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)
	assert.Nil(t, out)
}

func TestListByTag_TagIDInvalid(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByTag(t.Context(), 0, 1, 10, nil, nil)
	require.Error(t, err)
	assert.Equal(t, "tag id cannot be zero", err.Error())
	assert.Nil(t, out)
}

func TestListByTag_ServiceError(t *testing.T) {
	expected := errors.New("list failed")
	svc := &recordServiceStub{
		listByTagFn: func(_ context.Context, _, _ uint64, _ int, _ *string, _ *int64) ([]domain.Record, error) {
			return nil, expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByTag(t.Context(), 2, 3, 10, nil, nil)
	require.ErrorIs(t, err, expected)
	assert.Nil(t, out)
}

func TestListByTag_Success(t *testing.T) {
	when := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	svc := &recordServiceStub{
		listByTagFn: func(_ context.Context, tagID uint64, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error) {
			require.Equal(t, uint64(5), tagID)
			require.Equal(t, uint64(6), userID)
			require.Equal(t, 3, limit)
			require.Nil(t, afterEventTime)
			require.Nil(t, afterID)
			return []domain.Record{{ID: 1, UserID: userID, TagID: tagID, EventTime: when, CreatedAt: when, UpdatedAt: when}}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByTag(t.Context(), 5, 6, 3, nil, nil)
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "1", out[0].ID)
}

func TestCountByTag_Success(t *testing.T) {
	svc := &recordServiceStub{
		countByTagFn: func(_ context.Context, tagID uint64, userID uint64) (int64, error) {
			require.Equal(t, uint64(5), tagID)
			require.Equal(t, uint64(6), userID)
			return 10, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.CountByTag(t.Context(), 5, 6)
	require.NoError(t, err)
	assert.EqualValues(t, 10, out)
}
