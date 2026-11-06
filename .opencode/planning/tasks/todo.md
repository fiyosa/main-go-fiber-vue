# Todo: Minimal Frontend — Hello World + Vite Favicon

Scope: one page that renders `hello world` and one browser-tab icon from `public/vite.svg`.
No routing, API calls, state management, tests directory, CI, or docs edits.

Plan: `.opencode/planning/tasks/plan.md`

## Phase 1 — Hello World Slice

- [x] Task 1: Buildable Vue hello-world frontend served by Go
  - [x] `web/package.json` with Vue 3 + Vite + TypeScript; `dev` (:5173), `build`, and `test` scripts
  - [x] `web/vite.config.ts`: `outDir: '../public/build'`, input `src/main.ts`,
        `entryFileNames: 'index.js'`, CSS -> `index.css`
  - [x] `web/index.html` for the dev server only; build emits no `index.html`
  - [x] `web/src/main.ts`, `web/src/App.vue` (renders exactly `hello world`), `web/src/style.css`
  - [x] `web/src/App.test.ts` (Vitest + Vue Test Utils) proves the render
  - [x] `.gitignore` ignores `web/node_modules/` and `web/dist/`
  - [x] Verify: `npm --prefix web install` and `npm --prefix web run build` succeed
  - [x] Verify: `public/build/index.js` + `index.css` exist; no `index.html`
  - [x] Verify: `go run ./cmd/app`; browser at `:8080/` shows `hello world`
  - [x] Verify: `go build ./...`, `go vet ./...`, `gofmt -l` clean

### Checkpoint — Hello World

- [x] Frontend builds; Go serves it; page shows `hello world`
- [ ] Human confirms layout before favicon work

## Phase 2 — Browser Tab Icon

- [ ] Task 2: Serve `public/vite.svg` as the favicon
  - [ ] Add `public/vite.svg` (official Vite logo)
  - [ ] Add `<link rel="icon" type="image/svg+xml" href="/vite.svg" />` to `spa.html`
  - [ ] Add the same link to `web/index.html`
  - [ ] Add `GET /vite.svg` in `internal/routes/web.go` before the `/*` fallback
  - [ ] Verify: `curl -i localhost:8080/vite.svg` -> 200 `image/svg+xml`
  - [ ] Verify: unknown `/api/*` still returns JSON 404; other paths still return the SPA shell
  - [ ] Verify: browser tab shows the icon
  - [ ] Verify: `go build ./...`, `go vet ./...`, `go test ./...`

### Checkpoint — Complete

- [ ] `hello world` + favicon both work through the Go server
- [ ] No changes under `docs/`, `tests/`, or `.github/`
- [ ] Human approves; ready for `/build`

## Deferred

- [ ] Frontend routing, API client, auth screens
- [ ] Frontend unit/e2e tests
- [ ] CI build step for `web/`
- [ ] CDN/static deployment of `public/build`

## Standing Rules

- Keep the production shell in `spa.html`; the Vite build must not emit `index.html`.
- Build output stays flat: `index.js`, `index.css`, `assets/*` in `public/build`.
- Frontend code and comments in English.
- Do not edit `docs/`.
- Do not create `tests/` or `.github/`.
