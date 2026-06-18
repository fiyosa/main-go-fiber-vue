package lang

import (
	"fmt"
	"go-fiber-svelte/internal/config"
	"go-fiber-svelte/internal/lang/locales"
	"strings"
)

var T t

type t struct{}

func (*t) Convert(msg string, args ...map[string]any) string {
	if len(args) == 0 || args[0] == nil {
		return msg
	}

	newMsg := msg
	for key, value := range args[0] {
		newMsg = strings.ReplaceAll(newMsg, ":"+key, fmt.Sprintf("%v", value))
	}
	return newMsg
}

func (*t) Get() locales.ILang {
	return locale[config.APP_LOCALE]
}

var locale = map[string]locales.ILang{
	"en": locales.EN,
	"id": locales.ID,
}
