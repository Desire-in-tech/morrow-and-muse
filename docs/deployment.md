# Deployment readiness

This is a configuration review for a non-commercial portfolio demo of Morrow & Muse. The budget is $0. The hosts in use are:

- Supabase Free for Postgres
- Render Free for the Go API, at <https://morrow-and-muse-api.onrender.com>
- Vercel Hobby for the storefront, at <https://morrow-and-muse.vercel.app>

Nothing in this repository creates those accounts, projects, or services, and nothing here deploys. The services above already exist. The tables and the [order of operations](#order-of-operations) below are how to reproduce that setup. Do not add a payment method. Do not enable paid add-ons, including Supabase’s IPv4 add-on. Do not commit a database password, a connection string, or a Supabase host copied from the Connect dialog.

## Current deployed configuration

These are the public origins and the files in git. Dashboard screens were not opened for this note. Requests from 19 October 2026 are in [Checks after hosting exists](#checks-after-hosting-exists). There is no uptime guarantee.

| Piece | Deployed configuration |
| --- | --- |
| Storefront | `https://morrow-and-muse.vercel.app` |
| API | `https://morrow-and-muse-api.onrender.com` |
| Storefront config in git | [`frontend/vercel.json`](../frontend/vercel.json) rewrites to `/index.html` |
| API config in git | [`render.yaml`](../render.yaml): free Go web service `morrow-and-muse-api`, branch `main`, root `backend`, health check `/health` |
| API origin the storefront calls | `https://morrow-and-muse-api.onrender.com`, inlined at build time from `VITE_API_BASE_URL`. It is not written in `vercel.json` or in source. |
| CORS | A preflight from `https://morrow-and-muse.vercel.app` was allowed. The rest of the Render allowlist was not read. |

## What is already true in the code

The API listens on `":" + PORT` (`backend/cmd/api/main.go`). In Go, an empty host means every interface, which is the `0.0.0.0` bind Render requires. `PORT` defaults to `8080` when unset. Render sets `PORT` itself (the platform default is `10000`). No listen-address change was required.

`GET /health` does not touch the database. `GET /ready` pings it. The API does not run migrations on startup.

`DATABASE_URL` is required. `pgxpool.ParseConfig` reads it after trimming whitespace (`backend/internal/dbpool/pool.go`). pgx accepts `postgres://` and `postgresql://` URLs. Startup errors do not include the URL.

The API pool sets `DefaultQueryExecMode` to `QueryExecModeSimpleProtocol`, so catalogue queries do not use prepared statements. The pool allows at most 4 connections.

Migrations are a separate process: `go run ./cmd/migrate` in `backend/`. That command uses `database/sql` with the pgx stdlib driver and does not set an exec mode. pgx’s default exec mode is `QueryExecModeCacheStatement` (prepared statements). goose runs each SQL migration in a transaction unless the file contains `-- +goose NO TRANSACTION`. These migration files do not contain that annotation. goose’s session lock is off unless `WithSessionLocker` is set, and `cmd/migrate` does not set it, so goose does not take a Postgres advisory lock. Prepared statements are the session feature that matters: Supabase transaction mode does not support them.

CORS is an exact allowlist (`CORS_ALLOWED_ORIGINS`), comma-separated. There is no wildcard and no credential header. Allowed methods are `GET` and `OPTIONS`. An entry must be `http` or `https`, must have a host, and must not have userinfo, a path, a query, or a fragment. `https://example.vercel.app/` is rejected because of the trailing slash, and the process will not start. The browser’s `Origin` must match the configured string exactly, including the scheme.

`VITE_API_BASE_URL` is read from `import.meta.env` in `frontend/src/api/catalogClient.ts`. Vite inlines it at build time. Unset or blank means the storefront makes no catalogue request and keeps the bundled catalogue.

## Render

A `render.yaml` blueprint is in the repository root because it records the free plan, the monorepo root, the build and start commands, the health check, and the secret names without values. Prose alone would leave `plan: free` as a click that is easy to miss. Syncing the blueprint in the Render dashboard is what creates or updates the service. The service `morrow-and-muse-api` already exists. This file does not sync itself and does not deploy. Sync the blueprint only when you intend to create or update that free service from `main`.

| Setting | Value |
| --- | --- |
| Name | `morrow-and-muse-api` |
| Runtime | `go` |
| Plan | `free` |
| Branch | `main` |
| Root directory | `backend` |
| Build command | `go build -o bin/api ./cmd/api` |
| Start command | `./bin/api` |
| Health check path | `/health` |

Render runs the build and start commands with the root directory as the working directory, and files outside that directory are not available at build or runtime ([Monorepo Support](https://render.com/docs/monorepo-support)). The Go module lives entirely under `backend/`, so `rootDir: backend` is the right root. `bin/` is gitignored; the binary is produced on Render, not committed.

`healthCheckPath` must start with `/` ([Blueprint spec](https://render.com/docs/blueprint-spec)). `/health` is used so a paused database does not fail the platform health check. Use `/ready` yourself when you want to know that Postgres answered.

`DATABASE_URL` and `CORS_ALLOWED_ORIGINS` are `sync: false`, which makes the dashboard ask for a value and keeps that value out of the file ([Blueprint spec](https://render.com/docs/blueprint-spec)). `PORT` is omitted so Render can inject it. The service must bind host `0.0.0.0` on that port ([Web Services](https://render.com/docs/web-services)).

### Go version

`backend/go.mod` says `go 1.27.2`. Render’s native Go runtime tracks the latest stable Go 1.x and does not let you pin a version unless you deploy a Docker image ([Language Support](https://render.com/docs/language-support)). This phase does not add an image. With `GOTOOLCHAIN=auto`, which is the Go default, `go build` downloads the version in `go.mod` when the builder’s Go is older but still new enough to read the toolchain line. Whether Render’s builder is already on 1.27.x is unverified until a real build log exists. Do not set `GO_VERSION` expecting it to pin Go.

### Environment variables on Render

Enter these in the Render dashboard (the blueprint prompts for them). Do not put the values in git.

| Variable | Where | Value |
| --- | --- | --- |
| `DATABASE_URL` | Render → service → Environment | Session-pooler URI from Supabase, with `sslmode=require`. See below. |
| `CORS_ALLOWED_ORIGINS` | Render → service → Environment | Exact origins, no trailing slash. The deployed storefront origin is `https://morrow-and-muse.vercel.app`. For a new setup, start with `http://localhost:5173`, then replace it with the production origin and save so Render restarts the service. The rest of the dashboard value was not read. |
| `PORT` | Injected by Render | Do not set it. |

## Supabase

Create the project yourself, on the Free plan, in the Supabase dashboard. Choose the smallest free settings the dashboard offers. Do not enable the IPv4 add-on. Do not create the project from this repository.

Connection facts below are from [Connect to your database](https://supabase.com/docs/guides/database/connecting-to-postgres). Hostnames are copied from the dashboard’s Connect dialog. Do not invent a pooler host: the shared-pooler host contains an index that is not implied by the region name.

| Mode | Host and port | Free plan | Use it for |
| --- | --- | --- | --- |
| Direct | `db.[PROJECT-REF].supabase.co:5432` | IPv6 | A machine that actually has IPv6. Not the default for Render. |
| Session pooler | `aws-[INDEX]-[REGION].pooler.supabase.com:5432` | IPv4 | Migrations, and the Render API. |
| Transaction pooler | same pooler host, port `6543` | IPv4 | Not migrations. Optional for the API only, because this API disables prepared statements. |

The shared pooler username is `postgres.[PROJECT-REF]`. The direct username is `postgres`. The pooler password is the database password. Percent-encode reserved characters in the password (`&`, `#`, `?`, space, and so on). Append `sslmode=require` if the copied URI does not already have it. Supabase documents `sslmode=require` as the setting that refuses a plaintext fallback ([SSL](https://supabase.com/docs/guides/database/connecting-to-postgres#connecting-with-ssl)).

Supabase’s own table says a persistent backend on an IPv4-only network should use the shared pooler in session mode, and that migrations prefer a direct connection. The same page says to use session mode instead of direct when the network is IPv4-only and the project does not have the IPv4 add-on. Whether a Render free web service can open IPv6 to `db.[PROJECT-REF].supabase.co` was not tested. The $0 choice that does not depend on that answer is the session pooler.

Transaction mode does not support prepared statements or other session state ([Transaction mode limitations](https://supabase.com/docs/guides/database/connecting-to-postgres)). goose, through pgx’s default exec mode, does use prepared statements, and it wraps each of these migrations in a transaction. Use port `5432` (session pooler, or direct if the laptop has IPv6) for `cmd/migrate`. Do not point goose at port `6543`.

The running API can speak transaction mode because `dbpool` forces the simple protocol. Session mode is still the recommendation for Render: it matches Supabase’s persistent-backend row, and the same class of URI is valid for migrations. A later change that turned prepared statements back on would break a transaction-mode URL and would not break session mode.

### Migrate from a laptop

Run this only when you have a real Supabase project. The connection string is a secret. A shell history file, the process list (`ps` and `/proc/<pid>/cmdline`), a terminal log, and a shared screen all keep what you type on a command line. `DATABASE_URL='postgresql://...' go run ...` and `psql 'postgresql://...'` put the URL there. Type it at a hidden prompt, leave it in the environment for the one command, then remove it. Do not paste it into an issue, a CI log, or a file in this repository.

`cmd/migrate` accepts `up` or `down-to` and reads the URL with `os.Getenv("DATABASE_URL")`. There is no URL flag. A successful run prints nothing. On failure it prints `migrate:` and the error text. This command does not log the URL. goose is given an open database handle, not the URL string. If pgx cannot parse the URL, its error text redacts a recognizable password and can still include the host; a malformed URL can be too ambiguous to redact completely. Do not paste that output into an issue.

#### macOS and Linux (bash or zsh)

`read -rs` is raw (`-r`, so a backslash in the password stays) and silent (`-s`). Bash accepts the combined flag. Zsh does too: flags that take no argument combine, `-r` is raw mode, and `-s` does not echo when the input comes from the terminal ([zsh `read`](https://zsh.sourceforge.io/Doc/Release/Shell-Builtin-Commands.html)). With one variable, both shells assign the whole line, including any spaces. Percent-encode reserved characters in the password before you paste. Turn xtrace off first (`set +x` in bash, `unsetopt xtrace` in zsh); xtrace prints assignments.

```bash
cd backend
printf 'DATABASE_URL (input hidden): '
read -rs DATABASE_URL
printf '\n'
export DATABASE_URL
go run ./cmd/migrate up
unset DATABASE_URL
```

Paste the session-pooler URI at the hidden prompt: port `5432` and `sslmode=require`. The paste is not a shell command, so it is not written to history. `export` and `go run` do not receive the URL as an argument. `unset` removes it from this shell. Run this from `backend/`, where `go.mod` and `cmd/migrate` live. CI and the Migrate database workflow use that same working directory, and the workflow runs `up` only.

`go run ./cmd/migrate down-to 0` is not `up`. After a full `up`, goose runs the Down sections of `00003`, then `00002`, then `00001`, and stops at version 0. Version 2 deletes the catalogue rows. Version 1 drops `featured_products`, `product_highlights`, `products`, `collections`, and `catalog_assets`, their triggers, and the three trigger functions. `goose_db_version` remains. Use that only on a disposable local database. Never run it against the hosted database to experiment.

`go run ./cmd/migrate down-to 2` rolls back only version 3. The tables and the seed stay. The Down section of `00003_lock_down_public.sql` disables row level security on the six tables; where `anon` and `authenticated` exist, grants them `ALL` on those tables, `ALL` on sequences in `public`, and `EXECUTE` on the three trigger functions, and restores their schema-public default privileges; grants `EXECUTE` on those functions to `PUBLIC`; and replaces the functions with the version-1 bodies (`search_path` reset). It does not grant table or sequence privileges back to `PUBLIC`. That direction is also only for a disposable local database.

The bash sequence above was run against a disposable local Postgres (migrate exited 0, then `unset` cleared the variable). Zsh was not installed in that environment; the zsh flags are from the manual linked above.

#### Windows PowerShell

`Read-Host -AsSecureString` masks each typed character with `*`. The URL is not part of the command, so PSReadLine does not store it. This cmdlet exists in Windows PowerShell 5.1 and in PowerShell 7.

PowerShell 7.0 and later can recover the text with `-AsPlainText` ([ConvertFrom-SecureString](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.security/convertfrom-securestring?view=powershell-7.5)). Windows PowerShell 5.1 has no `-AsPlainText`; it needs the Marshal calls below. Assign `$env:DATABASE_URL` only. Do not call `[Environment]::SetEnvironmentVariable` with a `User` or `Machine` target, which would keep the URL after the window closes. These snippets were checked against that documentation and were not executed here.

PowerShell 7:

```powershell
Set-Location backend
$secure = Read-Host -AsSecureString -Prompt 'DATABASE_URL'
$env:DATABASE_URL = ConvertFrom-SecureString -SecureString $secure -AsPlainText
go run ./cmd/migrate up
Remove-Item Env:DATABASE_URL
```

Windows PowerShell 5.1:

```powershell
Set-Location backend
$secure = Read-Host -AsSecureString -Prompt 'DATABASE_URL'
$ptr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure)
try {
  $env:DATABASE_URL = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($ptr)
} finally {
  [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($ptr)
}
go run ./cmd/migrate up
Remove-Item Env:DATABASE_URL
```

`Remove-Item Env:DATABASE_URL` drops the variable from this process. Close the window afterwards.

#### Check with psql

Do not pass the URI to `psql`. Set the non-secret libpq variables and let `psql` prompt for the password. `-W` forces the prompt. The password is not an argument. Do not set `PGPASSWORD`; it is visible in the process list.

```bash
export PGHOST='aws-[INDEX]-[REGION].pooler.supabase.com'
export PGPORT=5432
export PGUSER='postgres.[PROJECT-REF]'
export PGDATABASE=postgres
export PGSSLMODE=require
psql -W \
  -c 'select version_id, is_applied from goose_db_version order by id;' \
  -c 'select count(*) from collections;' \
  -c 'select count(*) from products;'
unset PGHOST PGPORT PGUSER PGDATABASE PGSSLMODE
```

Copy `PGHOST` from the Connect dialog. `PGSSLMODE=require` is the minimum. `verify-full`, with `PGSSLROOTCERT` pointing at the root certificate from Database settings in the Supabase dashboard, also checks the server name ([SSL](https://supabase.com/docs/guides/database/connecting-to-postgres#connecting-with-ssl)). Do not use `disable`, `allow`, or `prefer`.

`goose_db_version` has `id`, `version_id`, `is_applied`, and `tstamp`. goose inserts `version_id` 0 when it creates that table. This module then applies version 1 (schema), version 2 (seed), and version 3 (public-schema lock-down). Expect `0`, `1`, `2`, and `3`, each with `is_applied` true, plus 4 collections and 15 products. Version 3 does not change those rows. The queries contain no secrets.

To skip the prompt, one line in `~/.pgpass` (Windows: `%APPDATA%\postgresql\pgpass.conf`) is enough:

```text
[POOLER-HOST]:5432:postgres:postgres.[PROJECT-REF]:[PASSWORD]
```

On macOS and Linux, `chmod 600` that file. Escape `:` and `\` in the password with a backslash. Delete the file when you are finished.

### Migrate from GitHub Actions

Use this when the computer you are on cannot open TCP port 5432 to Supabase. `.github/workflows/db-migrate.yml` is manual. It does not run on push or pull request.

Add the secret once, in the repository: Settings → Secrets and variables → Actions → New repository secret. Name it `DATABASE_URL`. The value is the session-pooler URI, port 5432, with `?sslmode=require`. Do not use port 6543. If that password was pasted into a chat, a shell, or an issue, reset it in the Supabase dashboard first, then store the new URI in this secret and in Render.

Run it from the Actions tab: Migrate database → Run workflow → branch `main` → type `migrate` in the confirm box. Type that word only. The run can show inputs, so the box is not the place for the connection string. The job rejects anything except that word, checks that the secret is a `postgres` or `postgresql` URL with a host, a user, a password, and `sslmode` of `require`, `verify-ca`, or `verify-full`, refuses port 6543, warns if the port is not 5432, and runs `go run ./cmd/migrate up` in `backend/`. The URL is passed through the environment and is not printed. A green run means goose finished.

Verify in the Supabase SQL editor, or with the `psql` `PG*` commands above, that `goose_db_version` has versions 0, 1, 2, and 3 applied, and that the counts are 4 collections and 15 products.

### Public schema lock-down

The storefront reads the catalogue through the Go API. Migration `00003_lock_down_public.sql` enables row level security on `collections`, `products`, `product_highlights`, `featured_products`, `catalog_assets`, and `goose_db_version`, and it creates no policies. `anon` and `authenticated` lose every privilege on those tables, on sequences in `public`, and on the trigger functions `set_updated_at`, `products_have_three_highlights`, and `product_highlights_stay_three`. `PUBLIC` loses `EXECUTE` on those functions.

The API connects as the table owner (`postgres` on Supabase). Row level security is enabled without `FORCE ROW LEVEL SECURITY`, so the owner bypasses it and goose can still write `goose_db_version`. On Supabase, `postgres` also has `BYPASSRLS`. `service_role` is unchanged: it is server-only and already bypasses row level security.

`supabase_admin` default privileges cannot be changed by the migration role. Objects this project creates are owned by `postgres`. Version 3 revokes that role's defaults that would grant new tables, sequences, and functions in `public` to `anon` and `authenticated`. It leaves the built-in `EXECUTE` grant to `PUBLIC` on functions created later, because a database-wide revoke would cover every schema. Each later migration that creates a table in `public` enables row level security, and each later migration that creates a function in `public` revokes `EXECUTE` from `PUBLIC`, `anon`, and `authenticated`. `TestSecurityHardening` fails if a migrated object in `public` skips that. A function created in `public` outside goose, for example from the Supabase SQL editor, stays executable by `PUBLIC` until it is revoked. `service_role` keeps its default grants.

Running Down for version 3 does what [Migrate from a laptop](#migrate-from-a-laptop) describes for `down-to 2`. Leave that direction unused on the hosted database.

Version 3 is on `main`. A hosted database that already has versions 0, 1, and 2 picks it up only when someone runs the [Migrate database](#migrate-from-github-actions) workflow (confirm input `migrate`). The workflow file runs `up` only. After that run, the Supabase security advisor should clear `rls_disabled_in_public` for these six tables and `function_search_path_mutable` for the three trigger functions. This repository cannot rerun that advisor. Catalogue responses show the schema and the seed. They do not show privileges.

A project-run read-only check through the Supabase management connection on 19 October 2026 is a historical record, not a result of this review, and not an independent audit. This review did not re-query the database. That check reported goose versions 0, 1, 2, and 3 applied, with version 3 at 18:04:47 EAT (15:04:47 UTC); row counts of 4 collections, 15 products, 45 product highlights, 4 featured products, and 2 catalogue assets; all 36 Cloudinary URLs matching `frontend/src/data/catalog.ts`; row level security enabled and not forced on the six tables, with no policies; `PUBLIC`, `anon`, and `authenticated` revoked on the tables, sequences, and the three trigger functions; those functions having `search_path` set to empty. The advisors then showed `rls_enabled_no_policy` (INFO, expected) and `unused_index`. The live API was not part of that check. The HTTP checks from the same date are in [Checks after hosting exists](#checks-after-hosting-exists).

GitHub Actions records a separate fact: the Migrate database workflow, `workflow_dispatch` on commit `135578461356cbf1d4940dd8e36a1300ed6cf849`, finished successfully at 15:04:51 UTC on 19 October 2026. That job runs `up` and does not print the database host. It is not a copy of the SQL result above.

Limits that check did not remove, and that the migration does not remove: `anon` still has `USAGE` on schema `public`; `supabase_admin` default privileges are unchanged; a function created by hand in the SQL editor is executable by `PUBLIC` until it is revoked.

### Separate read-only login (not built)

The API uses one `DATABASE_URL` and connects as the owner. Supabase’s connection guide says a custom role through the shared pooler uses the username `[ROLE].[PROJECT-REF]`. Creating that role is SQL, not a paid add-on named on that page. This review did not create one.

A non-owner role does not bypass row level security unless it has `BYPASSRLS`. With row level security enabled and no policy, that role sees no rows even if it holds `SELECT`. Assigning `BYPASSRLS` was not tested here. The design that matches the current policies (there are none) is a `SELECT` policy, not that attribute.

The catalogue queries select from `collections`, `products`, and `product_highlights` only (`backend/internal/catalog/postgres.go`). They do not read `featured_products`, `catalog_assets`, or `goose_db_version`. `/ready` pings the pool. It does not need a table grant.

A later change would be a new goose migration that creates the login, grants `USAGE` on `public` and `SELECT` on those three tables, and adds a `SELECT` policy on each; tests that the role can read those tables and cannot write or read the others; a Render `DATABASE_URL` whose user is `[ROLE].[PROJECT-REF]`, session pooler, port 5432, `sslmode=require`; and the existing migrate secret left on the owner. Do not point goose at the reader role. Do not do that work by editing the hosted database outside a migration.

## Vercel

The Hobby project is already served at `https://morrow-and-muse.vercel.app`. To reproduce it from this Git repository, create the project on the Hobby plan. In Project Settings:

| Setting | Value |
| --- | --- |
| Framework Preset | Vite |
| Root Directory | `frontend` |
| Install Command | `npm ci` |
| Build Command | `npm run build` (`tsc -b && vite build`) |
| Output Directory | `dist` |
| Node.js Version | 24.x, matching `frontend/package.json` `engines` (`24.22.0`) |

The repo `.nvmrc` is at the git root. With Root Directory set to `frontend`, that file is outside the Vercel project root, so set the Node.js version in Project Settings.

`frontend/vercel.json` is the project config once the root directory is `frontend`. It is the SPA rewrite from Vercel’s own example ([rewrites](https://vercel.com/docs/project-configuration/vercel-json#rewrites)):

```json
{
  "rewrites": [{ "source": "/(.*)", "destination": "/index.html" }]
}
```

Vercel gives filesystem routes precedence over rewrites, so a request that matches a built file (`/assets/…`, `/favicon.svg`) is not rewritten to `index.html` ([rewrites](https://vercel.com/docs/project-configuration/vercel-json#rewrites): “precedence is given to the filesystem prior to rewrites being applied”). Direct loads of `/shop`, `/collections/:slug`, and `/products/:slug` receive `index.html`, and React Router continues from there.

A Framework Preset of Other can stop this rewrite from taking effect ([Why is my deployed project giving 404?](https://vercel.com/kb/guide/why-is-my-deployed-project-giving-404)). Leave the preset on Vite.

Unknown client routes also receive `index.html` with HTTP 200. The not-found page is still rendered in the browser. A real HTTP 404 for unknown URLs is unchanged and is a later SEO concern.

Do not set a custom domain in this phase. Do not put a Render hostname in `vercel.json` or in source. The production bundle inlines `VITE_API_BASE_URL`, so the public API origin is present in the built JavaScript. On 19 October 2026 that bundle contained the string `morrow-and-muse-api.onrender.com`. The asset filename changes with each build.

### `VITE_API_BASE_URL`

Add it under Vercel → Project → Settings → Environment Variables, for the Production environment, before the production build:

`VITE_API_BASE_URL` = `https://morrow-and-muse-api.onrender.com`

No trailing slash. The client strips trailing slashes, but the value is still baked into the bundle. Changing it requires a new Vercel deployment. If the variable is absent or blank, the storefront never calls the API and serves the bundled catalogue.

Preview deployments are a different origin from production. Leave the variable unset on Preview unless you have also added that exact preview origin to `CORS_ALLOWED_ORIGINS`.

## CORS

The production origin is `https://morrow-and-muse.vercel.app`. Set Render’s `CORS_ALLOWED_ORIGINS` to that origin and save so the service restarts. On 19 October 2026 a preflight from that origin returned 204 and echoed it.

Exact format: `https://<project>.vercel.app`

No trailing slash, no path, scheme `https`. `http://` and `https://` are different origins. A Vercel preview host such as `https://<project>-git-<branch>-<scope>.vercel.app` does not match the production origin and is not allowed unless you add that full origin as another comma-separated entry. The code rejects `*`.

Credentials stay off. The storefront only sends `GET`. An allowed preflight returns 204 with `Access-Control-Allow-Origin` set to the request origin. An unlisted origin gets no allow-origin header; a preflight from an unlisted origin is 403.

## Order of operations

The demo at the two origins above already exists. These steps are how to reproduce that setup.

1. Create the Supabase Free project. Copy the session-pooler URI (port 5432) and add `sslmode=require`. Do not enable the IPv4 add-on.
2. From `backend/`, run the goose migrations with `go run ./cmd/migrate up`. Enter `DATABASE_URL` at the hidden prompt in [Migrate from a laptop](#migrate-from-a-laptop), then unset it. Use the session pooler on port 5432 with `sslmode=require`, and check the tables as that section describes. When the laptop cannot open TCP 5432, use the [Migrate database](#migrate-from-github-actions) workflow instead (`workflow_dispatch`, confirm input `migrate`).
3. Create the Render free web service from `render.yaml` (or enter the same settings by hand). Set `DATABASE_URL` to that session URI. Set `CORS_ALLOWED_ORIGINS` to `http://localhost:5173`.
4. Create the Vercel Hobby project with the settings in the table above. Set `VITE_API_BASE_URL` to the Render origin (`https://<service>.onrender.com`, no trailing slash) and deploy.
5. Replace Render’s `CORS_ALLOWED_ORIGINS` with `https://<project>.vercel.app` (no trailing slash) and save so Render restarts.

## Checks after hosting exists

Run on 19 October 2026 against the deployed hosts. These are single requests from that day, not an uptime or production-readiness result.

```bash
curl -sS -o /dev/null -w '%{http_code}\n' "https://morrow-and-muse.vercel.app/"
curl -sS -o /dev/null -w '%{http_code}\n' "https://morrow-and-muse.vercel.app/shop"
curl -sS -o /dev/null -w '%{http_code}\n' "https://morrow-and-muse.vercel.app/products/everyday-ceramic-mug"
curl -sS -o /dev/null -w '%{http_code}\n' "https://morrow-and-muse.vercel.app/this-route-does-not-exist"
curl -sS "https://morrow-and-muse-api.onrender.com/health"
curl -sS "https://morrow-and-muse-api.onrender.com/ready"
curl -sS "https://morrow-and-muse-api.onrender.com/api/v1/collections"
curl -sS "https://morrow-and-muse-api.onrender.com/api/v1/products"
curl -sS "https://morrow-and-muse-api.onrender.com/api/v1/products/everyday-ceramic-mug"
curl -sS -D - -o /dev/null -X OPTIONS \
  "https://morrow-and-muse-api.onrender.com/api/v1/collections" \
  -H "Origin: https://morrow-and-muse.vercel.app" \
  -H "Access-Control-Request-Method: GET"
```

That day:

- The storefront returned HTTP 200 for `/`, `/shop`, `/products/everyday-ceramic-mug`, and `/this-route-does-not-exist`. The unknown path is the SPA rewrite. The client still renders the not-found page. A collection path was not requested in this pass.
- `/health` returned HTTP 200 `{"status":"ok"}`. The free service had been idle, and that first request took about 12 seconds.
- `/ready` returned HTTP 200 `{"status":"ready"}`.
- `/api/v1/collections` returned HTTP 200 and 4 collections. `sortOrder` was 1 `morning-ritual`, 2 `workday-edit`, 3 `everyday-carry`, 4 `quiet-space`.
- `/api/v1/products` returned HTTP 200 and 15 products. Each product had 3 highlights. `everyday-ceramic-mug` had `priceCents` 2800.
- `/api/v1/products/everyday-ceramic-mug` returned HTTP 200.
- `/api/v1/products/Bad_Slug` returned HTTP 400 `{"error":{"code":"bad_request","message":"Slug is invalid."}}`.
- `/api/v1/products/not-a-real-product` returned HTTP 404 `{"error":{"code":"not_found","message":"Product not found."}}`.
- `POST /api/v1/products` returned HTTP 405 `{"error":{"code":"method_not_allowed","message":"Method not allowed."}}`.
- The allowed preflight returned HTTP 204, `access-control-allow-origin: https://morrow-and-muse.vercel.app`, methods `GET, OPTIONS`, headers `Accept, Content-Type`, and `access-control-max-age: 600`.

A free Render service that has been idle can show Render’s loading page for about a minute on the first request. The wake observed above was about 12 seconds. That is one sample, not a limit. Wait and retry `/health` before treating a slow first response as an application failure.

## Public deploy records

Checked on 20 October 2026 from GitHub’s API and from public HTTP. No dashboard was opened. No build log was read. No hosted SQL was run.

`main` is `135578461356cbf1d4940dd8e36a1300ed6cf849`. GitHub deployments for the environment `main - morrow-and-muse-api`, with environment URL `https://morrow-and-muse-api.onrender.com`, record a successful deploy of that commit at 20:40:43 UTC on 19 October 2026. The same commit also has a success at 15:56:53 UTC and a failure at 20:30:08 UTC. The Vercel commit status on that commit is success, “Deployment has completed”, updated at 20:12:02 UTC. The deployments list has no Vercel deployment object.

On 20 October 2026, `GET /health` and `GET /ready` returned `{"status":"ok"}` and `{"status":"ready"}`. Neither body includes a commit or a Go version. The API response identified Render behind Cloudflare. The storefront response identified Vercel. The HTML did not contain the commit. Its script URL was `/assets/index-BhBuGD7D.js`, and that file contained `morrow-and-muse-api.onrender.com`. The filename is build output, not a version pin. The Go toolchain on Render’s builder remains unverified.

## Free-tier limits

Citations are the pages fetched for this review. Numbers that those pages do not state are not repeated here.

**Render Free** ([Deploy for Free](https://render.com/docs/free), [Web Services](https://render.com/docs/web-services)). A free web service spins down after 15 minutes without inbound HTTP or WebSocket traffic. The next request spins it back up, which takes about a minute, and browsers see a loading page. Each workspace gets 750 free instance hours per calendar month. Spun-down services do not consume those hours. Using all of them suspends free web services until the next month. Outbound bandwidth and build pipeline minutes also have a monthly included amount; without a payment method, exhausting bandwidth suspends free services for the rest of the month, and exhausting pipeline minutes disables new builds. The free compute plan is 0.1 CPU and 512 MB RAM. The filesystem is ephemeral. Render may restart a free web service at any time. Free web services have no persistent disk, no SSH, and no scaling past one instance. Do not add a payment method. This demo uses Supabase, not Render Postgres (a free Render database expires 30 days after creation).

**Supabase Free** ([Project pausing](https://supabase.com/docs/guides/platform/free-project-pausing)). Free projects that do not get enough user database activity over a 7-day period are paused. Supabase says a few user requests to the database each day over the previous week is typically enough to avoid that. The owner gets a warning email about one week before the pause, and a confirmation after it. A paused project can be restored from the dashboard for up to 1 year. Paid plans are not paused. A paused database makes `/ready` fail and makes the storefront fall back to the bundled catalogue if the API cannot answer. The IPv4 add-on is paid; leave it off.

**Vercel Hobby** ([Hobby plan](https://vercel.com/docs/plans/hobby), [Terms of Service §4](https://vercel.com/legal/terms)). The Hobby plan is free and, in Vercel’s terms, only for personal or non-commercial use. The fair-use note on the Hobby page says the same. This storefront is a fictional boutique with no checkout. Do not use the deployment for payments, ads, or client work. Vercel can change Hobby limits or shut a Hobby project down. Exceeding included Hobby usage generally waits out the usage window rather than billing a card; do not add one.

## Still unverified

- Whether Render’s current native Go is already 1.27.x, or that `GOTOOLCHAIN=auto` successfully downloads 1.27.2 on that builder. No build log was read. A successful GitHub deployment of commit `1355784` does not show the toolchain. `/health` does not report it. Do not set `GO_VERSION`.
- Whether a Render free instance is IPv4-only. The session pooler is the choice that stays valid either way, without the paid IPv4 add-on.
- The Supabase session URI and the pooler host index. They were not copied into this review.
- The hosted privilege state after 19 October 2026. The project-run check from that date is the historical record above. This review did not re-query the database. `/ready` only shows that the owner ping succeeded. That check was not an independent security audit.
- The Vercel and Render dashboard screens (framework preset, Node version, plan, and Render’s health-check setting). `/health` answering is what the API does. It does not by itself show the dashboard health-check path. The storefront responses match the rewrite in `frontend/vercel.json`. The settings tables are the configuration this guide records.
