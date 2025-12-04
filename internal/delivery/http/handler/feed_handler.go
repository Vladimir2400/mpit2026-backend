package handler

import (
	"net/http"

	"github.com/gdugdh24/mpit2026-backend/internal/usecase/feed"
	"github.com/gin-gonic/gin"
)

type FeedHandler struct {
	feedUseCase *feed.FeedUseCase
}

func NewFeedHandler(feedUseCase *feed.FeedUseCase) *FeedHandler {
	return &FeedHandler{
		feedUseCase: feedUseCase,
	}
}

// GetNextUser handles GET /feed/next
// @Summary Get next user in feed
// @Description Get the next user to show in feed based on preferences
// @Tags feed
// @Security BearerAuth
// @Produce json
// @Success 200 {object} feed.FeedUserResponse
// @Success 204 "No more users in feed"
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /feed/next [get]
func (h *FeedHandler) GetNextUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "unauthorized",
		})
		return
	}

	user, err := h.feedUseCase.GetNextUser(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to get next user",
		})
		return
	}

	// No more users
	if user == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "no more users in feed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// ResetDislikes handles POST /feed/reset-dislikes
// @Summary Reset all dislikes
// @Description Delete all dislikes to refresh the feed
// @Tags feed
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /feed/reset-dislikes [post]
func (h *FeedHandler) ResetDislikes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "unauthorized",
		})
		return
	}

	count, err := h.feedUseCase.ResetDislikes(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to reset dislikes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "dislikes reset successfully",
		"reset_count": count,
	})
}
