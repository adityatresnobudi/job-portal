package usecase

import (
	"context"
	"time"

	"github.com/adityatresnobudi/job-portal/dto"
	"github.com/adityatresnobudi/job-portal/model"
	"github.com/adityatresnobudi/job-portal/repository"
	"github.com/adityatresnobudi/job-portal/shared"
)

type jobUsecase struct {
	jobRepo repository.JobRepository
}

type JobUsecase interface {
	GetAvailableJobs(ctx context.Context, name string) ([]dto.JobsDTO, error)
	GetJobsByID(ctx context.Context, jobId int) (dto.CloseJobsResponse, error)
	CreateJobs(ctx context.Context, newJob dto.JobsPayload, jobPosterId uint) (dto.JobsResponse, error)
	CloseJob(ctx context.Context, closeJob dto.CloseJobsResponse, jobPosterId uint) (dto.CloseJobsResponse, error)
	UpdateQuota(ctx context.Context, updateJob dto.CloseJobsResponse, quota int, jobPosterId uint) (dto.CloseJobsResponse, error)
	UpdateExpDate(ctx context.Context, updateJob dto.CloseJobsResponse, expDate string, jobPosterId uint) (dto.CloseJobsResponse, error)
}

func NewJobUsecase(jobRepo repository.JobRepository) JobUsecase {
	return &jobUsecase{
		jobRepo: jobRepo,
	}
}

func (ju *jobUsecase) GetAvailableJobs(ctx context.Context, name string) ([]dto.JobsDTO, error) {
	jobs := []dto.JobsDTO{}
	job := dto.JobsDTO{}
	jobList, err := ju.jobRepo.FindAll(ctx, name)
	if err != nil {
		return nil, shared.ErrGettingJobs
	}

	for _, j := range jobList {
		job.ID = j.ID
		job.JobPosterId = j.JobPosterId
		job.JobName = j.JobName
		job.JobDesc = j.JobDesc
		job.Quota = j.Quota
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (ju *jobUsecase) GetJobsByID(ctx context.Context, jobId int) (dto.CloseJobsResponse, error) {
	closeJob := dto.CloseJobsResponse{}
	cj, err := ju.jobRepo.FindById(ctx, jobId)
	if err != nil {
		return dto.CloseJobsResponse{}, shared.ErrGettingJobs
	}
	if cj.ID == 0 {
		return dto.CloseJobsResponse{}, shared.ErrJobNotFound
	}

	closeJob.ID = cj.ID
	closeJob.JobPosterId = cj.JobPosterId
	closeJob.JobName = cj.JobName
	closeJob.Quota = cj.Quota
	closeJob.IsOpen = cj.IsOpen
	closeJob.ExpiryDate = TimeToStrConv(cj.ExpiryDate)

	return closeJob, nil
}

func (ju *jobUsecase) CreateJobs(ctx context.Context, newJob dto.JobsPayload, jobPosterId uint) (dto.JobsResponse, error) {
	if newJob.JobPosterId != jobPosterId {
		return dto.JobsResponse{}, shared.ErrUnauthorized
	}

	job := model.Jobs{
		ID:          newJob.ID,
		JobPosterId: newJob.JobPosterId,
		JobName:     newJob.JobName,
		JobDesc:     newJob.JobDesc,
		Quota:       newJob.Quota,
		IsOpen:      true,
		ExpiryDate:  StrToTimeConv(newJob.ExpiryDate),
	}

	modelJob, err := ju.jobRepo.Create(ctx, job)
	if err != nil {
		return dto.JobsResponse{}, shared.ErrCreatingJobs
	}

	response := dto.JobsResponse{
		ID:          modelJob.ID,
		JobPosterId: modelJob.JobPosterId,
		JobName:     modelJob.JobName,
		JobDesc:     modelJob.JobDesc,
		Quota:       modelJob.Quota,
		ExpiryDate:  TimeToStrConv(modelJob.ExpiryDate),
	}

	return response, nil
}

func (ju *jobUsecase) CloseJob(ctx context.Context, closeJob dto.CloseJobsResponse, jobPosterId uint) (dto.CloseJobsResponse, error) {
	if closeJob.JobPosterId != jobPosterId {
		return dto.CloseJobsResponse{}, shared.ErrUnauthorized
	}

	modelJob := model.Jobs{
		ID:          closeJob.ID,
		JobPosterId: closeJob.JobPosterId,
		JobName:     closeJob.JobName,
		JobDesc:     closeJob.JobDesc,
		Quota:       closeJob.Quota,
		ExpiryDate:  StrToTimeConv(closeJob.ExpiryDate),
	}

	job, err := ju.jobRepo.Delete(ctx, modelJob)
	if err != nil {
		return dto.CloseJobsResponse{}, shared.ErrFindingJobs
	}

	response := dto.CloseJobsResponse{
		ID:          job.ID,
		JobPosterId: job.JobPosterId,
		JobName:     job.JobName,
		JobDesc:     job.JobDesc,
		Quota:       job.Quota,
		IsOpen:      job.IsOpen,
		ExpiryDate:  TimeToStrConv(job.ExpiryDate),
	}

	return response, nil
}

func (ju *jobUsecase) UpdateQuota(ctx context.Context, updateJob dto.CloseJobsResponse, quota int, jobPosterId uint) (dto.CloseJobsResponse, error) {
	if updateJob.JobPosterId != jobPosterId {
		return dto.CloseJobsResponse{}, shared.ErrUnauthorized
	}

	modelJob := model.Jobs{
		ID:          updateJob.ID,
		JobPosterId: updateJob.JobPosterId,
		JobName:     updateJob.JobName,
		JobDesc:     updateJob.JobDesc,
		Quota:       updateJob.Quota,
		ExpiryDate:  StrToTimeConv(updateJob.ExpiryDate),
	}

	if quota < 0 {
		return dto.CloseJobsResponse{}, shared.ErrMinusQuota
	}

	job, err := ju.jobRepo.UpdateQuota(ctx, modelJob, quota)
	if err != nil {
		return dto.CloseJobsResponse{}, shared.ErrFindingJobs
	}

	response := dto.CloseJobsResponse{
		ID:          job.ID,
		JobPosterId: job.JobPosterId,
		JobName:     job.JobName,
		JobDesc:     job.JobDesc,
		Quota:       job.Quota,
		IsOpen:      job.IsOpen,
		ExpiryDate:  TimeToStrConv(job.ExpiryDate),
	}

	return response, nil
}

func (ju *jobUsecase) UpdateExpDate(ctx context.Context, updateJob dto.CloseJobsResponse, expDate string, jobPosterId uint) (dto.CloseJobsResponse, error) {
	if updateJob.JobPosterId != jobPosterId {
		return dto.CloseJobsResponse{}, shared.ErrUnauthorized
	}

	modelJob := model.Jobs{
		ID:          updateJob.ID,
		JobPosterId: updateJob.JobPosterId,
		JobName:     updateJob.JobName,
		JobDesc:     updateJob.JobDesc,
		Quota:       updateJob.Quota,
		ExpiryDate:  StrToTimeConv(updateJob.ExpiryDate),
	}

	job, err := ju.jobRepo.UpdateExpDate(ctx, modelJob, StrToTimeConv(expDate))
	if err != nil {
		return dto.CloseJobsResponse{}, shared.ErrFindingJobs
	}

	response := dto.CloseJobsResponse{
		ID:          job.ID,
		JobPosterId: job.JobPosterId,
		JobName:     job.JobName,
		JobDesc:     job.JobDesc,
		Quota:       job.Quota,
		IsOpen:      job.IsOpen,
		ExpiryDate:  TimeToStrConv(job.ExpiryDate),
	}

	return response, nil
}

func TimeToStrConv(dateTime time.Time) string {
	return dateTime.Format("2006-01-02 15:04:05")
}

func StrToTimeConv(dateString string) time.Time {
	layoutFormat := "2006-01-02 15:04:05"
	date, _ := time.Parse(layoutFormat, dateString)
	return date
}
