package validation

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

// Validator represents a validator instance.
type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

// NewValidator creates a new instance of the validator.
func NewValidator() *Validator {
	v := &Validator{
		validator: validator.New(),
	}

	// Create a new 'en' locale
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)

	// Get a translator for the 'en' locale
	trans, _ := uni.GetTranslator("en")

	// Register the translator
	v.translator = trans

	// Optionally, you can register translations for error messages.
	err := enTranslations.RegisterDefaultTranslations(v.validator, v.translator)
	if err != nil {
		panic(err) // Handle error appropriately
	}
	return v
}

// ValidateStruct validates a struct against its defined validation rules.
func (v *Validator) ValidateStruct(s interface{}) error {
	if err := v.validator.Struct(s); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("field '%s' failed validation for tag '%s'", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation failed: %s", validationErrors)
	}
	return nil
}

// RegisterCustomValidationFunc registers a custom validation function with the validator.
func (v *Validator) RegisterCustomValidationFunc(tag string, fn validator.Func, msg string) {
	// Register the custom validation function
	v.validator.RegisterValidation(tag, fn)

	// Register translation for custom validation error message
	v.validator.RegisterTranslation(tag, v.translator, func(ut ut.Translator) error {
		return ut.Add(tag, msg, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}
