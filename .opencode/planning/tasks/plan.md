# Implementation Plan: Minimal Frontend — Hello World + Vite Favicon

## Overview

Deliver the smallest working Vue 3 + Vite frontend for the existing Go/Fiber app. The page
renders only `hello world`, and the browser tab uses a favicon served from `public/vite.svg`.

The Go side already owns the HTML shell (`internal/app/presentation/core/view/spa.html`, embedded
as `view.SPA`) and serves `/build/*` from `public/build`. This plan fills the empty `web/` source
folder, configures Vite to emit the flat artifact layout the Go server expects, and adds one
static route for the root-level favicon.

Out of scope: routing, API calls, state management, styling system, tests directory, CI, auth UI.

## Architecture Decisions

- **Vue 3 + TypeScript + Vite in `web/`**, per `docs/13-frontend-delivery.md`.
- **Vite output goes straight to `../public/build`** with fixed entry names:
  `index.js`, `index.css`, and hashed extras under `assets/`. The Go shell links those names.
- **`web/index.html` exists only for the Vite dev server** on `:5173`; the production shell stays
  in `spa.html`, so the Vite build must not emit its own `index.html`.
- **Build input is `src/main.ts`, not `index.html`**, which prevents an emitted `index.html`
  from overwriting the Go-owned shell contract.
- **Favicon lives at `public/vite.svg`** (repo root, not `web/public`) and is referenced as
  `/vite.svg` from both HTML shells.
- **Go serves `/vite.svg` explicitly** via `c.SendFile(filepath.Join(publicDir, "vite.svg"))`,
  registered before the `/*` SPA fallback. No broad root static mount is added.
- **No docs edits.** `docs/` is read-only per the spec; this plan follows the documented layout.
- **`public/build` stays gitignored**; the UI must be built locally or by a future CI step.
  Add `web/node_modules/` and `web/dist/` to `.gitignore`.

## Dependency Graph

```text
web/package.json + tsconfig
        │
        └──> vite.config.ts (outDir ../public/build, flat names, :5173 proxy)
                 │
                 └──> src/main.ts -> App.vue -> style.css
                          │
                          └──> npm run build
                                   │
                                   └──> public/build/index.js + index.css (+ assets/*)
                                            │
                                            └──> Go RegisterWeb (/build static)
                                                     │
                                                     └──> embedded spa.html
                                                              │
                                                              └──> browser shows "hello world"

public/vite.svg ──> spa.html + web/index.html <link rel="icon" href="/vite.svg">
                        │
                        └──> Go route GET /vite.svg ──> browser tab icon
```

## Task List

### Phase 1: Hello World Slice

- [x] Task 1: Buildable Vue hello-world frontend served by Go
  - Description: Scaffold `web/` with Vue 3 + TypeScript + Vite. Configure Vite to build
    `src/main.ts` into `../public/build` as `index.js` and `index.css`. The app renders exactly
    `hello world`. Run the production build and confirm the existing Go server serves it through
    the embedded shell.
  - Acceptance criteria:
    - [x] `web/package.json` defines `dev` (port 5173) and `build` scripts; deps are `vue`,
      `vite`, `@vitejs/plugin-vue`, `typescript`, `vue-tsc`, plus `vitest`, `@vue/test-utils`,
      `jsdom` for the render test.
    - [x] `npm --prefix web run build` exits 0 and writes `public/build/index.js`,
      `public/build/index.css`, and any extras under `public/build/assets/`.
    - [x] Build does **not** write `public/build/index.html`.
    - [x] `go run ./cmd/app` then `GET /` returns the SPA shell, and the loaded app shows
      `hello world` with no console errors.
    - [x] `go build ./...`, `go vet ./...`, and `gofmt -l` remain clean.
  - Verification:
    - [x] `npm --prefix web install`
    - [x] `npm --prefix web test`
    - [x] `npm --prefix web run build`
    - [x] `go run ./cmd/app` + browser at `http://localhost:8080/`
    - [x] `curl -s http://localhost:8080/build/index.js | head`
  - Dependencies: None.
  - Files likely touched:
    - `web/package.json`
    - `web/vite.config.ts`
    - `web/tsconfig.json`
    - `web/index.html` (dev server only)
    - `web/src/main.ts`
    - `web/src/App.vue`
    - `web/src/App.test.ts`
    - `web/src/env.d.ts`
    - `web/src/style.css`
    - `.gitignore`
  - Estimated scope: Medium (5–8 files).

### Checkpoint: Hello World

- [ ] Frontend builds without errors.
- [ ] Go server serves the built app.
- [ ] Browser displays `hello world`.
- [ ] Human confirms the page and build layout before adding the favicon.

### Phase 2: Browser Tab Icon

- [ ] Task 2: Serve `public/vite.svg` as the browser tab icon
  - Description: Add the Vite logo SVG at `public/vite.svg`, reference it from the embedded
    production shell and the Vite dev HTML, and register an explicit `GET /vite.svg` route in
    `internal/routes/web.go` before the SPA fallback.
  - Acceptance criteria:
    - [ ] `public/vite.svg` exists and is valid SVG.
    - [ ] `spa.html` and `web/index.html` include
      `<link rel="icon" type="image/svg+xml" href="/vite.svg" />`.
    - [ ] `GET /vite.svg` returns `200` with `image/svg+xml` and the SVG body.
    - [ ] `GET /api/v1/...` behavior is unchanged; unknown `/api/*` still returns JSON 404.
    - [ ] `GET /some/unknown/path` still returns the SPA shell.
    - [ ] Browser tab shows the Vite icon at `http://localhost:8080/`.
  - Verification:
    - [ ] `curl -i http://localhost:8080/vite.svg`
    - [ ] `curl -i http://localhost:8080/api/v1/nope`
    - [ ] `curl -s http://localhost:8080/some/unknown/path | head`
    - [ ] `go build ./...`, `go vet ./...`, `go test ./...`
    - [ ] Manual browser check of the tab icon.
  - Dependencies: Task 1.
  - Files likely touched:
    - `public/vite.svg` (new)
    - `internal/app/presentation/core/view/spa.html`
    - `web/index.html`
    - `internal/routes/web.go`
  - Estimated scope: Small (3–4 files).

### Checkpoint: Complete

- [ ] Both tasks' acceptance criteria met.
- [ ] `hello world` renders from the Go-served production shell.
- [ ] Favicon loads from `/vite.svg` in the browser tab.
- [ ] No changes under `docs/`, `tests/`, or `.github/`.
- [ ] Human reviews before `/build` starts.

## Risks and Mitigations

| Risk | Impact | Mitigation |
| --- | --- | --- |
| `public/vite.svg` does not exist yet | High | Task 2 creates it with the official Vite logo SVG; confirm with the human first |
| Vite default build emits hashed `assets/index-*.js` and an `index.html` | High | Set `rollupOptions.input = 'src/main.ts'` and fixed `entryFileNames`/`assetFileNames` |
| `/vite.svg` currently hits the `/*` SPA fallback and returns HTML | High | Register the explicit static route before the fallback |
| Dev server (`:5173`) cannot read root `public/vite.svg` | Medium | Proxy `/vite.svg` to the Go API in `server.proxy`, or accept the icon only in the Go-served app |
| `public/build` is gitignored, so a fresh clone has no assets | Medium | Document/keep the `npm --prefix web run build` step; CI is deferred |
| Node/npm version mismatch | Low | Toolchain is Node 24 / npm 11, matching the spec |

## Open Questions

1. `./public/vite.svg` is not present in the repo. May Task 2 create it using the official Vite
   logo SVG? (Assumed yes; otherwise provide the file.)
2. Is plain unstyled `hello world` text enough, or should it get minimal centering/typography?
   (Plan assumes plain.)
3. Should the favicon also appear in the Vite dev server at `:5173` via a proxy? (Plan assumes yes.)

## Parallelization Opportunities

- Task 1 and the creation of `public/vite.svg` can happen in parallel; the favicon link and Go
  route (Task 2) depend on both.
- Do not parallelize changes to `internal/routes/web.go` with other routing work.
