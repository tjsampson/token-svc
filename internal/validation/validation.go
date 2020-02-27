package validation

import (
	"fmt"
	"log"

	"github.com/tjsampson/token-svc/internal/errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type provider struct {
	validate *validator.Validate
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
		validate: v,
		translator: trans,
	}
}

func (p *provider) Validate(model interface{}) error {
	if err := p.validate.Struct(model); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errs := []string{}

		// errMsg := ""
		for _, vErr := range validationErrors {
			// errMsg = errMsg + fmt.Sprintf("%s %s (%s) ", vErr.Field(), vErr.Tag(), vErr.Param())
			// errMsg = errMsg + vErr.Value().(string)
			// errMsg = errMsg + vErr.Tag()
			errs = append(errs, vErr.Translate(p.translator))
			// errMsg = errMsg + "\n" + vErr.Translate(p.translator)
		}

		return &errors.RestError{
			Code:          400,
			Message:       fmt.Sprintf("validation error(s)"),
			OriginalError: err,
			Messages: errs,
		}
	}
	return nil
}
