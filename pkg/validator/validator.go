package validator

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate

	// Email regex pattern
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// Phone regex pattern (international format)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

func init() {
	validate = validator.New()

	// Register custom validators
	_ = validate.RegisterValidation("adult", validateAdult)
	_ = validate.RegisterValidation("validgender", validateGender)
	_ = validate.RegisterValidation("validmbti", validateMBTI)
}

// Validate validates a struct using validator tags
func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		// Format validation errors
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		var errMsg string
		for _, err := range err.(validator.ValidationErrors) {
			errMsg += fmt.Sprintf("Field '%s' failed validation '%s'; ", err.Field(), err.Tag())
		}
		return fmt.Errorf("%s", errMsg)
	}
	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// ValidatePassword checks password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if len(password) > 128 {
		return fmt.Errorf("password must be at most 128 characters long")
	}
	return nil
}

// validateAdult checks if birth date indicates user is 18+ years old
func validateAdult(fl validator.FieldLevel) bool {
	birthDate, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	age := time.Since(birthDate).Hours() / 24 / 365.25
	return age >= 18
}

// validateGender validates gender field
func validateGender(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	return gender == "male" || gender == "female"
}

// validateMBTI validates MBTI type
func validateMBTI(fl validator.FieldLevel) bool {
	mbti := fl.Field().String()
	validTypes := []string{
		"INTJ", "INTP", "ENTJ", "ENTP",
		"INFJ", "INFP", "ENFJ", "ENFP",
		"ISTJ", "ISFJ", "ESTJ", "ESFJ",
		"ISTP", "ISFP", "ESTP", "ESFP",
	}

	for _, valid := range validTypes {
		if mbti == valid {
			return true
		}
	}
	return false
}

// ValidateAge checks if age is within acceptable range
func ValidateAge(birthDate time.Time) error {
	age := time.Since(birthDate).Hours() / 24 / 365.25

	if age < 18 {
		return fmt.Errorf("user must be at least 18 years old")
	}
	if age > 120 {
		return fmt.Errorf("invalid birth date")
	}
	return nil
}
