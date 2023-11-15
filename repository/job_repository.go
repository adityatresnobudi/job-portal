package repository

import (
	"context"
	"time"

	"github.com/adityatresnobudi/job-portal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type jobRepository struct {
	db *gorm.DB
}

type JobRepository interface {
	FindAll(ctx context.Context, name string) ([]model.Jobs, error)
	FindById(ctx context.Context, jobId int) (model.Jobs, error)
	Create(ctx context.Context, newJob model.Jobs) (model.Jobs, error)
	Delete(ctx context.Context, job model.Jobs) (model.Jobs, error)
	UpdateQuota(ctx context.Context, job model.Jobs, quota int) (model.Jobs, error)
	UpdateExpDate(ctx context.Context, job model.Jobs, expDate time.Time) (model.Jobs, error)
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{
		db: db,
	}
}

func (j *jobRepository) FindAll(ctx context.Context, name string) ([]model.Jobs, error) {
	jobs := []model.Jobs{}

	err := j.db.WithContext(ctx).
		Model(&model.Jobs{}).
		Where("job_name ILIKE ? AND expiry_date > NOW() AND is_open IS TRUE", "%"+name+"%").
		Find(&jobs).Error
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (j *jobRepository) FindById(ctx context.Context, jobId int) (model.Jobs, error) {
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

func (j *jobRepository) Create(ctx context.Context, newJob model.Jobs) (model.Jobs, error) {
	err := j.db.WithContext(ctx).Model(&model.Jobs{}).Create(&newJob).Error
	if err != nil {
		return model.Jobs{}, err
	}

	return newJob, nil
}

func (j *jobRepository) Delete(ctx context.Context, job model.Jobs) (model.Jobs, error) {
	err := j.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		job := model.Jobs{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&job, job.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&job).Update("is_open", false).Error; err != nil {
			return err
		}
		return nil
	})
	return job, err
}

func (j *jobRepository) UpdateQuota(ctx context.Context, job model.Jobs, quota int) (model.Jobs, error) {
	err := j.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&model.Jobs{}, job.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Jobs{}).Where("id = ?", job.ID).Update("quota", quota).Error; err != nil {
			return err
		}
		return nil
	})
	return job, err
}

func (j *jobRepository) UpdateExpDate(ctx context.Context, job model.Jobs, expDate time.Time) (model.Jobs, error) {
	err := j.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&model.Jobs{}, job.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Jobs{}).Where("id = ?", job.ID).Update("expiry_date", expDate).Error; err != nil {
			return err
		}
		return nil
	})
	return job, err
}
