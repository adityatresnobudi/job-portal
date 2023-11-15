package repository

import (
	"context"

	"github.com/adityatresnobudi/job-portal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userJobRepository struct {
	db *gorm.DB
}

type UserJobRepository interface {
	FindById(ctx context.Context, jobId int) (model.Jobs, error)
	Create(ctx context.Context, newApply model.UserJobs) (model.UserJobs, error)
	UpdateMinusOneQuota(ctx context.Context, job model.Jobs) (model.Jobs, error)
	FindByJobIdUserId(ctx context.Context, jobId int, userId int) ([]model.UserJobs, error)
}

func NewUserJobRepository(db *gorm.DB) UserJobRepository {
	return &userJobRepository{
		db: db,
	}
}

func (j *userJobRepository) FindById(ctx context.Context, jobId int) (model.Jobs, error) {
	job := model.Jobs{}

	err := j.db.WithContext(ctx).
		Model(&model.Jobs{}).
		Where("id = ? AND expiry_date > NOW() AND is_open IS TRUE", jobId).
		First(&job).Error
	if err != nil {
		return model.Jobs{}, err
	}

	return job, nil
}

func (uj *userJobRepository) FindByJobIdUserId(ctx context.Context, jobId int, userId int) ([]model.UserJobs, error) {
	jobs := []model.UserJobs{}

	err := uj.db.WithContext(ctx).
		Model(&model.UserJobs{}).
		Where("job_id = ? AND user_id = ?", jobId, userId).
		Find(&jobs).Error
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (uj *userJobRepository) Create(ctx context.Context, newApply model.UserJobs) (model.UserJobs, error) {
	err := uj.db.WithContext(ctx).Model(&model.UserJobs{}).Create(&newApply).Error
	if err != nil {
		return model.UserJobs{}, err
	}

	return newApply, nil
}

func (uj *userJobRepository) UpdateMinusOneQuota(ctx context.Context, job model.Jobs) (model.Jobs, error) {
	err := uj.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&model.Jobs{}, job.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Jobs{}).Where("id = ?", job.ID).Update("quota", gorm.Expr("quota-1")).Error; err != nil {
			return err
		}
		return nil
	})
	return job, err
}
