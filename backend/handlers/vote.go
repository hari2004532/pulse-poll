package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/harih/pulse-poll/backend/services"
)

type VoteHandler struct {
	Service *services.VoteService
}

type VoteRequest struct {
	OptionID string `json:"optionId"`
}

func NewVoteHandler(service *services.VoteService) *VoteHandler {
	return &VoteHandler{
		Service: service,
	}
}

func (h *VoteHandler) Vote(c *gin.Context) {
	pollID := c.Param("id")

	var req VoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if req.OptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "optionId is required",
		})
		return
	}

	voterID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	results, err := h.Service.Vote(
		c.Request.Context(),
		pollID,
		req.OptionID,
		voterID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "vote recorded",
		"results": results,
	})
}
