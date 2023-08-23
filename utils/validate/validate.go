package validate

import (
	"errors"
	"fmt"
	"log"
	"unicode"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	validate := validator.New()

	registerAlphaNumSpace(validate, trans)

	return &Validator{
		Validator:  validate,
		Translator: trans,
	}
}

func (v *Validator) Validate(i interface{}) error {
	err := v.Validator.Struct(i)

	if err != nil {
		object, _ := err.(validator.ValidationErrors)

		for _, key := range object {
			return errors.New(key.Translate(v.Translator))
		}
	}

	return nil
}

// alphaNumSpace is a custom validation function that
// allow aplhanumeric and spaces
func alphaNumSpace(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	for _, char := range value {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && !unicode.IsSpace(char) {
			return false
		}
	}

	return true
}

func registerAlphaNumSpace(validate *validator.Validate, trans ut.Translator) {
	var (
		tag = "alphanumspace"
	)

	transaltionFn := func(ut ut.Translator, fe validator.FieldError) string {
		return fmt.Sprintf("%s can only contain alphanumeric characters", fe.Field())
	}

	alphaNumSpaceErr := func(ut ut.Translator) error {
		return nil
	}

	validate.RegisterValidation(tag, alphaNumSpace)
	validate.RegisterTranslation(tag, trans, alphaNumSpaceErr, transaltionFn)
}
