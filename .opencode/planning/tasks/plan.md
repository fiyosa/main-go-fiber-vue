# Implementation Plan: Scalar API Docs di `/api/docs`

> Plan baru (2026-11-09), menggantikan plan Neo-Brutalism yang sudah selesai
> (diarsipkan di `tasks/archive/`).

## 1. Objective

Tambahkan halaman dokumentasi API interaktif berbasis **Scalar** pada route
`GET /api/docs` di aplikasi Go + Fiber ini. Halaman tersebut berisi UI
referensi API (Scalar) yang membedah sebuah dokumen **OpenAPI** yang
mendokumentasikan API publik aplikasi ini.

**Keluar dari scope (explicit):**
- Route MCP (`/api/mcp`) **tidak** ikut didokumentasikan di Scalar. MCP
  memakai transport Streamable HTTP dengan protokolnya sendiri (JSON-RPC,
  endpoint dinamis per tool), auth bearer sendiri, dan bukan bagian dari
  API publik `/api/v1`. Satu-satunya pengguna Scalar adalah endpoint publik.
- Tidak ada endpoint/UI baru lain selain `/api/docs`.
- Tidak ada koneksi `openapi` tambahan di komponen `api-reference` Scalar
  (terintegrasi dengan OpenAPI generation pasca `/build` — lihat Noted di
  bawah).

## 2. Target endpoint yang akan didokumentasikan

Daftar lengkap saat ini (dihitung dari `internal/routes/`):

### `/api/v1`
| Method | Path | Deskripsi |
|---|---|---|
| GET | `/api/v1/healthz` | Liveness check → envelope sukses |
| GET | `/api/v1/readyz` | Readiness check (ping DB) → 200 atau 503 |
| POST | `/api/v1/auth/login` | Login; token dikirim via cookie `access_token` + `refresh_token` (HttpOnly, SameSite=Lax), body hanya `{success, message}` |

Response envelope umum: `{ success, data?, message }` untuk sukses;
`{ success, message, errors? }` untuk error; 422 dengan `errors` per-field.

Login tidak memakai header `Authorization` — token lewat cookie. Ini harus
tercermin di OpenAPI (login tanpa `security`, dan endpoint lain kelak pakai
`cookieAuth`).

## 3. Arsitektur & pendekatan

### 3.1 URL dokumen OpenAPI vs UI

Dua route, dua peran:

1. **`GET /api/docs/openapi.json`** — melayani dokumen OpenAPI 3.1 statis
   (embedded, di-`go:embed` dari `internal/platform/docs/openapi.json`).
   Server-side: tidak ada dinamika, CORS cocok, mudah di-curl/digabung tool.
2. **`GET /api/docs`** — melayani halaman HTML. Mengambil
   `@scalar/api-reference` dari CDN jsDelivr dan dirender dari
   `url: '/api/docs/openapi.json'` (relative — ikut host/port, include cookie
   untuk nanti).

Keputusan: **CDN ESM** (pendekatan resmi Scalar: `createApiReference('#app',
{ url })`), bukan npm. Alasan:
- docs bukan bagian dari SPA `web/` — tidak punya bundle Vite sendiri.
- `/api/docs` dirender oleh Fiber (server-rendered HTML), bukan Vue.
- Versi CDN memakai `@scalar/api-reference@1.77.0` (dicek via `npm view`).
- Tidak ada build step tambahan → `go:embed` sederhana, asset self-contained di binary.

### 3.2 Layout file baru

```
internal/platform/docs/
  openapi.go        // go:embed doc HTML + openapi.json, route handlers
  openapi.json      // spec OpenAPI 3.1 (ditulis manual sesuai endpoint nyata)
  index.html        // shell Scalar: creates api-reference + injects auth token
internal/routes/
  docs.go           // RegisterDocs(app): GET /api/docs + /api/docs/openapi.json
```

`docs.go` dieksekusi dari `bootstrap/app.go` — **tidak** di dalam group
`/api/v1`, karena docs bukan bagian dari versi API; `/api/docs` berdiri di
level root `/api`.

Urutan mount di `app.go` (penting — lihat `RegisterWeb` fallback):
```
routes.RegisterMCP(...)   // paling dulu: /api/mcp
routes.RegisterAPI(...)   // /api/v1/*
routes.RegisterDocs(...)  // /api/docs*  → JANGAN DI-SHADOW fallback SPA
routes.RegisterWeb(...)   // terakhir: fallback /* (skip /api)
```
Fallback SPA di `web.go:34` sudah `case strings.HasPrefix(path, "/api")` →
melempar `c.Next()`, sehingga `/api/docs` tidak ketimpa SPA selama
`RegisterDocs` terpasang sebelum `RegisterWeb`. Ini sudah terbukti jalan untuk
`/api/mcp` & `/api/v1`.

### 3.3 Skema HTML `/api/docs`

```html
<!doctype html>
<html>
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>app-api − API Reference</title>
</head>
<body>
  <div id="app"></div>
  <script type="module">
    import { createApiReference } from 'https://cdn.jsdelivr.net/npm/@scalar/api-reference@1.77.0/esm.js'
    {
      const token = document.cookie
        ?.split('; ')
        ?.find((row) => row.startsWith('access_token='))
        ?.split('=')
        ?.at(1)
      createApiReference('#app', {
        // Relative: resolves on the same host/port the browser is using, so it
        // works in dev (localhost:8080) and behind a proxy in prod.
        url: '/api/docs/openapi.json',
        ...(token ? { authentication: { token } } : {}),
        proxyUrl: 'https://proxy.scalar.com',
      })
    }
  </script>
</body>
</html>
```

Catatan auth: endpoint saat ini tidak butuh token. School of thought: tetap
ambil `access_token` dari cookie yang sudah HttpOnly — dikirim saat browser
meminta `url`. Ini jadi **fungsional nanti bila API punya endpoint
ber-auth**, tanpa double-work sekarang. Hapus blok ini kalau ternyata
masalah security/readability (easy to revert).

### 3.4 openapi.json

Dokumen **OpenAPI 3.1 manual.** Isi:

```json
{
  "openapi": "3.1.0",
  "info": { "title": "app-api", "version": "0.1.0", "description": "..." },
  "servers": [ { "url": "" } ],   // relative — Scalar resolves ke host request
  "paths": {
    "/api/v1/healthz":  { "get": {...200 envelope...} },
    "/api/v1/readyz":   { "get": {...200/503...} },
    "/api/v1/auth/login": {
      "post": {
        "description": "...token dikirim lewat cookie...",
        "requestBody": { "email", "password" },
        "responses": {
          "200": body {success, message} + header Set-Cookie dua cookie,
          "422": validation errors
        }
      }
    }
  },
  "components": {
    "schemas": { "Envelope", "DataEnvelope", "ErrorEnvelope", "ValidationEnvelope",
                 "LoginRequest", "LoginResponse", "ErrorResponse" },
    "securitySchemes": {
      "cookieAuth": { "type": "apiKey", "in": "cookie", "name": "access_token" }
    }
  }
}
```

`servers: [{url:""}]` → dibuat relatif agar prod + dev/port ajaib work. Scalar
dan sebagian tools OpenAPI support empty string (UE request). Jangan
hardcode `localhost:8080`.

**JANGAN sertakan** `/api/mcp` (route MCP) di spec. Ini sesuai requirement:
MCP dianggap internal/otomatis, bukan API publik.

Catatan: `security: [{cookieAuth: []}]` TIDAK dipakai di path login (login
method pakai cookie, tidak authenticate). Bila endpoint yang butuh auth
ditambahkan nanti, baru path itu diberi `security`.

### 3.5 Kode

`internal/routes/docs.go`:

```go
package routes

import (
    "github.com/gofiber/fiber/v3"
    "go-fiber-vue/internal/platform/docs"
)

// RegisterDocs mounts the Scalar UI and its OpenAPI document under /api/docs*.
// It must be registered before RegisterWeb so the SPA fallback doesn't
// shadow it.
func RegisterDocs(app *fiber.App) {
    if app == nil { return }
    app.Get("/api/docs", docs.Handler)
    app.Get("/api/docs/openapi.json", docs.SpecHandler)
}
```

`internal/platform/docs/openapi.go`:

```go
package docs

import (
    "embed"
    "mime"
    "github.com/gofiber/fiber/v3"
)

//go:embed index.html openapi.json
var docsFS embed.FS

// Handler serves the Scalar API reference shell.
func Handler(c fiber.Ctx) error {
    body, err := docsFS.ReadFile("index.html")
    if err != nil { return fiber.ErrNotFound }
    c.Set(fiber.HeaderContentType, mime.TypeByExtension(".html"))
    c.Set(fiber.HeaderCacheControl, "no-cache")
    return c.Send(body)
}

// SpecHandler serves the public OpenAPI 3.1 document (MCP endpoints excluded).
func SpecHandler(c fiber.Ctx) error {
    body, err := docsFS.ReadFile("openapi.json")
    if err != nil { return fiber.ErrNotFound }
    c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
    c.Set(fiber.HeaderCacheControl, "public, max-age=3600")
    return c.Send(body)
}
```

Perhatian: `go:embed` membutuhkan file di path yang sama dengan package —
`index.html` & `openapi.json` ditaruh di `internal/platform/docs/`.

### 3.6 Bootstrap

`internal/bootstrap/app.go` — `Router()` (baris 51-53) dan `New()` (baris
102-104): tambahkan `routes.RegisterDocs(app)` setelah `RegisterAPI` dan
sebelum `RegisterWeb`. Tidak perlu parameter tambahan (spec statis).

## 4. Dependency graph

```
openapi.json (manual, source of truth: routes + response envelope)
        │
        ├── internal/platform/docs (go:embed index.html + openapi.json)
        │           │
        │           └── internal/routes/docs.go  ← RegisterDocs(app)
        │                       │
        └── internal/bootstrap/app.go: Router() & New()
                                    └── (order) → RegisterWeb fallback
```
- ~~koneksi eksternal Scalar CDN~~ — bukan dependency graph; runtime-only,
  di-render browser setelah halaman flush.
- Tidak ada relasi dengan `web/` (SPA) atau MCP.

## 5. Tasks (vertical slicing)

**T1 → M1:** `internal/platform/docs/openapi.json` — OpenAPI 3.1 untuk
`/api/v1/healthz`, `/api/v1/readyz`, `/api/v1/auth/login` (skema envelope,
login cookie, validation 422). Bukan MCP.
**V:T1** `openapi.json` valid (jq / python yaml parser); 3 paths; tak ada
path mcp.

**T2 → M2:** `internal/platform/docs/` — `openapi.go` (go:embed + 2 handler),
`index.html` (createApiReference CDN, relative url `/api/docs/openapi.json`,
token injection). 
**V:T2** `go build ./...` (output ke `./tmp`), `go vet ./...`.

**T3 → M3:** `internal/routes/docs.go` (`RegisterDocs(app)`) + wire ke
`bootstrap/app.go` (`Router()` & `New()`; urutan sebelum `RegisterWeb`).
**V:T3** `go build ./...`, `go vet ./...`, `go test ./...` (test baru bila ada
route-check). Manual:
- `GET /api/docs/openapi.json` (curl) → 200 JSON aplikasi
- `GET /api/docs` → 200 HTML berisi `createApiReference`
- `/api/docs` saat `index.html` tak ada → 404 (fallback)
- `/api/mcp` tetap 404/pass (MCP off untuk default)

## 6. Checkpoints

- **CP1 (usai T1):** human review `openapi.json` — path/cookie/status code
  sudah sesuai endpoint nyata? Belum lanjut kalau masih ada yang mau diubah.
- **CP2 (usai T2):** human review `index.html` — arah CDN & token-injection
  (boleh drop kalau tidak mau).
- **CP3 (usai T3, sebelum merge/review):** keluar draft plan — cek mount
  order & SPA fallback di browser dev (Vite 5173 → proxy ke 8080? atau akses
  langsung `localhost:8080/api/docs`).

## 7. Risks & mitigations

| Risk | Likelihood | Mitigation |
|---|---|---|
| Fallback SPA `/*` menimpa `/api/docs` | Low (sudah skip `/api` prefix) | mount `RegisterDocs` sebelum `RegisterWeb` + CP3 |
| CDN Scalar kena CSP/blocked | Medium | `proxyUrl` tersedia; fallback manual via npm jika perlu |
| `proxyUrl` tambahan CORS | Low | Hanya dipakai kalau perlu; minim request |
| Scalar tidak render karena `url` relatif kosong | Low | `url:'/api/docs/openapi.json'` relatif host — path absolute, bukan server-relative kosong |
| `go:embed` path salah | Very low | file di package dir, build gagal jelas |
| UI docs kaya aksesibilitas | — | pakai `<meta viewport>`, fallback teks di `index.html` |

## 8. Out of scope / follow-up (pasca plan ini)

- Inject `authentication.token` dari cookie akses → kapan endpoint ber-auth
  ditambahkan.
- Integrasi OpenAPI generation (`gen` dari handler/struct) bila jumlah
  endpoint bertambah banyak.
- Auth pada `/api/docs` itu sendiri (mis. paksa login) — belum diminta.
- Route MCP tetap di luar spec.
