package bigfive

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/gdugdh24/mpit2026-backend/internal/domain"
	"github.com/gdugdh24/mpit2026-backend/internal/repository"
)

type BigFiveUseCase struct {
	bigFiveRepo repository.BigFiveRepository
}

func NewBigFiveUseCase(bigFiveRepo repository.BigFiveRepository) *BigFiveUseCase {
	return &BigFiveUseCase{
		bigFiveRepo: bigFiveRepo,
	}
}

// TIPIQuestion represents a TIPI test question
type TIPIQuestion struct {
	ID   int               `json:"id"`
	Text string            `json:"text"`
	Traits map[string]float64 `json:"-"`
}

// TIPIQuestions - 10 вопросов TIPI теста
var TIPIQuestions = []TIPIQuestion{
	{ID: 1, Text: "Экстраверт, энергичный", Traits: map[string]float64{"E": 1.0}},
	{ID: 2, Text: "Критичный, склонный к спорам", Traits: map[string]float64{"A": -1.0}},
	{ID: 3, Text: "Надёжный, дисциплинированный", Traits: map[string]float64{"C": 1.0}},
	{ID: 4, Text: "Тревожный, легко расстраиваюсь", Traits: map[string]float64{"N": 1.0}},
	{ID: 5, Text: "Открытый новому, со сложным внутренним миром", Traits: map[string]float64{"O": 1.0}},
	{ID: 6, Text: "Сдержанный, тихий", Traits: map[string]float64{"E": -1.0}},
	{ID: 7, Text: "Отзывчивый, тёплый", Traits: map[string]float64{"A": 1.0}},
	{ID: 8, Text: "Неорганизованный, беспечный", Traits: map[string]float64{"C": -1.0}},
	{ID: 9, Text: "Спокойный, эмоционально стабильный", Traits: map[string]float64{"N": -1.0}},
	{ID: 10, Text: "Консервативный, не склонный к творчеству", Traits: map[string]float64{"O": -1.0}},
}

// TIPIAnswersRequest represents answers to TIPI test
type TIPIAnswersRequest struct {
	Answers map[int]int `json:"answers" binding:"required"` // question_id -> score (1-7)
}

// GetQuestions returns all TIPI questions
func (uc *BigFiveUseCase) GetQuestions() []TIPIQuestion {
	// Return without traits for response
	questions := make([]TIPIQuestion, len(TIPIQuestions))
	for i, q := range TIPIQuestions {
		questions[i] = TIPIQuestion{
			ID:   q.ID,
			Text: q.Text,
		}
	}
	return questions
}

// SubmitAnswers calculates Big Five traits from TIPI answers and saves to DB
func (uc *BigFiveUseCase) SubmitAnswers(ctx context.Context, userID int, req *TIPIAnswersRequest) (*domain.BigFiveResult, error) {
	// Check if test already completed
	existing, err := uc.bigFiveRepo.GetByUserID(ctx, userID)
	if err == nil && existing != nil {
		return nil, errors.New("test already completed")
	}

	// Validate answers
	if len(req.Answers) != 10 {
		return nil, errors.New("must answer all 10 questions")
	}

	for qid, score := range req.Answers {
		if qid < 1 || qid > 10 {
			return nil, errors.New("invalid question id")
		}
		if score < 1 || score > 7 {
			return nil, errors.New("scores must be between 1 and 7")
		}
	}

	// Calculate Big Five traits using TIPI algorithm
	traits := uc.calculateTIPI(req.Answers)

	// Create result
	result := &domain.BigFiveResult{
		UserID:            userID,
		Extraversion:      traits["E"],
		Agreeableness:     traits["A"],
		Conscientiousness: traits["C"],
		Neuroticism:       traits["N"],
		Openness:          traits["O"],
		CompletedAt:       time.Now(),
	}

	if err := uc.bigFiveRepo.Create(ctx, result); err != nil {
		return nil, err
	}

	return result, nil
}

// calculateTIPI implements TIPI scoring algorithm
func (uc *BigFiveUseCase) calculateTIPI(answers map[int]int) map[string]float64 {
	// TIPI uses pairs of questions for each trait
	// E: 1, 6R (R = reversed)
	// A: 2R, 7
	// C: 3, 8R
	// N: 4, 9R
	// O: 5, 10R

	traits := map[string][2]float64{
		"E": {float64(answers[1]), float64(8 - answers[6])},
		"A": {float64(8 - answers[2]), float64(answers[7])},
		"C": {float64(answers[3]), float64(8 - answers[8])},
		"N": {float64(answers[4]), float64(8 - answers[9])},
		"O": {float64(answers[5]), float64(8 - answers[10])},
	}

	result := make(map[string]float64)
	for trait, scores := range traits {
		// Average the two scores (1-7 range)
		avg := (scores[0] + scores[1]) / 2.0
		// Normalize to 0.00-1.00 range
		normalized := (avg - 1.0) / 6.0
		// Round to 2 decimal places
		result[trait] = math.Round(normalized*100) / 100
	}

	return result
}

// GetMyResults returns current user's Big Five results
func (uc *BigFiveUseCase) GetMyResults(ctx context.Context, userID int) (*domain.BigFiveResult, error) {
	result, err := uc.bigFiveRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("test not completed yet")
	}
	return result, nil
}

// GetUserResults returns another user's Big Five results
func (uc *BigFiveUseCase) GetUserResults(ctx context.Context, targetUserID int) (*domain.BigFiveResult, error) {
	result, err := uc.bigFiveRepo.GetByUserID(ctx, targetUserID)
	if err != nil {
		return nil, errors.New("test results not found")
	}
	return result, nil
}
