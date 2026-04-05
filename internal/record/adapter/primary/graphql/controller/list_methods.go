package controller

import (
	"context"
	"strconv"
	"time"

	gmodel "github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ListByTag fetches all Records for a specific tag for the authenticated user.
func (h *controller) ListByTag(ctx context.Context, tagID, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListByTag)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListByTag),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.Int(AttrLimit, limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	if tagID == 0 {
		span.SetStatus(codes.Error, ErrTagIDCannotBeZero.Error())
		h.Logger.ErrorwCtx(ctx, ErrTagIDCannotBeZero.Error(), commonkeys.TagID, tagID)
		return nil, ErrTagIDCannotBeZero
	}

	records, err := h.RecordService.ListByTag(ctx, tagID, userID, limit, afterEventTime, afterID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgListByTagError)
		h.Logger.ErrorwCtx(ctx, MsgListByTagError, commonkeys.Error, err.Error(), commonkeys.TagID, tagID, commonkeys.UserID, userID)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int(AttrCount, len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// CountByTag returns the number of active records for a specific tag for the authenticated user.
func (h *controller) CountByTag(ctx context.Context, tagID, userID uint64) (int64, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "record.graphql.CountByTag")
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, "record.graphql.CountByTag"),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return 0, ErrUserIDNotFound
	}

	if tagID == 0 {
		span.SetStatus(codes.Error, ErrTagIDCannotBeZero.Error())
		h.Logger.ErrorwCtx(ctx, ErrTagIDCannotBeZero.Error(), commonkeys.TagID, tagID)
		return 0, ErrTagIDCannotBeZero
	}

	count, err := h.RecordService.CountByTag(ctx, tagID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to count records by tag")
		h.Logger.ErrorwCtx(ctx, "failed to count records by tag", commonkeys.Error, err.Error(), commonkeys.TagID, tagID, commonkeys.UserID, userID)
		return 0, err
	}

	span.SetAttributes(attribute.Int64(AttrCount, count))
	span.SetStatus(codes.Ok, StatusFetched)
	return count, nil
}

// ListByCategory fetches all Records for a specific category for the authenticated user.
// Records are retrieved via JOIN (records → tags → categories).
func (h *controller) ListByCategory(ctx context.Context, categoryID, userID uint64, limit int) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListByCategory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListByCategory),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.Int(AttrLimit, limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	if categoryID == 0 {
		span.SetStatus(codes.Error, ErrCategoryIDCannotBeZero.Error())
		h.Logger.ErrorwCtx(ctx, ErrCategoryIDCannotBeZero.Error(), commonkeys.CategoryID, categoryID)
		return nil, ErrCategoryIDCannotBeZero
	}

	records, err := h.RecordService.ListByCategory(ctx, categoryID, userID, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgListByCategoryError)
		h.Logger.ErrorwCtx(ctx, MsgListByCategoryError, commonkeys.Error, err.Error(), commonkeys.CategoryID, categoryID, commonkeys.UserID, userID)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListByDay fetches all Records for a specific day for the authenticated user.
func (h *controller) ListByDay(ctx context.Context, userID uint64, dateStr string) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListByDay)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListByDay),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(AttrDate, dateStr),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	date, err := parseRecordsDayQuery(dateStr, time.Now())
	if err != nil {
		span.SetStatus(codes.Error, MsgInvalidDateFormat)
		h.Logger.ErrorwCtx(ctx, MsgInvalidDateFormat, AttrDate, dateStr, commonkeys.Error, err.Error())
		return nil, err
	}

	records, err := h.RecordService.ListByDay(ctx, userID, date)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgListByDayError)
		h.Logger.ErrorwCtx(ctx, MsgListByDayError, commonkeys.Error, err.Error(), AttrDate, dateStr, commonkeys.UserID, userID)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListAllUntil fetches Records with event_time up to (and including) the given timestamp.
func (h *controller) ListAllUntil(ctx context.Context, userID uint64, untilStr string, limit int) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListAllUntil)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListAllUntil),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(AttrUntil, untilStr),
		attribute.Int(AttrLimit, limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	until, err := time.Parse(time.RFC3339, untilStr)
	if err != nil {
		span.SetStatus(codes.Error, MsgInvalidUntilTimestamp)
		h.Logger.ErrorwCtx(ctx, MsgInvalidUntilTimestamp, AttrUntil, untilStr, commonkeys.Error, err.Error())
		return nil, ErrInvalidUntilTimestamp
	}

	records, err := h.RecordService.ListAllUntil(ctx, userID, until, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgListUntilError)
		h.Logger.ErrorwCtx(ctx, MsgListUntilError, commonkeys.Error, err.Error(), AttrUntil, untilStr, commonkeys.UserID, userID)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListAllBetween fetches Records with event_time within the specified date range.
func (h *controller) ListAllBetween(ctx context.Context, userID uint64, startDateStr, endDateStr string, limit int) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListAllBetween)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListAllBetween),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(AttrStartDate, startDateStr),
		attribute.String(AttrEndDate, endDateStr),
		attribute.Int(AttrLimit, limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		span.SetStatus(codes.Error, MsgInvalidStartDate)
		h.Logger.ErrorwCtx(ctx, MsgInvalidStartDate, AttrStartDate, startDateStr, commonkeys.Error, err.Error())
		return nil, ErrInvalidStartDate
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		span.SetStatus(codes.Error, MsgInvalidEndDate)
		h.Logger.ErrorwCtx(ctx, MsgInvalidEndDate, AttrEndDate, endDateStr, commonkeys.Error, err.Error())
		return nil, ErrInvalidEndDate
	}

	records, err := h.RecordService.ListAllBetween(ctx, userID, startDate, endDate, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgListBetweenError)
		h.Logger.ErrorwCtx(
			ctx,
			MsgListBetweenError,
			commonkeys.Error,
			err.Error(),
			AttrStartDate,
			startDateStr,
			AttrEndDate,
			endDateStr,
			commonkeys.UserID,
			userID,
		)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListByUser fetches records for the authenticated user with optional cursors.
func (h *controller) ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListAll)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListAll),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int(AttrLimit, limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	records, err := h.RecordService.ListByUser(ctx, userID, limit, afterEventTime, afterID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrFailedToListRecords.Error())
		h.Logger.ErrorwCtx(ctx, ErrFailedToListRecords.Error(), commonkeys.Error, err.Error())
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}

// ListLatest fetches the N most recent records for the authenticated user.
func (h *controller) ListLatest(ctx context.Context, userID uint64, limit int) ([]*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListLatest)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListLatest),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int(AttrLimit, limit),
	)

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound.Error(), commonkeys.UserID, userID)
		return nil, ErrUserIDNotFound
	}

	if limit <= 0 || limit > 100 {
		limit = 10 // default limit for latest
	}

	records, err := h.RecordService.ListLatest(ctx, userID, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgListLatestError)
		h.Logger.ErrorwCtx(ctx, MsgListLatestError, commonkeys.Error, err.Error(), commonkeys.UserID, userID)
		return nil, err
	}

	out := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		out[i] = toModelOut(rec)
	}

	span.SetAttributes(attribute.Int("count", len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	return out, nil
}
