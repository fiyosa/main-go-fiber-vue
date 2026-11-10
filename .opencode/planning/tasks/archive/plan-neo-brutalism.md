# Implementation Plan: Neo-Brutalism Frontend UI

> **Revision 2 (2026-11-08) supersedes the scope decision below.** Neo-Brutalism
> is now confined to `web/src/pages/guest`. See **Revision 2** at the end of
> this file for the active plan. Revision 1 is kept as implementation history.

## Overview

Restyle the existing Vue 3 scaffold under `web/` with a **Neo-Brutalism** design
system sourced from ui-ux-pro-max. The work is UI-only: routing, Pinia stores,
axios wiring, and the Go API contract stay unchanged. The app currently has three
routes (`/`, `/about`, 404) and a minimal shell; this plan turns them into a
complete, consistent, accessible brutalist interface.

Source design references (ui-ux-pro-max CSV data; the Python CLI is unavailable
on this machine, so the rows were read directly):

- `styles.csv` row 38 — **Neubrutalism**: hard black borders, offset shadows,
  flat saturated colors, no gradients/blur.
- `styles.csv` row 77 — **Neo Brutalism (Mobile)**: cream background, sticker
  collage, mechanical press. Dark mode is a project extension: invert the
  semantic surfaces to near-black while keeping accents constant.
- `typography.csv` rows 54/66 — display sans (Space Grotesk / Lexend Mega) and
  row 17 — Space Mono for labels/code.

Existing spec (`.opencode/planning/SPEC.md`) is backend Phase 0 and does not
define frontend visuals; this plan adds the missing UI design layer without
changing backend scope or the `/api/v1` contract.

## Design System (Locked Tokens)

Semantic tokens flip between themes; accent colors stay constant. Every
component reads the semantic tokens, never raw hex.

| Semantic token | Light | Dark | Use |
| --- | --- | --- | --- |
| `--color-bg` | `#FFFDF5` | `#0F172A` | Page background |
| `--color-surface` | `#FFFFFF` | `#1B263B` | Card surface |
| `--color-fg` | `#000000` | `#F8FAFC` | Text, borders, shadows |
| `--color-muted-bg` | `#F2EEDF` | `#27334A` | Subtle blocks |

Dark palette revision: the first pass (`#0B0B0B` / `#171717` / `#232323` /
`#FFFDF5`) read as one flat black slab, then a warm charcoal pass still felt
dull. Final direction is deep navy (`#0F172A` / `#1B263B` / `#27334A` /
`#F8FAFC`) with three clearly separated surface levels and a very faint local
dot grid (inline SVG data URI, offline, opacity 0.03) on `body`.

| Accent token | Value | Use |
| --- | --- | --- |
| `--color-hot-red` | `#FF6B6B` | Danger / accent block |
| `--color-vivid-yellow` | `#FFD93D` | Primary CTA / highlight |
| `--color-soft-violet` | `#C4B5FD` | Secondary surface |
| `--color-mint` | `#4ECDC4` | Success / info |
| `--color-sky` | `#2196F3` | Link / focus accent |

- **Fonts**: `Space Grotesk` (variable, 300–700) for display and `Space Mono`
  400/700 for labels and code, system fallback. **Self-hosted** woff2 files live
  in `web/src/assets/resources/` and are declared with `@font-face` in `main.css`;
  no shell links to an external font host.
- **Borders**: `3px solid var(--color-fg)` on every card, button, input, badge.
- **Shadows**: `6px 6px 0 var(--color-fg)` (large), `4px 4px 0 var(--color-fg)`
  (small). No blur.
- **Radius**: `0`.
- **Motion**: 100–150ms mechanical press; active state `translate(3px, 3px)` and
  reduced shadow. Disabled for `prefers-reduced-motion`.
- **Rules**: no gradients, no blur, no glassmorphism; black text on all bright
  accents in both themes; all shadow/border colors come from `--color-fg`.


## Architecture Decisions

- **Tailwind v4 `@theme`** in `web/src/main.css` is the single token source. No
  new styling dependencies, no `tailwind.config.js` (v4 CSS-first).
- **Components** live under `web/src/components/ui/` (PascalCase SFCs) and
  `web/src/components/layout/`. Each exports through a small `index.ts` barrel,
  matching the existing `lib/index.ts` convention.
- **Shell**: `App.vue` becomes the layout (header + `<RouterView />` + footer);
  the old bare `<RouterView />` is replaced. Header owns navigation and the
  active-link state.
- **Theme**: light + dark. `useUIStore` keeps `theme: 'light' | 'dark'`, plus
  `setTheme`/`toggleTheme` persisted to `localStorage`. Default follows
  `prefers-color-scheme`; `class="dark"` on `<html>` drives the Tailwind v4
  `@custom-variant dark`. Both HTML shells include a tiny inline bootstrap
  script so the correct theme applies before first paint (no flash).
- **Fonts are self-hosted** as woff2 files in `web/src/assets/resources/` and
  declared with `@font-face` in `main.css`. The app makes **no external network
  requests**; Google Fonts CDN links are forbidden. OFL license files ship
  alongside the fonts. System fallbacks remain declared.
- **Production asset base**: Vite `base` is `/build/` for `vite build` (dev
  server stays `/`) because `internal/routes/web.go` serves the bundle only
  under `/build`; this also fixes lazy-route chunk URLs.
- **Icons**: inline SVG only (no emoji, no new icon dependency) to keep the
  bundle small. Sun/moon SVGs for the theme toggle.
- **Hero asset**: stop using the gradient `hero.svg`; replace hero with CSS
  color blocks and sticker badges.
- **No backend changes** except the production HTML shell (`spa.html`) for font
  links and `theme-color`.

## Dependency Graph

```text
main.css tokens + fonts  ──────────────┐
web/index.html + spa.html shells       │
                                       ▼
                          components/ui primitives
                                       │
                                       ▼
                          components/layout/AppHeader
                                       │
                                       ▼
                                  App.vue shell
                                       │
          ┌────────────────────────────┼────────────────────────────┐
          ▼                            ▼                            ▼
   pages/guest/Home.vue        pages/guest/About.vue       pages/guest/404.vue
          │                            │
          └───────────────┬────────────┘
                          ▼
              responsive + a11y + motion pass
```

Implementation order is bottom-up: tokens → primitives → shell → pages → audit.
Pages are the vertical slices; every page task delivers a complete, usable route.

## Task List

### Phase 1 — Foundation

- [x] Task 1: Neo-Brutalism tokens, fonts, and light/dark global styles
  - Description: Rewrite `web/src/main.css` using Tailwind v4 `@theme` with the
    semantic light/dark tokens and accent tokens, `@custom-variant dark`, base
    body styles, `:focus-visible` ring, selection color, and
    `prefers-reduced-motion` rules. Self-host Space Grotesk and Space Mono as
    woff2 in `web/src/assets/resources/` with `@font-face` (no external requests),
    remove CDN links from both HTML shells, set `theme-color` for both themes,
    and add the inline theme bootstrap script (localStorage +
    `prefers-color-scheme`) to both shells. Extend `useUIStore` with persisted
    `theme`, `setTheme`, and `toggleTheme`. Fix Vite `base` to `/build/` for
    production builds.
  - Acceptance criteria:
    - [x] Semantic tokens (`bg`, `surface`, `fg`, `muted-bg`) and accents
      (`hot-red`, `vivid-yellow`, `soft-violet`, `mint`, `sky`) exist in
      `main.css` for both themes.
    - [x] `class="dark"` on `<html>` flips all surfaces, text, borders, and
      shadows; accents stay constant with black text.
    - [x] Body uses the theme background and text colors with no default margin.
    - [x] Focus ring is a visible 3px outline with offset in both themes.
    - [x] Fonts are self-hosted in the repo (woff2 + OFL licenses), declared in
      `main.css`, and no shell references an external font host.
    - [x] Theme toggle state persists across reload; first paint has no
      wrong-theme flash.
    - [x] `prefers-reduced-motion` disables press transforms/transitions.
    - [x] No gradient or blur utility remains in global styles.
  - Verification:
    - [x] `npm --prefix web run build` succeeds (`vue-tsc` + Vite).
    - [x] Built `public/build/index.css` contains the token variables for both
      themes (unused accent tokens are tree-shaken until Task 2 consumes them).
    - [x] Built CSS references only `/build/assets/*.woff2`; no `fonts.*`
      domain appears in the bundle.
    - [x] `go build -o ./tmp/ ./...` and `go vet ./...` pass.
    - [ ] Manual dev check: toggle light/dark, reload, focus ring visible via
      keyboard tab. (deferred to human review at the Foundation checkpoint)
  - Dependencies: None.
  - Files likely touched:
    - `web/src/main.css`, `web/index.html`, `web/vite.config.ts`
    - `web/src/assets/resources/*` (woff2 + licenses)
    - `internal/app/presentation/core/view/spa.html`
    - `web/src/store/ui/index.ts`
  - Estimated scope: Medium (4 files).

- [x] Task 2: Brutal UI primitives
  - Description: Create `BrutalButton.vue` (variants `primary`, `secondary`,
    `ghost`, `danger`; sizes `sm`/`md`; optional block), `BrutalCard.vue`
    (header/body slots, optional accent), and `BrutalBadge.vue` (sticker badge
    with rotation prop). Add `components/ui/index.ts` barrel. All visuals use
    the tokens: 3px border, offset shadow, mechanical press.
  - Acceptance criteria:
    - [x] Button renders correct visual per variant and supports `disabled`.
    - [x] Active/pressed state translates by 3px and reduces the shadow without
      shifting surrounding layout at rest.
    - [x] Card and badge expose named slots/props and render on surface/muted
      backgrounds in both themes.
    - [x] Keyboard `Tab`/`Enter`/`Space` works; focus ring visible.
    - [x] No emoji icons; any icon is inline SVG.
  - Verification:
    - [x] `npm --prefix web run build` succeeds.
    - [ ] Manual dev check: all variants, disabled state, keyboard activation.
      (deferred to human review at the Foundation checkpoint)
  - Notes: press rules moved from `@layer components` to `@layer utilities`
    because the button's `shadow-brutal-sm` utility otherwise won the cascade
    and cancelled the pressed shadow.
  - Dependencies: Task 1.
  - Files likely touched:
    - `web/src/components/ui/BrutalButton.vue`
    - `web/src/components/ui/BrutalCard.vue`
    - `web/src/components/ui/BrutalBadge.vue`
    - `web/src/components/ui/index.ts`
  - Estimated scope: Medium (4 files).

- [ ] Task 3: App shell and header
  - Description: Replace the bare `App.vue` with a shell composed of
    `components/layout/AppHeader.vue` and a minimal footer. Header shows brand,
    `Home`/`About` links with active-route styling, a light/dark toggle wired
    to `useUIStore`, and a menu toggle backed by `useUIStore.sidebarOpen` on
    small screens.
  - Acceptance criteria:
    - [ ] Every route renders inside the shell with a consistent max-width and
      page padding.
    - [ ] Active route is visually distinct (inverted block), not color-only.
    - [ ] Theme toggle switches light/dark, persists, and updates `aria-label`
      + `aria-pressed` correctly in both themes.
    - [ ] Mobile menu opens/closes and is keyboard operable; no horizontal
      scroll at 320px.
    - [ ] Header/footer use only design tokens and inline SVG.
  - Verification:
    - [ ] `npm --prefix web run build` succeeds.
    - [ ] Manual check at 320, 375, 768, and 1280px; tab order is logical.
  - Dependencies: Tasks 1–2.
  - Files likely touched:
    - `web/src/App.vue`
    - `web/src/components/layout/AppHeader.vue`
    - `web/src/components/layout/AppFooter.vue`
  - Estimated scope: Small–Medium (3 files).

### Checkpoint: Foundation

- [ ] `npm --prefix web run build` passes.
- [ ] Tokens and fonts verified in dev and production shells.
- [ ] Primitives and shell usable; keyboard navigation works.
- [ ] Review with human before restyling pages.

### Phase 2 — Vertical Page Slices

- [x] Task 4: Home route (`/`)
  - Description: Rebuild `pages/guest/Home.vue` as a brutalist landing: sticker
    hero (badges, oversized headline, CTA), a three-card stats row, the Pinia
    counter and toggle demos, and the API health card. Keep existing behavior
    (counter/toggle actions, `getHealth()` on mount, loading/error/success
    states). Remove the gradient hero image.
  - Acceptance criteria:
    - [ ] Hero shows the product name, one-line pitch, and a primary CTA linking
      to `/about`; all styling is flat with hard borders/shadows.
    - [ ] Counter and toggle still mutate Pinia state.
    - [ ] Health card renders `checking`, error message, and `ok` states with
      distinct text + shape (not color only).
    - [ ] Layout works at 320/375/768/1280px with no overflow in both themes.
    - [ ] Headings follow one H1 → H2 order.
  - Verification:
    - [ ] `npm --prefix web run build` succeeds.
    - [ ] Manual dev check: increment/decrement/reset, toggle, health states
      (valid and invalid API base).
  - Dependencies: Tasks 1–3.
  - Files likely touched:
    - `web/src/pages/guest/Home.vue`
    - `web/src/assets/img/index.ts` (drop unused hero export if removed)
  - Estimated scope: Medium (2 files).

- [x] Task 5: About route (`/about`)
  - Description: Rebuild `pages/guest/About.vue` as a brutalist documentation
    page: typography specimen, color swatch grid, and a live gallery of the
    primitives (buttons, cards, badges). Includes copy describing the folder
    conventions and stack.
  - Acceptance criteria:
    - [ ] Page documents the token palette and type scale accurately.
    - [ ] Primitive gallery demonstrates every button variant and badge.
    - [ ] Copy is real content (no lorem ipsum) and explains the FE structure.
    - [ ] Contrast is black-on-bright; no white text on yellow/red.
    - [ ] Responsive grid collapses cleanly at 320px.
  - Verification:
    - [ ] `npm --prefix web run build` succeeds.
    - [ ] Manual check: all gallery elements render; links work; keyboard focus
      visible.
  - Dependencies: Tasks 1–3.
  - Files likely touched:
    - `web/src/pages/guest/About.vue`
  - Estimated scope: Small–Medium (1 file).

- [x] Task 6: 404 route
  - Description: Rebuild `pages/guest/404.vue` with an oversized sticker-style
    `404`, a short explanation, and a primary CTA back to `/`. Keep the
    catch-all route in `pages/routes.ts` unchanged.
  - Acceptance criteria:
    - [ ] Unknown routes still resolve to the 404 page.
    - [ ] Big numeral uses the display font and does not overflow at 320px.
    - [ ] CTA returns to Home; focus is visible after navigation.
    - [ ] Route title remains readable and uses semantic `h1` in both themes.
  - Verification:
    - [ ] `npm --prefix web run build` succeeds.
    - [ ] Manual dev check: visit `/does-not-exist`, click CTA, verify Home.
  - Dependencies: Tasks 1–3.
  - Files likely touched:
    - `web/src/pages/guest/404.vue`
  - Estimated scope: Small (1 file).

### Checkpoint: Pages

- [ ] All three routes render fully styled and responsive.
- [ ] Existing interactions (router, Pinia, health request) still work.
- [ ] No console errors or layout overflow at 320–1280px.
- [ ] Human reviews the visual direction before polish.

### Phase 3 — Polish

- [ ] Task 7: Accessibility, motion, and consistency pass
  - Description: Audit every page and shared component against the ui-ux-pro-max
    pre-delivery checklist: contrast, focus order, 44px touch targets, labels,
    reduced motion, heading order, and interaction states — independently in
    light and dark themes. Fix inconsistencies and remove dead code/assets.
  - Acceptance criteria:
    - [ ] Primary text contrast ≥ 4.5:1; secondary ≥ 3:1 on all surfaces in
      both light and dark themes.
    - [ ] Every interactive element is keyboard reachable with a visible focus
      ring; touch targets ≥ 44px.
    - [ ] `prefers-reduced-motion` removes transforms and long transitions.
    - [ ] No emoji structural icons, no gradients/blur, no ad-hoc hex values
      outside `main.css` tokens.
    - [ ] `npm --prefix web run format:check` and `npm --prefix web run build`
      pass.
  - Verification:
    - [ ] Manual keyboard-only walkthrough of all routes.
    - [ ] Manual contrast spot-check on accent surfaces.
    - [ ] `npm --prefix web run format:check`; `npm --prefix web run build`.
  - Dependencies: Tasks 1–6.
  - Files likely touched:
    - `web/src/main.css` and any page/component with findings.
  - Estimated scope: Medium (audit-driven).

### Checkpoint: Complete

- [ ] Tasks 1–7 acceptance criteria met.
- [ ] Build and format checks pass.
- [ ] Visual direction approved by human.
- [ ] No backend behavior or `/api/v1` contract changed.

## Risks and Mitigations

| Risk | Impact | Mitigation |
| --- | --- | --- |
| Production HTML shell (`spa.html`) differs from `web/index.html`, so fonts may load in dev but not prod | High | Task 1 updates both shells and verifies the built CSS/HTML |
| External runtime resources (fonts/CDN) | High | All assets are vendored in the repo; builds must contain no external `http(s)` font URLs |
| Production asset paths (`/assets` vs `/build/assets`) | High | Vite `base = '/build/'` on build; verified in emitted CSS/JS |
| Bright accents (yellow/red) fail WCAG contrast | High | Black text only on bright surfaces; contrast spot-check in Task 7 |
| Mechanical press transforms cause layout jitter | Medium | Use fixed-size buttons and transform only; no margin changes |
| Dark theme contrast on accents and muted surfaces | High | Audit each theme independently; black text only on accents; Task 7 gate |
| Wrong-theme flash before the app hydrates | Medium | Inline bootstrap script in both HTML shells; reload check in Task 1 |
| `prettier-plugin-tailwindcss` reorders classes and breaks custom syntax | Low | Run `npm run format` after each task; build is the gate |
| No automated visual tests | Medium | Manual route checklist at 320/375/768/1280px at each checkpoint |

## Resolved Decisions

- **Theme**: light + dark, persisted per browser. `useUIStore` owns `theme`.
- **Fonts**: Space Grotesk (display) + Space Mono (labels/code), self-hosted in
  the repo. No external font/CDN requests at runtime.
- **Scope**: apply Neo-Brutalism to **all existing UI** — shell, Home, About,
  404, and shared primitives. No new routes.
- **Tokens**: semantic light/dark tokens + constant accents.
- **Assets**: gradient hero image dropped in favor of flat color blocks.

## Open Questions

1. New routes (login/dashboard) are out of scope; they can reuse the primitives
   in a later plan.

## Parallelization Opportunities

- Tasks 4, 5, and 6 (pages) are independent once Phase 1 is done and can run in
  parallel if separate sessions edit disjoint files.
- Task 7 must wait for all pages and is not parallelizable.
- Do not edit `main.css` concurrently with page tasks once tokens are locked.

---

# Revision 2 — Confine Neo-Brutalism to `web/src/pages/guest`

> Status: **proposed, awaiting human approval.** No code changes yet.
> Supersedes the Revision 1 scope ("all existing UI"). Revision 1 remains as
> implementation history; its visuals are the parity target.

## Context

Every current route is a guest route, and the brutalist shell (header, footer,
tokens, fonts, primitives) is global. To confine the style to
`web/src/pages/guest` **without changing a single pixel**, the style system must
move from the app shell into a guest-owned wrapper. The shell becomes a neutral
host; the guest layout owns chrome and tokens. Later non-guest areas can then be
built with their own design language and cannot accidentally inherit brutal
utilities.

Constraints carried over:

- Offline only: no CDN, no external font/asset requests.
- FE-only: no `/api/v1` contract, router URL, or route-name changes.
- Do not edit `docs/`; do not create `tests/` or `.github/`.
- Never read or edit `.env` files.
- Go verification only via `go vet ./...` and `go build -o ./tmp/ ./...`.
- Theme bootstrap in both HTML shells stays global (non-visual infrastructure).

## Scope

**In scope**

- Move brutal tokens, fonts, utilities, primitives, header, and footer under
  `web/src/pages/guest/`.
- Introduce `GuestLayout.vue` as the single brutalist wrapper.
- Nest the three existing routes under `GuestLayout` with unchanged `path` and
  `name`.
- Reduce `App.vue` and `main.css` to a neutral host.

**Out of scope**

- New routes, pages, or features.
- Visual redesign — parity with the current build is the acceptance bar.
- Non-guest styling (there is none yet).
- Backend, API contract, or router URL/behavior changes.
- HTML-shell theme bootstrap and `theme-color` metas (stay global).

## Target Tree

```text
web/src/
  App.vue                        # neutral: <RouterView />
  main.css                       # neutral host only
  main.ts                        # unchanged (imports ./main.css)
  pages/
    routes.ts                    # guest routes nested under GuestLayout
    guest/
      GuestLayout.vue            # .guest-brutal wrapper: skip link, header, main, footer
      Home.vue
      About.vue
      404.vue
      style/
        brutal.css               # @font-face, @theme inline, tokens, scoped base/utilities
      assets/
        fonts/                   # woff2 + OFL licenses (moved from src/assets/resources)
      components/
        layout/
          AppHeader.vue
          AppFooter.vue
          index.ts
        ui/
          BrutalButton.vue
          BrutalCard.vue
          BrutalBadge.vue
          index.ts
```

Removed after the move: `web/src/components/` and
`web/src/assets/resources/` (only if empty after moving fonts/licenses).

## Token Strategy

1. **Neutral host** — `web/src/main.css` keeps only what the document needs and
   what cannot be scoped:
   - `@import 'tailwindcss'`.
   - `@custom-variant dark (&:where(.dark, .dark *))` (the `.dark` class lives on
     `<html>`, outside the guest wrapper).
   - Base host rules: `margin: 0`, `min-height: 100vh`,
     `-webkit-text-size-adjust`, `color-scheme: light`/`.dark { color-scheme: dark }`,
     and the global `prefers-reduced-motion` reset.
   - A neutral body fallback (`background-color: Canvas; color: CanvasText`)
     so a non-guest route renders sanely if added later.
   - **No** `--color-*`, `--font-*`, `--shadow-brutal*`, `.brutal-press`,
     selection color, or dot grid.

2. **Guest tokens** — `web/src/pages/guest/style/brutal.css`:
   ```css
   @font-face { /* Space Grotesk / Space Mono, URL relative to guest/assets/fonts */ }

   @theme inline {
     --font-display: var(--guest-font-display);
     --font-mono: var(--guest-font-mono);
     --color-bg: var(--guest-bg);
     --color-surface: var(--guest-surface);
     --color-fg: var(--guest-fg);
     --color-muted-bg: var(--guest-muted-bg);
     --color-hot-red: var(--guest-hot-red);
     /* ...other accents... */
     --shadow-brutal: var(--guest-shadow-brutal);
     --shadow-brutal-sm: var(--guest-shadow-brutal-sm);
     --shadow-brutal-press: var(--guest-shadow-brutal-press);
   }
   ```
   `@theme inline` makes generated utilities emit `var(--guest-*)` instead of
   copying values, so utilities used outside `.guest-brutal` resolve to an
   undefined custom property and paint nothing.

3. **Values live only on the wrapper**:
   ```css
   .guest-brutal { --guest-bg: #fffdf5; /* ...light values... */ }
   .dark .guest-brutal { --guest-bg: #0f172a; /* ...navy values... */ }
   ```
   Font-family vars resolve inside the wrapper only; `@font-face` stays a
   global registration but its files live under `pages/guest/assets/fonts/`.

4. **Guest-scoped base rules** (inside `brutal.css`):
   - `.guest-brutal { background-color: var(--guest-bg); color: var(--guest-fg); font-family: var(--guest-font-display); -webkit-font-smoothing: antialiased; }`
   - `body:has(.guest-brutal) { background-color: var(--guest-bg); }` keeps the
     overscroll/body paint identical. Fallback when `:has` is unsupported: the
     neutral `Canvas` body is acceptable (not covered by the visual parity gate).
   - `.dark .guest-brutal` dot-grid background.
   - `.guest-brutal ::selection` and `.guest-brutal :focus-visible`.
   - `.brutal-press` and `.brutal-press:active` unchanged, plus its
     reduced-motion override. Keep these selectors under `.guest-brutal` or let
     the class name itself be the contract (only guest markup uses it).

5. **Fallback** if `@theme inline` misbehaves (Tailwind version quirks, missing
   utilities): keep plain custom props on `.guest-brutal` and define the handful
   of needed utilities manually inside `@layer utilities`, scoped by
   `.guest-brutal`. Parity gate decides.

## Dependency Graph

```text
R1 (style module + neutral host + GuestLayout + nested routes)
 └─> R2 (move primitives/chrome/fonts into guest, delete legacy dirs)
      ├─> R3 (Home slice verification)
      └─> R4 (About + 404 slice verification)
           └─> R5 (isolation gate, cleanup, a11y audit)
```

R3 and R4 may run in parallel (disjoint page files) only after R2 lands.
R5 must wait for all pages.

## Tasks

### R1 — Guest style module, neutral host, layout, nested routes

**Description.** Create `pages/guest/style/brutal.css` with `@font-face`,
`@theme inline` mappings, and token values on `.guest-brutal` / `.dark
.guest-brutal`. Strip tokens, fonts, utilities, selection, focus, and dot grid
from `main.css`, leaving the neutral host. Replace `App.vue` with a bare
`<RouterView />`. Add `GuestLayout.vue` that renders the skip link, `AppHeader`,
`<main id="main-content">`, `AppFooter`, and the `.guest-brutal` wrapper.
Convert `routes.ts` to a single parent route with `GuestLayout` and three
children keeping identical `path` and `name`. Update imports in Home/About/404
only if required by the layout change.

**Acceptance criteria**

- [ ] Tokens (`--guest-*` and `--color-*` definitions) exist only under
      `web/src/pages/guest` and inside the `@theme inline` mapping block.
- [ ] `App.vue` contains no layout, chrome, or brutal class.
- [ ] All three routes keep their URL, `name`, and lazy-loaded component.
- [ ] Guest pages render inside `.guest-brutal`; header/footer/main DOM order
      and headings are unchanged.
- [ ] Light and dark themes, theme persistence (`app.theme`), and the
      anti-flash bootstrap behave exactly as before.
- [ ] `npm --prefix web run format:check` and `npm --prefix web run build` pass.

**Verification**

- [ ] Compare `/`, `/about`, and an unknown path against the pre-change build at
      320/768/1280px in both themes; no visible diff.
- [ ] Toggle theme, reload, confirm no flash and persisted choice.
- [ ] Grep: `--color-bg|--color-fg|--color-surface|--color-muted-bg|--color-hot-red|--shadow-brutal|--font-display` outside `pages/guest` returns only the
      `@theme inline` mapping block in `brutal.css` (allowed) and no hits in
      `main.css`.

**Files likely touched:** `web/src/main.css`, `web/src/App.vue`,
`web/src/pages/routes.ts`, new `web/src/pages/guest/GuestLayout.vue`, new
`web/src/pages/guest/style/brutal.css`.

**Estimate:** Medium.

### R2 — Move primitives, chrome, and fonts into guest

**Description.** `git mv` (or copy+delete) `components/ui/*`,
`components/layout/*`, and `assets/resources/*.woff2` + licenses into
`pages/guest/components/{ui,layout}` and `pages/guest/assets/fonts`. Update all
imports (`Home.vue`, `About.vue`, `404.vue`, `GuestLayout.vue`) to relative guest
paths. Keep barrels (`index.ts`) so page imports stay tidy. Delete
`web/src/components/` and `web/src/assets/resources/` when empty. Update the
`@font-face` URLs in `brutal.css`.

**Acceptance criteria**

- [ ] `web/src/components/` no longer exists.
- [ ] `web/src/assets/resources/` no longer exists; fonts and licenses live under
      `pages/guest/assets/fonts/`.
- [ ] No file outside `web/src/pages/guest` imports a brutal component.
- [ ] Barrel exports resolve; no deep relative paths crossing the guest boundary.
- [ ] Fonts still load offline from the built bundle.
- [ ] Format and build pass.

**Verification**

- [ ] `npm --prefix web run build`; inspect emitted CSS for `/build/assets/*.woff2`
      and confirm no external `http(s)` font URL.
- [ ] Dev server: Network panel shows woff2 served locally, zero external
      requests.
- [ ] Spot-check Home/About/404 render after the move.

**Files likely touched:** `web/src/pages/guest/components/**`,
`web/src/pages/guest/assets/fonts/**`, `web/src/pages/guest/**/*.vue`,
`web/src/pages/guest/style/brutal.css`, deletions under `web/src/components/`
and `web/src/assets/resources/`.

**Estimate:** Medium.

### R3 — Home slice verification

**Description.** Verify the Home page after the structural move: hero, cards,
buttons, badges, header nav, active-link inversion, press states, dark mode, and
responsive behavior. Fix only regressions caused by R1/R2 — no redesign.

**Acceptance criteria**

- [ ] Home is pixel-parity with the pre-change build in light and dark.
- [ ] Active nav state, theme toggle, and mobile menu (if present) work.
- [ ] Keyboard order reaches skip link, nav, theme toggle, CTAs.
- [ ] No horizontal scroll at 320px.
- [ ] Build and format pass.

**Verification:** manual route checklist at 320/375/768/1280px in both themes;
keyboard-only walkthrough.

**Dependencies:** R2. **Estimate:** Small.

### R4 — About + 404 slice verification

**Description.** Same parity pass for About and the catch-all 404 route,
including direct navigation to an unknown URL and header/footer consistency.

**Acceptance criteria**

- [ ] About and 404 are pixel-parity with the pre-change build in both themes.
- [ ] Unknown URLs still resolve to `not-found`; no redirect or URL change.
- [ ] 404 CTA returns to `/` without a full reload.
- [ ] No horizontal scroll at 320px.
- [ ] Build and format pass.

**Verification:** manual route checklist at 320/375/768/1280px in both themes;
keyboard-only walkthrough.

**Dependencies:** R2. **Estimate:** Small.

### R5 — Isolation gate, cleanup, accessibility audit

**Description.** Enforce the boundary and finish the residual Revision 1 Task 7
audit. Remove dead CSS/asset references left by the move, then run the isolation
grep gate. Audit contrast, focus visibility, 44px targets, labels, heading order,
reduced motion, and interaction states in both themes.

**Acceptance criteria**

- [ ] Isolation grep gate returns no brutal token utility outside
      `web/src/pages/guest` (see gate below).
- [ ] `main.css` contains no brutal color/font/shadow definitions, no dot grid,
      no `.brutal-press`.
- [ ] Both HTML shells still contain the theme bootstrap and `theme-color` and
      need no change.
- [ ] Task 7 audit items (contrast ≥ 4.5:1 body / 3:1 large, visible focus,
      ≥ 44px touch targets, reduced motion, heading order) pass in both themes.
- [ ] No dead imports, orphan CSS, or unused font files.
- [ ] `npm --prefix web run format:check`, `npm --prefix web run build`,
      `go vet ./...`, and `go build -o ./tmp/ ./...` pass.

**Verification**

```bash
# Isolation gate: brutal utilities/tokens must not appear outside guest.
grep -rnE 'bg-(bg|surface|fg|muted-bg|hot-red|vivid-yellow|soft-violet|mint|sky)|text-(bg|surface|fg|muted-bg|hot-red|vivid-yellow|soft-violet|mint|sky)|border-(bg|surface|fg|muted-bg)|shadow-brutal|brutal-press' \
  web/src --include='*.vue' --include='*.ts' --include='*.css' \
  | grep -v 'web/src/pages/guest'
```

Expected: no output.

**Dependencies:** R2, R3, R4. **Estimate:** Medium.

## Checkpoints

- **C1 — Boundary established (after R1):** neutral host + guest wrapper; routes
  unchanged; no visual diff; theme persistence intact.
- **C2 — Structure isolated (after R2):** only `pages/guest` owns brutal code,
  components, and fonts; legacy dirs deleted.
- **C3 — Parity proven (after R3/R4):** all three routes pixel-parity at
  320/375/768/1280px in both themes.
- **C4 — Enforced (after R5):** isolation grep gate empty; a11y audit passed;
  build/format/Go checks green.

## Risks and Mitigations

| Risk | Impact | Mitigation |
| --- | --- | --- |
| `@theme inline` emits concrete values, leaking tokens globally | High | Inspect emitted CSS in R1; fallback to scoped custom props + manual utilities |
| Tailwind drops utilities whose `@theme` values are `var()` references it cannot see | High | Build after R1; if classes disappear, use the fallback strategy |
| `.dark` lives on `<html>`, so guest tokens need `.dark .guest-brutal` (a plain selector, not the `dark:` variant) | Medium | Author dark overrides in plain CSS; keep `@custom-variant dark` global |
| `body`/overscroll background loses the cream/navy paint | Medium | `body:has(.guest-brutal)` in the guest stylesheet; neutral `Canvas` fallback documented |
| Fonts registered globally via `@font-face` even though files are guest-owned | Low | Accepted: registration is inert without use; grep confirms no non-guest `font-display`/`font-mono` utilities |
| Nested routes change `route.matched` shape and break active-link logic | Medium | `AppHeader` compares `route.name`; verify both nav states in R1/R3 |
| Delete moves drop license files or leave orphan woff2 | Low | Move licenses with fonts; R5 dead-asset check |
| Prettier class reordering touches parity-sensitive markup | Low | Run `npm run format` before comparison; diffs limited to order |
| Manual-only visual parity has no automated guard | Medium | Fixed viewport checklist at each checkpoint; document diffs explicitly |

## Resolved Decisions (Revision 2)

- **Isolation level:** strict — tokens, primitives, chrome, and fonts all live
  under `web/src/pages/guest`.
- **Parity:** visuals must not change; this is a structural move.
- **Token mechanism:** `@theme inline` mapping to `--guest-*`, values declared
  only on `.guest-brutal` and `.dark .guest-brutal`.
- **Routing:** one `GuestLayout` parent with three children; URLs and route
  names unchanged.
- **Global exceptions:** HTML-shell theme bootstrap, `theme-color` metas, and the
  `@custom-variant dark` declaration remain global.
- **`App.vue`:** neutral `<RouterView />` only.

## Open Questions

1. Should the guest style module be imported by `GuestLayout.vue` (loaded with the
   layout chunk) or by `main.ts` (guaranteed but global)? Plan assumes
   `GuestLayout.vue`.
2. Are later non-guest pages expected soon? If yes, R5 can add a short
   `pages/README` boundary note; otherwise the isolation gate is the contract.
