package controller_test

import (
	"context"
	"testing"
	"time"

	gmodel "github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardWidgetCatalog_ReturnsCanonicalCatalog(t *testing.T) {
	h, ctrl := newRecordController(t, &recordServiceStub{})
	defer ctrl.Finish()

	out, err := h.DashboardWidgetCatalog(t.Context())

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, int32(domain.MaxLargeWidgetsPerDashboard), out.MaxLargeWidgets)
	assert.Equal(t, []gmodel.DashboardWidgetSize{
		gmodel.DashboardWidgetSizeSmall,
		gmodel.DashboardWidgetSizeMedium,
		gmodel.DashboardWidgetSizeLarge,
	}, out.Sizes)
	assert.Equal(t, []gmodel.DashboardWidgetType{
		gmodel.DashboardWidgetTypeKpiNumber,
		gmodel.DashboardWidgetTypeGoalProgress,
		gmodel.DashboardWidgetTypeTrendLine,
		gmodel.DashboardWidgetTypeChecklist,
		gmodel.DashboardWidgetTypeInsightFeed,
	}, out.Types)
}

func TestUpsertDashboardWidget_MapsGraphQLEnumsToUsecaseContract(t *testing.T) {
	var captured input.UpsertDashboardWidgetCommand
	title := "Radar operacional"
	config := `{"layoutVersion":2,"gridW":8,"gridH":3}`
	isActive := true
	now := time.Date(2026, 3, 18, 4, 0, 0, 0, time.UTC)

	svc := &recordServiceStub{
		upsertWidgetFn: func(_ context.Context, userID uint64, cmd input.UpsertDashboardWidgetCommand) (domain.DashboardWidget, error) {
			require.Equal(t, uint64(999), userID)
			captured = cmd
			return domain.DashboardWidget{
				ID:            44,
				UserID:        userID,
				ViewID:        10,
				WidgetType:    domain.DashboardWidgetTypeInsightFeed,
				Size:          domain.DashboardWidgetSizeLarge,
				OrderIndex:    6,
				TitleOverride: &title,
				ConfigJSON:    config,
				IsActive:      true,
				CreatedAt:     now,
				UpdatedAt:     now,
			}, nil
		},
	}

	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.UpsertDashboardWidget(t.Context(), 999, gmodel.UpsertDashboardWidgetInput{
		ViewID:        "10",
		WidgetType:    gmodel.DashboardWidgetTypeInsightFeed,
		Size:          gmodel.DashboardWidgetSizeLarge,
		OrderIndex:    ptrInt32(6),
		TitleOverride: &title,
		ConfigJSON:    &config,
		IsActive:      &isActive,
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, uint64(10), captured.ViewID)
	assert.Equal(t, domain.DashboardWidgetTypeInsightFeed, captured.WidgetType)
	assert.Equal(t, domain.DashboardWidgetSizeLarge, captured.Size)
	require.NotNil(t, captured.OrderIndex)
	assert.Equal(t, 6, *captured.OrderIndex)
	assert.Equal(t, config, captured.ConfigJSON)
	require.NotNil(t, captured.IsActive)
	assert.True(t, *captured.IsActive)
	assert.Zero(t, captured.MetricDefinitionID)
	assert.Equal(t, gmodel.DashboardWidgetTypeInsightFeed, out.WidgetType)
	assert.Equal(t, gmodel.DashboardWidgetSizeLarge, out.Size)
	require.NotNil(t, out.ConfigJSON)
	assert.Equal(t, config, *out.ConfigJSON)
}

func TestUpdateDashboardView_MapsRenameCommand(t *testing.T) {
	var capturedViewID uint64
	var captured input.UpdateDashboardViewCommand
	now := time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC)

	svc := &recordServiceStub{
		updateViewFn: func(_ context.Context, userID uint64, viewID uint64, cmd input.UpdateDashboardViewCommand) (domain.DashboardView, error) {
			require.Equal(t, uint64(999), userID)
			capturedViewID = viewID
			captured = cmd
			return domain.DashboardView{
				ID:        viewID,
				UserID:    userID,
				Name:      cmd.Name,
				IsDefault: false,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}

	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.UpdateDashboardView(t.Context(), 999, gmodel.UpdateDashboardViewInput{
		ViewID: "15",
		Name:   "Foco semanal",
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, uint64(15), capturedViewID)
	assert.Equal(t, "Foco semanal", captured.Name)
	assert.Equal(t, "Foco semanal", out.Name)
}

func TestDeleteDashboardView_MapsDeleteRequest(t *testing.T) {
	var capturedViewID uint64

	svc := &recordServiceStub{
		deleteViewFn: func(_ context.Context, userID uint64, viewID uint64) error {
			require.Equal(t, uint64(999), userID)
			capturedViewID = viewID
			return nil
		},
	}

	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	err := h.DeleteDashboardView(t.Context(), 999, "21")

	require.NoError(t, err)
	assert.Equal(t, uint64(21), capturedViewID)
}

func TestCreateMetricAndWidget_MapsMetricAndWidgetCommands(t *testing.T) {
	var captured input.CreateMetricAndWidgetCommand
	title := "Hidratacao"
	config := `{"layoutVersion":2,"gridW":5,"gridH":2}`
	now := time.Date(2026, 3, 18, 4, 30, 0, 0, time.UTC)

	svc := &recordServiceStub{
		createMetricAndWidgetFn: func(_ context.Context, userID uint64, cmd input.CreateMetricAndWidgetCommand) (domain.DashboardWidget, error) {
			require.Equal(t, uint64(999), userID)
			captured = cmd
			return domain.DashboardWidget{
				ID:                 55,
				UserID:             userID,
				ViewID:             10,
				MetricDefinitionID: 77,
				WidgetType:         domain.DashboardWidgetTypeTrendLine,
				Size:               domain.DashboardWidgetSizeMedium,
				OrderIndex:         2,
				TitleOverride:      &title,
				ConfigJSON:         config,
				IsActive:           true,
				CreatedAt:          now,
				UpdatedAt:          now,
			}, nil
		},
	}

	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.CreateMetricAndWidget(t.Context(), 999, gmodel.CreateMetricAndWidgetInput{
		Metric: &gmodel.UpsertMetricDefinitionInput{
			MetricKey:   "water.count",
			DisplayName: "Agua",
			CategoryID:  ptrString("3"),
			TagID:       "7",
			TagIds:      []string{"7", "8"},
			ValueSource: ptrString("count"),
			Aggregation: ptrString("sum"),
			Unit:        ptrString("count"),
			IsActive:    ptrBool(true),
		},
		Widget: &gmodel.UpsertDashboardWidgetInput{
			ViewID:        "10",
			WidgetType:    gmodel.DashboardWidgetTypeTrendLine,
			Size:          gmodel.DashboardWidgetSizeMedium,
			OrderIndex:    ptrInt32(2),
			TitleOverride: &title,
			ConfigJSON:    &config,
			IsActive:      ptrBool(true),
		},
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "water.count", captured.Metric.MetricKey)
	assert.Equal(t, "Agua", captured.Metric.DisplayName)
	require.NotNil(t, captured.Metric.CategoryID)
	assert.Equal(t, uint64(3), *captured.Metric.CategoryID)
	assert.Equal(t, uint64(7), captured.Metric.TagID)
	assert.Equal(t, []uint64{7, 8}, captured.Metric.TagIDs)
	assert.Equal(t, "count", captured.Metric.ValueSource)
	assert.Equal(t, "sum", captured.Metric.Aggregation)
	assert.Equal(t, "count", captured.Metric.Unit)
	assert.Equal(t, uint64(10), captured.Widget.ViewID)
	assert.Equal(t, domain.DashboardWidgetTypeTrendLine, captured.Widget.WidgetType)
	assert.Equal(t, domain.DashboardWidgetSizeMedium, captured.Widget.Size)
	assert.Equal(t, config, captured.Widget.ConfigJSON)
	assert.Equal(t, gmodel.DashboardWidgetTypeTrendLine, out.WidgetType)
	assert.Equal(t, gmodel.DashboardWidgetSizeMedium, out.Size)
}

func ptrString(v string) *string { return &v }

func ptrBool(v bool) *bool { return &v }

func ptrInt32(v int32) *int32 { return &v }
