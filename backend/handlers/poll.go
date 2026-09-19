package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/harih/pulse-poll/backend/services"
)

type PollHandler struct {
	Service *services.PollService
}

type CreatePollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

func NewPollHandler(service *services.PollService) *PollHandler {
	return &PollHandler{
		Service: service,
	}
}

func (h *PollHandler) CreatePoll(c *gin.Context) {

	var req CreatePollRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	poll, err := h.Service.CreatePoll(
		c.Request.Context(),
		req.Question,
		req.Options,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "poll created",
		"poll":    poll,
	})
}

func (h *PollHandler) GetPoll(c *gin.Context) {
	pollID := c.Param("id")

	poll, err := h.Service.GetPoll(c.Request.Context(), pollID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "poll not found",
		})
		return
	}

	results, err := h.Service.GetPollResults(c.Request.Context(), pollID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load poll results",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"poll":    poll,
		"results": results,
	})
}
