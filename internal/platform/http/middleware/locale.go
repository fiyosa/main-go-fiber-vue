package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/platform/lang"
)

// Locale resolves the request locale (?lang=, then Accept-Language, then the
// configured default) and stores the translator plus locale in Fiber locals.
func Locale(translator *lang.Translator) fiber.Handler {
	return func(c fiber.Ctx) error {
		locale := resolveLocale(c, translator)
		lang.SetRequest(c, translator, locale)
		return c.Next()
	}
}

func resolveLocale(c fiber.Ctx, translator *lang.Translator) string {
	fallback := lang.DefaultLocale
	if translator != nil {
		fallback = translator.DefaultLocale()
	}

	if query := strings.TrimSpace(c.Query("lang")); lang.IsSupported(query) {
		return query
	}

	for _, candidate := range parseAcceptLanguage(c.Get("Accept-Language")) {
		if lang.IsSupported(candidate) {
			return candidate
		}
	}
	return fallback
}

func parseAcceptLanguage(header string) []string {
	if header == "" {
		return nil
	}

	parts := strings.Split(header, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		tag := strings.TrimSpace(strings.SplitN(part, ";", 2)[0])
		if tag == "" || tag == "*" {
			continue
		}
		if index := strings.Index(tag, "-"); index > 0 {
			tag = tag[:index]
		}
		out = append(out, strings.ToLower(tag))
	}
	return out
}
