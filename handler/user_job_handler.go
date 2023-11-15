package handler

import (
	"log"
	"net/http"

	"github.com/adityatresnobudi/job-portal/dto"
	"github.com/adityatresnobudi/job-portal/shared"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ApplyJob(c *gin.Context) {
	ctx := c.Request.Context()
	userId := c.GetUint("id")
	newUserJob := dto.UserJobsPayload{}

	if err := c.ShouldBindJSON(&newUserJob); err != nil {
		log.Println(err)
		c.Error(shared.ErrInvalidRequestBody)
		return
	}

	users, err := h.UserJobUsecase.ApplyJob(ctx, newUserJob, int(userId))
	if err != nil {
		log.Println(err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.JsonResponse{Data: users})
}
