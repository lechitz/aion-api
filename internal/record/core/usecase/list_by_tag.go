package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ListByTag returns records filtered by tag for the authenticated user.
func (s *Service) ListByTag(ctx context.Context, tagID uint64, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListByTag)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListByTag),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.Int("limit", limit),
	)

	if limit <= 0 || limit > 100 {
		limit = 50 // default limit
	}

	span.AddEvent(EventRepositoryList)
	records, err := s.RecordRepository.ListByTag(ctx, tagID, userID, limit, afterEventTime, afterID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToListRecords)
		s.Logger.ErrorwCtx(ctx, FailedToListRecords,
			commonkeys.TagID, tagID,
			commonkeys.UserID, userID,
			commonkeys.Error, err,
		)
		return nil, fmt.Errorf("%s: %w", FailedToListRecords, err)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusListedAll)
	s.Logger.InfowCtx(ctx, "records listed by tag successfully",
		commonkeys.TagID, tagID,
		commonkeys.UserID, userID,
		"count", len(records),
	)

	return records, nil
}

// CountByTag returns the number of active records for a tag and user.
func (s *Service) CountByTag(ctx context.Context, tagID uint64, userID uint64) (int64, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "record.usecase.CountByTag")
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, "record.usecase.CountByTag"),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
	)

	count, err := s.RecordRepository.CountByTag(ctx, tagID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToListRecords)
		s.Logger.ErrorwCtx(ctx, "failed to count records by tag",
			commonkeys.TagID, tagID,
			commonkeys.UserID, userID,
			commonkeys.Error, err,
		)
		return 0, fmt.Errorf("failed to count records by tag: %w", err)
	}

	span.SetAttributes(attribute.Int64("count", count))
	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusListedAll)
	return count, nil
}
