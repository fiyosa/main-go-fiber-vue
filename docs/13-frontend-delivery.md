# 13. Frontend Delivery

`public/` is a plain static-file directory. The frontend source lives in `web/` (Vue 3 + TypeScript + Vite) and is built into `public/build`.

## 13.1 Production

- Build output goes to `public/build`, produced by Vite.
- `public/build` contains only the flat build artifacts — `index.js`, `index.css`, and `assets/*` — with the remaining Vite output files placed directly inside it, so no extra nested folders are allowed.
- Serve `public/` as-is next to the Go binary. Frontend build assets are not embedded; only the HTML shell (`internal/app/presentation/core/view/spa.html`) is embedded with `//go:embed` as `view.SPA`.
- Serve hashed assets with long cache headers.
- `internal/routes/web.go` owns the `/*` fallback: `RegisterWeb` serves the embedded SPA shell for any non-`/api` path and `/build/*` from `public/build` for `index.js`, `index.css`, and `assets/*`; `internal/bootstrap/app.go` calls it after `routes.RegisterAPI`.
- Keep `/admin/login` returning the SPA shell.
- Keep public media and file routes unauthenticated.

## 13.2 Development

- Run Vite dev server from `web/` on `:5173`.
- Proxy `/api`, `/media`, and `/file-upload` to the Go API on `:8080`.
- CORS allowlist: `http://localhost:5173` and the production origin.

## 13.3 Deployment

- The single Go binary serves API from `public/` (static files + build) and runs the cron scheduler plus job goroutines in the same process.
- Alternatively, serve `public/build` from a CDN while keeping the API and scheduler in one Go process. Keep the same origin if possible to avoid CORS and cookie issues.
