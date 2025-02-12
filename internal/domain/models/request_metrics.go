package models

import (
	"time"
)

type RequestMetrics struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Method    string    `json:"method" gorm:"not null;uniqueIndex:idx_method_path"`
	Path      string    `json:"path" gorm:"not null;uniqueIndex:idx_method_path"`
	Count     int64     `json:"count" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
