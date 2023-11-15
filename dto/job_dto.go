package dto

type JobsDTO struct {
	ID          uint   `json:"id"`
	JobPosterId uint   `json:"job_poster_id"`
	JobName     string `json:"job_name"`
	JobDesc     string `json:"job_desc"`
	Quota       int    `json:"quota"`
}

type JobsPayload struct {
	ID          uint   `json:"id"`
	JobPosterId uint   `json:"job_poster_id" binding:"required"`
	JobName     string `json:"job_name" binding:"required"`
	JobDesc     string `json:"job_desc" binding:"required"`
	Quota       int    `json:"quota" binding:"required"`
	ExpiryDate  string `json:"expiry_date" binding:"required"`
}

type JobsResponse struct {
	ID          uint   `json:"id"`
	JobPosterId uint   `json:"job_poster_id"`
	JobName     string `json:"job_name"`
	JobDesc     string `json:"job_desc"`
	Quota       int    `json:"quota"`
	ExpiryDate  string `json:"expiry_date"`
}

type CloseJobsResponse struct {
	ID          uint   `json:"id"`
	JobPosterId uint   `json:"job_poster_id"`
	JobName     string `json:"job_name"`
	JobDesc     string `json:"job_desc"`
	Quota       int    `json:"quota"`
	IsOpen      bool   `json:"is_open"`
	ExpiryDate  string `json:"expiry_date"`
}
