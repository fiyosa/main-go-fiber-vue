# Portfolio `/about` — Spec, Plan, Tasks

Status: proposed, awaiting approval.
Source: `C:\Users\fys\Downloads\Documents\CV_Fiter_Yogie_Saputra.md` (ATS CV).
Scope: replace the design-system catalog on `/about` with a portfolio page.
Route `path: 'about'` and `name: 'about'` stay unchanged.

## 1. Spec

### Goal

`/about` becomes a personal portfolio for Fiter Yogie Saputra. Content is
English and derived only from the supplied CV. No invented facts (no phone, no
LinkedIn, no availability claims, no metrics absent from the CV).

### Sections (top to bottom)

1. **Hero** — badge `Fullstack Developer`, `h1` name, one-line positioning,
   location badge `Bandung, Indonesia`, primary `Email me` button
   (`mailto:fiterys@gmail.com`), secondary anchor `See projects` scrolling to
   the projects section.
2. **Summary** — the CV summary paragraph, set at readable measure.
3. **Experience** — two timeline entries, newest first:
   - `Full Stack Software Engineer`, PT Sigma Cipta Caraka, Jan 2022 – present.
   - `Backend Engineer`, PT Daya Desa Asasta, Mar 2021 – Oct 2021.
   Each shows company blurb + the CV bullet list.
4. **Projects** — `KYPAS` and `SSO`, each with description + bullets.
5. **Skills** — six groups: Languages, Frontend, Backend, Databases &
   Infrastructure, DevOps, Product & Integrations.
6. **Education** — State Polytechnic of Bandung, Diploma IV Electronics
   Engineering, GPA 3.60 / 4.00, Aug 2016 – Nov 2020.
7. **Contact CTA** — closing block: `Email me` mailto button + `Bandung,
   Indonesia` badge. No Download CV button.

### Design direction

- Existing Neo-Brutalism system only: `BrutalCard`, `BrutalBadge`,
  `BrutalButton`, semantic tokens, `brutal-press`.
- Oversized headline, sticker badges, hard 3px borders, offset shadows.
- Experience as a vertical timeline: `border-l-[3px] border-fg` rail with
  square markers, one card per role.
- Projects and skills in responsive card grids (`md:grid-cols-2`,
  `lg:grid-cols-3` for skills).
- Single `h1`; each section uses `h2` + `aria-labelledby`.
- No images, icons-only decoration, gradients, or ad-hoc hex.
- Works at 320px without horizontal scroll; both light and dark themes.
- `prefers-reduced-motion` already handled by `brutal.css`.

### Content module

`web/src/pages/guest/content/profile.ts` — typed, frozen data:

```ts
export interface ProfileProject { name: string; description: string; bullets: string[] }
export interface ProfileRole { role: string; company: string; period: string; blurb: string; bullets: string[] }
export interface ProfileSkillGroup { group: string; items: string[] }
export const profile: Profile
```

This keeps `About.vue` presentational and lets tests assert real content.

### Navigation

`AppHeader.vue` links array: label `About` → `Profile`. `name`/`to` unchanged.

### Boundaries

- No `/api/v1` contract, router URL, or route-name change.
- No change outside `web/src/pages/guest` except nothing (nav label lives in
  the guest layout already).
- Brutal tokens remain guest-only; do not edit `main.css` or `brutal.css`.
- Offline: no external assets or fonts.
- `docs/` untouched; no `tests/` or `.github/` directories.

### Testing

- `web/src/pages/guest/About.test.ts`
  - renders the name as the only `h1`
  - has section headings: Experience, Projects, Skills, Education, Contact
  - `Email me` links to `mailto:fiterys@gmail.com`
  - shows `Bandung, Indonesia`
  - shows education period `Aug 2016 – Nov 2020`
  - no design-system catalog copy (`Design system`, `Palette`, `Primitives`)
- `web/src/pages/guest/components/layout/AppHeader.test.ts`
  - nav renders `Profile` and no `About` label
  - `/about` link still targets `/about`

### Success criteria

- `npm --prefix web run test` green (existing 5 + new tests).
- `npm --prefix web run build` passes.
- `npm --prefix web run format:check` passes.
- Isolation grep gate from `plan.md` still returns nothing.
- `/about` renders all sections at 320/768/1280px in light and dark, no
  horizontal scroll, keyboard focus visible.
- No external `http(s)` asset request.

## 2. Plan

1. Add `content/profile.ts` with the CV data and explicit types.
2. Rewrite `About.vue` using `profile` + brutal primitives; keep the page
   self-contained (no `RouterLink` dependency).
3. Rename the nav label in `AppHeader.vue`.
4. Add the two local test files.
5. Verify, then update `todo.md` R4 note and commit scoped files.

Rollback: revert the `About.vue` rewrite and the nav label line; the content
module and tests are additive.

## 3. Tasks

- [x] T1: `content/profile.ts` with hero, summary, experience, projects,
      skills, education, contact data.
- [x] T2: Hero + summary sections (hero "Email me" CTA removed per user).
- [x] T3: Experience timeline section.
- [x] T4: Projects + skills sections.
- [x] T5: Education + contact CTA (email link, lowercase, no uppercase transform).
- [x] T6: Nav label stays `About` (profile tab rename reverted per user).
- [x] T7: Tests (page + header).
- [x] T8: Verify test/build/format + isolation gate, update trackers, commit.

## 4. Open Questions

None. Resolved by the user: full replacement, Nov 2020, nav rename + contact
CTA, no Download CV button.
