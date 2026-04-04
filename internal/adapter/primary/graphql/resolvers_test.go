//nolint:testpackage // requires access to generated internal resolver types.
package graphql

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gmodel "github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/category/core/domain"
	catinput "github.com/lechitz/aion-api/internal/category/core/ports/input"
	chatdomain "github.com/lechitz/aion-api/internal/chat/core/domain"
	chatinput "github.com/lechitz/aion-api/internal/chat/core/ports/input"
	"github.com/lechitz/aion-api/internal/platform/config"
	recorddomain "github.com/lechitz/aion-api/internal/record/core/domain"
	recordinput "github.com/lechitz/aion-api/internal/record/core/ports/input"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	tagdomain "github.com/lechitz/aion-api/internal/tag/core/domain"
	taginput "github.com/lechitz/aion-api/internal/tag/core/ports/input"
	userdomain "github.com/lechitz/aion-api/internal/user/core/domain"
	userinput "github.com/lechitz/aion-api/internal/user/core/ports/input"
	"github.com/stretchr/testify/require"
)

type gqlLoggerStub struct{}

func (gqlLoggerStub) Infof(string, ...any)                      {}
func (gqlLoggerStub) Errorf(string, ...any)                     {}
func (gqlLoggerStub) Debugf(string, ...any)                     {}
func (gqlLoggerStub) Warnf(string, ...any)                      {}
func (gqlLoggerStub) Infow(string, ...any)                      {}
func (gqlLoggerStub) Errorw(string, ...any)                     {}
func (gqlLoggerStub) Debugw(string, ...any)                     {}
func (gqlLoggerStub) Warnw(string, ...any)                      {}
func (gqlLoggerStub) InfowCtx(context.Context, string, ...any)  {}
func (gqlLoggerStub) ErrorwCtx(context.Context, string, ...any) {}
func (gqlLoggerStub) WarnwCtx(context.Context, string, ...any)  {}
func (gqlLoggerStub) DebugwCtx(context.Context, string, ...any) {}

type categorySvcStub struct{}

func (categorySvcStub) Create(context.Context, catinput.CreateCategoryCommand) (domain.Category, error) {
	return domain.Category{ID: 1, UserID: 1, Name: "cat", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (categorySvcStub) Update(context.Context, catinput.UpdateCategoryCommand) (domain.Category, error) {
	return domain.Category{ID: 1, UserID: 1, Name: "cat2", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (categorySvcStub) GetByID(context.Context, uint64, uint64) (domain.Category, error) {
	return domain.Category{ID: 1, UserID: 1, Name: "cat", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (categorySvcStub) GetByName(context.Context, string, uint64) (domain.Category, error) {
	return domain.Category{ID: 1, UserID: 1, Name: "cat", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (categorySvcStub) ListAll(context.Context, uint64) ([]domain.Category, error) {
	return []domain.Category{{ID: 1, UserID: 1, Name: "cat", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
}
func (categorySvcStub) SoftDelete(context.Context, uint64, uint64) error { return nil }

type tagSvcStub struct{}

func (tagSvcStub) Create(context.Context, taginput.CreateTagCommand) (tagdomain.Tag, error) {
	return tagdomain.Tag{ID: 1, UserID: 1, CategoryID: 1, Name: "tag", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (tagSvcStub) Update(context.Context, taginput.UpdateTagCommand) (tagdomain.Tag, error) {
	return tagdomain.Tag{ID: 1, UserID: 1, CategoryID: 1, Name: "tag2", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (tagSvcStub) GetByID(context.Context, uint64, uint64) (tagdomain.Tag, error) {
	return tagdomain.Tag{ID: 1, UserID: 1, CategoryID: 1, Name: "tag", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (tagSvcStub) GetByName(context.Context, string, uint64) (tagdomain.Tag, error) {
	return tagdomain.Tag{ID: 1, UserID: 1, CategoryID: 1, Name: "tag", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (tagSvcStub) GetByCategoryID(context.Context, uint64, uint64) ([]tagdomain.Tag, error) {
	return []tagdomain.Tag{{ID: 1, UserID: 1, CategoryID: 1, Name: "tag", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
}

func (tagSvcStub) GetAll(context.Context, uint64) ([]tagdomain.Tag, error) {
	return []tagdomain.Tag{{ID: 1, UserID: 1, CategoryID: 1, Name: "tag", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
}
func (tagSvcStub) SoftDelete(context.Context, uint64, uint64) error { return nil }

type recordSvcStub struct{}

func (recordSvcStub) Create(context.Context, recordinput.CreateRecordCommand) (recorddomain.Record, error) {
	now := time.Now().UTC()
	return recorddomain.Record{ID: 1, UserID: 1, TagID: 1, EventTime: now, CreatedAt: now, UpdatedAt: now}, nil
}

func (recordSvcStub) GetByID(context.Context, uint64, uint64) (recorddomain.Record, error) {
	now := time.Now().UTC()
	return recorddomain.Record{ID: 1, UserID: 1, TagID: 1, EventTime: now, CreatedAt: now, UpdatedAt: now}, nil
}

func (recordSvcStub) GetProjectedByID(context.Context, uint64, uint64) (recorddomain.RecordProjection, error) {
	now := time.Now().UTC()
	return recorddomain.RecordProjection{
		RecordID:          1,
		UserID:            1,
		TagID:             1,
		EventTimeUTC:      now,
		LastEventID:       "evt-1",
		LastEventType:     "record.created",
		LastEventVersion:  "v1",
		LastKafkaTopic:    "aion.record.events.v1",
		LastConsumedAtUTC: now,
		CreatedAtUTC:      now,
		UpdatedAtUTC:      now,
	}, nil
}

func (recordSvcStub) ListByUser(context.Context, uint64, int, *string, *int64) ([]recorddomain.Record, error) {
	return []recorddomain.Record{}, nil
}

func (recordSvcStub) ListProjectedPage(context.Context, uint64, int, *string, *int64) ([]recorddomain.RecordProjection, error) {
	return []recorddomain.RecordProjection{}, nil
}

func (recordSvcStub) ListByTag(context.Context, uint64, uint64, int) ([]recorddomain.Record, error) {
	return []recorddomain.Record{}, nil
}

func (recordSvcStub) ListByCategory(context.Context, uint64, uint64, int) ([]recorddomain.Record, error) {
	return []recorddomain.Record{}, nil
}

func (recordSvcStub) ListByDay(context.Context, uint64, time.Time) ([]recorddomain.Record, error) {
	return []recorddomain.Record{}, nil
}

func (recordSvcStub) ListAllUntil(context.Context, uint64, time.Time, int) ([]recorddomain.Record, error) {
	return []recorddomain.Record{}, nil
}

func (recordSvcStub) ListAllBetween(context.Context, uint64, time.Time, time.Time, int) ([]recorddomain.Record, error) {
	return []recorddomain.Record{}, nil
}

func (recordSvcStub) ListLatest(context.Context, uint64, int) ([]recorddomain.Record, error) {
	return []recorddomain.Record{}, nil
}

func (recordSvcStub) ListProjectedLatest(context.Context, uint64, int) ([]recorddomain.RecordProjection, error) {
	return []recorddomain.RecordProjection{}, nil
}

func (recordSvcStub) Update(context.Context, uint64, uint64, recordinput.UpdateRecordCommand) (recorddomain.Record, error) {
	now := time.Now().UTC()
	return recorddomain.Record{ID: 1, UserID: 1, TagID: 1, EventTime: now, CreatedAt: now, UpdatedAt: now}, nil
}
func (recordSvcStub) Delete(context.Context, uint64, uint64) error { return nil }
func (recordSvcStub) DeleteAll(context.Context, uint64) error      { return nil }
func (recordSvcStub) SearchRecords(context.Context, uint64, recorddomain.SearchFilters) ([]recorddomain.Record, error) {
	return []recorddomain.Record{}, nil
}

func (recordSvcStub) DashboardSnapshot(context.Context, uint64, recordinput.DashboardSnapshotQuery) (recorddomain.DashboardSnapshot, error) {
	return recorddomain.DashboardSnapshot{}, nil
}

func (recordSvcStub) InsightFeed(context.Context, uint64, recordinput.InsightFeedQuery) ([]recorddomain.InsightCard, error) {
	return []recorddomain.InsightCard{}, nil
}

func (recordSvcStub) AnalyticsSeries(context.Context, uint64, recordinput.AnalyticsSeriesQuery) (recorddomain.AnalyticsSeriesResult, error) {
	return recorddomain.AnalyticsSeriesResult{}, nil
}

func (recordSvcStub) ListMetricDefinitions(context.Context, uint64) ([]recorddomain.MetricDefinition, error) {
	return []recorddomain.MetricDefinition{}, nil
}

func (recordSvcStub) UpsertMetricDefinition(context.Context, uint64, recordinput.UpsertMetricDefinitionCommand) (recorddomain.MetricDefinition, error) {
	return recorddomain.MetricDefinition{}, nil
}

func (recordSvcStub) UpsertGoalTemplate(context.Context, uint64, recordinput.UpsertGoalTemplateCommand) (recorddomain.GoalTemplate, error) {
	return recorddomain.GoalTemplate{}, nil
}
func (recordSvcStub) DeleteGoalTemplate(context.Context, uint64, uint64) error { return nil }
func (recordSvcStub) ListDashboardViews(context.Context, uint64) ([]recorddomain.DashboardView, error) {
	return []recorddomain.DashboardView{}, nil
}

func (recordSvcStub) GetDashboardView(context.Context, uint64, uint64) (recorddomain.DashboardView, error) {
	return recorddomain.DashboardView{}, nil
}

func (recordSvcStub) CreateDashboardView(context.Context, uint64, recordinput.CreateDashboardViewCommand) (recorddomain.DashboardView, error) {
	return recorddomain.DashboardView{}, nil
}

func (recordSvcStub) UpdateDashboardView(context.Context, uint64, uint64, recordinput.UpdateDashboardViewCommand) (recorddomain.DashboardView, error) {
	return recorddomain.DashboardView{}, nil
}

func (recordSvcStub) SetDefaultDashboardView(context.Context, uint64, uint64) (recorddomain.DashboardView, error) {
	return recorddomain.DashboardView{}, nil
}

func (recordSvcStub) DeleteDashboardView(context.Context, uint64, uint64) error { return nil }

func (recordSvcStub) UpsertDashboardWidget(context.Context, uint64, recordinput.UpsertDashboardWidgetCommand) (recorddomain.DashboardWidget, error) {
	return recorddomain.DashboardWidget{}, nil
}

func (recordSvcStub) ReorderDashboardWidgets(context.Context, uint64, recordinput.ReorderDashboardWidgetsCommand) ([]recorddomain.DashboardWidget, error) {
	return []recorddomain.DashboardWidget{}, nil
}
func (recordSvcStub) DeleteDashboardWidget(context.Context, uint64, uint64) error { return nil }
func (recordSvcStub) CreateMetricAndWidget(context.Context, uint64, recordinput.CreateMetricAndWidgetCommand) (recorddomain.DashboardWidget, error) {
	return recorddomain.DashboardWidget{}, nil
}

func (recordSvcStub) SuggestMetricDefinitions(context.Context, uint64, int) ([]recorddomain.MetricDefinitionSuggestion, error) {
	return []recorddomain.MetricDefinitionSuggestion{}, nil
}

type chatSvcStub struct{}

func (chatSvcStub) ProcessMessage(context.Context, uint64, string, map[string]interface{}, *chatdomain.RuntimeSelection) (*chatdomain.ChatResult, error) {
	return &chatdomain.ChatResult{}, nil
}

func (chatSvcStub) SaveChatHistory(context.Context, uint64, string, string, int, map[string]string) error {
	return nil
}

func (chatSvcStub) GetChatHistory(context.Context, uint64, int, int) ([]chatdomain.ChatHistory, error) {
	now := time.Now().UTC()
	return []chatdomain.ChatHistory{{ChatID: 1, UserID: 1, Message: "m", Response: "r", CreatedAt: now, UpdatedAt: now}}, nil
}

func (chatSvcStub) GetLatestChatHistory(context.Context, uint64, int) ([]chatdomain.ChatHistory, error) {
	return []chatdomain.ChatHistory{}, nil
}

func (chatSvcStub) GetChatContext(context.Context, uint64) (*chatdomain.ChatContext, error) {
	now := time.Now().UTC()
	return &chatdomain.ChatContext{RecentChats: []chatdomain.ChatHistory{{ChatID: 1, UserID: 1, Message: "m", Response: "r", CreatedAt: now, UpdatedAt: now}}}, nil
}

type userSvcStub struct{}

func (userSvcStub) Create(context.Context, userinput.CreateUserCommand) (userdomain.User, error) {
	return userdomain.User{}, nil
}

func (userSvcStub) StartRegistration(context.Context, userinput.StartRegistrationCommand) (userdomain.RegistrationSession, error) {
	return userdomain.RegistrationSession{}, nil
}

func (userSvcStub) UpdateRegistrationProfile(context.Context, string, userinput.UpdateRegistrationProfileCommand) (userdomain.RegistrationSession, error) {
	return userdomain.RegistrationSession{}, nil
}

func (userSvcStub) UpdateRegistrationAvatar(context.Context, string, userinput.UpdateRegistrationAvatarCommand) (userdomain.RegistrationSession, error) {
	return userdomain.RegistrationSession{}, nil
}

func (userSvcStub) CompleteRegistration(context.Context, string) (userdomain.User, error) {
	return userdomain.User{}, nil
}

func (userSvcStub) UploadAvatar(context.Context, userinput.UploadAvatarCommand) (string, string, int64, error) {
	return "", "", 0, nil
}

func (userSvcStub) GetByID(context.Context, uint64) (userdomain.User, error) {
	return userdomain.User{}, nil
}

func (userSvcStub) GetUserByUsername(context.Context, string) (userdomain.User, error) {
	return userdomain.User{}, nil
}
func (userSvcStub) ListAll(context.Context) ([]userdomain.User, error) { return nil, nil }
func (userSvcStub) GetUserStats(context.Context, uint64) (userdomain.UserStats, error) {
	return userdomain.UserStats{
		TotalRecords:     10,
		TotalCategories:  3,
		TotalTags:        4,
		RecordsThisWeek:  2,
		RecordsThisMonth: 5,
		MostUsedCategory: &userdomain.UsageCount{ID: 1, Name: "cat", Count: 7},
		MostUsedTag:      &userdomain.UsageCount{ID: 2, Name: "tag", Count: 8},
	}, nil
}

func (userSvcStub) UpdateUser(context.Context, uint64, userinput.UpdateUserCommand) (userdomain.User, error) {
	return userdomain.User{}, nil
}

func (userSvcStub) RemoveAvatar(context.Context, uint64) (userdomain.User, error) {
	return userdomain.User{}, nil
}

func (userSvcStub) UpdatePassword(context.Context, uint64, string, string) (string, error) {
	return "", nil
}
func (userSvcStub) SoftDeleteUser(context.Context, uint64) error { return nil }

func newResolver() *Resolver {
	return &Resolver{
		CategoryService: categorySvcStub{},
		TagService:      tagSvcStub{},
		RecordService:   recordSvcStub{},
		ChatService:     chatSvcStub{},
		UserService:     userSvcStub{},
		Logger:          gqlLoggerStub{},
	}
}

func userCtx(t *testing.T) context.Context {
	return context.WithValue(t.Context(), ctxkeys.UserID, uint64(1))
}

func TestRootResolversAndControllerFactories(t *testing.T) {
	r := newResolver()
	require.NotNil(t, r.Mutation())
	require.NotNil(t, r.Query())
	require.NotNil(t, r.CategoryController())
	require.NotNil(t, r.TagController())
	require.NotNil(t, r.RecordController())
	require.NotNil(t, r.ChatController())

	m, ok := r.Mutation().(*mutationResolver)
	require.True(t, ok)
	q, ok := r.Query().(*queryResolver)
	require.True(t, ok)
	_, err := m.Empty(t.Context())
	require.Error(t, err)
	_, err = q.Empty(t.Context())
	require.Error(t, err)
}

func TestCategoryAndTagResolvers(t *testing.T) {
	r := newResolver()
	m, ok := r.Mutation().(*mutationResolver)
	require.True(t, ok)
	q, ok := r.Query().(*queryResolver)
	require.True(t, ok)
	ctx := userCtx(t)

	_, err := m.CreateCategory(ctx, gmodel.CreateCategoryInput{Name: "n"})
	require.NoError(t, err)
	_, err = m.UpdateCategory(ctx, gmodel.UpdateCategoryInput{ID: "1"})
	require.NoError(t, err)
	deleted, err := m.SoftDeleteCategory(ctx, gmodel.DeleteCategoryInput{ID: "1"})
	require.NoError(t, err)
	require.True(t, deleted)
	_, err = m.SoftDeleteCategory(ctx, gmodel.DeleteCategoryInput{ID: "bad"})
	require.Error(t, err)

	_, err = q.Categories(ctx)
	require.NoError(t, err)
	_, err = q.CategoryByID(ctx, "1")
	require.NoError(t, err)
	_, err = q.CategoryByID(ctx, "bad")
	require.Error(t, err)
	_, err = q.CategoryByName(ctx, "x")
	require.NoError(t, err)

	_, err = m.CreateTag(ctx, gmodel.CreateTagInput{Name: "t", CategoryID: "1"})
	require.NoError(t, err)
	_, err = m.UpdateTag(ctx, gmodel.UpdateTagInput{ID: "1"})
	require.NoError(t, err)
	deleted, err = m.SoftDeleteTag(ctx, gmodel.DeleteTagInput{ID: "1"})
	require.NoError(t, err)
	require.True(t, deleted)
	_, err = m.SoftDeleteTag(ctx, gmodel.DeleteTagInput{ID: "bad"})
	require.Error(t, err)

	_, err = q.TagByName(ctx, "x")
	require.NoError(t, err)
	_, err = q.TagByID(ctx, "1")
	require.NoError(t, err)
	_, err = q.TagByID(ctx, "bad")
	require.Error(t, err)
	_, err = q.TagsByCategoryID(ctx, "1")
	require.NoError(t, err)
	_, err = q.TagsByCategoryID(ctx, "bad")
	require.Error(t, err)
	_, err = q.Tags(ctx)
	require.NoError(t, err)
}

func TestRecordResolvers_SuccessPaths(t *testing.T) {
	r := newResolver()
	m, ok := r.Mutation().(*mutationResolver)
	require.True(t, ok)
	q, ok := r.Query().(*queryResolver)
	require.True(t, ok)
	ctx := userCtx(t)
	limit := int32(5)
	date := "2026-02-14"
	until := "2026-02-14T00:00:00Z"
	start := "2026-02-01T00:00:00Z"
	end := "2026-02-14T00:00:00Z"
	after := "1"
	bad := "bad"

	_, err := m.CreateRecord(ctx, gmodel.CreateRecordInput{TagID: "1"})
	require.NoError(t, err)
	_, err = q.RecordByID(ctx, "1")
	require.NoError(t, err)
	_, err = q.RecordsByTag(ctx, "1", &limit)
	require.NoError(t, err)
	_, err = q.RecordsByDay(ctx, &date)
	require.NoError(t, err)
	_, err = q.RecordsUntil(ctx, until, &limit)
	require.NoError(t, err)
	_, err = q.RecordsBetween(ctx, start, end, &limit)
	require.NoError(t, err)
	_, err = q.Records(ctx, &limit, &date, &after)
	require.NoError(t, err)
	_, err = q.Records(ctx, nil, nil, &bad)
	require.NoError(t, err)
	_, err = q.RecordsLatest(ctx, &limit)
	require.NoError(t, err)
	_, err = q.RecordsByCategory(ctx, "1", &limit)
	require.NoError(t, err)
	_, err = m.UpdateRecord(ctx, gmodel.UpdateRecordInput{ID: "1"})
	require.NoError(t, err)
	deleted, err := m.SoftDeleteRecord(ctx, gmodel.DeleteRecordInput{ID: "1"})
	require.NoError(t, err)
	require.True(t, deleted)
	deleted, err = m.SoftDeleteAllRecords(ctx)
	require.NoError(t, err)
	require.True(t, deleted)
	_, err = q.SearchRecords(ctx, gmodel.SearchFilters{Query: "q"})
	require.NoError(t, err)
	_, err = q.RecordStats(ctx, nil)
	require.NoError(t, err)
}

func TestRecordResolvers_ParseErrors(t *testing.T) {
	r := newResolver()
	m, ok := r.Mutation().(*mutationResolver)
	require.True(t, ok)
	q, ok := r.Query().(*queryResolver)
	require.True(t, ok)
	ctx := userCtx(t)
	bad := "bad"

	_, err := q.RecordByID(ctx, bad)
	require.Error(t, err)
	_, err = q.RecordsByTag(ctx, bad, nil)
	require.Error(t, err)
	_, err = q.RecordsByCategory(ctx, bad, nil)
	require.Error(t, err)
	_, err = m.SoftDeleteRecord(ctx, gmodel.DeleteRecordInput{ID: bad})
	require.Error(t, err)
}

func TestChatResolversAndUserStats(t *testing.T) {
	r := newResolver()
	q, ok := r.Query().(*queryResolver)
	require.True(t, ok)
	ctx := userCtx(t)
	limit := int32(3)
	offset := int32(1)

	_, err := q.ChatHistory(ctx, &limit, &offset)
	require.NoError(t, err)
	_, err = q.ChatHistory(ctx, nil, nil)
	require.NoError(t, err)
	_, err = q.ChatContext(ctx)
	require.NoError(t, err)
	_, err = q.ChatDataPack(ctx, &limit, false)
	require.NoError(t, err)
	_, err = q.ChatDataPack(ctx, &limit, true)
	require.NoError(t, err)

	stats, err := q.UserStats(ctx)
	require.NoError(t, err)
	require.NotNil(t, stats)
	require.NotNil(t, stats.MostUsedCategory)
	require.NotNil(t, stats.MostUsedTag)

	r.UserService = userSvcErrStub{}
	_, err = q.UserStats(ctx)
	require.Error(t, err)
}

type userSvcErrStub struct{ userSvcStub }

func (userSvcErrStub) GetUserStats(context.Context, uint64) (userdomain.UserStats, error) {
	return userdomain.UserStats{}, errors.New("stats failed")
}

func TestNewGraphqlHandler(t *testing.T) {
	h, err := NewGraphqlHandler(nil, categorySvcStub{}, tagSvcStub{}, recordSvcStub{}, chatSvcStub{}, userSvcStub{}, gqlLoggerStub{}, &config.Config{})
	require.NoError(t, err)
	require.NotNil(t, h)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	require.NotEqual(t, 500, rw.Code)
}

var (
	_ catinput.CategoryService  = categorySvcStub{}
	_ taginput.TagService       = tagSvcStub{}
	_ recordinput.RecordService = recordSvcStub{}
	_ chatinput.ChatService     = chatSvcStub{}
	_ userinput.UserService     = userSvcStub{}
)
