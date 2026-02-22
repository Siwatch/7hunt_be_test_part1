package validator

import (
	"7hunt-be-rest-api/utils"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator *validator.Validate
	Trans     ut.Translator
}

func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New(validator.WithRequiredStructEnabled())
	validatorInstance := &Validator{
		Validator: validate,
		Trans:     trans,
	}

	return validatorInstance
}

func (v *Validator) ValidateRequest(req interface{}) *string {
	var errResult *string
	err := v.Validator.Struct(req)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	for _, e := range errs {
		translatedErr := fmt.Errorf("%s", e.Translate(v.Trans))
		if errResult == nil {
			errResult = utils.SetPtr(translatedErr.Error())
		} else {
			errResult = utils.SetPtr(*errResult + " | " + translatedErr.Error())
		}
	}

	return errResult
}
