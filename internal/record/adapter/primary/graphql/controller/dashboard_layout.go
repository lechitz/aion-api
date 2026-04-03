package controller

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
)

func (c *controller) ListDashboardViews(ctx context.Context, userID uint64) ([]*model.DashboardView, error) {
	out, err := c.RecordService.ListDashboardViews(ctx, userID)
	if err != nil {
		return nil, err
	}

	items := make([]*model.DashboardView, 0, len(out))
	for _, view := range out {
		detailed, err := c.RecordService.GetDashboardView(ctx, userID, view.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, toGraphQLDashboardView(detailed, detailed.Widgets))
	}
	return items, nil
}

func (c *controller) GetDashboardView(ctx context.Context, userID uint64, viewID string) (*model.DashboardView, error) {
	id := mustParseID(viewID)
	out, err := c.RecordService.GetDashboardView(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	return toGraphQLDashboardView(out, out.Widgets), nil
}

func (c *controller) CreateDashboardView(ctx context.Context, userID uint64, in model.CreateDashboardViewInput) (*model.DashboardView, error) {
	out, err := c.RecordService.CreateDashboardView(ctx, userID, input.CreateDashboardViewCommand{
		Name:      in.Name,
		IsDefault: in.IsDefault,
	})
	if err != nil {
		return nil, err
	}
	return toGraphQLDashboardView(out, nil), nil
}

func (c *controller) UpdateDashboardView(ctx context.Context, userID uint64, in model.UpdateDashboardViewInput) (*model.DashboardView, error) {
	out, err := c.RecordService.UpdateDashboardView(ctx, userID, mustParseID(in.ViewID), input.UpdateDashboardViewCommand{
		Name: in.Name,
	})
	if err != nil {
		return nil, err
	}
	return toGraphQLDashboardView(out, nil), nil
}

func (c *controller) SetDefaultDashboardView(ctx context.Context, userID uint64, viewID string) (*model.DashboardView, error) {
	out, err := c.RecordService.SetDefaultDashboardView(ctx, userID, mustParseID(viewID))
	if err != nil {
		return nil, err
	}
	return toGraphQLDashboardView(out, nil), nil
}

func (c *controller) DeleteDashboardView(ctx context.Context, userID uint64, viewID string) error {
	return c.RecordService.DeleteDashboardView(ctx, userID, mustParseID(viewID))
}

func (c *controller) UpsertDashboardWidget(ctx context.Context, userID uint64, in model.UpsertDashboardWidgetInput) (*model.DashboardWidget, error) {
	cmd := input.UpsertDashboardWidgetCommand{
		ViewID:        mustParseID(in.ViewID),
		WidgetType:    normalizeGQLEnum(string(in.WidgetType)),
		Size:          normalizeGQLEnum(string(in.Size)),
		OrderIndex:    toIntPtr(in.OrderIndex),
		TitleOverride: in.TitleOverride,
		ConfigJSON:    toStringValue(in.ConfigJSON),
		IsActive:      in.IsActive,
	}
	if in.MetricDefinitionID != nil {
		cmd.MetricDefinitionID = mustParseID(*in.MetricDefinitionID)
	}
	if in.ID != nil {
		id := mustParseID(*in.ID)
		cmd.ID = &id
	}

	out, err := c.RecordService.UpsertDashboardWidget(ctx, userID, cmd)
	if err != nil {
		return nil, err
	}
	return toGraphQLDashboardWidget(out), nil
}

func (c *controller) ReorderDashboardWidgets(ctx context.Context, userID uint64, in model.ReorderDashboardWidgetsInput) ([]*model.DashboardWidget, error) {
	items := make([]input.ReorderDashboardWidgetsItem, 0, len(in.Items))
	for _, item := range in.Items {
		items = append(items, input.ReorderDashboardWidgetsItem{
			WidgetID:   mustParseID(item.ID),
			OrderIndex: int(item.OrderIndex),
		})
	}

	out, err := c.RecordService.ReorderDashboardWidgets(ctx, userID, input.ReorderDashboardWidgetsCommand{
		ViewID: mustParseID(in.ViewID),
		Items:  items,
	})
	if err != nil {
		return nil, err
	}

	res := make([]*model.DashboardWidget, 0, len(out))
	for _, item := range out {
		res = append(res, toGraphQLDashboardWidget(item))
	}
	return res, nil
}

func (c *controller) DeleteDashboardWidget(ctx context.Context, userID uint64, widgetID string) error {
	return c.RecordService.DeleteDashboardWidget(ctx, userID, mustParseID(widgetID))
}

func (c *controller) CreateMetricAndWidget(ctx context.Context, userID uint64, in model.CreateMetricAndWidgetInput) (*model.DashboardWidget, error) {
	if in.Metric == nil {
		return nil, errors.New("metric is required")
	}
	if in.Widget == nil {
		return nil, errors.New("widget is required")
	}
	metric := *in.Metric
	widget := *in.Widget

	metricCmd := input.UpsertMetricDefinitionCommand{
		MetricKey:   metric.MetricKey,
		DisplayName: metric.DisplayName,
		TagID:       mustParseID(metric.TagID),
		TagIDs:      parseIDs(metric.TagIds),
		ValueSource: toPtrValue(metric.ValueSource),
		Aggregation: toPtrValue(metric.Aggregation),
		Unit:        toPtrValue(metric.Unit),
		GoalDefault: metric.GoalDefault,
		IsActive:    metric.IsActive,
	}
	if metric.ID != nil {
		id := mustParseID(*metric.ID)
		metricCmd.ID = &id
	}
	if metric.CategoryID != nil {
		cid := mustParseID(*metric.CategoryID)
		metricCmd.CategoryID = &cid
	}

	widgetCmd := input.UpsertDashboardWidgetCommand{
		ViewID:        mustParseID(widget.ViewID),
		WidgetType:    normalizeGQLEnum(string(widget.WidgetType)),
		Size:          normalizeGQLEnum(string(widget.Size)),
		OrderIndex:    toIntPtr(widget.OrderIndex),
		TitleOverride: widget.TitleOverride,
		ConfigJSON:    toStringValue(widget.ConfigJSON),
		IsActive:      widget.IsActive,
	}
	if widget.MetricDefinitionID != nil {
		widgetCmd.MetricDefinitionID = mustParseID(*widget.MetricDefinitionID)
	}
	if widget.ID != nil {
		id := mustParseID(*widget.ID)
		widgetCmd.ID = &id
	}

	out, err := c.RecordService.CreateMetricAndWidget(ctx, userID, input.CreateMetricAndWidgetCommand{
		Metric: metricCmd,
		Widget: widgetCmd,
	})
	if err != nil {
		return nil, err
	}
	return toGraphQLDashboardWidget(out), nil
}

func (c *controller) DashboardWidgetCatalog(_ context.Context) (*model.DashboardWidgetCatalog, error) {
	return &model.DashboardWidgetCatalog{
		MaxLargeWidgets: int32(domain.MaxLargeWidgetsPerDashboard),
		Sizes: []model.DashboardWidgetSize{
			model.DashboardWidgetSizeSmall,
			model.DashboardWidgetSizeMedium,
			model.DashboardWidgetSizeLarge,
		},
		Types: []model.DashboardWidgetType{
			model.DashboardWidgetTypeKpiNumber,
			model.DashboardWidgetTypeGoalProgress,
			model.DashboardWidgetTypeTrendLine,
			model.DashboardWidgetTypeChecklist,
			model.DashboardWidgetTypeInsightFeed,
		},
	}, nil
}

func (c *controller) SuggestMetricDefinitions(ctx context.Context, userID uint64, limit *int32) ([]*model.MetricDefinitionSuggestion, error) {
	lim := 8
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}
	out, err := c.RecordService.SuggestMetricDefinitions(ctx, userID, lim)
	if err != nil {
		return nil, err
	}

	items := make([]*model.MetricDefinitionSuggestion, 0, len(out))
	for _, item := range out {
		row := &model.MetricDefinitionSuggestion{
			MetricKey:   item.MetricKey,
			DisplayName: item.DisplayName,
			TagIds:      formatIDs(item.TagIDs),
			ValueSource: item.ValueSource,
			Aggregation: item.Aggregation,
			Unit:        item.Unit,
			Reason:      item.Reason,
		}
		if item.CategoryID != nil {
			value := strconv.FormatUint(*item.CategoryID, 10)
			row.CategoryID = &value
		}
		items = append(items, row)
	}
	return items, nil
}

func toGraphQLDashboardView(in domain.DashboardView, widgets []domain.DashboardWidget) *model.DashboardView {
	out := &model.DashboardView{
		ID:        strconv.FormatUint(in.ID, 10),
		Name:      in.Name,
		IsDefault: in.IsDefault,
		CreatedAt: in.CreatedAt.Format(timeLayout),
		UpdatedAt: in.UpdatedAt.Format(timeLayout),
		Widgets:   make([]*model.DashboardWidget, 0, len(widgets)),
	}
	for _, widget := range widgets {
		out.Widgets = append(out.Widgets, toGraphQLDashboardWidget(widget))
	}
	return out
}

func toGraphQLDashboardWidget(in domain.DashboardWidget) *model.DashboardWidget {
	var metricDefinitionID *string
	if in.MetricDefinitionID != 0 {
		value := strconv.FormatUint(in.MetricDefinitionID, 10)
		metricDefinitionID = &value
	}
	return &model.DashboardWidget{
		ID:                 strconv.FormatUint(in.ID, 10),
		ViewID:             strconv.FormatUint(in.ViewID, 10),
		MetricDefinitionID: metricDefinitionID,
		WidgetType:         toGraphQLWidgetType(in.WidgetType),
		Size:               toGraphQLWidgetSize(in.Size),
		OrderIndex:         safeInt32(in.OrderIndex),
		TitleOverride:      in.TitleOverride,
		ConfigJSON:         toStringPtr(in.ConfigJSON),
		IsActive:           in.IsActive,
		CreatedAt:          in.CreatedAt.Format(timeLayout),
		UpdatedAt:          in.UpdatedAt.Format(timeLayout),
	}
}

func toGraphQLWidgetType(v string) model.DashboardWidgetType {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case domain.DashboardWidgetTypeGoalProgress:
		return model.DashboardWidgetTypeGoalProgress
	case domain.DashboardWidgetTypeTrendLine:
		return model.DashboardWidgetTypeTrendLine
	case domain.DashboardWidgetTypeChecklist:
		return model.DashboardWidgetTypeChecklist
	case domain.DashboardWidgetTypeInsightFeed:
		return model.DashboardWidgetTypeInsightFeed
	default:
		return model.DashboardWidgetTypeKpiNumber
	}
}

func toGraphQLWidgetSize(v string) model.DashboardWidgetSize {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case domain.DashboardWidgetSizeMedium:
		return model.DashboardWidgetSizeMedium
	case domain.DashboardWidgetSizeLarge:
		return model.DashboardWidgetSizeLarge
	default:
		return model.DashboardWidgetSizeSmall
	}
}

func normalizeGQLEnum(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

func parseIDs(values []string) []uint64 {
	out := make([]uint64, 0, len(values))
	for _, value := range values {
		id := mustParseID(value)
		if id == 0 {
			continue
		}
		out = append(out, id)
	}
	return out
}

const timeLayout = "2006-01-02T15:04:05Z07:00"

func toStringValue(v *string) string {
	if v == nil || strings.TrimSpace(*v) == "" {
		return "{}"
	}
	return *v
}

func toStringPtr(v string) *string {
	value := strings.TrimSpace(v)
	if value == "" {
		return nil
	}
	return &value
}

func toIntPtr(v *int32) *int {
	if v == nil {
		return nil
	}
	value := int(*v)
	return &value
}

func safeInt32(value int) int32 {
	if value > math.MaxInt32 {
		return math.MaxInt32
	}
	if value < math.MinInt32 {
		return math.MinInt32
	}
	return int32(value)
}
