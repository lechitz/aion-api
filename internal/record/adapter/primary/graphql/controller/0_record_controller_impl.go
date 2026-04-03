package controller

import (
	"context"

	"github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
)

// RecordController is the contract used by GraphQL resolvers.
type RecordController interface {
	Create(ctx context.Context, in model.CreateRecordInput, userID uint64) (*model.Record, error)
	GetByID(ctx context.Context, recordID, userID uint64) (*model.Record, error)
	GetProjectedByID(ctx context.Context, recordID, userID uint64) (*model.RecordProjection, error)
	ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]*model.Record, error)
	ListProjectedPage(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]*model.RecordProjection, error)
	ListLatest(ctx context.Context, userID uint64, limit int) ([]*model.Record, error)
	ListProjectedLatest(ctx context.Context, userID uint64, limit int) ([]*model.RecordProjection, error)
	ListByTag(ctx context.Context, tagID, userID uint64, limit int) ([]*model.Record, error)
	ListByCategory(ctx context.Context, categoryID, userID uint64, limit int) ([]*model.Record, error)
	ListByDay(ctx context.Context, userID uint64, date string) ([]*model.Record, error)
	ListAllUntil(ctx context.Context, userID uint64, until string, limit int) ([]*model.Record, error)
	ListAllBetween(ctx context.Context, userID uint64, startDate, endDate string, limit int) ([]*model.Record, error)
	DashboardSnapshot(ctx context.Context, userID uint64, date string, timezone *string) (*model.DashboardSnapshot, error)
	InsightFeed(
		ctx context.Context,
		userID uint64,
		window model.InsightWindow,
		limit *int32,
		date *string,
		timezone *string,
		categoryID *string,
		tagIDs []string,
	) ([]*model.InsightCard, error)
	AnalyticsSeries(
		ctx context.Context,
		userID uint64,
		seriesKey string,
		window model.InsightWindow,
		date *string,
		timezone *string,
		categoryID *string,
		tagIDs []string,
	) (*model.AnalyticsSeriesResult, error)
	ListMetricDefinitions(ctx context.Context, userID uint64) ([]*model.MetricDefinition, error)
	Update(ctx context.Context, in model.UpdateRecordInput, userID uint64) (*model.Record, error)
	SoftDelete(ctx context.Context, recordID, userID uint64) error
	SoftDeleteAll(ctx context.Context, userID uint64) error
	SearchRecords(ctx context.Context, filters model.SearchFilters, userID uint64) ([]*model.Record, error)
	RecordStats(ctx context.Context, filters *model.RecordStatsFilters, userID uint64) (*model.RecordStats, error)
	UpsertMetricDefinition(ctx context.Context, userID uint64, in model.UpsertMetricDefinitionInput) (*model.MetricDefinition, error)
	UpsertGoalTemplate(ctx context.Context, userID uint64, in model.UpsertGoalTemplateInput) (*model.GoalTemplate, error)
	DeleteGoalTemplate(ctx context.Context, userID uint64, id uint64) error
	ListDashboardViews(ctx context.Context, userID uint64) ([]*model.DashboardView, error)
	GetDashboardView(ctx context.Context, userID uint64, viewID string) (*model.DashboardView, error)
	CreateDashboardView(ctx context.Context, userID uint64, in model.CreateDashboardViewInput) (*model.DashboardView, error)
	UpdateDashboardView(ctx context.Context, userID uint64, in model.UpdateDashboardViewInput) (*model.DashboardView, error)
	SetDefaultDashboardView(ctx context.Context, userID uint64, viewID string) (*model.DashboardView, error)
	DeleteDashboardView(ctx context.Context, userID uint64, viewID string) error
	UpsertDashboardWidget(ctx context.Context, userID uint64, in model.UpsertDashboardWidgetInput) (*model.DashboardWidget, error)
	ReorderDashboardWidgets(ctx context.Context, userID uint64, in model.ReorderDashboardWidgetsInput) ([]*model.DashboardWidget, error)
	DeleteDashboardWidget(ctx context.Context, userID uint64, widgetID string) error
	CreateMetricAndWidget(ctx context.Context, userID uint64, in model.CreateMetricAndWidgetInput) (*model.DashboardWidget, error)
	DashboardWidgetCatalog(ctx context.Context) (*model.DashboardWidgetCatalog, error)
	SuggestMetricDefinitions(ctx context.Context, userID uint64, limit *int32) ([]*model.MetricDefinitionSuggestion, error)
}

// controller is the controller for the record service.
type controller struct {
	RecordService input.RecordService
	Logger        logger.ContextLogger
}

// NewController wires dependencies and returns a Controller.
func NewController(svc input.RecordService, logger logger.ContextLogger) RecordController {
	return &controller{
		RecordService: svc,
		Logger:        logger,
	}
}
