package handler

import (
	"net/http"
	"strconv"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
	"github.com/gdugdh24/mpit2026-backend/internal/usecase/swipe"
	"github.com/gin-gonic/gin"
)

type SwipeHandler struct {
	swipeUseCase *swipe.SwipeUseCase
}

func NewSwipeHandler(swipeUseCase *swipe.SwipeUseCase) *SwipeHandler {
	return &SwipeHandler{
		swipeUseCase: swipeUseCase,
	}
}

// CreateSwipe handles POST /swipe
// @Summary Create a swipe (like/dislike)
// @Description Swipe on a user and check if it's a match
// @Tags swipe
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body swipe.SwipeRequest true "Swipe data"
// @Success 200 {object} swipe.SwipeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /swipe [post]
func (h *SwipeHandler) CreateSwipe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "unauthorized",
		})
		return
	}

	var req swipe.SwipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	result, err := h.swipeUseCase.CreateSwipe(c.Request.Context(), userID.(int), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "failed to create swipe"

		switch err {
		case domain.ErrCannotSwipeSelf:
			statusCode = http.StatusBadRequest
			message = "cannot swipe yourself"
		case domain.ErrSwipeAlreadyExists:
			statusCode = http.StatusConflict
			message = "swipe already exists"
		}

		c.JSON(statusCode, ErrorResponse{
			Error: message,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetLikesReceived handles GET /swipe/likes-received
// @Summary Get likes received
// @Description Get list of users who liked current user
// @Tags swipe
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /swipe/likes-received [get]
func (h *SwipeHandler) GetLikesReceived(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "unauthorized",
		})
		return
	}

	// Parse query params
	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	likes, total, err := h.swipeUseCase.GetLikesReceived(c.Request.Context(), userID.(int), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to get likes received",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"likes": likes,
		"total": total,
	})
}
