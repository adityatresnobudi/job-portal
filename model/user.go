package model

import "time"

type Users struct {
	ID          uint      `gorm:"primary_key;column:id"`
	Name        string    `gorm:"column:user_name"`
	Email       string    `gorm:"column:email"`
	Phone       string    `gorm:"column:phone"`
	Password    string    `gorm:"column:user_password"`
	CurrentJob  string    `gorm:"column:current_job"`
	Age         uint      `gorm:"column:user_age"`
	IsJobPoster bool      `gorm:"column:is_job_poster"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"-"`
	DeletedAt   time.Time `gorm:"column:deleted_at" json:"-"`
}
