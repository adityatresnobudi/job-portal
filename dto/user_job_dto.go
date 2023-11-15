package dto

type UserJobsDTO struct {
	JobId uint `json:"job_id"`
	Status string `json:"status"`
	Message string `json:"message"`
	AppliedAt string `json:"applied_at"`
}

type UserJobsPayload struct {
	UserId uint `json:"user_id" binding:"required"`
	JobId uint `json:"job_id" binding:"required"`
}