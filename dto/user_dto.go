package dto

type UserPayload struct {
	ID          uint   `json:"id"`
	Name        string `json:"user_name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Password    string `json:"user_password" binding:"required"`
	CurrentJob  string `json:"current_job"`
	Age         uint   `json:"user_age"`
	IsJobPoster bool   `json:"is_job_poster" binding:"required"`
}

type UserResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"user_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	CurrentJob  string `json:"current_job,omitempty"`
	Age         uint   `json:"user_age,omitempty"`
	IsJobPoster bool   `json:"is_job_poster"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}