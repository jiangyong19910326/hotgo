# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HotGo is a full-stack rapid development framework with a **Go backend** (GoFrame 2.x) and **Vue 3 frontend** (TypeScript + NaiveUI). The backend lives in `server/`, the frontend in `web/`.

## Common Commands

All backend commands run from `server/` via `make`:

```bash
# Development (with hot reload via gf run)
make all        # Start HTTP + Queue + Cron together
make http       # HTTP server only
make queue      # Message queue processor only
make cron       # Scheduled tasks only

# Build
make build      # Compile frontend, copy to server/resource/public/, compile Go binary

# Code generation (requires DB connection)
make dao        # Regenerate DAO/DO/Entity from DB schema
make service    # Regenerate service interfaces from logic

# Linting
make lint       # Run golangci-lint

# Casbin
make refresh    # Refresh Casbin permissions
```

Frontend commands from `web/`:

```bash
pnpm install
pnpm run dev              # Dev server (proxies /admin to localhost:8000)
pnpm run build            # Production build → outputs to server/resource/public/admin/
pnpm run lint:eslint      # Fix ESLint
pnpm run lint:prettier    # Format code
```

## Architecture

### Backend (`server/`)

**Entry & Commands** — `main.go` → `internal/cmd/` defines runnable modes: `Http`, `Queue`, `Cron`, `Auth`, `Tools`, `Upgrade`. Each mode is a separate process.

**Layer stack:**
1. `api/` — Request/response struct definitions (input validation via GoFrame tags)
2. `internal/router/` — Route registration mapping URLs to controllers
3. `internal/controller/` — Thin HTTP handlers; delegates to service
4. `internal/service/` — **Auto-generated** interface files (`gf gen service`); do not edit manually
5. `internal/logic/` — Business logic implementing service interfaces; this is where feature code lives
6. `internal/dao/` — **Auto-generated** DAO objects (`gf gen dao`); do not edit manually
7. `internal/model/` — Input/output structs shared across layers

**Key subsystems:**
- `internal/library/` — Utility packages (JWT, Casbin wrapper, queue clients, storage drivers, etc.)
- `internal/global/` — App-wide singleton initialization
- `internal/middleware/` (under `logic/middleware/`) — Auth, CORS, rate limiting, logging
- `internal/websocket/` — Real-time push notifications
- `internal/crons/` — Cron job definitions
- `internal/queues/` — Message queue consumer handlers

**Plugin system** — `addons/` contains independently-deployable feature modules. Each addon registers its own routes under `admin/`, `api/`, `home/`, and `websocket/` sub-entries via `addons/modules/`.

**Config** — Copy `manifest/config/config.example.yaml` → `manifest/config/config.yaml`. Sections cover DB, Redis, JWT, storage, payments, queues, etc.

### Frontend (`web/`)

**Stack:** Vue 3 + TypeScript + Vite + NaiveUI + Pinia + Vue Router

**Structure:**
- `src/api/` — One file per backend feature area; each wraps axios calls
- `src/views/` — Page components; directory structure mirrors the admin menu tree
- `src/components/` — Shared UI components (tables, forms, upload, etc.)
- `src/store/` — Pinia stores: auth, user info, app/UI state
- `src/router/` — Static routes + dynamic menu-driven routes loaded after login
- `src/hooks/` — Reusable composition functions (pagination, filters, etc.)
- `src/locale/` — i18n strings (zh-CN, zh-TW, en)

**Dev proxy** — `.env.development` proxies `/admin` to `http://localhost:8000/admin`.

### Code Generation Workflow

When adding a new CRUD feature:
1. Create/alter the DB table
2. Run `make dao` to regenerate DAO/model files
3. Add API structs in `api/admin/<feature>/`
4. Add logic in `internal/logic/admin/<feature>/`
5. Run `make service` to regenerate the service interface
6. Register controller + routes in `internal/router/`
7. Generate or write the corresponding Vue page in `web/src/views/`

Alternatively, use the built-in code generator UI (System → Code Generation) to scaffold all layers automatically from the DB schema.

### RBAC & Permissions

Permissions are enforced via **Casbin** (policy stored in DB). Middleware checks JWT → extracts user → verifies Casbin policy for the route. After changing role/menu permissions in the UI, run `make refresh` to reload Casbin policies without restarting.

### Database

Supports MySQL and PostgreSQL. ORM is GoFrame's built-in `gdb`. DAO files are generated and should never be hand-edited. Model structs live in `internal/model/entity/` (generated) and `internal/model/input/` (hand-written input/output types).
