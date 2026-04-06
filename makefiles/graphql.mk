# GraphQL Documentation and Shared Queries
.PHONY: graphql.schema graphql.queries graphql.manifest graphql.validate graphql.check-dirty graphql.docs graphql.setup graphql.clean

ROOT_DIR := $(shell pwd)
DOCS_DIR := $(ROOT_DIR)/docs/graphql
SCHEMA_DIR := $(ROOT_DIR)/internal/adapter/primary/graphql/schema
SCHEMA_OUT := $(DOCS_DIR)/schema.graphql
CONTRACTS_DIR := $(ROOT_DIR)/contracts/graphql
MUTATIONS_DIR := $(CONTRACTS_DIR)/mutations
QUERIES_DIR := $(ROOT_DIR)/contracts/graphql/queries
MANIFEST_OUT := $(CONTRACTS_DIR)/manifest.json
MANIFEST_SCRIPT := $(ROOT_DIR)/hack/tools/graphql-contract-manifest.sh
VALIDATE_SCRIPT := $(ROOT_DIR)/hack/tools/graphql-contract-validate.sh

graphql.schema:
	@echo "Exporting schema..."
	@mkdir -p "$(DOCS_DIR)"
	@echo "# Aion GraphQL Schema" > "$(SCHEMA_OUT)"
	@find "$(SCHEMA_DIR)" -name "*.graphqls" | sort | xargs cat >> "$(SCHEMA_OUT)"
	@echo "✅ Schema: $(SCHEMA_OUT)"

graphql.queries:
	@echo "Creating shared GraphQL operations..."
	@mkdir -p "$(QUERIES_DIR)/categories" "$(QUERIES_DIR)/tags" "$(QUERIES_DIR)/records" "$(QUERIES_DIR)/chat" "$(QUERIES_DIR)/user" "$(QUERIES_DIR)/dashboard"
	@mkdir -p "$(MUTATIONS_DIR)/categories" "$(MUTATIONS_DIR)/tags" "$(MUTATIONS_DIR)/records" "$(MUTATIONS_DIR)/dashboard"
	@printf 'query ListCategories { categories { id name description colorHex icon } }\n' > "$(QUERIES_DIR)/categories/list.graphql"
	@printf 'query CategoryById($$id: ID!) { categoryById(id: $$id) { id userId name description colorHex icon } }\n' > "$(QUERIES_DIR)/categories/by-id.graphql"
	@printf 'query CategoryByName($$name: String!) { categoryByName(name: $$name) { id userId name description colorHex icon } }\n' > "$(QUERIES_DIR)/categories/by-name.graphql"
	@printf 'query ListTags { tags { id userId name categoryId description icon createdAt updatedAt } }\n' > "$(QUERIES_DIR)/tags/list.graphql"
	@printf 'query TagById($$id: ID!) { tagById(id: $$id) { id userId name categoryId description icon createdAt updatedAt } }\n' > "$(QUERIES_DIR)/tags/by-id.graphql"
	@printf 'query TagByName($$name: String!) { tagByName(name: $$name) { id userId name categoryId description icon createdAt updatedAt } }\n' > "$(QUERIES_DIR)/tags/by-name.graphql"
	@printf 'query TagsByCategoryId($$categoryId: ID!) { tagsByCategoryId(categoryId: $$categoryId) { id userId name categoryId description icon createdAt updatedAt } }\n' > "$(QUERIES_DIR)/tags/by-category-id.graphql"
	@printf 'query ListRecords($$limit: Int) { records(limit: $$limit) { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } }\n' > "$(QUERIES_DIR)/records/list.graphql"
	@printf 'query RecordById($$id: ID!) { recordById(id: $$id) { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } }\n' > "$(QUERIES_DIR)/records/by-id.graphql"
	@printf 'query RecordsLatest($$limit: Int) { recordsLatest(limit: $$limit) { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } }\n' > "$(QUERIES_DIR)/records/latest.graphql"
	@printf 'query RecordProjectionById($$id: ID!) { recordProjectionById(id: $$id) { recordId userId tagId description eventTimeUTC recordedAtUTC durationSeconds value source timezone status createdAtUTC updatedAtUTC lastEventType } }\n' > "$(QUERIES_DIR)/records/projection-by-id.graphql"
	@printf 'query RecordProjectionsLatest($$limit: Int) { recordProjectionsLatest(limit: $$limit) { recordId userId tagId description eventTimeUTC recordedAtUTC durationSeconds value source timezone status createdAtUTC updatedAtUTC lastEventType } }\n' > "$(QUERIES_DIR)/records/projections-latest.graphql"
	@printf 'query RecordProjections($$limit: Int, $$afterEventTime: String, $$afterId: ID) { recordProjections(limit: $$limit, afterEventTime: $$afterEventTime, afterId: $$afterId) { recordId userId tagId description eventTimeUTC recordedAtUTC durationSeconds value source timezone status createdAtUTC updatedAtUTC lastEventType } }\n' > "$(QUERIES_DIR)/records/projections.graphql"
	@printf 'query RecordsByTag($$tagId: ID!, $$limit: Int) { recordsByTag(tagId: $$tagId, limit: $$limit) { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } }\n' > "$(QUERIES_DIR)/records/by-tag.graphql"
	@printf 'query RecordsByCategory($$categoryId: ID!, $$limit: Int) { recordsByCategory(categoryId: $$categoryId, limit: $$limit) { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } }\n' > "$(QUERIES_DIR)/records/by-category.graphql"
	@printf 'query RecordsByDay($$date: String!) { recordsByDay(date: $$date) { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } }\n' > "$(QUERIES_DIR)/records/by-day.graphql"
	@printf 'query RecordsUntil($$until: String!, $$limit: Int) { recordsUntil(until: $$until, limit: $$limit) { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } }\n' > "$(QUERIES_DIR)/records/until.graphql"
	@printf 'query RecordsBetween($$startDate: String!, $$endDate: String!, $$limit: Int) { recordsBetween(startDate: $$startDate, endDate: $$endDate, limit: $$limit) { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } }\n' > "$(QUERIES_DIR)/records/between.graphql"
	@printf 'query SearchRecords($$filters: SearchFilters!) { searchRecords(filters: $$filters) { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } }\n' > "$(QUERIES_DIR)/records/search.graphql"
	@printf 'query RecordStats($$filters: RecordStatsFilters) { recordStats(filters: $$filters) { totalRecords recordsWithValue totalDurationSeconds sumValue avgValue avgDurationSeconds minValue maxValue } }\n' > "$(QUERIES_DIR)/records/stats.graphql"
	@printf 'query ChatHistory($$limit: Int, $$offset: Int) { chatHistory(limit: $$limit, offset: $$offset) { id userId message response tokensUsed functionCalls createdAt updatedAt } }\n' > "$(QUERIES_DIR)/chat/history.graphql"
	@printf 'query ChatContext { chatContext { recentChats { id userId message response tokensUsed functionCalls createdAt updatedAt } totalRecords totalCategories totalTags } }\n' > "$(QUERIES_DIR)/chat/context.graphql"
	@printf 'query ChatDataPack($$limitRecords: Int, $$includeStats: Boolean!) { chatDataPack(limitRecords: $$limitRecords, includeStats: $$includeStats) { categories { id userId name description colorHex icon } tags { id userId name categoryId description icon createdAt updatedAt } recentRecords { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } userStats @include(if: $$includeStats) { totalRecords totalCategories totalTags recordsThisWeek recordsThisMonth mostUsedCategory { id name count } mostUsedTag { id name count } } } }\n' > "$(QUERIES_DIR)/chat/data-pack.graphql"
	@printf 'query UserStats { userStats { totalRecords totalCategories totalTags recordsThisWeek recordsThisMonth mostUsedCategory { id name count } mostUsedTag { id name count } } }\n' > "$(QUERIES_DIR)/user/stats.graphql"
	@printf 'query DashboardSnapshot($$date: String!, $$timezone: String) { dashboardSnapshot(date: $$date, timezone: $$timezone) { date timezone metrics { metricKey label value unit target progressPct status } goals { goalId title metricKey currentValue targetValue progressPct status } } }\n' > "$(QUERIES_DIR)/dashboard/snapshot.graphql"
	@printf 'query InsightFeed($$window: InsightWindow!, $$limit: Int, $$date: String, $$timezone: String, $$categoryId: ID, $$tagIds: [ID!]) { insightFeed(window: $$window, limit: $$limit, date: $$date, timezone: $$timezone, categoryId: $$categoryId, tagIds: $$tagIds) { id type title summary status window confidence metricKeys recommendedAction evidence { label value kind } generatedAt } }\n' > "$(QUERIES_DIR)/dashboard/insight-feed.graphql"
	@printf 'query AnalyticsSeries($$seriesKey: String!, $$window: InsightWindow!, $$date: String, $$timezone: String, $$categoryId: ID, $$tagIds: [ID!]) { analyticsSeries(seriesKey: $$seriesKey, window: $$window, date: $$date, timezone: $$timezone, categoryId: $$categoryId, tagIds: $$tagIds) { seriesKey window points { timestamp value label } summary } }\n' > "$(QUERIES_DIR)/dashboard/analytics-series.graphql"
	@printf 'query MetricDefinitions { metricDefinitions { id metricKey displayName categoryId tagId tagIds valueSource aggregation unit goalDefault isActive } }\n' > "$(QUERIES_DIR)/dashboard/metric-definitions.graphql"
	@printf 'query DashboardViews { dashboardViews { id name isDefault widgets { id viewId metricDefinitionId widgetType size orderIndex titleOverride configJson isActive createdAt updatedAt } createdAt updatedAt } }\n' > "$(QUERIES_DIR)/dashboard/views.graphql"
	@printf 'query DashboardView($$id: ID!) { dashboardView(id: $$id) { id name isDefault widgets { id viewId metricDefinitionId widgetType size orderIndex titleOverride configJson isActive createdAt updatedAt } createdAt updatedAt } }\n' > "$(QUERIES_DIR)/dashboard/view.graphql"
	@printf 'query DashboardWidgetCatalog { dashboardWidgetCatalog { maxLargeWidgets sizes types } }\n' > "$(QUERIES_DIR)/dashboard/widget-catalog.graphql"
	@printf 'query SuggestMetricDefinitions($$limit: Int) { suggestMetricDefinitions(limit: $$limit) { metricKey displayName categoryId tagIds valueSource aggregation unit reason } }\n' > "$(QUERIES_DIR)/dashboard/suggest-metric-definitions.graphql"
	@printf 'mutation CreateCategory($$input: CreateCategoryInput!) { createCategory(input: $$input) { id userId name description colorHex icon } }\n' > "$(MUTATIONS_DIR)/categories/create.graphql"
	@printf 'mutation UpdateCategory($$input: UpdateCategoryInput!) { updateCategory(input: $$input) { id userId name description colorHex icon } }\n' > "$(MUTATIONS_DIR)/categories/update.graphql"
	@printf 'mutation SoftDeleteCategory($$input: DeleteCategoryInput!) { softDeleteCategory(input: $$input) }\n' > "$(MUTATIONS_DIR)/categories/delete.graphql"
	@printf 'mutation CreateTag($$input: CreateTagInput!) { createTag(input: $$input) { id userId name categoryId description icon createdAt updatedAt } }\n' > "$(MUTATIONS_DIR)/tags/create.graphql"
	@printf 'mutation UpdateTag($$input: UpdateTagInput!) { updateTag(input: $$input) { id userId name categoryId description icon createdAt updatedAt } }\n' > "$(MUTATIONS_DIR)/tags/update.graphql"
	@printf 'mutation SoftDeleteTag($$input: DeleteTagInput!) { softDeleteTag(input: $$input) }\n' > "$(MUTATIONS_DIR)/tags/delete.graphql"
	@printf 'mutation CreateRecord($$input: CreateRecordInput!) { createRecord(input: $$input) { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } }\n' > "$(MUTATIONS_DIR)/records/create.graphql"
	@printf 'mutation UpdateRecord($$input: UpdateRecordInput!) { updateRecord(input: $$input) { id userId tagId description eventTime recordedAt durationSeconds value source timezone status createdAt updatedAt } }\n' > "$(MUTATIONS_DIR)/records/update.graphql"
	@printf 'mutation SoftDeleteRecord($$input: DeleteRecordInput!) { softDeleteRecord(input: $$input) }\n' > "$(MUTATIONS_DIR)/records/delete.graphql"
	@printf 'mutation SoftDeleteAllRecords { softDeleteAllRecords }\n' > "$(MUTATIONS_DIR)/records/delete-all.graphql"
	@printf 'mutation UpsertMetricDefinition($$input: UpsertMetricDefinitionInput!) { upsertMetricDefinition(input: $$input) { id metricKey displayName categoryId tagId tagIds valueSource aggregation unit goalDefault isActive } }\n' > "$(MUTATIONS_DIR)/dashboard/upsert-metric-definition.graphql"
	@printf 'mutation UpsertGoalTemplate($$input: UpsertGoalTemplateInput!) { upsertGoalTemplate(input: $$input) { id metricKey title targetValue comparison period isActive } }\n' > "$(MUTATIONS_DIR)/dashboard/upsert-goal-template.graphql"
	@printf 'mutation DeleteGoalTemplate($$input: DeleteGoalTemplateInput!) { deleteGoalTemplate(input: $$input) }\n' > "$(MUTATIONS_DIR)/dashboard/delete-goal-template.graphql"
	@printf 'mutation CreateDashboardView($$input: CreateDashboardViewInput!) { createDashboardView(input: $$input) { id name isDefault widgets { id viewId metricDefinitionId widgetType size orderIndex titleOverride configJson isActive createdAt updatedAt } createdAt updatedAt } }\n' > "$(MUTATIONS_DIR)/dashboard/create-view.graphql"
	@printf 'mutation SetDefaultDashboardView($$input: SetDefaultDashboardViewInput!) { setDefaultDashboardView(input: $$input) { id name isDefault widgets { id viewId metricDefinitionId widgetType size orderIndex titleOverride configJson isActive createdAt updatedAt } createdAt updatedAt } }\n' > "$(MUTATIONS_DIR)/dashboard/set-default-view.graphql"
	@printf 'mutation UpsertDashboardWidget($$input: UpsertDashboardWidgetInput!) { upsertDashboardWidget(input: $$input) { id viewId metricDefinitionId widgetType size orderIndex titleOverride configJson isActive createdAt updatedAt } }\n' > "$(MUTATIONS_DIR)/dashboard/upsert-widget.graphql"
	@printf 'mutation ReorderDashboardWidgets($$input: ReorderDashboardWidgetsInput!) { reorderDashboardWidgets(input: $$input) { id viewId metricDefinitionId widgetType size orderIndex titleOverride configJson isActive createdAt updatedAt } }\n' > "$(MUTATIONS_DIR)/dashboard/reorder-widgets.graphql"
	@printf 'mutation DeleteDashboardWidget($$input: DeleteDashboardWidgetInput!) { deleteDashboardWidget(input: $$input) }\n' > "$(MUTATIONS_DIR)/dashboard/delete-widget.graphql"
	@printf 'mutation CreateMetricAndWidget($$input: CreateMetricAndWidgetInput!) { createMetricAndWidget(input: $$input) { id viewId metricDefinitionId widgetType size orderIndex titleOverride configJson isActive createdAt updatedAt } }\n' > "$(MUTATIONS_DIR)/dashboard/create-metric-and-widget.graphql"
	@echo "✅ Operations: $(CONTRACTS_DIR)/"

graphql.manifest:
	@echo "Generating GraphQL contracts manifest..."
	@"$(MANIFEST_SCRIPT)"
	@echo "✅ Manifest: $(MANIFEST_OUT)"

graphql.validate:
	@echo "Validating GraphQL contracts against schema..."
	@"$(VALIDATE_SCRIPT)"

graphql.check-dirty:
	@tracked_diff=0; \
	git diff --quiet -- ':(glob)contracts/graphql/queries/**/*.graphql' ':(glob)contracts/graphql/mutations/**/*.graphql' "$(MANIFEST_OUT)" || tracked_diff=1; \
	untracked=$$(git ls-files --others --exclude-standard -- "$(QUERIES_DIR)" "$(MUTATIONS_DIR)" "$(MANIFEST_OUT)" | grep -E '(\.graphql$$|manifest\.json$$)' || true); \
	if [ $$tracked_diff -ne 0 ] || [ -n "$$untracked" ]; then \
		echo "GraphQL contracts out-of-date. Run 'make graphql.queries graphql.manifest'."; \
		exit 1; \
	fi

graphql.docs: graphql.schema
	@printf '# GraphQL Documentation\n\nSchema: schema.graphql\nPlayground: http://localhost:5001/aion/api/v1/graphql/playground\n' > "$(DOCS_DIR)/README.md"
	@echo "✅ Docs ready"

graphql.clean:
	@rm -rf "$(DOCS_DIR)" "$(QUERIES_DIR)" "$(MUTATIONS_DIR)" "$(MANIFEST_OUT)"
	@echo "Cleaned"

graphql.setup: graphql.schema graphql.queries graphql.manifest graphql.validate graphql.docs
	@echo ""
	@echo "✅ GraphQL documentation setup complete!"
	@echo "   Schema:  $(SCHEMA_OUT)"
	@echo "   Queries: $(QUERIES_DIR)/"
	@echo "   Manifest: $(MANIFEST_OUT)"
	@echo ""
