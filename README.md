# morrow-and-muse

Morrow & Muse is a fictional premium lifestyle boutique.

Thoughtfully chosen. Beautifully lived.

## Structure

- `frontend/` — React, TypeScript, and Vite. Brand tokens live in `frontend/src/index.css`. The storefront (`/`, `/shop`, `/collections/:slug`, `/products/:slug`) is a static catalogue with proposed display prices. It does not call the API. Unknown addresses render the not-found page in the client. A static host still returns HTTP 200 for those URLs; mapping them to a real 404 belongs to a later deployment and SEO phase.
- `backend/` — Go API. Entry point `backend/cmd/api`, router and `GET /health` in `backend/internal/httpapi`.
- `.github/workflows/ci.yml` — frontend typecheck, lint, format, test, and build; backend gofmt, vet, and test.

## Requirements

- Node.js **24.22.0** (Active LTS, Krypton), pinned in `.nvmrc` and `frontend/package.json` `engines`.
- npm **11.19.1**, the npm that ships with Node 24.22.0. The lockfile is `frontend/package-lock.json`.
- Go **1.27.2**, pinned by the `go` line in `backend/go.mod`.

On this cloud agent VM, `/exec-daemon/node` (v22.14.0) is first on `PATH`. That build is older than React Router 8 (Node >= 22.22.0) and jsdom 30 (Node ^22.22.2 or ^24.15.0). Run `nvm use` from the repo root so Node 24.22.0 comes first, then confirm with `node -v`.

The VM image includes Go 1.22.2. Current supported Go releases are 1.26 and 1.27, so the module targets 1.27.2. With `GOTOOLCHAIN=auto` (the default), the Go command downloads that toolchain on first use. CI installs the same version via `actions/setup-go` and `go-version-file: backend/go.mod`.

## Fonts

Headings use DM Serif Display and UI text uses DM Sans, self-hosted with `@fontsource` (`latin` subsets: serif weight 400, sans weights 400, 500, and 600). Those packages keep the CSS family names `DM Serif Display` and `DM Sans`. The variable DM Sans package publishes the face as `DM Sans Variable`, so it is not used. Nothing is requested from Google Fonts at runtime.

## Setup

From the repository root:

```bash
nvm install
nvm use
node -v
```

`node -v` should print `v24.22.0`.

### Frontend

The frontend does not call the API. From the repository root:

```bash
cd frontend
npm ci
npm run dev
```

Vite prints the local URL. The default is `http://localhost:5173/`.

Quality checks, still in `frontend/`:

```bash
npm run typecheck
npm run lint
npm run format:check
npm run test
npm run build
```

### API

From the repository root:

```bash
cd backend
go run ./cmd/api
```

With that process running, the health command is:

```bash
curl -sS http://127.0.0.1:8080/health
```

`GET /health` returns `{"status":"ok"}` with `Content-Type: application/json`.

Backend checks, from `backend/`:

```bash
gofmt -l .
go vet ./...
go test ./...
```

The API reads `PORT` from the environment and listens on `8080` when it is unset. `.env.example` documents that name. The process does not load `.env` files.

### Database

The catalogue schema lives in `backend/db/migrations`. [goose](https://github.com/pressly/goose) `v3.29.0` is pinned in `backend/go.mod` and applied by the `db` package, not by a global binary:

```bash
cd backend
DATABASE_URL='postgres://USER:PASSWORD@HOST:5432/DBNAME?sslmode=require' go run ./cmd/migrate up
DATABASE_URL='postgres://USER:PASSWORD@HOST:5432/DBNAME?sslmode=require' go run ./cmd/migrate down-to 0
```

`00001_schema.sql` creates the tables. `00002_seed.sql` inserts the catalogue. Both have Up and Down. The seed is generated from `frontend/src/data/catalog.ts` and is applied once, as a goose version. Regenerate it, without editing the catalogue file, from the repository root:

```bash
node --experimental-strip-types scripts/export-catalog.ts
node --experimental-strip-types scripts/export-catalog.ts --check
```

Integration tests compare the database with `backend/db/testdata/catalog_export.json`, which that script also writes. They skip when `TEST_DATABASE_URL` is unset, and fail in CI if it is unset. CI uses a disposable Postgres 17 container, the same major version Supabase currently defaults to. No hosted database is required for this step.
