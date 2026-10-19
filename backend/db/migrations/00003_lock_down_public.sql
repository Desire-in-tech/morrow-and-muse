-- Lock the public schema.
--
-- The storefront reads only through the Go API. The API and goose both
-- connect as the table owner (postgres on Supabase, the migration role
-- in CI and local Postgres). This migration enables row level security
-- and adds no policies.
--
-- Row level security is not forced. Without FORCE ROW LEVEL SECURITY the
-- table owner bypasses it, so catalogue reads and goose_db_version writes
-- keep working. On Supabase the postgres role also has BYPASSRLS, which
-- bypasses row level security even when it is forced. Forcing it with no
-- policies would refuse a non-superuser owner that does not bypass row
-- level security, including a local migration role. ENABLE without FORCE
-- is the setting that keeps that owner working.
--
-- service_role is left unchanged. It is server-only, it bypasses row
-- level security, and Supabase already grants it full access. The path
-- this closes is anon, authenticated, and PUBLIC.
--
-- Default privileges are changed for the current role, which is the role
-- goose runs as. That is postgres on Supabase. CI and local Postgres use
-- a different role name, so the statements are not written FOR ROLE
-- postgres. supabase_admin default privileges cannot be changed by this
-- role; objects this project creates are owned by the migration role.
-- A later migration that creates a table still has to enable row level
-- security itself. service_role keeps its default grants.
--
-- PostgreSQL merges schema default privileges onto the built-in function
-- default, and that built-in default includes EXECUTE for PUBLIC. Revoking
-- PUBLIC inside the schema does not remove it when the stored entry has
-- no PUBLIC grant. The database-wide revoke below does. Down grants that
-- built-in default back, which drops the row.
--
-- On PostgreSQL 17, firing an existing trigger does not require the
-- invoking role to hold EXECUTE on the trigger function. EXECUTE is still
-- revoked from PUBLIC, anon, and authenticated. The owner keeps EXECUTE,
-- and the migration tests check that the owner's writes still fire the
-- triggers.
--
-- Down restores the previous grants and turns row level security off.
-- That exposes the catalogue to anon and authenticated again. Do not run
-- Down for this version on the hosted database.
--
-- +goose Up

alter table public.collections enable row level security;
alter table public.products enable row level security;
alter table public.product_highlights enable row level security;
alter table public.featured_products enable row level security;
alter table public.catalog_assets enable row level security;
alter table public.goose_db_version enable row level security;

revoke all on table
    public.collections,
    public.products,
    public.product_highlights,
    public.featured_products,
    public.catalog_assets,
    public.goose_db_version
from public;

revoke all on all sequences in schema public from public;

revoke all on function
    public.set_updated_at(),
    public.products_have_three_highlights(),
    public.product_highlights_stay_three()
from public;

-- +goose StatementBegin
do $lock$
declare
    role_name text;
begin
    foreach role_name in array array['anon', 'authenticated']
    loop
        if exists (select 1 from pg_roles where rolname = role_name) then
            execute format(
                'revoke all on table public.collections, public.products, public.product_highlights, public.featured_products, public.catalog_assets, public.goose_db_version from %I',
                role_name
            );
            execute format(
                'revoke all on all sequences in schema public from %I',
                role_name
            );
            execute format(
                'revoke all on function public.set_updated_at(), public.products_have_three_highlights(), public.product_highlights_stay_three() from %I',
                role_name
            );
            execute format(
                'alter default privileges in schema public revoke all on tables from %I',
                role_name
            );
            execute format(
                'alter default privileges in schema public revoke all on sequences from %I',
                role_name
            );
            execute format(
                'alter default privileges in schema public revoke all on functions from %I',
                role_name
            );
        end if;
    end loop;
end
$lock$;
-- +goose StatementEnd

alter default privileges in schema public revoke execute on functions from public;
alter default privileges revoke execute on functions from public;

-- +goose StatementBegin
create or replace function public.set_updated_at()
returns trigger
language plpgsql
set search_path = ''
as $fn$
begin
    new.updated_at = pg_catalog.clock_timestamp();
    return new;
end;
$fn$;
-- +goose StatementEnd

alter function public.set_updated_at() set search_path = '';

-- +goose StatementBegin
create or replace function public.products_have_three_highlights()
returns trigger
language plpgsql
set search_path = ''
as $fn$
declare
    highlight_count integer;
begin
    if not exists (select 1 from public.products where id = new.id) then
        return null;
    end if;
    select count(*) into highlight_count
    from public.product_highlights
    where product_id = new.id;
    if highlight_count <> 3 then
        raise exception 'product % must have exactly 3 highlights', new.id;
    end if;
    return null;
end;
$fn$;
-- +goose StatementEnd

alter function public.products_have_three_highlights() set search_path = '';

-- coalesce() is SQL grammar, not a pg_catalog function, so it stays unqualified.
-- +goose StatementBegin
create or replace function public.product_highlights_stay_three()
returns trigger
language plpgsql
set search_path = ''
as $fn$
declare
    pid text;
    highlight_count integer;
begin
    pid = coalesce(new.product_id, old.product_id);
    if not exists (select 1 from public.products where id = pid) then
        return null;
    end if;
    select count(*) into highlight_count
    from public.product_highlights
    where product_id = pid;
    if highlight_count <> 3 then
        raise exception 'product % must have exactly 3 highlights', pid;
    end if;
    return null;
end;
$fn$;
-- +goose StatementEnd

alter function public.product_highlights_stay_three() set search_path = '';

-- +goose Down

-- +goose StatementBegin
create or replace function public.set_updated_at()
returns trigger
language plpgsql
as $fn$
begin
    new.updated_at = clock_timestamp();
    return new;
end;
$fn$;
-- +goose StatementEnd

alter function public.set_updated_at() reset search_path;

-- +goose StatementBegin
create or replace function public.products_have_three_highlights()
returns trigger
language plpgsql
as $fn$
declare
    highlight_count integer;
begin
    if not exists (select 1 from products where id = new.id) then
        return null;
    end if;
    select count(*) into highlight_count
    from product_highlights
    where product_id = new.id;
    if highlight_count <> 3 then
        raise exception 'product % must have exactly 3 highlights', new.id;
    end if;
    return null;
end;
$fn$;
-- +goose StatementEnd

alter function public.products_have_three_highlights() reset search_path;

-- +goose StatementBegin
create or replace function public.product_highlights_stay_three()
returns trigger
language plpgsql
as $fn$
declare
    pid text;
    highlight_count integer;
begin
    pid = coalesce(new.product_id, old.product_id);
    if not exists (select 1 from products where id = pid) then
        return null;
    end if;
    select count(*) into highlight_count
    from product_highlights
    where product_id = pid;
    if highlight_count <> 3 then
        raise exception 'product % must have exactly 3 highlights', pid;
    end if;
    return null;
end;
$fn$;
-- +goose StatementEnd

alter function public.product_highlights_stay_three() reset search_path;

alter table public.collections disable row level security;
alter table public.products disable row level security;
alter table public.product_highlights disable row level security;
alter table public.featured_products disable row level security;
alter table public.catalog_assets disable row level security;
alter table public.goose_db_version disable row level security;

grant execute on function
    public.set_updated_at(),
    public.products_have_three_highlights(),
    public.product_highlights_stay_three()
to public;

-- +goose StatementBegin
do $lock$
declare
    role_name text;
begin
    foreach role_name in array array['anon', 'authenticated']
    loop
        if exists (select 1 from pg_roles where rolname = role_name) then
            execute format(
                'grant all on table public.collections, public.products, public.product_highlights, public.featured_products, public.catalog_assets, public.goose_db_version to %I',
                role_name
            );
            execute format(
                'grant all on all sequences in schema public to %I',
                role_name
            );
            execute format(
                'grant execute on function public.set_updated_at(), public.products_have_three_highlights(), public.product_highlights_stay_three() to %I',
                role_name
            );
            execute format(
                'alter default privileges in schema public grant all on tables to %I',
                role_name
            );
            execute format(
                'alter default privileges in schema public grant all on sequences to %I',
                role_name
            );
            execute format(
                'alter default privileges in schema public grant execute on functions to %I',
                role_name
            );
        end if;
    end loop;
end
$lock$;
-- +goose StatementEnd

-- Inverse of the database-wide revoke in Up. When the result matches
-- PostgreSQL's built-in function default, the pg_default_acl row is
-- removed. Granting PUBLIC inside schema public would add an explicit
-- entry that was not there before Up.
alter default privileges grant execute on functions to public;
