package handler

import "github.com/adityatresnobudi/job-portal/usecase"

type Handler struct {
	JobUsecase usecase.JobUsecase
	UserUsecase usecase.UserUsecase
	UserJobUsecase usecase.UserJobUsecase
}

func NewHandler(JobUsecase usecase.JobUsecase, UserUsecase usecase.UserUsecase, UserJobUsecase usecase.UserJobUsecase) *Handler {
	return &Handler{
		JobUsecase: JobUsecase,
		UserUsecase: UserUsecase,
		UserJobUsecase: UserJobUsecase,
	}
}
