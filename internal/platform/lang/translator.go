package lang

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// Locales supported by the API. APP_LOCALE must be one of these.
var supportedLocales = []string{"en", "id"}

const (
	DefaultLocale     = "en"
	translatorLocals  = "translator"
	requestLocaleKey  = "locale"
	catalogFileSuffix = ".json"
)

// Translator resolves translation keys against the locale catalogues loaded at
// startup. It is read-only after construction and safe for concurrent use.
type Translator struct {
	catalog       map[string]map[string]string
	defaultLocale string
}

// Load reads every <locale>.json file in dir and builds a Translator.
func Load(dir, fallback string) (*Translator, error) {
	if fallback == "" {
		fallback = DefaultLocale
	}
	if !IsSupported(fallback) {
		return nil, fmt.Errorf("unsupported default locale %q", fallback)
	}

	translator := &Translator{
		catalog:       make(map[string]map[string]string),
		defaultLocale: fallback,
	}

	for _, locale := range supportedLocales {
		path := filepath.Join(dir, locale+catalogFileSuffix)
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read locale %q: %w", locale, err)
		}

		var nested map[string]any
		if err := json.Unmarshal(raw, &nested); err != nil {
			return nil, fmt.Errorf("parse locale %q: %w", locale, err)
		}

		flat := make(map[string]string)
		flatten("", nested, flat)
		translator.catalog[locale] = flat
	}

	if err := translator.validateKeyParity(); err != nil {
		return nil, err
	}
	return translator, nil
}

// IsSupported reports whether locale is part of the allowlist.
func IsSupported(locale string) bool {
	return slices.Contains(supportedLocales, locale)
}

// DefaultLocale returns the configured fallback locale.
func (t *Translator) DefaultLocale() string {
	return t.defaultLocale
}

// T resolves key for locale, falling back to the default locale. A missing key
// returns the key itself and logs a warning; it never returns an empty string.
func (t *Translator) T(locale, key string, params map[string]string) string {
	if t == nil {
		return key
	}
	if locale == "" || !IsSupported(locale) {
		locale = t.defaultLocale
	}

	message, ok := t.lookup(locale, key)
	if !ok && locale != t.defaultLocale {
		message, ok = t.lookup(t.defaultLocale, key)
	}
	if !ok {
		slog.Warn("missing translation", "locale", locale, "key", key)
		return key
	}
	return replacePlaceholders(message, params)
}

// Has reports whether key exists for locale. Used to resolve attribute labels
// without logging a warning for optional translations.
func (t *Translator) Has(locale, key string) bool {
	if t == nil {
		return false
	}
	if locale == "" || !IsSupported(locale) {
		locale = t.defaultLocale
	}
	_, ok := t.lookup(locale, key)
	return ok
}

// SetRequest stores the resolved translator and locale in Fiber locals.
func SetRequest(c fiber.Ctx, translator *Translator, locale string) {
	c.Locals(translatorLocals, translator)
	c.Locals(requestLocaleKey, locale)
}

// From returns the translator stored in Fiber locals by the locale middleware.
func From(c fiber.Ctx) *Translator {
	if value, ok := c.Locals(translatorLocals).(*Translator); ok {
		return value
	}
	return nil
}

// Locale returns the request locale resolved by the locale middleware.
func Locale(c fiber.Ctx) string {
	if value, ok := c.Locals(requestLocaleKey).(string); ok && value != "" {
		return value
	}
	return DefaultLocale
}

func (t *Translator) lookup(locale, key string) (string, bool) {
	messages, ok := t.catalog[locale]
	if !ok {
		return "", false
	}
	message, ok := messages[key]
	if !ok || message == "" {
		return "", false
	}
	return message, true
}

func (t *Translator) validateKeyParity() error {
	reference := keysOf(t.catalog[DefaultLocale])
	for locale, messages := range t.catalog {
		if locale == DefaultLocale {
			continue
		}
		current := keysOf(messages)
		if missing, extra := diff(reference, current), diff(current, reference); len(missing) > 0 || len(extra) > 0 {
			return fmt.Errorf("locale %q key mismatch: missing=%v extra=%v", locale, missing, extra)
		}
	}
	return nil
}

func keysOf(messages map[string]string) []string {
	keys := make([]string, 0, len(messages))
	for key := range messages {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func diff(left, right []string) []string {
	seen := make(map[string]struct{}, len(right))
	for _, value := range right {
		seen[value] = struct{}{}
	}
	var out []string
	for _, value := range left {
		if _, ok := seen[value]; !ok {
			out = append(out, value)
		}
	}
	return out
}

func flatten(prefix string, value any, out map[string]string) {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			next := key
			if prefix != "" {
				next = prefix + "." + key
			}
			flatten(next, child, out)
		}
	case string:
		if typed != "" {
			out[prefix] = typed
		}
	}
}

func replacePlaceholders(message string, params map[string]string) string {
	for name, value := range params {
		message = strings.ReplaceAll(message, ":"+name, value)
	}
	return message
}
