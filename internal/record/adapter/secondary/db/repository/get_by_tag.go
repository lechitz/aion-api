package repository

import (
	"context"

	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/record/core/domain"
)

// ListByTag returns records filtered by tag for a given user.
func (r *RecordRepository) ListByTag(
	ctx context.Context,
	tagID uint64,
	userID uint64,
	limit int,
	afterEventTime *string,
	afterID *int64,
) ([]domain.Record, error) {
	var recordsDB []model.Record

	q := r.db.WithContext(ctx).
		Where("tag_id = ? AND user_id = ? AND deleted_at IS NULL", tagID, userID).
		Order("event_time DESC, id DESC").
		Limit(limit)

	if afterEventTime != nil && afterID != nil {
		q = q.Where("event_time < ? OR (event_time = ? AND id < ?)", *afterEventTime, *afterEventTime, *afterID)
	}

	if err := q.Find(&recordsDB).Error(); err != nil {
		return nil, err
	}

	return mapper.RecordsFromDB(recordsDB), nil
}

// CountByTag returns the number of active records for a tag and user.
func (r *RecordRepository) CountByTag(ctx context.Context, tagID uint64, userID uint64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Record{}).
		Where("tag_id = ? AND user_id = ? AND deleted_at IS NULL", tagID, userID).
		Count(&count).
		Error(); err != nil {
		return 0, err
	}

	return count, nil
}
