# Todo: Neo-Brutalism Frontend UI

> **Revision 2 (2026-11-08) supersedes the scope line below.** Neo-Brutalism is
> confined to `web/src/pages/guest`; the shell becomes a neutral host. See the
> **Revision 2** section at the end of this file. Phase 1–3 are retained as
> completed implementation history and parity baseline.

Scope: restyle **all existing UI** in `web/` (shell, Home, About, 404, shared
primitives) with the Neo-Brutalism design system from ui-ux-pro-max, including
**light + dark** themes. No API, router-behavior, or backend changes (except the
production HTML shell for fonts and theme bootstrap).

**Offline rule:** every runtime asset (fonts, icons, images, CSS) must be a file
inside this repo. No CDN, no Google Fonts, no external `http(s)` requests.

Plan: `.opencode/planning/tasks/plan.md`

Design source: `styles.csv` #38 (Neubrutalism) + #77 (Neo Brutalism Mobile),
`typography.csv` #54/#66 + #17 (Space Grotesk + Space Mono).

Resolved: Space Grotesk + Space Mono (self-hosted woff2); light + dark; all
existing UI restyled; no new routes.

## Phase 1 — Foundation

- [x] Task 1: Tokens, fonts, light + dark global styles
  - [x] `main.css`: Tailwind v4 `@theme` semantic tokens `bg/surface/fg/muted-bg` for light + dark
  - [x] Accent tokens: hot-red/vivid-yellow/soft-violet/mint/sky (constant)
  - [x] `@custom-variant dark` driven by `class="dark"` on `<html>`
  - [x] Self-hosted fonts in `web/src/assets/resources/`: Space Grotesk variable woff2 + Space Mono 400/700 woff2 + OFL license files
  - [x] `@font-face` declarations in `main.css`; CDN `<link>`/`preconnect` removed from both shells
  - [x] Theme `theme-color` metas in `web/index.html` and `internal/app/presentation/core/view/spa.html`
  - [x] Inline theme bootstrap in both shells (localStorage + `prefers-color-scheme`, no flash)
  - [x] `useUIStore`: persisted `theme`, `setTheme`, `toggleTheme`; keep `sidebarOpen`
  - [x] Body theme background/text, zero margin, `:focus-visible` 3px outline
  - [x] Dark palette revision: deep navy `#0F172A/#1B263B/#27334A/#F8FAFC` + very faint dot grid (inline SVG, opacity 0.03, offline); root `App.vue` no longer paints `bg-bg` over the body texture
  - [x] `prefers-reduced-motion` disables press transforms/transitions
  - [x] Vite `base = '/build/'` for production build (fixes lazy chunks + font URLs)
  - [x] Verify: `npm --prefix web run build` pass; fonts emitted to `public/build/assets`; no `fonts.*` domain in bundle; `format:check` pass; `go build`/`go vet` pass
  - [ ] Manual: reload keeps theme, no first-paint flash (Foundation checkpoint)

- [x] Task 2: Brutal UI primitives
  - [x] `components/ui/BrutalButton.vue`: primary/secondary/ghost/danger, sm/md, disabled
  - [x] `components/ui/BrutalCard.vue`: header/body/footer slots, tone + accent props
  - [x] `components/ui/BrutalBadge.vue`: sticker badge, tone + rotation props
  - [x] `components/ui/index.ts` barrel
  - [x] 3px border + offset shadow + mechanical press, both themes
  - [x] Verify: build pass; press shadow no longer overridden by shadow utility
  - [ ] Manual: variants, disabled, keyboard activation, light + dark (Foundation checkpoint)

- [ ] Task 3: App shell and header
  - [ ] `App.vue` layout: header + RouterView + footer
  - [ ] `components/layout/AppHeader.vue`: brand, Home/About, active state, mobile menu
  - [ ] Light/dark toggle wired to `useUIStore` (`aria-label`, `aria-pressed`)
  - [ ] `components/layout/AppFooter.vue`
  - [ ] `useUIStore.sidebarOpen` drives the mobile menu
  - [ ] Verify: build; 320/375/768/1280px; no horizontal scroll; tab order; both themes

### Checkpoint — Foundation

- [ ] Build passes
- [ ] Tokens/fonts verified in dev and production shells
- [ ] Theme toggle persists with no first-paint flash
- [ ] Primitives and shell usable by keyboard in both themes
- [ ] Human review before pages

## Phase 2 — Page Slices

- [x] Task 4: Home (`/`)
  - [x] Sticker hero: badges, oversized headline, CTA to `/about`
  - [x] Stats cards row
  - [x] Pinia counter + toggle demos restyled
  - [x] Health card: checking / error / ok (text + shape, not color only)
  - [x] Remove gradient hero image usage
  - [x] Verify: build; counter/toggle/health flows; 320–1280px; both themes

- [x] Task 5: About (`/about`)
  - [x] Typography specimen
  - [x] Color swatch grid (light + dark values)
  - [x] Live primitive gallery (all variants + badges)
  - [x] Real copy describing FE structure and stack
  - [x] Verify: build; contrast black-on-bright; 320px collapse; both themes

- [x] Task 6: 404
  - [x] Oversized sticker `404`
  - [x] Short explanation + CTA to Home
  - [x] Catch-all route unchanged
  - [x] Verify: build; visit unknown route; CTA works; both themes

### Checkpoint — Pages

- [ ] All routes styled, responsive, and functional
- [ ] Router, Pinia, and health request still work
- [ ] No console errors or overflow at 320–1280px
- [ ] Human review of visual direction in light + dark

## Phase 3 — Polish

- [ ] Task 7: Accessibility, motion, consistency
  - [ ] Contrast: primary ≥ 4.5:1, secondary ≥ 3:1 — light and dark separately
  - [ ] Keyboard reachable, visible focus ring, touch targets ≥ 44px
  - [ ] Reduced motion respected
  - [ ] No emoji icons, gradients, blur, or ad-hoc hex outside tokens
  - [ ] Remove dead assets/code
  - [ ] Verify: `npm --prefix web run format:check`; `npm --prefix web run build`

### Checkpoint — Complete

- [ ] Tasks 1–7 accepted
- [ ] Build + format checks pass
- [ ] Visual direction approved in both themes
- [ ] No backend behavior or `/api/v1` contract changed

## Deferred

- [ ] High-contrast / color-blind theme variant
- [ ] Login or dashboard pages (not in current FE scope)
- [ ] Automated visual regression tests (no browser-test setup in repo)

## Standing Rules

- **Everything local**: all runtime assets live in the repo. No CDN/Google
  Fonts/external `http(s)` URLs in HTML, CSS, or JS.
- Fonts: woff2 in `web/src/assets/resources/` + `@font-face` in `main.css`; keep the
  OFL license files with the fonts.
- Single token source: `web/src/main.css` (`@theme`). No ad-hoc hex in components.
- Components use semantic tokens only; accents stay constant across themes.
- No gradients, blur, or glassmorphism.
- Black text on bright accents in both themes.
- Inline SVG icons; never emoji as structural icons.
- Update both HTML shells when adding head assets or bootstrap scripts.
- Keep Vite `base = '/build/'` for production builds so `/build/assets/*` URLs
  resolve (dev server stays `/`).
- Verify with `npm --prefix web run build`; run `npm run format` after edits.
- Keep `docs/` untouched; do not create `tests/` or `.github/`.

---

## Revision 2 — Confine Neo-Brutalism to `web/src/pages/guest`

Plan: `.opencode/planning/tasks/plan.md` → "Revision 2".
Goal: structural move with **zero visual change**. Only `pages/guest` may own
brutal tokens, fonts, primitives, or chrome. Status: proposed, awaiting approval.

### R1 — Style module, neutral host, layout, nested routes

- [x] Create `web/src/pages/guest/style/brutal.css`
  - [x] `@font-face` for Space Grotesk + Space Mono pointing into guest fonts
  - [x] `@theme inline` mapping `--color-*`, `--font-*`, `--shadow-brutal*` to `var(--guest-*)`
  - [x] Light token values on `.guest-brutal`; navy overrides on `.dark .guest-brutal`
  - [x] Guest-scoped dot grid, `::selection`, `:focus-visible`, `.brutal-press`, reduced-motion override
  - [x] `body:has(.guest-brutal)` background parity; neutral `Canvas` fallback
- [x] Strip tokens/fonts/utilities/dot grid from `web/src/main.css`; keep neutral host reset, `@custom-variant dark`, `color-scheme`, reduced-motion reset
- [x] Reduce `web/src/App.vue` to `<RouterView />`
- [x] Add `web/src/pages/guest/GuestLayout.vue` with `.guest-brutal`, skip link, header, `main#main-content`, footer
- [x] Nest routes in `web/src/pages/routes.ts` under `GuestLayout`; keep paths and names identical
- [x] Verify: no brutal token definitions outside guest; DOM order + route names covered by tests
- [x] Verify: `npm --prefix web run format:check`; `npm --prefix web run build`
- [ ] Manual browser parity check at 320/768/1280px in both themes (theme persistence + no flash)
- [ ] Checkpoint C1: boundary established, zero visual diff

### R2 — Move primitives, chrome, fonts into guest

- [x] Move `components/ui/*` → `pages/guest/components/ui/` (keep barrel)
- [x] Move `components/layout/*` → `pages/guest/components/layout/` (keep barrel)
- [x] Move `assets/resources/*.woff2` + licenses → `pages/guest/assets/fonts/`
- [x] Update imports in Home, About, 404, `GuestLayout`
- [x] Update `@font-face` URLs in `brutal.css`
- [x] Delete `web/src/components/` and `web/src/assets/resources/`
- [x] Verify: no imports of brutal components outside guest; fonts load offline in built bundle; no external `http(s)` URLs
- [x] Verify: `npm --prefix web run format:check`; `npm --prefix web run build`
- [x] Checkpoint C2: structure isolated

### R3 — Home slice parity

- [ ] Visual parity, light + dark: hero, cards, buttons, badges, nav, press states
- [ ] Active nav inversion + theme toggle + mobile menu work
- [ ] Keyboard walkthrough; no horizontal scroll at 320px
- [ ] Verify: `npm --prefix web run format:check`; `npm --prefix web run build`
- [ ] Checkpoint C3a: Home accepted

### R4 — About + 404 slice parity

- [ ] Visual parity, light + dark for `/about` and catch-all 404
- [ ] Unknown URL still resolves to `not-found`; no redirect or URL change
- [ ] 404 CTA returns to `/` without full reload
- [ ] Keyboard walkthrough; no horizontal scroll at 320px
- [ ] Verify: `npm --prefix web run format:check`; `npm --prefix web run build`
- [ ] Checkpoint C3b: About + 404 accepted

### R5 — Isolation gate, cleanup, accessibility audit

- [ ] Isolation grep gate returns nothing outside `web/src/pages/guest`
  - [ ] Run gate from plan; paste output (expected empty)
- [ ] `main.css` has no brutal token, dot grid, or `.brutal-press`
- [ ] Both HTML shells still have theme bootstrap + `theme-color` (no change needed)
- [ ] A11y audit both themes: contrast, focus ring, 44px targets, labels, heading order, reduced motion, interaction states
- [ ] Cleanup: no dead imports, orphan CSS, or unused font files
- [ ] Verify: `npm --prefix web run format:check`; `npm --prefix web run build`
- [ ] Verify: `go vet ./...`; `go build -o ./tmp/ ./...`
- [ ] Checkpoint C4: enforced + green

### Revision 2 standing rules

- Brutal tokens resolve only inside `.guest-brutal`; utilities outside it must
  paint nothing.
- No file outside `web/src/pages/guest` may reference a brutal token utility or
  `brutal-press` (grep gate is the contract).
- HTML-shell theme bootstrap, `theme-color` metas, and `@custom-variant dark`
  stay global by exception.
- `web/src/pages/guest/style/brutal.css` is the single brutal token source.
- Zero visual change: parity beats refactor elegance.
