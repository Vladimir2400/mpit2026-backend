package handler

import (
	"net/http"
	"strconv"

	"github.com/gdugdh24/mpit2026-backend/internal/usecase/bigfive"
	"github.com/gin-gonic/gin"
)

type BigFiveHandler struct {
	bigFiveUseCase *bigfive.BigFiveUseCase
}

func NewBigFiveHandler(bigFiveUseCase *bigfive.BigFiveUseCase) *BigFiveHandler {
	return &BigFiveHandler{
		bigFiveUseCase: bigFiveUseCase,
	}
}

// GetQuestions handles GET /big-five/questions
// @Summary Get TIPI test questions
// @Description Get all 10 TIPI personality test questions
// @Tags big-five
// @Produce json
// @Success 200 {array} bigfive.TIPIQuestion
// @Router /big-five/questions [get]
func (h *BigFiveHandler) GetQuestions(c *gin.Context) {
	questions := h.bigFiveUseCase.GetQuestions()
	c.JSON(http.StatusOK, gin.H{
		"questions": questions,
		"instruction": "Оцените, насколько каждое утверждение описывает вас, по шкале от 1 (совершенно не согласен) до 7 (полностью согласен)",
	})
}

// SubmitAnswers handles POST /big-five/submit
// @Summary Submit TIPI test answers
// @Description Submit answers to TIPI test and calculate Big Five traits
// @Tags big-five
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body bigfive.TIPIAnswersRequest true "TIPI answers"
// @Success 201 {object} domain.BigFiveResult
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /big-five/submit [post]
func (h *BigFiveHandler) SubmitAnswers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "unauthorized",
		})
		return
	}

	var req bigfive.TIPIAnswersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	result, err := h.bigFiveUseCase.SubmitAnswers(c.Request.Context(), userID.(int), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := err.Error()

		switch err.Error() {
		case "test already completed":
			statusCode = http.StatusConflict
		case "must answer all 10 questions", "invalid question id", "scores must be between 1 and 7":
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ErrorResponse{
			Error: message,
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetMyResults handles GET /big-five/my-results
// @Summary Get my Big Five results
// @Description Get current user's Big Five test results
// @Tags big-five
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.BigFiveResult
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /big-five/my-results [get]
func (h *BigFiveHandler) GetMyResults(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "unauthorized",
		})
		return
	}

	result, err := h.bigFiveUseCase.GetMyResults(c.Request.Context(), userID.(int))
	if err != nil {
		if err.Error() == "test not completed yet" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "test not completed yet",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to get results",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserResults handles GET /big-five/user/:user_id
// @Summary Get user Big Five results
// @Description Get another user's Big Five test results
// @Tags big-five
// @Security BearerAuth
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} domain.BigFiveResult
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /big-five/user/{user_id} [get]
func (h *BigFiveHandler) GetUserResults(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "unauthorized",
		})
		return
	}

	targetUserIDStr := c.Param("user_id")
	targetUserID, err := strconv.Atoi(targetUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid user_id",
		})
		return
	}

	result, err := h.bigFiveUseCase.GetUserResults(c.Request.Context(), targetUserID)
	if err != nil {
		if err.Error() == "test results not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "test results not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to get results",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
