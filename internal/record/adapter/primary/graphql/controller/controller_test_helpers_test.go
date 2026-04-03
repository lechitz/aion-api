package controller_test

import (
	"context"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/record/adapter/primary/graphql/controller"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
	"github.com/lechitz/aion-api/tests/mocks"
	"github.com/lechitz/aion-api/tests/setup"
	"go.uber.org/mock/gomock"
)

type recordServiceStub struct {
	createFn                func(context.Context, input.CreateRecordCommand) (domain.Record, error)
	getByIDFn               func(context.Context, uint64, uint64) (domain.Record, error)
	getProjectedByIDFn      func(context.Context, uint64, uint64) (domain.RecordProjection, error)
	listByUserFn            func(context.Context, uint64, int, *string, *int64) ([]domain.Record, error)
	listProjectedPageFn     func(context.Context, uint64, int, *string, *int64) ([]domain.RecordProjection, error)
	listByTagFn             func(context.Context, uint64, uint64, int) ([]domain.Record, error)
	listByCatFn             func(context.Context, uint64, uint64, int) ([]domain.Record, error)
	listByDayFn             func(context.Context, uint64, time.Time) ([]domain.Record, error)
	listAllUntilFn          func(context.Context, uint64, time.Time, int) ([]domain.Record, error)
	listAllBetweenFn        func(context.Context, uint64, time.Time, time.Time, int) ([]domain.Record, error)
	listLatestFn            func(context.Context, uint64, int) ([]domain.Record, error)
	listProjectedLatestFn   func(context.Context, uint64, int) ([]domain.RecordProjection, error)
	updateFn                func(context.Context, uint64, uint64, input.UpdateRecordCommand) (domain.Record, error)
	deleteFn                func(context.Context, uint64, uint64) error
	deleteAllFn             func(context.Context, uint64) error
	searchFn                func(context.Context, uint64, domain.SearchFilters) ([]domain.Record, error)
	dashboardFn             func(context.Context, uint64, input.DashboardSnapshotQuery) (domain.DashboardSnapshot, error)
	insightFeedFn           func(context.Context, uint64, input.InsightFeedQuery) ([]domain.InsightCard, error)
	analyticsSeriesFn       func(context.Context, uint64, input.AnalyticsSeriesQuery) (domain.AnalyticsSeriesResult, error)
	listMetricFn            func(context.Context, uint64) ([]domain.MetricDefinition, error)
	upsertMetricFn          func(context.Context, uint64, input.UpsertMetricDefinitionCommand) (domain.MetricDefinition, error)
	upsertGoalFn            func(context.Context, uint64, input.UpsertGoalTemplateCommand) (domain.GoalTemplate, error)
	deleteGoalFn            func(context.Context, uint64, uint64) error
	listViewsFn             func(context.Context, uint64) ([]domain.DashboardView, error)
	getViewFn               func(context.Context, uint64, uint64) (domain.DashboardView, error)
	createViewFn            func(context.Context, uint64, input.CreateDashboardViewCommand) (domain.DashboardView, error)
	updateViewFn            func(context.Context, uint64, uint64, input.UpdateDashboardViewCommand) (domain.DashboardView, error)
	setDefaultViewFn        func(context.Context, uint64, uint64) (domain.DashboardView, error)
	deleteViewFn            func(context.Context, uint64, uint64) error
	upsertWidgetFn          func(context.Context, uint64, input.UpsertDashboardWidgetCommand) (domain.DashboardWidget, error)
	reorderWidgetFn         func(context.Context, uint64, input.ReorderDashboardWidgetsCommand) ([]domain.DashboardWidget, error)
	deleteWidgetFn          func(context.Context, uint64, uint64) error
	createMetricAndWidgetFn func(context.Context, uint64, input.CreateMetricAndWidgetCommand) (domain.DashboardWidget, error)
	suggestMetricFn         func(context.Context, uint64, int) ([]domain.MetricDefinitionSuggestion, error)
}

func newRecordController(t *testing.T, svc input.RecordService) (controller.RecordController, *gomock.Controller) {
	t.Helper()
	ctrl := gomock.NewController(t)
	log := mocks.NewMockContextLogger(ctrl)
	setup.ExpectLoggerDefaultBehavior(log)
	return controller.NewController(svc, log), ctrl
}

func (s *recordServiceStub) Create(ctx context.Context, cmd input.CreateRecordCommand) (domain.Record, error) {
	if s.createFn == nil {
		panic("unexpected Create call")
	}
	return s.createFn(ctx, cmd)
}

func (s *recordServiceStub) GetByID(ctx context.Context, recordID uint64, userID uint64) (domain.Record, error) {
	if s.getByIDFn == nil {
		panic("unexpected GetByID call")
	}
	return s.getByIDFn(ctx, recordID, userID)
}

func (s *recordServiceStub) GetProjectedByID(ctx context.Context, recordID uint64, userID uint64) (domain.RecordProjection, error) {
	if s.getProjectedByIDFn == nil {
		panic("unexpected GetProjectedByID call")
	}
	return s.getProjectedByIDFn(ctx, recordID, userID)
}

func (s *recordServiceStub) ListByUser(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.Record, error) {
	if s.listByUserFn == nil {
		panic("unexpected ListByUser call")
	}
	return s.listByUserFn(ctx, userID, limit, afterEventTime, afterID)
}

func (s *recordServiceStub) ListProjectedPage(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.RecordProjection, error) {
	if s.listProjectedPageFn == nil {
		panic("unexpected ListProjectedPage call")
	}
	return s.listProjectedPageFn(ctx, userID, limit, afterEventTime, afterID)
}

func (s *recordServiceStub) ListByTag(ctx context.Context, tagID uint64, userID uint64, limit int) ([]domain.Record, error) {
	if s.listByTagFn == nil {
		panic("unexpected ListByTag call")
	}
	return s.listByTagFn(ctx, tagID, userID, limit)
}

func (s *recordServiceStub) ListByCategory(ctx context.Context, categoryID uint64, userID uint64, limit int) ([]domain.Record, error) {
	if s.listByCatFn == nil {
		panic("unexpected ListByCategory call")
	}
	return s.listByCatFn(ctx, categoryID, userID, limit)
}

func (s *recordServiceStub) ListByDay(ctx context.Context, userID uint64, date time.Time) ([]domain.Record, error) {
	if s.listByDayFn == nil {
		panic("unexpected ListByDay call")
	}
	return s.listByDayFn(ctx, userID, date)
}

func (s *recordServiceStub) ListAllUntil(ctx context.Context, userID uint64, until time.Time, limit int) ([]domain.Record, error) {
	if s.listAllUntilFn == nil {
		panic("unexpected ListAllUntil call")
	}
	return s.listAllUntilFn(ctx, userID, until, limit)
}

func (s *recordServiceStub) ListAllBetween(ctx context.Context, userID uint64, startDate time.Time, endDate time.Time, limit int) ([]domain.Record, error) {
	if s.listAllBetweenFn == nil {
		panic("unexpected ListAllBetween call")
	}
	return s.listAllBetweenFn(ctx, userID, startDate, endDate, limit)
}

func (s *recordServiceStub) ListLatest(ctx context.Context, userID uint64, limit int) ([]domain.Record, error) {
	if s.listLatestFn == nil {
		panic("unexpected ListLatest call")
	}
	return s.listLatestFn(ctx, userID, limit)
}

func (s *recordServiceStub) ListProjectedLatest(ctx context.Context, userID uint64, limit int) ([]domain.RecordProjection, error) {
	if s.listProjectedLatestFn == nil {
		panic("unexpected ListProjectedLatest call")
	}
	return s.listProjectedLatestFn(ctx, userID, limit)
}

func (s *recordServiceStub) Update(ctx context.Context, recordID uint64, userID uint64, cmd input.UpdateRecordCommand) (domain.Record, error) {
	if s.updateFn == nil {
		panic("unexpected Update call")
	}
	return s.updateFn(ctx, recordID, userID, cmd)
}

func (s *recordServiceStub) Delete(ctx context.Context, recordID uint64, userID uint64) error {
	if s.deleteFn == nil {
		panic("unexpected Delete call")
	}
	return s.deleteFn(ctx, recordID, userID)
}

func (s *recordServiceStub) DeleteAll(ctx context.Context, userID uint64) error {
	if s.deleteAllFn == nil {
		panic("unexpected DeleteAll call")
	}
	return s.deleteAllFn(ctx, userID)
}

func (s *recordServiceStub) SearchRecords(ctx context.Context, userID uint64, filters domain.SearchFilters) ([]domain.Record, error) {
	if s.searchFn == nil {
		panic("unexpected SearchRecords call")
	}
	return s.searchFn(ctx, userID, filters)
}

func (s *recordServiceStub) DashboardSnapshot(ctx context.Context, userID uint64, query input.DashboardSnapshotQuery) (domain.DashboardSnapshot, error) {
	if s.dashboardFn == nil {
		panic("unexpected DashboardSnapshot call")
	}
	return s.dashboardFn(ctx, userID, query)
}

func (s *recordServiceStub) InsightFeed(ctx context.Context, userID uint64, query input.InsightFeedQuery) ([]domain.InsightCard, error) {
	if s.insightFeedFn == nil {
		panic("unexpected InsightFeed call")
	}
	return s.insightFeedFn(ctx, userID, query)
}

func (s *recordServiceStub) AnalyticsSeries(ctx context.Context, userID uint64, query input.AnalyticsSeriesQuery) (domain.AnalyticsSeriesResult, error) {
	if s.analyticsSeriesFn == nil {
		panic("unexpected AnalyticsSeries call")
	}
	return s.analyticsSeriesFn(ctx, userID, query)
}

func (s *recordServiceStub) ListMetricDefinitions(ctx context.Context, userID uint64) ([]domain.MetricDefinition, error) {
	if s.listMetricFn == nil {
		panic("unexpected ListMetricDefinitions call")
	}
	return s.listMetricFn(ctx, userID)
}

func (s *recordServiceStub) UpsertMetricDefinition(ctx context.Context, userID uint64, cmd input.UpsertMetricDefinitionCommand) (domain.MetricDefinition, error) {
	if s.upsertMetricFn == nil {
		panic("unexpected UpsertMetricDefinition call")
	}
	return s.upsertMetricFn(ctx, userID, cmd)
}

func (s *recordServiceStub) UpsertGoalTemplate(ctx context.Context, userID uint64, cmd input.UpsertGoalTemplateCommand) (domain.GoalTemplate, error) {
	if s.upsertGoalFn == nil {
		panic("unexpected UpsertGoalTemplate call")
	}
	return s.upsertGoalFn(ctx, userID, cmd)
}

func (s *recordServiceStub) DeleteGoalTemplate(ctx context.Context, userID uint64, goalTemplateID uint64) error {
	if s.deleteGoalFn == nil {
		panic("unexpected DeleteGoalTemplate call")
	}
	return s.deleteGoalFn(ctx, userID, goalTemplateID)
}

func (s *recordServiceStub) ListDashboardViews(ctx context.Context, userID uint64) ([]domain.DashboardView, error) {
	if s.listViewsFn == nil {
		panic("unexpected ListDashboardViews call")
	}
	return s.listViewsFn(ctx, userID)
}

func (s *recordServiceStub) GetDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error) {
	if s.getViewFn == nil {
		panic("unexpected GetDashboardView call")
	}
	return s.getViewFn(ctx, userID, viewID)
}

func (s *recordServiceStub) CreateDashboardView(ctx context.Context, userID uint64, cmd input.CreateDashboardViewCommand) (domain.DashboardView, error) {
	if s.createViewFn == nil {
		panic("unexpected CreateDashboardView call")
	}
	return s.createViewFn(ctx, userID, cmd)
}

func (s *recordServiceStub) UpdateDashboardView(ctx context.Context, userID uint64, viewID uint64, cmd input.UpdateDashboardViewCommand) (domain.DashboardView, error) {
	if s.updateViewFn == nil {
		panic("unexpected UpdateDashboardView call")
	}
	return s.updateViewFn(ctx, userID, viewID, cmd)
}

func (s *recordServiceStub) SetDefaultDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error) {
	if s.setDefaultViewFn == nil {
		panic("unexpected SetDefaultDashboardView call")
	}
	return s.setDefaultViewFn(ctx, userID, viewID)
}

func (s *recordServiceStub) DeleteDashboardView(ctx context.Context, userID uint64, viewID uint64) error {
	if s.deleteViewFn == nil {
		panic("unexpected DeleteDashboardView call")
	}
	return s.deleteViewFn(ctx, userID, viewID)
}

func (s *recordServiceStub) UpsertDashboardWidget(ctx context.Context, userID uint64, cmd input.UpsertDashboardWidgetCommand) (domain.DashboardWidget, error) {
	if s.upsertWidgetFn == nil {
		panic("unexpected UpsertDashboardWidget call")
	}
	return s.upsertWidgetFn(ctx, userID, cmd)
}

func (s *recordServiceStub) ReorderDashboardWidgets(ctx context.Context, userID uint64, cmd input.ReorderDashboardWidgetsCommand) ([]domain.DashboardWidget, error) {
	if s.reorderWidgetFn == nil {
		panic("unexpected ReorderDashboardWidgets call")
	}
	return s.reorderWidgetFn(ctx, userID, cmd)
}

func (s *recordServiceStub) DeleteDashboardWidget(ctx context.Context, userID uint64, widgetID uint64) error {
	if s.deleteWidgetFn == nil {
		panic("unexpected DeleteDashboardWidget call")
	}
	return s.deleteWidgetFn(ctx, userID, widgetID)
}

func (s *recordServiceStub) CreateMetricAndWidget(ctx context.Context, userID uint64, cmd input.CreateMetricAndWidgetCommand) (domain.DashboardWidget, error) {
	if s.createMetricAndWidgetFn == nil {
		panic("unexpected CreateMetricAndWidget call")
	}
	return s.createMetricAndWidgetFn(ctx, userID, cmd)
}

func (s *recordServiceStub) SuggestMetricDefinitions(ctx context.Context, userID uint64, limit int) ([]domain.MetricDefinitionSuggestion, error) {
	if s.suggestMetricFn == nil {
		panic("unexpected SuggestMetricDefinitions call")
	}
	return s.suggestMetricFn(ctx, userID, limit)
}
