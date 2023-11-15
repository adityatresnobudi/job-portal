package model

import "time"

type UserJobs struct {
	ID        uint      `gorm:"primary_key;column:id"`
	JobId     uint      `gorm:"column:job_id"`
	Jobs      Jobs      `gorm:"foreignKey:JobId" json:"jobs"`
	UserId    uint      `gorm:"column:user_id"`
	Users     Users     `gorm:"foreignKey:UserId" json:"users"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	DeletedAt time.Time `gorm:"column:deleted_at" json:"-"`
}
