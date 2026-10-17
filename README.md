# morrow-and-muse

Morrow & Muse is a fictional premium lifestyle boutique.

Thoughtfully chosen. Beautifully lived.

## Structure

- `frontend/` — React, TypeScript, and Vite. Brand tokens live in `frontend/src/index.css`. The storefront (`/`, `/shop`, `/collections/:slug`, `/products/:slug`) renders the bundled catalogue immediately. When `VITE_API_BASE_URL` is set, it replaces that catalogue only if both list requests succeed. Unknown addresses render the not-found page in the client. A static host still returns HTTP 200 for those URLs; mapping them to a real 404 belongs to a later deployment and SEO phase.
- `backend/` — Go API. Entry point `backend/cmd/api`. Routes live in `backend/internal/httpapi`. Migrations stay in `backend/cmd/migrate` and are not run by the API.
- `.github/workflows/ci.yml` — frontend typecheck, lint, format, test, and build; backend gofmt, vet, and test against a disposable Postgres 17 service.

## Deployment

Hosting notes for a $0 portfolio demo (Supabase Free, Render Free, Vercel Hobby) are in [docs/deployment.md](docs/deployment.md). That guide is configuration only. It does not create accounts or deploy anything.

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

From the repository root:

```bash
cd frontend
npm ci
npm run dev
```

Vite prints the local URL. The default is `http://localhost:5173/`.

Leave `VITE_API_BASE_URL` unset or blank to keep the bundled catalogue and make no catalogue requests. Set it to the API origin, such as `http://127.0.0.1:8080`, to load `GET /api/v1/collections` and `GET /api/v1/products` once per session. The storefront switches to that response only when both succeed and map onto the storefront types. If either request fails, or the JSON is the wrong shape, both lists stay on the bundle, so a collection list and a product list never come from different sources. Product pages still look up the slug in that one catalogue. Logo, hero, and the featured-slug order stay in the bundle, because the API does not serve them. An outage keeps the bundled pages. It does not show an error screen.

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

`GET /health` does not touch the database. It returns `{"status":"ok"}`.

`GET /ready` pings the database with a one-second timeout. It returns `{"status":"ready"}` or, when the database is unreachable, `503` and `{"error":{"code":"service_unavailable","message":"Not ready."}}`.

Catalogue responses use `{"data": ...}`. Errors use `{"error":{"code":"...","message":"..."}}`. Unknown paths are `404` `not_found`. Other methods are `405` `method_not_allowed`. The API accepts `GET` and CORS preflight `OPTIONS` only.

`GET /api/v1/collections` returns the four collections in display order. `name` is the display title. Each item has `slug`, `name`, `description`, `image` (`src`, `alt`), `seoTitle`, `metaDescription`, and `sortOrder`.

`GET /api/v1/products` returns all 15 products in display order. `GET /api/v1/products/{slug}` returns one. A slug must match `^[a-z0-9]+(-[a-z0-9]+)*$`; a malformed slug is `400` `bad_request` with message `Slug is invalid.` A well-formed unknown slug is `404` `not_found` with message `Product not found.` A product has `id`, `slug`, `name` (the display title), `priceCents`, `collection` (`slug`, `name`), `summary`, `description`, `highlights` (exactly three strings), `cardImage`, `detailImage`, and `sortOrder`.

```json
{"data":{"id":"mm-01","slug":"everyday-ceramic-mug","name":"Everyday Ceramic Mug","priceCents":2800,"collection":{"slug":"morning-ritual","name":"Morning Ritual"},"summary":"...","description":"...","highlights":["...","...","..."],"cardImage":{"src":"https://...","alt":"..."},"detailImage":{"src":"https://...","alt":"..."},"sortOrder":1}}
```

Unexpected database failures are `500` `{"error":{"code":"internal","message":"Something went wrong."}}`. The response and the log omit SQL and connection strings.

Browser callers send an `Origin`. `CORS_ALLOWED_ORIGINS` is a comma-separated allowlist with no wildcard. Local Vite is `http://localhost:5173`. An allowed origin receives `Access-Control-Allow-Origin` set to that origin, `Vary: Origin`, and methods `GET, OPTIONS`. Credentials are not allowed. A missing or unlisted origin gets no allow-origin header. There is no production domain in the code.

`DATABASE_URL` is required. The process exits on startup if it is missing or not a Postgres URL, and the message does not include the value. With `MinConns` at 0 the pool connects on the first query. `PORT` defaults to `8080`.

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
go run ./cmd/migrate up
go run ./cmd/migrate down-to 0
```

`cmd/migrate` reads `DATABASE_URL` from the environment. Do not put the URL on the command line. [docs/deployment.md](docs/deployment.md) shows how to type it without saving it in shell history, and how to check the tables with `psql` without passing the URL.

`00001_schema.sql` creates the tables. `00002_seed.sql` inserts the catalogue. Both have Up and Down. The seed is generated from `frontend/src/data/catalog.ts` and is applied once, as a goose version. Regenerate it, without editing the catalogue file, from the repository root:

```bash
node --experimental-strip-types scripts/export-catalog.ts
node --experimental-strip-types scripts/export-catalog.ts --check
```

Integration tests compare the database and the catalogue API with `backend/db/testdata/catalog_export.json`, which that script also writes. They skip when `TEST_DATABASE_URL` is unset, and fail in CI if it is unset. CI uses a disposable Postgres 17 container, the same major version Supabase currently defaults to. No hosted database is required for this step.

The pool uses pgx simple protocol (`QueryExecModeSimpleProtocol`), so it does not prepare statements. That works with a direct Postgres connection and with Supabase's transaction-mode pooler. Prepared statements need a direct connection or Supabase session mode; this API does not use them. The pool allows at most 4 connections, keeps none idle at startup, closes idle connections after 5 minutes, and limits connection lifetime to 30 minutes. Catalogue queries time out after 3 seconds.
