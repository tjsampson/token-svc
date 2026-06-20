package validation

import (
	"fmt"
	"log"

	"github.com/tjsampson/token-svc/internal/errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type provider struct {
	validate   *validator.Validate
	translator ut.Translator
}

// Provider is the validation interface
type Provider interface {
	Validate(model interface{}) error
}

// New returns a new Validation Provider
func New(v *validator.Validate) Provider {

	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	return &provider{
		validate:   v,
		translator: trans,
	}
}

func (p *provider) Validate(model interface{}) error {
	if err := p.validate.Struct(model); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errs := []string{}

		for _, vErr := range validationErrors {
			errs = append(errs, vErr.Translate(p.translator))
		}

		return &errors.RestError{
			Code:          400,
			Message:       fmt.Sprintf("validation error(s)"),
			OriginalError: err,
			Messages:      errs,
		}
	}
	return nil
}
