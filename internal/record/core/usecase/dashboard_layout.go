package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
)

// ListDashboardViews lists all user dashboard views, creating defaults when empty.
func (s *Service) ListDashboardViews(ctx context.Context, userID uint64) ([]domain.DashboardView, error) {
	if userID == 0 {
		return nil, ErrUserIDIsRequired
	}
	return s.ensureDashboardViews(ctx, userID)
}

// GetDashboardView retrieves one dashboard view by ID.
func (s *Service) GetDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error) {
	if userID == 0 {
		return domain.DashboardView{}, ErrUserIDIsRequired
	}
	if viewID == 0 {
		return domain.DashboardView{}, errors.New(ErrDashboardViewIDRequired)
	}
	return s.RecordRepository.GetDashboardView(ctx, userID, viewID)
}

// CreateDashboardView creates a new dashboard view.
func (s *Service) CreateDashboardView(ctx context.Context, userID uint64, cmd input.CreateDashboardViewCommand) (domain.DashboardView, error) {
	if userID == 0 {
		return domain.DashboardView{}, ErrUserIDIsRequired
	}

	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		name = DefaultDashboardViewName
	}
	isDefault := cmd.IsDefault != nil && *cmd.IsDefault

	view, err := s.RecordRepository.CreateDashboardView(ctx, domain.DashboardView{
		UserID:    userID,
		Name:      name,
		IsDefault: isDefault,
	})
	if err != nil {
		return domain.DashboardView{}, err
	}

	if isDefault {
		return s.RecordRepository.SetDefaultDashboardView(ctx, userID, view.ID)
	}
	return view, nil
}

// UpdateDashboardView renames one dashboard view.
func (s *Service) UpdateDashboardView(ctx context.Context, userID uint64, viewID uint64, cmd input.UpdateDashboardViewCommand) (domain.DashboardView, error) {
	if userID == 0 {
		return domain.DashboardView{}, ErrUserIDIsRequired
	}
	if viewID == 0 {
		return domain.DashboardView{}, errors.New(ErrDashboardViewIDRequired)
	}

	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return domain.DashboardView{}, errors.New(ErrDashboardViewNameRequired)
	}

	return s.RecordRepository.UpdateDashboardView(ctx, userID, viewID, name)
}

// SetDefaultDashboardView sets the user's default dashboard view.
func (s *Service) SetDefaultDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error) {
	if userID == 0 {
		return domain.DashboardView{}, ErrUserIDIsRequired
	}
	if viewID == 0 {
		return domain.DashboardView{}, errors.New(ErrDashboardViewIDRequired)
	}
	return s.RecordRepository.SetDefaultDashboardView(ctx, userID, viewID)
}

// DeleteDashboardView removes one dashboard view while preserving at least one remaining view.
func (s *Service) DeleteDashboardView(ctx context.Context, userID uint64, viewID uint64) error {
	if userID == 0 {
		return ErrUserIDIsRequired
	}
	if viewID == 0 {
		return errors.New(ErrDashboardViewIDRequired)
	}

	views, err := s.ensureDashboardViews(ctx, userID)
	if err != nil {
		return err
	}

	if len(views) <= 1 {
		return errors.New(ErrDashboardLastViewDeleteBlocked)
	}

	return s.RecordRepository.DeleteDashboardView(ctx, userID, viewID)
}

// UpsertDashboardWidget creates or updates a dashboard widget.
func (s *Service) UpsertDashboardWidget(ctx context.Context, userID uint64, cmd input.UpsertDashboardWidgetCommand) (domain.DashboardWidget, error) {
	if userID == 0 {
		return domain.DashboardWidget{}, ErrUserIDIsRequired
	}
	if cmd.ViewID == 0 {
		return domain.DashboardWidget{}, errors.New(ErrDashboardViewIDRequired)
	}
	widgetType := normalizeWidgetType(cmd.WidgetType)
	if cmd.MetricDefinitionID == 0 && widgetType != domain.DashboardWidgetTypeInsightFeed {
		return domain.DashboardWidget{}, errors.New(ErrDashboardMetricDefinitionIDRequired)
	}
	size := normalizeWidgetSize(cmd.Size)

	if size == domain.DashboardWidgetSizeLarge {
		var exclude *uint64
		if cmd.ID != nil && *cmd.ID != 0 {
			exclude = cmd.ID
		}
		countLarge, err := s.RecordRepository.CountLargeWidgetsInView(ctx, userID, cmd.ViewID, exclude)
		if err != nil {
			return domain.DashboardWidget{}, err
		}
		if countLarge >= domain.MaxLargeWidgetsPerDashboard {
			return domain.DashboardWidget{}, fmt.Errorf(ErrDashboardLimitLargeWidgets, domain.MaxLargeWidgetsPerDashboard)
		}
	}

	var orderIndex int
	if cmd.OrderIndex != nil {
		orderIndex = *cmd.OrderIndex
	} else {
		existing, err := s.RecordRepository.ListDashboardWidgetsByView(ctx, userID, cmd.ViewID)
		if err != nil {
			return domain.DashboardWidget{}, err
		}
		orderIndex = len(existing)
	}

	isActive := true
	if cmd.IsActive != nil {
		isActive = *cmd.IsActive
	}

	widget := domain.DashboardWidget{
		UserID:             userID,
		ViewID:             cmd.ViewID,
		MetricDefinitionID: cmd.MetricDefinitionID,
		WidgetType:         widgetType,
		Size:               size,
		OrderIndex:         orderIndex,
		TitleOverride:      cmd.TitleOverride,
		ConfigJSON:         strings.TrimSpace(cmd.ConfigJSON),
		IsActive:           isActive,
	}
	if widget.ConfigJSON == "" {
		widget.ConfigJSON = DefaultDashboardConfigJSON
	}
	if cmd.ID != nil {
		widget.ID = *cmd.ID
	}

	return s.RecordRepository.UpsertDashboardWidget(ctx, widget)
}

// ReorderDashboardWidgets updates ordering for dashboard widgets in a view.
func (s *Service) ReorderDashboardWidgets(ctx context.Context, userID uint64, cmd input.ReorderDashboardWidgetsCommand) ([]domain.DashboardWidget, error) {
	if userID == 0 {
		return nil, ErrUserIDIsRequired
	}
	if cmd.ViewID == 0 {
		return nil, errors.New(ErrDashboardViewIDRequired)
	}
	if len(cmd.Items) == 0 {
		return nil, errors.New(ErrDashboardItemsRequired)
	}

	items := make([]domain.DashboardWidget, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		if item.WidgetID == 0 {
			return nil, errors.New(ErrDashboardWidgetIDRequired)
		}
		items = append(items, domain.DashboardWidget{
			ID:         item.WidgetID,
			OrderIndex: item.OrderIndex,
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].OrderIndex < items[j].OrderIndex })
	return s.RecordRepository.ReorderDashboardWidgets(ctx, userID, cmd.ViewID, items)
}

// DeleteDashboardWidget disables a dashboard widget.
func (s *Service) DeleteDashboardWidget(ctx context.Context, userID uint64, widgetID uint64) error {
	if userID == 0 {
		return ErrUserIDIsRequired
	}
	if widgetID == 0 {
		return errors.New(ErrDashboardWidgetIDRequired)
	}
	return s.RecordRepository.DeleteDashboardWidget(ctx, userID, widgetID)
}

// CreateMetricAndWidget creates/updates a metric definition and then creates/updates its widget.
func (s *Service) CreateMetricAndWidget(ctx context.Context, userID uint64, cmd input.CreateMetricAndWidgetCommand) (domain.DashboardWidget, error) {
	metric, err := s.UpsertMetricDefinition(ctx, userID, cmd.Metric)
	if err != nil {
		return domain.DashboardWidget{}, err
	}

	widgetCmd := cmd.Widget
	if widgetCmd.MetricDefinitionID == 0 {
		widgetCmd.MetricDefinitionID = metric.ID
	}
	return s.UpsertDashboardWidget(ctx, userID, widgetCmd)
}

// SuggestMetricDefinitions proposes deterministic metric definitions from user tags.
func (s *Service) SuggestMetricDefinitions(ctx context.Context, userID uint64, limit int) ([]domain.MetricDefinitionSuggestion, error) {
	if userID == 0 {
		return nil, ErrUserIDIsRequired
	}
	if limit <= 0 {
		limit = DefaultDashboardSuggestionsLimit
	}
	if limit > MaxDashboardSuggestionsLimit {
		limit = MaxDashboardSuggestionsLimit
	}

	// Deterministic suggestions from existing active tags.
	tags, err := s.TagRepository.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return []domain.MetricDefinitionSuggestion{}, nil
	}

	out := make([]domain.MetricDefinitionSuggestion, 0, limit)
	seen := make(map[string]struct{}, limit)
	for _, tag := range tags {
		if len(out) >= limit {
			break
		}
		key := slugMetricKey(tag.Name)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		tagID := tag.ID
		out = append(out, domain.MetricDefinitionSuggestion{
			MetricKey:   key,
			DisplayName: strings.TrimSpace(tag.Name),
			CategoryID:  &tag.CategoryID,
			TagIDs:      []uint64{tagID},
			ValueSource: DashboardValueSourceCount,
			Aggregation: DashboardAggregationSum,
			Unit:        DashboardUnitCount,
			Reason:      DashboardSuggestionReasonTaxonomy,
		})
	}
	return out, nil
}

func (s *Service) ensureDashboardViews(ctx context.Context, userID uint64) ([]domain.DashboardView, error) {
	views, err := s.RecordRepository.ListDashboardViews(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(views) > 0 {
		return views, nil
	}

	defaultView, err := s.RecordRepository.CreateDashboardView(ctx, domain.DashboardView{
		UserID:    userID,
		Name:      FallbackDashboardViewName,
		IsDefault: true,
	})
	if err != nil {
		return nil, err
	}
	_ = defaultView

	views, err = s.RecordRepository.ListDashboardViews(ctx, userID)
	if err != nil {
		return nil, err
	}
	return views, nil
}

func normalizeWidgetType(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case domain.DashboardWidgetTypeGoalProgress:
		return domain.DashboardWidgetTypeGoalProgress
	case domain.DashboardWidgetTypeTrendLine:
		return domain.DashboardWidgetTypeTrendLine
	case domain.DashboardWidgetTypeChecklist:
		return domain.DashboardWidgetTypeChecklist
	case domain.DashboardWidgetTypeInsightFeed:
		return domain.DashboardWidgetTypeInsightFeed
	default:
		return domain.DashboardWidgetTypeKPINumber
	}
}

func normalizeWidgetSize(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case domain.DashboardWidgetSizeMedium:
		return domain.DashboardWidgetSizeMedium
	case domain.DashboardWidgetSizeLarge:
		return domain.DashboardWidgetSizeLarge
	default:
		return domain.DashboardWidgetSizeSmall
	}
}

func slugMetricKey(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		return ""
	}
	replacer := strings.NewReplacer(
		DashboardSlugSpace, DashboardSlugUnderscore,
		DashboardSlugHyphen, DashboardSlugUnderscore,
		DashboardSlugSlash, DashboardSlugUnderscore,
		DashboardSlugLeftParenthesis, "",
		DashboardSlugRightParenthesis, "",
	)
	v = replacer.Replace(v)
	return strings.Trim(v, DashboardSlugUnderscore)
}
