package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adityatresnobudi/job-portal/dto"
	"github.com/adityatresnobudi/job-portal/handler"
	"github.com/adityatresnobudi/job-portal/mocks"
	"github.com/adityatresnobudi/job-portal/router"
	"github.com/adityatresnobudi/job-portal/shared"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createJobsDTO() dto.JobsDTO {
	return dto.JobsDTO{
		ID:          1,
		JobPosterId: 2,
		JobName:     "test",
		JobDesc:     "test",
		Quota:       3,
	}
}

func createJobsPayload() dto.JobsPayload {
	return dto.JobsPayload{
		ID:          1,
		JobPosterId: 2,
		JobName:     "test",
		JobDesc:     "test",
		Quota:       3,
		ExpiryDate:  "2023-06-14 20:00:00",
	}
}

func createJobsResponse() dto.JobsResponse {
	return dto.JobsResponse{
		ID:          1,
		JobPosterId: 2,
		JobName:     "test",
		JobDesc:     "test",
		Quota:       3,
		ExpiryDate:  "2023-06-14 20:00:00",
	}
}

func createCloseJobsResponse() dto.CloseJobsResponse {
	return dto.CloseJobsResponse{
		ID:          1,
		JobPosterId: 2,
		JobName:     "test",
		JobDesc:     "test",
		Quota:       3,
		IsOpen:      true,
		ExpiryDate:  "2023-06-14 20:00:00",
	}
}

func MakeRequestBody(dto interface{}) *strings.Reader {
	payload, _ := json.Marshal(dto)
	return strings.NewReader(string(payload))
}

func TestJobHandler_GetJobs(t *testing.T) {
	t.Setenv("ENV_MODE", "testing")
	t.Run("should return status code 200 when job fetched", func(t *testing.T) {
		// 1. setup router
		mockJobUsecase := new(mocks.JobUsecase)
		mockUserUsecase := new(mocks.UserUsecase)
		mockUserJobUsecase := new(mocks.UserJobUsecase)
		h := handler.NewHandler(mockJobUsecase, mockUserUsecase, mockUserJobUsecase)
		jobs := []dto.JobsDTO{
			createJobsDTO(),
		}

		// 2. make request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("GET", "/jobs", nil)
		c.Request = req

		mockJobUsecase.On("GetAvailableJobs", c.Request.Context(), "").Return(jobs, nil)
		expectedResp, _ := json.Marshal(dto.JsonResponse{Data: jobs})
		h.GetJobs(c)

		// 3. assert
		assert.Equal(t, http.StatusOK, w.Code)
		str := strings.Trim(w.Body.String(), "\n")
		assert.Equal(t, string(expectedResp), str)
	})

	t.Run("should return status code 500 when job fetch failed", func(t *testing.T) {
		// 1. setup router
		mockJobUsecase := new(mocks.JobUsecase)
		mockUserUsecase := new(mocks.UserUsecase)
		mockUserJobUsecase := new(mocks.UserJobUsecase)
		h := handler.NewHandler(mockJobUsecase, mockUserUsecase, mockUserJobUsecase)
		router := router.NewRouter(h)
		mockJobUsecase.On("GetAvailableJobs", mock.Anything, "").Return(nil, shared.ErrGettingJobs)
		expectedResp, _ := json.Marshal(dto.JsonResponse{Message: shared.ErrGettingJobs.Message})

		// 2. make request
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/jobs", nil)
		router.ServeHTTP(rec, req)

		// 3. assert
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		str := strings.Trim(rec.Body.String(), "\n")
		assert.Equal(t, string(expectedResp), str)
	})
}

func TestJobHandler_CreateNewJobs(t *testing.T) {
	t.Setenv("ENV_MODE", "testing")
	t.Run("should return status code 201 when job created", func(t *testing.T) {
		// 1. setup router
		mockJobUsecase := new(mocks.JobUsecase)
		mockUserUsecase := new(mocks.UserUsecase)
		mockUserJobUsecase := new(mocks.UserJobUsecase)
		h := handler.NewHandler(mockJobUsecase, mockUserUsecase, mockUserJobUsecase)
		router.NewRouter(h)
		jobPayload := createJobsPayload()
		jobResponse := createJobsResponse()

		// 2. make request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("id", uint(1))
		req, _ := http.NewRequest("POST", "/jobs", MakeRequestBody(jobPayload))
		c.Request = req

		mockJobUsecase.On("CreateJobs", c.Request.Context(), jobPayload, uint(1)).Return(jobResponse, nil)
		expectedResp, _ := json.Marshal(dto.JsonResponse{Message: "successfully add post job with id 1", Data: jobResponse})
		h.CreateNewJobs(c)

		// 3. assert
		assert.Equal(t, http.StatusCreated, w.Code)
		str := strings.Trim(w.Body.String(), "\n")
		assert.Equal(t, string(expectedResp), str)
	})

	t.Run("should return status code 500 when job fetch failed", func(t *testing.T) {
		// 1. setup router
		mockJobUsecase := new(mocks.JobUsecase)
		mockUserUsecase := new(mocks.UserUsecase)
		mockUserJobUsecase := new(mocks.UserJobUsecase)
		h := handler.NewHandler(mockJobUsecase, mockUserUsecase, mockUserJobUsecase)
		router.NewRouter(h)
		jobPayload := createJobsPayload()

		// 2. make request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("id", uint(1))
		req, _ := http.NewRequest("POST", "/jobs", MakeRequestBody(jobPayload))
		c.Request = req

		mockJobUsecase.On("CreateJobs", c.Request.Context(), jobPayload, uint(1)).Return(dto.JobsResponse{}, shared.ErrCreatingJobs)
		h.CreateNewJobs(c)

		// 3. assert
		assert.Equal(t, c.Errors[0].Err, shared.ErrCreatingJobs)
	})

	t.Run("should return status code 500 when job fetch failed", func(t *testing.T) {
		// 1. setup router
		mockJobUsecase := new(mocks.JobUsecase)
		mockUserUsecase := new(mocks.UserUsecase)
		mockUserJobUsecase := new(mocks.UserJobUsecase)
		h := handler.NewHandler(mockJobUsecase, mockUserUsecase, mockUserJobUsecase)
		router.NewRouter(h)
		jobPayload := dto.JobsPayload{
			ID:         1,
			JobName:    "test",
			JobDesc:    "test",
			Quota:      3,
			ExpiryDate: "2023-06-14 20:00:00",
		}

		// 2. make request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("POST", "/jobs", MakeRequestBody(jobPayload))
		c.Request = req
		h.CreateNewJobs(c)

		// 3. assert
		assert.Equal(t, c.Errors[0].Err, shared.ErrInvalidRequestBody)
	})
}

func TestJobHandler_CloseJobs(t *testing.T) {
	t.Setenv("ENV_MODE", "testing")
	t.Run("should return status code 200 when job closed", func(t *testing.T) {
		// 1. setup router
		mockJobUsecase := new(mocks.JobUsecase)
		mockUserUsecase := new(mocks.UserUsecase)
		mockUserJobUsecase := new(mocks.UserJobUsecase)
		h := handler.NewHandler(mockJobUsecase, mockUserUsecase, mockUserJobUsecase)
		// router.NewRouter(h)
		jobCloseResponse := createCloseJobsResponse()

		// 2. make request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("id", uint(1))
		c.AddParam("id", "1")
		req, _ := http.NewRequest("POST", "/jobs/1/close", nil)
		c.Request = req

		mockJobUsecase.On("GetJobsByID", c.Request.Context(), 1).Return(jobCloseResponse, nil)
		mockJobUsecase.On("CloseJob", c.Request.Context(), jobCloseResponse, uint(1)).Return(jobCloseResponse, nil)
		expectedResp, _ := json.Marshal(dto.JsonResponse{Message: "successfully closed job with id 1", Data: jobCloseResponse})
		h.CloseJobs(c)

		// 3. assert
		assert.Equal(t, http.StatusOK, w.Code)
		str := strings.Trim(w.Body.String(), "\n")
		assert.Equal(t, string(expectedResp), str)
	})

	t.Run("should return status code 500 when job fetch failed", func(t *testing.T) {
		// 1. setup router
		mockJobUsecase := new(mocks.JobUsecase)
		mockUserUsecase := new(mocks.UserUsecase)
		mockUserJobUsecase := new(mocks.UserJobUsecase)
		h := handler.NewHandler(mockJobUsecase, mockUserUsecase, mockUserJobUsecase)
		router.NewRouter(h)
		jobPayload := createJobsPayload()

		// 2. make request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("id", uint(1))
		req, _ := http.NewRequest("POST", "/jobs", MakeRequestBody(jobPayload))
		c.Request = req

		mockJobUsecase.On("CreateJobs", c.Request.Context(), jobPayload, uint(1)).Return(dto.JobsResponse{}, shared.ErrCreatingJobs)
		h.CreateNewJobs(c)

		// 3. assert
		assert.Equal(t, c.Errors[0].Err, shared.ErrCreatingJobs)
	})

	t.Run("should return status code 500 when job fetch failed", func(t *testing.T) {
		// 1. setup router
		mockJobUsecase := new(mocks.JobUsecase)
		mockUserUsecase := new(mocks.UserUsecase)
		mockUserJobUsecase := new(mocks.UserJobUsecase)
		h := handler.NewHandler(mockJobUsecase, mockUserUsecase, mockUserJobUsecase)
		router.NewRouter(h)
		jobPayload := dto.JobsPayload{
			ID:         1,
			JobName:    "test",
			JobDesc:    "test",
			Quota:      3,
			ExpiryDate: "2023-06-14 20:00:00",
		}

		// 2. make request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("POST", "/jobs", MakeRequestBody(jobPayload))
		c.Request = req
		h.CreateNewJobs(c)

		// 3. assert
		assert.Equal(t, c.Errors[0].Err, shared.ErrInvalidRequestBody)
	})
}
