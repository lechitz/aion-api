package graphql

import (
	"context"
	"math"
	"strconv"

	"github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
)

func hydrateTagRecordCounts(_ context.Context, tags []*model.Tag, countFn func(tagID uint64) (int64, error)) error {
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		tagID, err := strconv.ParseUint(tag.ID, 10, 64)
		if err != nil {
			return err
		}
		count, err := countFn(tagID)
		if err != nil {
			return err
		}
		tag.RecordCount = safeRecordCount(count)
	}
	return nil
}

func safeRecordCount(count int64) int32 {
	if count >= math.MaxInt32 {
		return math.MaxInt32
	}
	return int32(count) // #nosec G115: safe due to preceding clamp
}

// CreateTag is the resolver for the createTag field.
func (m *mutationResolver) CreateTag(ctx context.Context, input model.CreateTagInput) (*model.Tag, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.TagController().Create(ctx, input, uid)
}

// UpdateTag is the resolver for the updateTag field.
func (m *mutationResolver) UpdateTag(ctx context.Context, input model.UpdateTagInput) (*model.Tag, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.TagController().Update(ctx, input, uid)
}

// SoftDeleteTag is the resolver for the softDeleteTag field.
func (m *mutationResolver) SoftDeleteTag(ctx context.Context, input model.DeleteTagInput) (bool, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	tagID, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		return false, err
	}
	if err := m.TagController().SoftDelete(ctx, tagID, uid); err != nil {
		return false, err
	}
	return true, nil
}

// TagByName is the resolve for the tagByName field.
func (q *queryResolver) TagByName(ctx context.Context, tagName string) (*model.Tag, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	tag, err := q.TagController().GetByName(ctx, tagName, uid)
	if err != nil || tag == nil {
		return tag, err
	}
	tagID, parseErr := strconv.ParseUint(tag.ID, 10, 64)
	if parseErr != nil {
		return nil, parseErr
	}
	count, countErr := q.RecordController().CountByTag(ctx, tagID, uid)
	if countErr != nil {
		return nil, countErr
	}
	tag.RecordCount = safeRecordCount(count)
	return tag, nil
}

// TagByID is the resolve for the tagByID field.
func (q *queryResolver) TagByID(ctx context.Context, tagID string) (*model.Tag, error) {
	id, err := strconv.ParseUint(tagID, 10, 64)
	if err != nil {
		return nil, err
	}
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	tag, getErr := q.TagController().GetByID(ctx, id, uid)
	if getErr != nil || tag == nil {
		return tag, getErr
	}
	count, countErr := q.RecordController().CountByTag(ctx, id, uid)
	if countErr != nil {
		return nil, countErr
	}
	tag.RecordCount = safeRecordCount(count)
	return tag, nil
}

// TagsByCategoryID is the resolver for the tagsByCategoryId field.
func (q *queryResolver) TagsByCategoryID(ctx context.Context, categoryID string) ([]*model.Tag, error) {
	id, err := strconv.ParseUint(categoryID, 10, 64)
	if err != nil {
		return nil, err
	}
	userID, _ := ctx.Value(ctxkeys.UserID).(uint64)
	tags, getErr := q.TagController().GetByCategoryID(ctx, id, userID)
	if getErr != nil {
		return nil, getErr
	}
	if err := hydrateTagRecordCounts(ctx, tags, func(tagID uint64) (int64, error) {
		return q.RecordController().CountByTag(ctx, tagID, userID)
	}); err != nil {
		return nil, err
	}
	return tags, nil
}

// Tags is the resolver for the tags field (list all tags for user).
func (q *queryResolver) Tags(ctx context.Context) ([]*model.Tag, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	tags, err := q.TagController().GetAll(ctx, uid)
	if err != nil {
		return nil, err
	}
	if err := hydrateTagRecordCounts(ctx, tags, func(tagID uint64) (int64, error) {
		return q.RecordController().CountByTag(ctx, tagID, uid)
	}); err != nil {
		return nil, err
	}
	return tags, nil
}
