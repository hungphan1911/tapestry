# Frontend architecture

A single-page app (Vite + React + TypeScript + Tailwind v4) that serves the landing page and every
module UI (personal finance, ...). One app, one build, static `dist/`. Production serving is nginx:
it serves `dist/` with an SPA fallback to `index.html` and proxies `/api` to the Go service.

## Stack

| Concern | Tool |
|---|---|
| Build / dev server | Vite (`/api` proxied to `localhost:8080`, mirroring nginx) |
| UI | React 19, TypeScript, Tailwind v4 (tokens in `src/index.css`), lucide-react icons |
| Routing | React Router (`src/App.tsx`, paths in `src/config/routes.ts`) |
| Server data | TanStack Query (client in `src/api/queryClient.ts`) |
| Tests | Vitest + Testing Library (`*.test.tsx` next to the code) |
| Lint / format | oxlint, prettier (with the Tailwind class-sorting plugin) |

Added when first needed (finance PR): React Hook Form + Zod (forms), TanStack Table (transactions grid),
Recharts (charts, later).

## Layout

```
src/
  main.tsx, App.tsx      entry point; providers + routes
  index.css              Tailwind import, theme tokens, global styles
  api/                   shared plumbing only: apiFetch() wrapper, ApiError, QueryClient
  assets/                images/svgs imported by code (static files go in public/)
  components/ui/         primitives: Button, Card, ... (no feature knowledge)
  components/layout/     AppShell, Sidebar, BottomNav, PageHeader
  config/                env access, route paths
  features/<name>/       everything for one module: api.ts, queries.ts, types.ts, components/, hooks/
  hooks/                 hooks shared by several features
  lib/                   cn(), formatters and other pure helpers
  pages/                 thin route-level components that compose features
  types/                 types shared across features
```

Rules:
- `features/` mirrors the backend's one-folder-per-module rule. **Features do not import each other**;
  anything shared moves to `components/`, `hooks/`, `lib/` or `types/`.
- `pages/` stay thin: layout plus feature components. Logic lives in `features/`.
- Import with the `@/` alias (`@/components/ui/Button`), not long relative paths.
- No global client store. Server state is TanStack Query; local state is `useState`. Add Zustand only if
  genuinely global client state appears.
- No auth (single user, Tailscale only).

## Hooks primer

- `useState`: local UI state (is a dialog open, what is typed in a field).
- `useQuery`: read server data. Gives `data`, `isPending`, `error`, and caches by `queryKey`.
- `useMutation`: write server data. On success, invalidate the related query so lists refresh.
- Custom hooks (`useX`): wrap one piece of logic so components stay small, e.g. `useTransactions(filters)`.
- Rules: call hooks only at the top level of a component or another hook, never in loops or conditions.

Data flow for a feature, bottom to top:
1. `features/<name>/api.ts`: plain async functions calling `apiFetch` (`/api/v1/<module>/...`).
2. `features/<name>/queries.ts`: `useQuery`/`useMutation` hooks wrapping those functions, with query keys defined once.
3. Components call the hooks. **Components never call `fetch` or `apiFetch` directly.**

## Look and feel

- Neutral Monarch-style: warm off-white canvas, white rounded cards, 1px soft borders, almost no shadow.
- **Primary buttons and active states are navy**. Green/red (`text-positive` / `text-negative`) are only for
  positive/negative amounts.
- Font: system UI stack (matches the claude.ai fallback stack; the `anthropic-sans` face itself is proprietary).
- Money uses `tabular-nums`.
- Colors, radius and font are tokens in `@theme` in `src/index.css`. Use token classes (`bg-canvas`,
  `bg-surface`, `border-line`, `text-ink`, `text-muted`, `bg-primary`), never raw hex values or arbitrary colors.

## Responsive: laptop-first

The base (unprefixed) classes are the laptop layout (>= 1024px). Phone and tablet are overrides using
Tailwind's `max-*` variants:

| Device | Width | Variant | Behavior |
|---|---|---|---|
| Laptop (and iPad landscape) | >= 1024px | none | sidebar with labels, multi-column grids |
| iPad portrait | 768-1023px | `max-lg:` | icon-rail sidebar, fewer grid columns |
| iPhone | < 768px | `max-md:` | bottom tab bar, single column, tables become card lists |

Write the laptop layout first, then the phone overrides, then the tablet ones. Use `min-h-dvh`
(not `vh`) and `pb-[env(safe-area-inset-bottom)]` for anything pinned to the bottom of the screen.

## Adding a page or feature

1. Add the path to `src/config/routes.ts` and a `<Route>` in `src/App.tsx`.
2. Add a nav entry in `src/components/layout/navItems.ts` (shows in both sidebar and bottom bar).
3. Create `src/pages/<Name>Page.tsx` using `PageHeader` and `Card`.
4. For data: `features/<name>/{types,api,queries}.ts` first, then components.
5. If the page should appear on the landing page, add it to `modules` in `src/pages/LandingPage.tsx`.

## Commands (run in `frontend/`)

```
npm run dev         dev server on :5173 (proxies /api to :8080)
npm run typecheck   tsc -b
npm run lint        oxlint
npm run test        vitest run
npm run build       typecheck + static build to dist/
npm run format      prettier --write src
```

## Configuration

`.env` is git-ignored; copy `.env.example`. `VITE_API_BASE_URL` is empty by default, meaning relative
`/api` URLs (Vite proxy in dev, nginx in production). Never hardcode hosts.
