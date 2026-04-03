package repository

import (
	"context"
	"errors"

	dbport "github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"gorm.io/gorm"
)

// ListDashboardViews returns all dashboard views for a user ordered by default flag and ID.
func (r *RecordRepository) ListDashboardViews(ctx context.Context, userID uint64) ([]domain.DashboardView, error) {
	var rows []model.DashboardView
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, id ASC").
		Find(&rows).Error(); err != nil {
		return nil, err
	}

	out := make([]domain.DashboardView, len(rows))
	for i := range rows {
		out[i] = mapper.DashboardViewFromDB(rows[i])
	}
	return out, nil
}

// GetDashboardView returns one dashboard view and its widgets for the given user.
func (r *RecordRepository) GetDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error) {
	var row model.DashboardView
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", viewID, userID).
		First(&row).Error(); err != nil {
		return domain.DashboardView{}, err
	}

	view := mapper.DashboardViewFromDB(row)
	widgets, err := r.ListDashboardWidgetsByView(ctx, userID, viewID)
	if err != nil {
		return domain.DashboardView{}, err
	}
	view.Widgets = widgets
	return view, nil
}

// CreateDashboardView persists a new dashboard view.
func (r *RecordRepository) CreateDashboardView(ctx context.Context, view domain.DashboardView) (domain.DashboardView, error) {
	row := mapper.DashboardViewToDB(view)
	if err := r.db.WithContext(ctx).Create(&row).Error(); err != nil {
		return domain.DashboardView{}, err
	}
	return mapper.DashboardViewFromDB(row), nil
}

// UpdateDashboardView renames one dashboard view.
func (r *RecordRepository) UpdateDashboardView(ctx context.Context, userID uint64, viewID uint64, name string) (domain.DashboardView, error) {
	if err := r.db.WithContext(ctx).
		Model(&model.DashboardView{}).
		Where("id = ? AND user_id = ?", viewID, userID).
		Update("name", name).Error(); err != nil {
		return domain.DashboardView{}, err
	}

	var row model.DashboardView
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", viewID, userID).
		First(&row).Error(); err != nil {
		return domain.DashboardView{}, err
	}
	return mapper.DashboardViewFromDB(row), nil
}

// SetDefaultDashboardView marks one dashboard view as default for the user.
func (r *RecordRepository) SetDefaultDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx dbport.DB) error {
		if err := tx.Model(&model.DashboardView{}).
			Where("user_id = ?", userID).
			Update("is_default", false).Error(); err != nil {
			return err
		}

		if err := tx.Model(&model.DashboardView{}).
			Where("id = ? AND user_id = ?", viewID, userID).
			Update("is_default", true).Error(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.DashboardView{}, err
	}

	var row model.DashboardView
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", viewID, userID).
		First(&row).Error(); err != nil {
		return domain.DashboardView{}, err
	}
	return mapper.DashboardViewFromDB(row), nil
}

// DeleteDashboardView removes one dashboard view and its widgets, promoting another default if needed.
func (r *RecordRepository) DeleteDashboardView(ctx context.Context, userID uint64, viewID uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx dbport.DB) error {
		var current model.DashboardView
		if err := tx.
			Where("id = ? AND user_id = ?", viewID, userID).
			First(&current).Error(); err != nil {
			return err
		}

		if err := tx.
			Where("user_id = ? AND view_id = ?", userID, viewID).
			Delete(&model.DashboardWidget{}).Error(); err != nil {
			return err
		}

		if err := tx.
			Where("id = ? AND user_id = ?", viewID, userID).
			Delete(&model.DashboardView{}).Error(); err != nil {
			return err
		}

		if !current.IsDefault {
			return nil
		}

		var fallback model.DashboardView
		if err := tx.
			Where("user_id = ?", userID).
			Order("id ASC").
			First(&fallback).Error(); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		return tx.
			Model(&model.DashboardView{}).
			Where("id = ? AND user_id = ?", fallback.ID, userID).
			Update("is_default", true).Error()
	})
}

// UpsertDashboardWidget creates or updates a dashboard widget.
func (r *RecordRepository) UpsertDashboardWidget(ctx context.Context, widget domain.DashboardWidget) (domain.DashboardWidget, error) {
	row := mapper.DashboardWidgetToDB(widget)
	if row.ConfigJSON == "" {
		row.ConfigJSON = "{}"
	}

	if row.ID != 0 {
		updateMap := map[string]interface{}{
			"view_id":        row.ViewID,
			"widget_type":    row.WidgetType,
			"size":           row.Size,
			"order_index":    row.OrderIndex,
			"title_override": row.TitleOverride,
			"config_json":    row.ConfigJSON,
			"is_active":      row.IsActive,
		}
		if row.MetricDefinitionID != nil {
			updateMap["metric_definition_id"] = *row.MetricDefinitionID
		} else {
			updateMap["metric_definition_id"] = nil
		}
		if err := r.db.WithContext(ctx).
			Model(&model.DashboardWidget{}).
			Where("id = ? AND user_id = ?", row.ID, row.UserID).
			Updates(updateMap).Error(); err != nil {
			return domain.DashboardWidget{}, err
		}
		if err := r.db.WithContext(ctx).
			Where("id = ? AND user_id = ?", row.ID, row.UserID).
			First(&row).Error(); err != nil {
			return domain.DashboardWidget{}, err
		}
		return mapper.DashboardWidgetFromDB(row), nil
	}

	if err := r.db.WithContext(ctx).Create(&row).Error(); err != nil {
		return domain.DashboardWidget{}, err
	}
	return mapper.DashboardWidgetFromDB(row), nil
}

// ListDashboardWidgetsByView returns active widgets for a specific dashboard view.
func (r *RecordRepository) ListDashboardWidgetsByView(ctx context.Context, userID uint64, viewID uint64) ([]domain.DashboardWidget, error) {
	var rows []model.DashboardWidget
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND view_id = ? AND is_active = ?", userID, viewID, true).
		Order("order_index ASC, id ASC").
		Find(&rows).Error(); err != nil {
		return nil, err
	}

	out := make([]domain.DashboardWidget, len(rows))
	for i := range rows {
		out[i] = mapper.DashboardWidgetFromDB(rows[i])
	}
	return out, nil
}

// ReorderDashboardWidgets updates the order index for widgets in a dashboard view.
func (r *RecordRepository) ReorderDashboardWidgets(ctx context.Context, userID uint64, viewID uint64, items []domain.DashboardWidget) ([]domain.DashboardWidget, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx dbport.DB) error {
		for _, item := range items {
			if err := tx.Model(&model.DashboardWidget{}).
				Where("id = ? AND user_id = ? AND view_id = ?", item.ID, userID, viewID).
				Update("order_index", item.OrderIndex).Error(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.ListDashboardWidgetsByView(ctx, userID, viewID)
}

// DeleteDashboardWidget performs a soft delete by marking a widget inactive.
func (r *RecordRepository) DeleteDashboardWidget(ctx context.Context, userID uint64, widgetID uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.DashboardWidget{}).
		Where("id = ? AND user_id = ?", widgetID, userID).
		Update("is_active", false).Error()
}

// CountLargeWidgetsInView counts active large widgets in a view, optionally excluding one widget.
func (r *RecordRepository) CountLargeWidgetsInView(ctx context.Context, userID uint64, viewID uint64, excludeWidgetID *uint64) (int64, error) {
	q := r.db.WithContext(ctx).
		Model(&model.DashboardWidget{}).
		Where("user_id = ? AND view_id = ? AND is_active = ? AND size = ?", userID, viewID, true, domain.DashboardWidgetSizeLarge)
	if excludeWidgetID != nil && *excludeWidgetID != 0 {
		q = q.Where("id <> ?", *excludeWidgetID)
	}

	var count int64
	if err := q.Count(&count).Error(); err != nil {
		return 0, err
	}
	return count, nil
}
