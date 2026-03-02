package models

import (
	"github.com/google/uuid"
)

type ProjectMetric struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProjectID   uuid.UUID `gorm:"type:uuid;not null" json:"project_id"`
	MetricKey   string    `gorm:"type:varchar(100)" json:"metric_key"`
	MetricLabel string    `gorm:"type:varchar(255)" json:"metric_label"`
	MetricValue float64   `gorm:"type:numeric" json:"metric_value"`
	MetricUnit  string    `gorm:"type:varchar(50)" json:"metric_unit"`
	Order       int       `gorm:"type:integer;default:0" json:"order"`
}