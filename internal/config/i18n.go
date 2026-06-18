package config

import (
	"fmt"
	"os"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

var (
	Translator ut.Translator
	Validate   *validator.Validate
)

func I18n() {
	locale := APP_LOCALE
	en := en.New()
	uni := ut.New(en, en, id.New())

	var found bool
	Translator, found = uni.GetTranslator(locale)
	if !found {
		fmt.Printf("Translator for locale %v not found", locale)
		os.Exit(1)
	}

	Validate = validator.New()
	var err error

	switch locale {
	case "en":
		err = en_translations.RegisterDefaultTranslations(Validate, Translator)
	case "id":
		err = id_translations.RegisterDefaultTranslations(Validate, Translator)
	default:
		err = en_translations.RegisterDefaultTranslations(Validate, Translator)
	}

	if err != nil {
		fmt.Printf("Error register translation: %v", err.Error())
		os.Exit(1)
	}
}
