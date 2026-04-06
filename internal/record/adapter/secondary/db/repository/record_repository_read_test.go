package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRecordReadBaseQueries(t *testing.T) {
	repo, dbMock := newRecordRepo(t)
	rec := sampleRecord()
	at := rec.EventTime.Format(time.RFC3339)
	afterID := int64(10)

	t.Run("get by id success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*model.Record)
			require.True(t, ok)
			*row = model.Record{ID: rec.ID, UserID: rec.UserID, TagID: rec.TagID}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)
		got, err := repo.GetByID(t.Context(), rec.ID, rec.UserID)
		require.NoError(t, err)
		require.Equal(t, rec.ID, got.ID)
	})

	t.Run("list latest error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(5).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("query fail"))
		_, err := repo.ListLatest(t.Context(), rec.UserID, 5)
		require.Error(t, err)
	})

	t.Run("list by day success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			rows, ok := dest.(*[]model.Record)
			require.True(t, ok)
			*rows = []model.Record{{ID: rec.ID, UserID: rec.UserID}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)
		got, err := repo.ListByDay(t.Context(), rec.UserID, rec.EventTime)
		require.NoError(t, err)
		require.Len(t, got, 1)
	})

	t.Run("list by user with cursor", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(10).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			rows, ok := dest.(*[]model.Record)
			require.True(t, ok)
			*rows = []model.Record{{ID: rec.ID}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)
		got, err := repo.ListByUser(t.Context(), rec.UserID, 10, &at, &afterID)
		require.NoError(t, err)
		require.Len(t, got, 1)
	})
}

func TestRecordReadFilterQueries(t *testing.T) {
	repo, dbMock := newRecordRepo(t)
	rec := sampleRecord()

	t.Run("list by tag success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(3).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			rows, ok := dest.(*[]model.Record)
			require.True(t, ok)
			*rows = []model.Record{{ID: rec.ID}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)
		got, err := repo.ListByTag(t.Context(), rec.TagID, rec.UserID, 3, nil, nil)
		require.NoError(t, err)
		require.Len(t, got, 1)
	})

	t.Run("count by tag success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Count(gomock.Any()).DoAndReturn(func(dest any) db.DB {
			value, ok := dest.(*int64)
			require.True(t, ok)
			*value = 10
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)
		got, err := repo.CountByTag(t.Context(), rec.TagID, rec.UserID)
		require.NoError(t, err)
		require.EqualValues(t, 10, got)
	})

	t.Run("list by category and range variants", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(4).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			rows, ok := dest.(*[]model.Record)
			require.True(t, ok)
			*rows = []model.Record{{ID: 1}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)
		_, err := repo.ListByCategory(t.Context(), 9, rec.UserID, 4, nil, nil)
		require.NoError(t, err)

		start := rec.EventTime.Add(-time.Hour)
		end := rec.EventTime.Add(time.Hour)
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(4).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)
		_, err = repo.ListAllBetween(t.Context(), rec.UserID, start, end, 4)
		require.NoError(t, err)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(4).Return(dbMock)
		dbMock.EXPECT().Find(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(nil)
		_, err = repo.ListAllUntil(t.Context(), rec.UserID, end, 4)
		require.NoError(t, err)
	})
}
