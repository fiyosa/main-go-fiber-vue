# 20. Localization

## 20.1 Goal

Every JSON API response must contain a non-empty `message`. The message is translated from the active locale, never written as a literal in a controller or handler.

- Supported locales: `en`, `id`.
- Default locale: `en`.
- The equivalent of Laravel `__()` is `Translator.T(locale, key, params)`.
- Unknown locale falls back to `en`.
- A missing key returns the key itself and writes a warning log. There is no fallback-locale configuration.

The catalogue in `docs/lang/` is the review copy. The Go runtime loads the same files from `lang/`.

## 20.2 Configuration

| Variable | Default | Meaning |
| --- | --- | --- |
| `APP_LOCALE` | `en` | Default locale when the request does not select one |

`APP_LOCALE` is the only locale setting.

Validation rules:

- `APP_LOCALE` must be one of `en`, `id`.
- An invalid `APP_LOCALE` fails startup instead of silently using another language.
- Per-request locale never changes the global default.

## 20.3 Catalogue format

Each locale is a nested JSON object. Dot notation resolves the path:

```text
lang.retrieved -> lang.en.json -> { "lang": { "retrieved": ":operator Retrieved" } }
```

Files:

```text
lang/
├── en.json
└── id.json
```

Rules:

- Both files must contain exactly the same key set.
- Values must be non-empty strings.
- Keys use lowercase snake case.
- Placeholders use Laravel syntax: `:operator`, `:attribute`, `:min`, `:email`.
- Do not add pluralization until the frontend contract requires it.
- Do not translate API field names, enum values, paths, or IDs.

Example:

```json
{
  "lang": {
    "retrieved": ":operator Retrieved",
    "created": ":operator Created"
  },
  "common": {
    "forbidden": "Forbidden"
  }
}
```

## 20.4 Translator

```go
// internal/platform/lang/translator.go
package lang

type Translator struct {
    catalog map[string]map[string]string
    defaultLocale string
}

func (t *Translator) T(locale, key string, params map[string]string) string {
    if locale == "" {
        locale = t.defaultLocale
    }

    message, ok := t.lookup(locale, key)
    if !ok {
        slog.Warn("missing translation", "locale", locale, "key", key)
        return key
    }

    return replacePlaceholders(message, params)
}
```

Rules:

- `replacePlaceholders` replaces every `:name` occurrence with `params["name"]`.
- Unknown placeholders are left untouched.
- The translator is built once at startup and is read-only after construction.
- Controllers receive the translator through their dependencies; no package-level global.
- The translator never returns an empty string.

Loading:

```text
cmd/app/main.go
  -> bootstrap.Config loads APP_LOCALE
  -> bootstrap.newTranslator loads lang/en.json and lang/id.json
  -> translator is stored in Dependencies
  -> helpers and controllers receive the same instance
```

The loader looks for the catalogues in this order and uses the first directory that exists: `lang`, `<binary dir>/lang`. Deployment must keep one of them next to the binary. Hot reload is not supported; restart the process after changing translations.

## 20.5 Locale resolution

The locale middleware resolves the request locale in this order:

1. `?lang=en` or `?lang=id`.
2. `Accept-Language` header, first supported language.
3. `APP_LOCALE` (`en` by default).

Unsupported values are ignored. The middleware stores the resolved locale in Fiber locals so controllers and response helpers do not parse headers themselves. It lives in `internal/platform/http/middleware/locale.go` and is registered before the maintenance middleware, so 503 responses are translated too.

Examples:

| Request | Locale |
| --- | --- |
| no header | `en` |
| `Accept-Language: id-ID,id;q=0.9,en;q=0.8` | `id` |
| `Accept-Language: fr-FR` | `en` |
| `?lang=id` | `id` |
| `?lang=fr` | `en` |

## 20.6 Response contract

`message` is mandatory for every JSON response. `Envelope.Message` must not use `omitempty`.

```go
type Envelope struct {
    Success bool   `json:"success"`
    Data    any    `json:"data,omitempty"`
    Message string `json:"message"`
}
```

Response helpers accept a translation key, not a literal:

```go
helpers.Success(c, data, "lang.retrieved", map[string]string{
    "operator": translator.T(locale, "lang.appointment", nil),
}, fiber.StatusOK)

helpers.Error(c, "common.forbidden", nil, fiber.StatusForbidden)
```

Required messages:

| Condition | Key |
| --- | --- |
| Login success | `lang.logged_in_successfully` |
| Login failure | `lang.invalid_credentials` |
| Logout success | `lang.logged_out_successfully` |
| Token missing, invalid, or expired | `auth.unauthenticated` / `auth.invalid_token` |
| Permission denied | `common.forbidden` |
| Resource missing | `lang.not_found` or `common.not_found` |
| Validation failure | `common.validation_failed` |
| Unexpected server error | `common.try_again_later` |
| Generic success | `common.success` |

`204 No Content` cannot carry a JSON body. JSON delete endpoints therefore return `200` with a `message`. This is an intentional difference from Laravel and is recorded in [Section 18](18-risks-and-open-decisions.md).

The existing Laravel quirk where a `403` envelope still has `success: true` is preserved separately; it does not change the `message` rule.

## 20.7 Validation messages

The validator adapter prefers keys under `validation.*`:

```text
validation.required
validation.email
validation.min.string
validation.max.string
validation.unique
validation.exists
validation.mimes
```

The response shape stays compatible:

```json
{
  "success": false,
  "message": "Data yang diberikan tidak valid",
  "errors": {
    "email": ["Kolom email wajib diisi."]
  }
}
```

`message` is the generic translated summary. Field messages come from the same catalogue and use `:attribute` and rule placeholders.

Validation attribute names are translated where a dedicated key exists. Do not expose raw database column names when a human-readable label is available.

## 20.8 Controller rules

- Never pass a literal to `helpers.Success` or `helpers.Error`.
- Pass a stable key, not user input.
- Build dynamic values through the translator; for example `lang.retrieved` with operator `lang.appointment`.
- Do not concatenate translated fragments. Translate the full sentence with placeholders.
- Do not return English text from a service, integration adapter, or database constraint. Services return typed errors or keys; the presentation layer translates.
- Every new endpoint must add its keys to both locale files in the same change.

## 20.9 Migration audit

Before a module is considered ported:

1. Find every `sendResponse(`, `sendSuccess(`, and `sendError(` call.
2. Replace every literal message with a key.
3. Confirm the key exists in both `en.json` and `id.json`.
4. Confirm the response is non-empty for `2xx`, `4xx`, and `5xx`.
5. Confirm `message` is still present when `data` is absent.
6. Confirm delete endpoints return a body or document an exception.

Known hard-coded messages to migrate include controller strings such as `Aux type created successfully`.

## 20.10 Tests

- `en.json` and `id.json` have identical key sets.
- Both catalogues parse as valid JSON and contain no empty values.
- `T("en", "lang.logged_in_successfully", nil)` returns the English message.
- `T("id", "lang.logged_in_successfully", nil)` returns the Indonesian message.
- Unsupported locale resolves to `en`.
- Missing key returns the key and logs a warning.
- Placeholder replacement covers `:operator`, `:attribute`, `:email`, `:min`, `:values`.
- Every API response, including `401`, `403`, `404`, `422`, and `500`, has a non-empty `message`.
- Delete endpoints return `200` with a message.
- Tests do not assert a hard-coded message string unless the locale is explicit.

## 20.11 Adding a language

1. Add `lang/<code>.json` with the same key set as `en.json`.
2. Add the code to the supported-locale allowlist.
3. Add it to the locale middleware tests.
4. A new language is complete only when its catalogue is complete.
