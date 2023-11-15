package model

import "time"

type Jobs struct {
	ID          uint      `gorm:"primary_key;column:id"`
	JobPosterId uint      `gorm:"column:job_poster_id"`
	JobPoster   Users     `gorm:"foreignKey:JobPosterId"`
	JobName     string    `gorm:"column:job_name"`
	JobDesc     string    `gorm:"column:job_desc"`
	Quota       int       `gorm:"column:quota"`
	IsOpen      bool      `gorm:"column:is_open"`
	ExpiryDate  time.Time `gorm:"column:expiry_date"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"-"`
	DeletedAt   time.Time `gorm:"column:deleted_at" json:"-"`
}
