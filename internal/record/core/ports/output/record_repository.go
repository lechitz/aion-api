// Package output defines persistence output ports for the record bounded context.
package output

import (
	"context"
	"time"

	"github.com/lechitz/aion-api/internal/record/core/domain"
)

// RecordRepository defines persistence operations for records.
type RecordRepository interface {
	Create(ctx context.Context, r domain.Record) (domain.Record, error)
	Update(ctx context.Context, r domain.Record) (domain.Record, error)
	GetByID(ctx context.Context, recordID uint64, userID uint64) (domain.Record, error)
	GetByUserCategoryDate(ctx context.Context, userID uint64, categoryID uint64, date time.Time) (domain.Record, error)

	ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error)
	ListByCategory(ctx context.Context, categoryID uint64, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error)
	ListByTag(ctx context.Context, tagID uint64, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error)
	CountByTag(ctx context.Context, tagID uint64, userID uint64) (int64, error)
	ListByDay(ctx context.Context, userID uint64, date time.Time) ([]domain.Record, error)
	ListLatest(ctx context.Context, userID uint64, limit int) ([]domain.Record, error)
	ListAllUntil(ctx context.Context, userID uint64, until time.Time, limit int) ([]domain.Record, error)
	ListAllBetween(ctx context.Context, userID uint64, startDate time.Time, endDate time.Time, limit int) ([]domain.Record, error)

	Delete(ctx context.Context, id uint64, userID uint64) error
	DeleteAllByUser(ctx context.Context, userID uint64) error

	// SearchRecords performs text search with filters
	SearchRecords(ctx context.Context, userID uint64, filters domain.SearchFilters) ([]domain.Record, error)

	// Dashboard semantic configuration
	ListMetricDefinitions(ctx context.Context, userID uint64) ([]domain.MetricDefinition, error)
	UpsertMetricDefinition(ctx context.Context, definition domain.MetricDefinition) (domain.MetricDefinition, error)

	ListGoalTemplates(ctx context.Context, userID uint64) ([]domain.GoalTemplate, error)
	UpsertGoalTemplate(ctx context.Context, template domain.GoalTemplate) (domain.GoalTemplate, error)
	DeleteGoalTemplate(ctx context.Context, userID uint64, goalTemplateID uint64) error

	// White-label dashboard layout persistence
	ListDashboardViews(ctx context.Context, userID uint64) ([]domain.DashboardView, error)
	GetDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error)
	CreateDashboardView(ctx context.Context, view domain.DashboardView) (domain.DashboardView, error)
	UpdateDashboardView(ctx context.Context, userID uint64, viewID uint64, name string) (domain.DashboardView, error)
	SetDefaultDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error)
	DeleteDashboardView(ctx context.Context, userID uint64, viewID uint64) error
	UpsertDashboardWidget(ctx context.Context, widget domain.DashboardWidget) (domain.DashboardWidget, error)
	ListDashboardWidgetsByView(ctx context.Context, userID uint64, viewID uint64) ([]domain.DashboardWidget, error)
	ReorderDashboardWidgets(ctx context.Context, userID uint64, viewID uint64, items []domain.DashboardWidget) ([]domain.DashboardWidget, error)
	DeleteDashboardWidget(ctx context.Context, userID uint64, widgetID uint64) error
	CountLargeWidgetsInView(ctx context.Context, userID uint64, viewID uint64, excludeWidgetID *uint64) (int64, error)
}
