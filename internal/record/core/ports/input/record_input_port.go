package input

import (
	"context"
	"time"

	"github.com/lechitz/aion-api/internal/record/core/domain"
)

// RecordCreator interface for creating a new record.
type RecordCreator interface {
	Create(ctx context.Context, cmd CreateRecordCommand) (domain.Record, error)
}

// RecordRetriever defines methods for retrieving record details.
type RecordRetriever interface {
	GetByID(ctx context.Context, recordID uint64, userID uint64) (domain.Record, error)
	ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error)
	ListByTag(ctx context.Context, tagID uint64, userID uint64, limit int) ([]domain.Record, error)
	ListByCategory(ctx context.Context, categoryID uint64, userID uint64, limit int) ([]domain.Record, error)
	ListByDay(ctx context.Context, userID uint64, date time.Time) ([]domain.Record, error)
	ListAllUntil(ctx context.Context, userID uint64, until time.Time, limit int) ([]domain.Record, error)
	ListAllBetween(ctx context.Context, userID uint64, startDate time.Time, endDate time.Time, limit int) ([]domain.Record, error)
	ListLatest(ctx context.Context, userID uint64, limit int) ([]domain.Record, error)
}

// RecordUpdater defines update operations for a record.
type RecordUpdater interface {
	Update(ctx context.Context, recordID uint64, userID uint64, cmd UpdateRecordCommand) (domain.Record, error)
}

// RecordDeleter defines deletion operations for records.
type RecordDeleter interface {
	Delete(ctx context.Context, recordID uint64, userID uint64) error
	DeleteAll(ctx context.Context, userID uint64) error
}

// RecordService defines the input port used by controllers/handlers to interact with record use cases.
type RecordService interface {
	RecordCreator
	RecordRetriever
	RecordProjectionRetriever
	RecordUpdater
	RecordDeleter

	// SearchRecords performs full-text search with filters
	SearchRecords(ctx context.Context, userID uint64, filters domain.SearchFilters) ([]domain.Record, error)
	// DashboardSnapshot computes deterministic metrics and goals for a specific day.
	DashboardSnapshot(ctx context.Context, userID uint64, query DashboardSnapshotQuery) (domain.DashboardSnapshot, error)
	// InsightFeed returns explainable deterministic insights for one analysis window.
	InsightFeed(ctx context.Context, userID uint64, query InsightFeedQuery) ([]domain.InsightCard, error)
	// AnalyticsSeries returns a compact series for one supported key/window pair.
	AnalyticsSeries(ctx context.Context, userID uint64, query AnalyticsSeriesQuery) (domain.AnalyticsSeriesResult, error)
	// ListMetricDefinitions retrieves active dashboard metric definitions.
	ListMetricDefinitions(ctx context.Context, userID uint64) ([]domain.MetricDefinition, error)
	// UpsertMetricDefinition creates/updates a metric definition.
	UpsertMetricDefinition(ctx context.Context, userID uint64, cmd UpsertMetricDefinitionCommand) (domain.MetricDefinition, error)
	// UpsertGoalTemplate creates/updates a goal template.
	UpsertGoalTemplate(ctx context.Context, userID uint64, cmd UpsertGoalTemplateCommand) (domain.GoalTemplate, error)
	// DeleteGoalTemplate soft deletes/removes a goal template.
	DeleteGoalTemplate(ctx context.Context, userID uint64, goalTemplateID uint64) error
	// Dashboard views/widgets (white-label dashboard layout)
	ListDashboardViews(ctx context.Context, userID uint64) ([]domain.DashboardView, error)
	GetDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error)
	CreateDashboardView(ctx context.Context, userID uint64, cmd CreateDashboardViewCommand) (domain.DashboardView, error)
	UpdateDashboardView(ctx context.Context, userID uint64, viewID uint64, cmd UpdateDashboardViewCommand) (domain.DashboardView, error)
	SetDefaultDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error)
	DeleteDashboardView(ctx context.Context, userID uint64, viewID uint64) error
	UpsertDashboardWidget(ctx context.Context, userID uint64, cmd UpsertDashboardWidgetCommand) (domain.DashboardWidget, error)
	ReorderDashboardWidgets(ctx context.Context, userID uint64, cmd ReorderDashboardWidgetsCommand) ([]domain.DashboardWidget, error)
	DeleteDashboardWidget(ctx context.Context, userID uint64, widgetID uint64) error
	CreateMetricAndWidget(ctx context.Context, userID uint64, cmd CreateMetricAndWidgetCommand) (domain.DashboardWidget, error)
	SuggestMetricDefinitions(ctx context.Context, userID uint64, limit int) ([]domain.MetricDefinitionSuggestion, error)
}
