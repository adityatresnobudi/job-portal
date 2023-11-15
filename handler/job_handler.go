package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/adityatresnobudi/job-portal/dto"
	"github.com/adityatresnobudi/job-portal/shared"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetJobs(c *gin.Context) {
	ctx := c.Request.Context()
	jobName := c.Query("name")
	jobs, err := h.JobUsecase.GetAvailableJobs(ctx, jobName)
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.JsonResponse{Data: jobs})
}

func (h *Handler) CreateNewJobs(c *gin.Context) {
	ctx := c.Request.Context()
	jobPosterId := c.GetUint("id")
	newJob := dto.JobsPayload{}
	if err := c.ShouldBindJSON(&newJob); err != nil {
		log.Println(err)
		c.Error(shared.ErrInvalidRequestBody)
		return
	}

	newJob.JobName = strings.TrimSpace(newJob.JobName)

	jobs, err := h.JobUsecase.CreateJobs(ctx, newJob, jobPosterId)
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}
	message := fmt.Sprintf("successfully add post job with id %d", jobs.ID)
	c.JSON(http.StatusCreated, dto.JsonResponse{Message: message, Data: jobs})
}

func (h *Handler) CloseJobs(c *gin.Context) {
	ctx := c.Request.Context()
	jobPosterId := c.GetUint("id")

	jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.Error(shared.ErrIdNotFound)
		return
	}

	cjUpdate, err := h.JobUsecase.GetJobsByID(ctx, jobId)
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}

	output, err := h.JobUsecase.CloseJob(ctx, cjUpdate, jobPosterId)
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}
	message := fmt.Sprintf("successfully closed job with id %d", output.ID)
	c.JSON(http.StatusOK, dto.JsonResponse{Message: message, Data: output})
}

func (h *Handler) ChangeJobs(c *gin.Context) {
	ctx := c.Request.Context()
	output := dto.CloseJobsResponse{}
	jobPosterId := c.GetUint("id")
	quota := c.Query("quota")
	expDate := c.Query("expDate")

	jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.Error(shared.ErrIdNotFound)
		return
	}

	q, err := strconv.Atoi(quota)
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}

	cjUpdate, err := h.JobUsecase.GetJobsByID(ctx, jobId)
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}

	if quota != "" {
		res, err := h.JobUsecase.UpdateQuota(ctx, cjUpdate, q, jobPosterId)
		if err != nil {
			log.Println(err)
			c.Error(err)
			return
		}
		output = res
	}

	if expDate != "" {
		res, err := h.JobUsecase.UpdateExpDate(ctx, cjUpdate, expDate, jobPosterId)
		if err != nil {
			log.Println(err)
			c.Error(err)
			return
		}
		output = res
	}

	message := fmt.Sprintf("successfully change job with id %d", jobId)
	c.JSON(http.StatusNoContent, dto.JsonResponse{Message: message, Data: output})
}
