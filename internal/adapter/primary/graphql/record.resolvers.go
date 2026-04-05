package graphql

import (
	"context"
	"strconv"

	"github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
)

// CreateRecord is the resolver for the createRecord field.
func (m *mutationResolver) CreateRecord(ctx context.Context, input model.CreateRecordInput) (*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().Create(ctx, input, uid)
}

// RecordByID is the resolver for the recordByID field.
func (q *queryResolver) RecordByID(ctx context.Context, recordID string) (*model.Record, error) {
	id, err := strconv.ParseUint(recordID, 10, 64)
	if err != nil {
		return nil, err
	}

	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().GetByID(ctx, id, uid)
}

// RecordProjectionByID is the resolver for the recordProjectionById field.
func (q *queryResolver) RecordProjectionByID(ctx context.Context, id string) (*model.RecordProjection, error) {
	recordID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().GetProjectedByID(ctx, recordID, uid)
}

// RecordsByTag is the resolver for the recordsByTag field.
func (q *queryResolver) RecordsByTag(ctx context.Context, tagID string, limit *int32, afterEventTime *string, afterID *string) ([]*model.Record, error) {
	tid, err := strconv.ParseUint(tagID, 10, 64)
	if err != nil {
		return nil, err
	}

	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)

	lim := 50 // default
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}

	var afterIDInt *int64
	if afterID != nil && *afterID != "" {
		if v, err := strconv.ParseInt(*afterID, 10, 64); err == nil {
			afterIDInt = &v
		}
	}

	return q.RecordController().ListByTag(ctx, tid, uid, lim, afterEventTime, afterIDInt)
}

// RecordsByDay is the resolver for the recordsByDay field.
func (q *queryResolver) RecordsByDay(ctx context.Context, date *string) ([]*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	dateStr := ""
	if date != nil {
		dateStr = *date
	}
	return q.RecordController().ListByDay(ctx, uid, dateStr)
}

// RecordsUntil is the resolver for the recordsUntil field.
func (q *queryResolver) RecordsUntil(ctx context.Context, until string, limit *int32) ([]*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)

	lim := 50 // default
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}

	return q.RecordController().ListAllUntil(ctx, uid, until, lim)
}

// RecordsBetween is the resolver for the recordsBetween field.
func (q *queryResolver) RecordsBetween(ctx context.Context, startDate string, endDate string, limit *int32) ([]*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)

	lim := 50 // default
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}

	return q.RecordController().ListAllBetween(ctx, uid, startDate, endDate, lim)
}

// Records is the resolver for the records field (list by user with optional cursors).
func (q *queryResolver) Records(ctx context.Context, limit *int32, afterEventTime *string, afterID *string) ([]*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	lim := 50
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}
	var afterIDInt *int64
	if afterID != nil && *afterID != "" {
		if v, err := strconv.ParseInt(*afterID, 10, 64); err == nil {
			afterIDInt = &v
		}
	}

	return q.RecordController().ListByUser(ctx, uid, lim, afterEventTime, afterIDInt)
}

// RecordProjections is the resolver for the recordProjections field.
func (q *queryResolver) RecordProjections(ctx context.Context, limit *int32, afterEventTime *string, afterID *string) ([]*model.RecordProjection, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	lim := 50
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}
	var afterIDInt *int64
	if afterID != nil && *afterID != "" {
		if v, err := strconv.ParseInt(*afterID, 10, 64); err == nil {
			afterIDInt = &v
		}
	}

	return q.RecordController().ListProjectedPage(ctx, uid, lim, afterEventTime, afterIDInt)
}

// RecordsLatest is the resolver for the recordsLatest field.
func (q *queryResolver) RecordsLatest(ctx context.Context, limit *int32) ([]*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)

	lim := 10 // default for latest
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}

	return q.RecordController().ListLatest(ctx, uid, lim)
}

// RecordProjectionsLatest is the resolver for the recordProjectionsLatest field.
func (q *queryResolver) RecordProjectionsLatest(ctx context.Context, limit *int32) ([]*model.RecordProjection, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)

	lim := 10
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}

	return q.RecordController().ListProjectedLatest(ctx, uid, lim)
}

// RecordsByCategory is the resolver for the recordsByCategory field.
func (q *queryResolver) RecordsByCategory(ctx context.Context, categoryID string, limit *int32) ([]*model.Record, error) {
	cid, err := strconv.ParseUint(categoryID, 10, 64)
	if err != nil {
		return nil, err
	}

	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)

	lim := 50 // default
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}

	return q.RecordController().ListByCategory(ctx, cid, uid, lim)
}

// UpdateRecord is the resolver for the updateRecord field.
func (m *mutationResolver) UpdateRecord(ctx context.Context, input model.UpdateRecordInput) (*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().Update(ctx, input, uid)
}

// SoftDeleteRecord is the resolver for the softDeleteRecord field.
func (m *mutationResolver) SoftDeleteRecord(ctx context.Context, input model.DeleteRecordInput) (bool, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	id, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		return false, err
	}
	if err := m.RecordController().SoftDelete(ctx, id, uid); err != nil {
		return false, err
	}
	return true, nil
}

// SoftDeleteAllRecords is the resolver for the softDeleteAllRecords field.
func (m *mutationResolver) SoftDeleteAllRecords(ctx context.Context) (bool, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	if err := m.RecordController().SoftDeleteAll(ctx, uid); err != nil {
		return false, err
	}
	return true, nil
}

// SearchRecords is the resolver for the searchRecords field.
func (q *queryResolver) SearchRecords(ctx context.Context, filters model.SearchFilters) ([]*model.Record, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().SearchRecords(ctx, filters, uid)
}

// RecordStats is the resolver for the recordStats field.
func (q *queryResolver) RecordStats(ctx context.Context, filters *model.RecordStatsFilters) (*model.RecordStats, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().RecordStats(ctx, filters, uid)
}

// DashboardSnapshot is the resolver for the dashboardSnapshot field.
func (q *queryResolver) DashboardSnapshot(ctx context.Context, date string, timezone *string) (*model.DashboardSnapshot, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().DashboardSnapshot(ctx, uid, date, timezone)
}

// InsightFeed is the resolver for the insightFeed field.
func (q *queryResolver) InsightFeed(
	ctx context.Context,
	window model.InsightWindow,
	limit *int32,
	date *string,
	timezone *string,
	categoryID *string,
	tagIDs []string,
) ([]*model.InsightCard, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().InsightFeed(ctx, uid, window, limit, date, timezone, categoryID, tagIDs)
}

// AnalyticsSeries is the resolver for the analyticsSeries field.
func (q *queryResolver) AnalyticsSeries(
	ctx context.Context,
	seriesKey string,
	window model.InsightWindow,
	date *string,
	timezone *string,
	categoryID *string,
	tagIDs []string,
) (*model.AnalyticsSeriesResult, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().AnalyticsSeries(ctx, uid, seriesKey, window, date, timezone, categoryID, tagIDs)
}

// MetricDefinitions is the resolver for the metricDefinitions field.
func (q *queryResolver) MetricDefinitions(ctx context.Context) ([]*model.MetricDefinition, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().ListMetricDefinitions(ctx, uid)
}

// UpsertMetricDefinition is the resolver for the upsertMetricDefinition field.
func (m *mutationResolver) UpsertMetricDefinition(ctx context.Context, input model.UpsertMetricDefinitionInput) (*model.MetricDefinition, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().UpsertMetricDefinition(ctx, uid, input)
}

// UpsertGoalTemplate is the resolver for the upsertGoalTemplate field.
func (m *mutationResolver) UpsertGoalTemplate(ctx context.Context, input model.UpsertGoalTemplateInput) (*model.GoalTemplate, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().UpsertGoalTemplate(ctx, uid, input)
}

// DeleteGoalTemplate is the resolver for the deleteGoalTemplate field.
func (m *mutationResolver) DeleteGoalTemplate(ctx context.Context, input model.DeleteGoalTemplateInput) (bool, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	id, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		return false, err
	}
	if err := m.RecordController().DeleteGoalTemplate(ctx, uid, id); err != nil {
		return false, err
	}
	return true, nil
}

// DashboardViews is the resolver for the dashboardViews field.
func (q *queryResolver) DashboardViews(ctx context.Context) ([]*model.DashboardView, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().ListDashboardViews(ctx, uid)
}

// DashboardView is the resolver for the dashboardView field.
func (q *queryResolver) DashboardView(ctx context.Context, id string) (*model.DashboardView, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().GetDashboardView(ctx, uid, id)
}

// DashboardWidgetCatalog is the resolver for the dashboardWidgetCatalog field.
func (q *queryResolver) DashboardWidgetCatalog(ctx context.Context) (*model.DashboardWidgetCatalog, error) {
	return q.RecordController().DashboardWidgetCatalog(ctx)
}

// SuggestMetricDefinitions is the resolver for the suggestMetricDefinitions field.
func (q *queryResolver) SuggestMetricDefinitions(ctx context.Context, limit *int32) ([]*model.MetricDefinitionSuggestion, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.RecordController().SuggestMetricDefinitions(ctx, uid, limit)
}

// CreateDashboardView is the resolver for the createDashboardView field.
func (m *mutationResolver) CreateDashboardView(ctx context.Context, input model.CreateDashboardViewInput) (*model.DashboardView, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().CreateDashboardView(ctx, uid, input)
}

// UpdateDashboardView is the resolver for the updateDashboardView field.
func (m *mutationResolver) UpdateDashboardView(ctx context.Context, input model.UpdateDashboardViewInput) (*model.DashboardView, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().UpdateDashboardView(ctx, uid, input)
}

// SetDefaultDashboardView is the resolver for the setDefaultDashboardView field.
func (m *mutationResolver) SetDefaultDashboardView(ctx context.Context, input model.SetDefaultDashboardViewInput) (*model.DashboardView, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().SetDefaultDashboardView(ctx, uid, input.ViewID)
}

// DeleteDashboardView is the resolver for the deleteDashboardView field.
func (m *mutationResolver) DeleteDashboardView(ctx context.Context, input model.DeleteDashboardViewInput) (bool, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	if err := m.RecordController().DeleteDashboardView(ctx, uid, input.ID); err != nil {
		return false, err
	}
	return true, nil
}

// UpsertDashboardWidget is the resolver for the upsertDashboardWidget field.
func (m *mutationResolver) UpsertDashboardWidget(ctx context.Context, input model.UpsertDashboardWidgetInput) (*model.DashboardWidget, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().UpsertDashboardWidget(ctx, uid, input)
}

// ReorderDashboardWidgets is the resolver for the reorderDashboardWidgets field.
func (m *mutationResolver) ReorderDashboardWidgets(ctx context.Context, input model.ReorderDashboardWidgetsInput) ([]*model.DashboardWidget, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().ReorderDashboardWidgets(ctx, uid, input)
}

// DeleteDashboardWidget is the resolver for the deleteDashboardWidget field.
func (m *mutationResolver) DeleteDashboardWidget(ctx context.Context, input model.DeleteDashboardWidgetInput) (bool, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	if err := m.RecordController().DeleteDashboardWidget(ctx, uid, input.ID); err != nil {
		return false, err
	}
	return true, nil
}

// CreateMetricAndWidget is the resolver for the createMetricAndWidget field.
func (m *mutationResolver) CreateMetricAndWidget(ctx context.Context, input model.CreateMetricAndWidgetInput) (*model.DashboardWidget, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.RecordController().CreateMetricAndWidget(ctx, uid, input)
}

// Additional unimplemented resolvers can be added below as needed.
