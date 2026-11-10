# Todo: Scalar API Docs di `/api/docs`

> Digenerate berdasarkan `plan.md` (2026-11-09). Menggantikan todo Neo-Brutalism
> (diarsipkan di `tasks/archive/todo-neo-brutalism.md`).

## Phase 1 — OpenAPI 3.1 Document

- [x] **Task 1: Buat `internal/platform/docs/openapi.json`**
  - Tulis manual spec OpenAPI 3.1.0 untuk `/api/v1`
  - Tag: `Health` (`GET /api/v1/healthz`, `GET /api/v1/readyz`) & `Auth` (`POST /api/v1/auth/login`)
  - Skema: `Envelope`, `ValidationEnvelope`, `LoginRequest`, `LoginResponse`
  - Autentikasi: `cookieAuth` (`in: cookie`, `name: access_token`), tapi jangan pasang `security` di endpoint `/login`
  - **Eksplisit JANGAN sertakan `/api/mcp`** (route MCP terpisah dari API publik)
  - *Acceptance:* Valid JSON, valid OpenAPI 3.1, server URL `""` (relatif), 3 path.

---

## Phase 2 — Handler & Embedding

- [x] **Task 2: Buat `internal/platform/docs/index.html`**
  - HTML 5 shell dengan `<meta charset>`, viewport, title `app-api − API Reference`
  - Inject CDN Scalar ESM: `@scalar/api-reference@1.77.0`
  - Inisialisasi `createApiReference('#app', { url: '/api/docs/openapi.json', proxyUrl: 'https://proxy.scalar.com' })`
  - Tambahkan snippet pembaca cookie `access_token` untuk token injection jika ada
  - *Acceptance:* Pure HTML valid, standalone (tanpa build Vue/Vite).

- [x] **Task 3: Buat `internal/platform/docs/openapi.go`**
  - Go package `docs` dengan `//go:embed index.html openapi.json`
  - Handler `Handler(c fiber.Ctx)` → return `index.html` dengan header `Content-Type: text/html; charset=utf-8`, `Cache-Control: no-cache`
  - Handler `SpecHandler(c fiber.Ctx)` → return `openapi.json` dengan header `Content-Type: application/json`, `Cache-Control: public, max-age=3600`
  - *Acceptance:* `go build ./...` (dengan output ke `./tmp`) dan `go vet ./...` lulus tanpa warning.

---

## Phase 3 — Routing & Wiring

- [x] **Task 4: Buat `internal/routes/docs.go` & daftarkan di `bootstrap/app.go`**
  - Function `RegisterDocs(app *fiber.App)` daftarkan route:
    - `GET /api/docs` → `docs.Handler`
    - `GET /api/docs/openapi.json` → `docs.SpecHandler`
  - Di `internal/bootstrap/app.go`:
    - Panggil `routes.RegisterDocs(app)` di `Router()` dan `New()`, **sebelum** `routes.RegisterWeb(app, cfg.PublicDir)` agar tidak ketutupan SPA fallback `/*`.
  - *Acceptance:*
    - `go build ./...` (output ke `./tmp`), `go vet ./...`, dan `go test ./...` lulus.
    - Test HTTP via Fiber / httptest:
      - `GET /api/docs` -> 200 OK + `text/html` + memuat script CDN Scalar
      - `GET /api/docs/openapi.json` -> 200 OK + `application/json` + memuat json yang valid
      - `GET /api/mcp` -> tidak muncul di spec OpenAPI (`openapi.json` tak ada kata `/api/mcp`)

---

## Phase 4 — Verification & Cleanup

- [x] **Task 5: Verifikasi menyeluruh & regression check**
  - Jalankan `go test ./...`
  - `gofmt -l` bersih pada file baru/yang diubah
  - Cek `git status` memastikan tidak ada file temporary atau binary tercecer di root.

---

## Checkpoints

1. **Pasca Task 1:** Tinjau `openapi.json` (apakah ada endpoint yang kurang / kualifikasi path salah).
2. **Pasca Task 3:** Tinjau `docs.Handler` & `index.html` (apakah CDN diizinkan & `go:embed` tepat).
3. **Pasca Task 4:** Test runtime HTTP handler untuk `/api/docs`, `/api/docs/openapi.json`, dan SPA fallback `/some-vue-route`.
