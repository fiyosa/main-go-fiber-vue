# AGENTS.md

Guidance for AI coding agents in this repository, using OpenCode.

## Global Rules

- Do not read any file that contains "env" in its name or path.
- Write code comments in English.
- This machine has low specifications. Do NOT run build output (binaries produced by `go build`) unless explicitly requested.
- To verify Go code, use only `go vet` or `go build ./...`. Never execute the resulting binary.
- When running `go build`, always write the output binary into `./tmp` (for example `go build -o ./tmp/app.exe .`). Never let build output land in the project root or the working directory.

## Progress Tracking (Todos)

- Post a `todowrite` list for any task with 3+ distinct steps, multiple files, or a debugging loop.
- Keep exactly one item `in_progress`; update it the moment a step starts or finishes.
- Narrate each phase in one line: what you are doing now and what the result was.
- Mark `completed` only after verification (tests, lint, command output). If blocked, state the cause and the next attempt instead of silently retrying.
- Add newly discovered follow-ups to the list; cancel stale items.
- For quick questions or a single trivial edit, skip the list but still say what you are checking.

## Intent → Skill

Map intent to a skill even without a slash command:

- Feature / new functionality → `spec-driven-development`, then `incremental-implementation` and `test-driven-development`
- Planning / breakdown → `planning-and-task-breakdown`
- Bug / unexpected behavior → `debugging-and-error-recovery`
- Code review → `code-review-and-quality`
- Refactoring / simplification → `code-simplification`
- API / interface design → `api-and-interface-design`
- UI/UX → see UI/UX Pro Max below

If a task matches a skill, invoke it. Never implement directly when a skill applies. Follow the skill's instructions exactly; do not partially apply them.

## Workflow

Lifecycle: DEFINE → `spec` · PLAN → `plan` · BUILD → `build` · VERIFY → `test` · REVIEW → `review` · SHIP → `ship`.

Default sequence: `/spec → /plan → /build → /test → /review → /ship`.

Do not skip ahead (no `/build` before a spec exists, no `/ship` before tests pass) unless the user explicitly asks. `/code-simplify` is optional after `/build` and `/test`. `/webperf` runs for browser-facing UI with meaningful visual or data load. `/ship` is mandatory unless the change is 2 files or fewer, under 50 lines, and touches no auth, payments, data access, or config/env.

Each task in `/build` must pass its own test before the next one starts. `/review` runs before every merge, regardless of change size — lighter weight than `/ship`, single persona.

## Rules

- Ignore these thoughts: "too small for a skill", "just implement quickly", "gather context first". Check for a skill first.
- Do not introduce other agent frameworks or workflows. This project uses only CodeGraph, Agent Skills, and UI/UX Pro Max.
- Do not build a router persona that decides which other persona to call; that is the job of slash commands and intent mapping.
- Skills are the _how_ (`.opencode/skills/<name>/SKILL.md`), slash commands are the _when_ (`.opencode/commands/*.md`), personas are review perspectives (`.opencode/agent-skills/*.md`).
- The user or a slash command orchestrates. Personas never invoke other personas; a persona may invoke skills. The only fan-out is `/ship` running `code-reviewer`, `security-auditor`, and `test-engineer` in sequence, then merging their reports.
- All skill, agent, and reference paths inside SKILL.md and command files resolve relative to `.opencode/`.
- Checklists live in `.opencode/references/`: `security`, `performance`, `accessibility`, `testing-patterns`, `observability`, `definition-of-done`, `orchestration-patterns`. Load them when a skill needs detail beyond its SKILL.md.

## CodeGraph

If `.codegraph/` exists at the repo root, use it before grep/find or reading files: the `codegraph_explore` MCP tool, or `codegraph explore "<symbols or question>"` in the shell. Skip CodeGraph entirely if `.codegraph/` is absent. Indexing is the user's decision — never run `codegraph init` yourself.

## UI/UX Pro Max

`ui-ux-pro-max` auto-activates on UI/UX requests ("build a landing page", "create a dashboard", "design a..."). It runs before or during `/spec`, and its design decisions belong in the spec before `/plan`. It does not replace the lifecycle: `frontend-ui-engineering` still guides implementation.

## Caveman

Caveman mode is active unless the user says "stop caveman" or "normal mode". Default mode is full; `/caveman lite`, `/caveman full`, `/caveman ultra`, `/caveman wenyan`, `/caveman wenyan-lite`, and `/caveman wenyan-ultra` switch intensity for the current session. Respond terse, keep technical substance exact, drop filler. Use normal clarity for code, commits, PR descriptions, destructive confirmations, and security findings.

Skills: `caveman`, `caveman-commit`, `caveman-review`, `caveman-help`, `caveman-compress`.
