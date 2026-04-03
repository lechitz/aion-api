package input

import "time"

// UpsertMetricDefinitionCommand contains input data to create/update a metric definition.
type UpsertMetricDefinitionCommand struct {
	ID          *uint64
	MetricKey   string
	DisplayName string
	CategoryID  *uint64
	TagID       uint64
	TagIDs      []uint64
	ValueSource string
	Aggregation string
	Unit        string
	GoalDefault *float64
	IsActive    *bool
}

// UpsertGoalTemplateCommand contains input data to create/update a goal template.
type UpsertGoalTemplateCommand struct {
	ID          *uint64
	MetricKey   string
	Title       string
	TargetValue float64
	Comparison  string
	Period      string
	IsActive    *bool
}

// DashboardSnapshotQuery contains input parameters for dashboard snapshot queries.
type DashboardSnapshotQuery struct {
	Date     time.Time
	Timezone string
}

// InsightFeedQuery contains input parameters for the canonical insight feed.
type InsightFeedQuery struct {
	Window     string
	Limit      int
	Date       time.Time
	Timezone   string
	CategoryID *uint64
	TagIDs     []uint64
}

// AnalyticsSeriesQuery contains input parameters for analytics series retrieval.
type AnalyticsSeriesQuery struct {
	SeriesKey  string
	Window     string
	Date       time.Time
	Timezone   string
	CategoryID *uint64
	TagIDs     []uint64
}

// CreateDashboardViewCommand contains input for creating a dashboard view.
type CreateDashboardViewCommand struct {
	Name      string
	IsDefault *bool
}

// UpdateDashboardViewCommand contains input for renaming one dashboard view.
type UpdateDashboardViewCommand struct {
	Name string
}

// UpsertDashboardWidgetCommand contains input data to create/update dashboard widget.
type UpsertDashboardWidgetCommand struct {
	ID                 *uint64
	ViewID             uint64
	MetricDefinitionID uint64
	WidgetType         string
	Size               string
	OrderIndex         *int
	TitleOverride      *string
	ConfigJSON         string
	IsActive           *bool
}

// ReorderDashboardWidgetsItem defines widget order item.
type ReorderDashboardWidgetsItem struct {
	WidgetID   uint64
	OrderIndex int
}

// ReorderDashboardWidgetsCommand contains input for bulk reorder.
type ReorderDashboardWidgetsCommand struct {
	ViewID uint64
	Items  []ReorderDashboardWidgetsItem
}

// CreateMetricAndWidgetCommand creates metric definition + widget in one operation.
type CreateMetricAndWidgetCommand struct {
	Metric UpsertMetricDefinitionCommand
	Widget UpsertDashboardWidgetCommand
}
