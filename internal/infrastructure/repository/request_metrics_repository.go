package repository

import (
	"github.com/guisithos/neoway/internal/domain/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RequestMetricsRepository interface {
	IncrementCount(method, path string) error
	GetAllMetrics() ([]models.RequestMetrics, error)
}

type requestMetricsRepository struct {
	db *gorm.DB
}

func NewRequestMetricsRepository(db *gorm.DB) RequestMetricsRepository {
	return &requestMetricsRepository{db: db}
}

func (r *requestMetricsRepository) IncrementCount(method, path string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		result := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "method"}, {Name: "path"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"count":      gorm.Expr("request_metrics.count + 1"),
				"updated_at": gorm.Expr("NOW()"),
			}),
		}).Create(&models.RequestMetrics{
			Method: method,
			Path:   path,
			Count:  1,
		})

		return result.Error
	})
}

func (r *requestMetricsRepository) GetAllMetrics() ([]models.RequestMetrics, error) {
	var metrics []models.RequestMetrics
	err := r.db.Find(&metrics).Error
	return metrics, err
}
