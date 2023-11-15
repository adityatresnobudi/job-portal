package usecase

import (
	"context"
	"time"

	"github.com/adityatresnobudi/job-portal/dto"
	"github.com/adityatresnobudi/job-portal/model"
	"github.com/adityatresnobudi/job-portal/repository"
	"github.com/adityatresnobudi/job-portal/shared"
)

type userJobUsecase struct {
	userJobRepo repository.UserJobRepository
}

type UserJobUsecase interface {
	ApplyJob(ctx context.Context, job dto.UserJobsPayload, userId int) (dto.UserJobsDTO, error)
}

func NewUserJobUsecase(userJobRepo repository.UserJobRepository) UserJobUsecase {
	return &userJobUsecase{
		userJobRepo: userJobRepo,
	}
}

func (uj *userJobUsecase) ApplyJob(ctx context.Context, job dto.UserJobsPayload, userId int) (dto.UserJobsDTO, error) {
	if job.UserId != uint(userId) {
		return dto.UserJobsDTO{}, shared.ErrUnauthorized
	}

	jobs, err := uj.userJobRepo.FindByJobIdUserId(ctx, int(job.JobId), userId)
	if err != nil {
		return dto.UserJobsDTO{}, shared.ErrGettingUserJob
	}

	if len(jobs) != 0 {
		return dto.UserJobsDTO{}, shared.ErrAlreadyApplied
	}

	j, err := uj.userJobRepo.FindById(ctx, int(job.JobId))
	if err != nil {
		return dto.UserJobsDTO{}, shared.ErrGettingUserJob
	}

	output, err := uj.userJobRepo.UpdateMinusOneQuota(ctx, j)
	if err != nil {
		return dto.UserJobsDTO{}, shared.ErrJobTransaction
	}

	modelUserJob := model.UserJobs{
		UserId: uint(userId),
		JobId:  output.ID,
	}

	res, err := uj.userJobRepo.Create(ctx, modelUserJob)
	if err != nil {
		return dto.UserJobsDTO{}, shared.ErrCreateApplyJob
	}

	userJobRes := dto.UserJobsDTO{
		JobId:     res.JobId,
		Status:    "Applied",
		Message:   "Application success",
		AppliedAt: TimeToStrConv(time.Now()),
	}

	return userJobRes, nil
}
