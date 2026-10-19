package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// TestSecurityHardening applies 00003 on a scratch database whose default
// privileges mimic Supabase, then checks the lock-down and a down/up cycle.
func TestSecurityHardening(t *testing.T) {
	baseURL := os.Getenv("TEST_DATABASE_URL")
	if baseURL == "" {
		if os.Getenv("CI") == "true" {
			t.Fatal("TEST_DATABASE_URL is unset; database integration tests must run in CI")
		}
		t.Skip("TEST_DATABASE_URL is unset; skipping database integration tests")
	}

	ctx := context.Background()
	admin, err := sql.Open("pgx", baseURL)
	if err != nil {
		t.Fatalf("open admin database: %v", err)
	}
	t.Cleanup(func() { _ = admin.Close() })

	createdRoles := createStandInRoles(t, ctx, admin)
	t.Cleanup(func() { dropStandInRoles(t, ctx, admin, createdRoles) })

	dbName := scratchDatabaseName(t)
	if _, err := admin.ExecContext(ctx, "create database "+dbName); err != nil {
		t.Fatalf("create scratch database: %v", err)
	}
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), "drop database "+dbName+" with (force)")
	})

	scratchURL := databaseURLWithName(t, baseURL, dbName)
	scratch, err := sql.Open("pgx", scratchURL)
	if err != nil {
		t.Fatalf("open scratch database: %v", err)
	}
	t.Cleanup(func() { _ = scratch.Close() })

	grantSupabaseLikeDefaults(t, ctx, scratch)

	provider, err := newProvider(scratchURL)
	if err != nil {
		t.Fatalf("goose provider: %v", err)
	}
	t.Cleanup(func() { _ = provider.Close() })

	if _, err := provider.UpTo(ctx, 2); err != nil {
		t.Fatalf("migrate up to 2: %v", err)
	}
	before := privilegeSnapshot(t, ctx, scratch)

	if _, err := provider.Up(ctx); err != nil {
		t.Fatalf("migrate up: %v", err)
	}

	t.Run("locked_down", func(t *testing.T) {
		assertLockedDown(t, ctx, scratch)
	})
	t.Run("anon_denied", func(t *testing.T) {
		assertAnonDenied(t, ctx, scratch)
	})
	t.Run("triggers", func(t *testing.T) {
		assertOwnerTriggers(t, ctx, scratch)
	})
	t.Run("down_restores_then_up", func(t *testing.T) {
		if _, err := provider.DownTo(ctx, 2); err != nil {
			t.Fatalf("migrate down to 2: %v", err)
		}
		afterDown := privilegeSnapshot(t, ctx, scratch)
		if before != afterDown {
			t.Fatalf("privilege snapshot after down\n--- before up ---\n%s\n--- after down ---\n%s", before, afterDown)
		}
		if _, err := provider.Up(ctx); err != nil {
			t.Fatalf("migrate up again: %v", err)
		}
		assertLockedDown(t, ctx, scratch)
		assertAnonDenied(t, ctx, scratch)
	})
}

func createStandInRoles(t *testing.T, ctx context.Context, admin *sql.DB) []string {
	t.Helper()
	names := []string{"anon", "authenticated", "service_role"}
	var created []string
	for _, name := range names {
		var exists bool
		if err := admin.QueryRowContext(ctx, "select exists (select 1 from pg_roles where rolname = $1)", name).Scan(&exists); err != nil {
			t.Fatalf("lookup role %s: %v", name, err)
		}
		if exists {
			t.Fatalf("role %s already exists; this test creates and drops its own stand-in roles", name)
		}
		if _, err := admin.ExecContext(ctx, "create role "+name+" nologin nobypassrls"); err != nil {
			t.Fatalf("create role %s: %v", name, err)
		}
		created = append(created, name)
		if _, err := admin.ExecContext(ctx, "grant "+name+" to current_user"); err != nil {
			t.Fatalf("grant %s to current user: %v", name, err)
		}
	}
	return created
}

func dropStandInRoles(t *testing.T, ctx context.Context, admin *sql.DB, names []string) {
	t.Helper()
	for _, name := range names {
		if _, err := admin.ExecContext(ctx, "revoke "+name+" from current_user"); err != nil {
			t.Errorf("revoke %s: %v", name, err)
		}
		if _, err := admin.ExecContext(ctx, "drop role "+name); err != nil {
			t.Errorf("drop role %s: %v", name, err)
		}
	}
}

func scratchDatabaseName(t *testing.T) string {
	t.Helper()
	var buf [6]byte
	if _, err := rand.Read(buf[:]); err != nil {
		t.Fatalf("random database name: %v", err)
	}
	return "mmsec" + hex.EncodeToString(buf[:])
}

func databaseURLWithName(t *testing.T, raw, name string) string {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse database URL: %v", err)
	}
	parsed.Path = "/" + name
	return parsed.String()
}

func grantSupabaseLikeDefaults(t *testing.T, ctx context.Context, scratch *sql.DB) {
	t.Helper()
	statements := []string{
		"alter default privileges in schema public grant all on tables to anon, authenticated, service_role",
		"alter default privileges in schema public grant all on sequences to anon, authenticated, service_role",
		"alter default privileges in schema public grant execute on functions to anon, authenticated, service_role",
		"grant usage on schema public to anon, authenticated, service_role",
	}
	for _, statement := range statements {
		if _, err := scratch.ExecContext(ctx, statement); err != nil {
			t.Fatalf("%s: %v", statement, err)
		}
	}
}

func assertLockedDown(t *testing.T, ctx context.Context, scratch *sql.DB) {
	t.Helper()

	var super, bypass bool
	if err := scratch.QueryRowContext(ctx, "select rolsuper, rolbypassrls from pg_roles where rolname = current_user").Scan(&super, &bypass); err != nil {
		t.Fatalf("role flags: %v", err)
	}
	t.Logf("migration role super=%t bypassrls=%t", super, bypass)

	tables := catalogueTables()
	var flags []string
	rows, err := scratch.QueryContext(ctx, `
		select c.relname, c.relrowsecurity, c.relforcerowsecurity
		from pg_class c
		join pg_namespace n on n.oid = c.relnamespace
		where n.nspname = 'public' and c.relkind = 'r'
		order by c.relname`)
	if err != nil {
		t.Fatalf("row security flags: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var enabled, forced bool
		if err := rows.Scan(&name, &enabled, &forced); err != nil {
			t.Fatalf("scan row security: %v", err)
		}
		flags = append(flags, fmt.Sprintf("%s enabled=%t forced=%t", name, enabled, forced))
		if !enabled || forced {
			t.Errorf("%s relrowsecurity=%t relforcerowsecurity=%t, want enabled and not forced", name, enabled, forced)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("row security rows: %v", err)
	}
	if len(flags) != len(tables) {
		t.Fatalf("public tables = %v, want %v", flags, tables)
	}

	var policies int
	if err := scratch.QueryRowContext(ctx, `
		select count(*)
		from pg_policy p
		join pg_class c on c.oid = p.polrelid
		join pg_namespace n on n.oid = c.relnamespace
		where n.nspname = 'public'`).Scan(&policies); err != nil {
		t.Fatalf("policies: %v", err)
	}
	if policies != 0 {
		t.Fatalf("public policies = %d, want 0", policies)
	}

	for _, role := range []string{"anon", "authenticated"} {
		for _, table := range tables {
			for _, priv := range tablePrivileges() {
				if hasTablePriv(t, ctx, scratch, role, table, priv) {
					t.Errorf("%s has %s on %s", role, priv, table)
				}
			}
		}
	}
	for _, table := range tables {
		for _, priv := range tablePrivileges() {
			if !hasTablePriv(t, ctx, scratch, "service_role", table, priv) {
				t.Errorf("service_role lost %s on %s", priv, table)
			}
		}
	}

	seqRows, err := scratch.QueryContext(ctx, `
		select c.relname
		from pg_class c
		join pg_namespace n on n.oid = c.relnamespace
		where n.nspname = 'public' and c.relkind = 'S'
		order by c.relname`)
	if err != nil {
		t.Fatalf("sequences: %v", err)
	}
	defer seqRows.Close()
	var sequences []string
	for seqRows.Next() {
		var name string
		if err := seqRows.Scan(&name); err != nil {
			t.Fatalf("scan sequence: %v", err)
		}
		sequences = append(sequences, name)
	}
	if err := seqRows.Err(); err != nil {
		t.Fatalf("sequence rows: %v", err)
	}
	if !containsAll(sequences, []string{"collections_id_seq", "goose_db_version_id_seq"}) {
		t.Fatalf("sequences = %v, want collections_id_seq and goose_db_version_id_seq", sequences)
	}
	for _, role := range []string{"anon", "authenticated"} {
		for _, seq := range sequences {
			for _, priv := range []string{"USAGE", "SELECT", "UPDATE"} {
				if hasSequencePriv(t, ctx, scratch, role, seq, priv) {
					t.Errorf("%s has %s on sequence %s", role, priv, seq)
				}
			}
		}
	}

	for _, role := range []string{"public", "anon", "authenticated"} {
		for _, fn := range triggerFunctions() {
			if hasFunctionPriv(t, ctx, scratch, role, fn) {
				t.Errorf("%s has EXECUTE on %s", role, fn)
			}
		}
	}
	for _, fn := range triggerFunctions() {
		if !hasFunctionPriv(t, ctx, scratch, "service_role", fn) {
			t.Errorf("service_role lost EXECUTE on %s", fn)
		}
	}

	var searchPath int
	if err := scratch.QueryRowContext(ctx, `
		select count(*)
		from pg_proc p
		join pg_namespace n on n.oid = p.pronamespace
		where n.nspname = 'public'
		  and p.proname in ('set_updated_at', 'products_have_three_highlights', 'product_highlights_stay_three')
		  and p.proconfig @> array['search_path=""']`).Scan(&searchPath); err != nil {
		t.Fatalf("proconfig: %v", err)
	}
	if searchPath != 3 {
		t.Fatalf("functions with search_path=\"\" = %d, want 3", searchPath)
	}

	assertNoDefaultGrantee(t, ctx, scratch, "anon")
	assertNoDefaultGrantee(t, ctx, scratch, "authenticated")
	for _, objType := range []string{"r", "S", "f"} {
		var n int
		if err := scratch.QueryRowContext(ctx, `
			select count(*)
			from pg_default_acl d
			join lateral aclexplode(d.defaclacl) a on true
			join pg_roles g on g.oid = a.grantee
			where d.defaclnamespace = 'public'::regnamespace
			  and d.defaclrole = (select oid from pg_roles where rolname = current_user)
			  and d.defaclobjtype = $1
			  and g.rolname = 'service_role'`, objType).Scan(&n); err != nil {
			t.Fatalf("service_role default privileges %s: %v", objType, err)
		}
		if n == 0 {
			t.Errorf("service_role default privileges for objtype %s were removed", objType)
		}
	}

	assertFutureFunctionDefaults(t, ctx, scratch)

	var collections, products, highlights, featured, assets int
	if err := scratch.QueryRowContext(ctx, `
		select
			(select count(*) from public.collections),
			(select count(*) from public.products),
			(select count(*) from public.product_highlights),
			(select count(*) from public.featured_products),
			(select count(*) from public.catalog_assets)`).Scan(&collections, &products, &highlights, &featured, &assets); err != nil {
		t.Fatalf("owner counts: %v", err)
	}
	if collections != 4 || products != 15 || highlights != 45 || featured != 4 || assets != 2 {
		t.Fatalf("owner counts = %d collections, %d products, %d highlights, %d featured, %d assets; want 4/15/45/4/2", collections, products, highlights, featured, assets)
	}
}

func assertAnonDenied(t *testing.T, ctx context.Context, scratch *sql.DB) {
	t.Helper()
	for _, table := range catalogueTables() {
		expectPermissionDenied(t, ctx, scratch, "select 1 from public."+table)
		expectPermissionDenied(t, ctx, scratch, "insert into public."+table+" default values")
	}
}

func assertOwnerTriggers(t *testing.T, ctx context.Context, scratch *sql.DB) {
	t.Helper()

	var updated string
	if err := scratch.QueryRowContext(ctx, `
		update public.collections
		set updated_at = timestamptz '1999-01-01+00'
		where id = (select min(id) from public.collections)
		returning updated_at::text`).Scan(&updated); err != nil {
		t.Fatalf("update collection timestamp: %v", err)
	}
	if strings.HasPrefix(updated, "1999-01-01") {
		t.Fatalf("updated_at = %s; set_updated_at did not replace it", updated)
	}

	tx, err := scratch.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin short highlight insert: %v", err)
	}
	if _, err := tx.ExecContext(ctx, `
		insert into public.products (
			id, slug, name, price_cents, collection_id, summary, description,
			card_image_url, card_image_alt, detail_image_url, detail_image_alt, sort_order
		) values (
			'mm-99', 'security-probe', 'Security probe', 100,
			(select min(id) from public.collections),
			'Probe summary', 'Probe description',
			'https://example.com/card.jpg', 'Card',
			'https://example.com/detail.jpg', 'Detail',
			99
		)`); err != nil {
		_ = tx.Rollback()
		t.Fatalf("insert probe product: %v", err)
	}
	err = tx.Commit()
	if err == nil || !strings.Contains(err.Error(), "exactly 3 highlights") {
		t.Fatalf("product without highlights: %v, want exactly 3 highlights", err)
	}

	tx, err = scratch.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin highlight delete: %v", err)
	}
	if _, err := tx.ExecContext(ctx, `
		delete from public.product_highlights
		where product_id = (select min(id) from public.products)
		  and position = 1`); err != nil {
		_ = tx.Rollback()
		t.Fatalf("delete highlight: %v", err)
	}
	err = tx.Commit()
	if err == nil || !strings.Contains(err.Error(), "exactly 3 highlights") {
		t.Fatalf("delete one highlight: %v, want exactly 3 highlights", err)
	}

	tx, err = scratch.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin valid product: %v", err)
	}
	if _, err := tx.ExecContext(ctx, `
		insert into public.products (
			id, slug, name, price_cents, collection_id, summary, description,
			card_image_url, card_image_alt, detail_image_url, detail_image_alt, sort_order
		) values (
			'mm-99', 'security-probe', 'Security probe', 100,
			(select min(id) from public.collections),
			'Probe summary', 'Probe description',
			'https://example.com/card.jpg', 'Card',
			'https://example.com/detail.jpg', 'Detail',
			99
		)`); err != nil {
		_ = tx.Rollback()
		t.Fatalf("insert valid probe product: %v", err)
	}
	if _, err := tx.ExecContext(ctx, `
		insert into public.product_highlights (product_id, position, body) values
			('mm-99', 1, 'one'),
			('mm-99', 2, 'two'),
			('mm-99', 3, 'three')`); err != nil {
		_ = tx.Rollback()
		t.Fatalf("insert three highlights: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit product with three highlights: %v", err)
	}
	if _, err := scratch.ExecContext(ctx, "delete from public.products where id = 'mm-99'"); err != nil {
		t.Fatalf("delete probe product: %v", err)
	}
}

func assertFutureFunctionDefaults(t *testing.T, ctx context.Context, scratch *sql.DB) {
	t.Helper()
	tx, err := scratch.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin function probe: %v", err)
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, "create function public.mm_sec_probe() returns integer language sql as 'select 1'"); err != nil {
		t.Fatalf("create probe function: %v", err)
	}
	for _, role := range []string{"public", "anon", "authenticated"} {
		var allowed bool
		if err := tx.QueryRowContext(ctx, "select has_function_privilege($1, 'public.mm_sec_probe()', 'EXECUTE')", role).Scan(&allowed); err != nil {
			t.Fatalf("future execute %s: %v", role, err)
		}
		if allowed {
			t.Errorf("future function is executable by %s", role)
		}
	}
	var serviceAllowed bool
	if err := tx.QueryRowContext(ctx, "select has_function_privilege('service_role', 'public.mm_sec_probe()', 'EXECUTE')").Scan(&serviceAllowed); err != nil {
		t.Fatalf("future execute service_role: %v", err)
	}
	if !serviceAllowed {
		t.Error("future function lost EXECUTE for service_role")
	}
}

func assertNoDefaultGrantee(t *testing.T, ctx context.Context, scratch *sql.DB, role string) {
	t.Helper()
	var n int
	if err := scratch.QueryRowContext(ctx, `
		select count(*)
		from pg_default_acl d
		join lateral aclexplode(d.defaclacl) a on true
		left join pg_roles g on g.oid = a.grantee
		where d.defaclrole = (select oid from pg_roles where rolname = current_user)
		  and (d.defaclnamespace = 'public'::regnamespace or d.defaclnamespace = 0)
		  and coalesce(g.rolname, 'public') = $1`, role).Scan(&n); err != nil {
		t.Fatalf("default privileges for %s: %v", role, err)
	}
	if n != 0 {
		t.Errorf("default privileges still grant %s (%d entries)", role, n)
	}
}

func expectPermissionDenied(t *testing.T, ctx context.Context, scratch *sql.DB, statement string) {
	t.Helper()
	tx, err := scratch.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, "set local role anon"); err != nil {
		t.Fatalf("set local role anon: %v", err)
	}
	_, err = tx.ExecContext(ctx, statement)
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "42501" {
		t.Fatalf("%s: got %v, want SQLSTATE 42501", statement, err)
	}
}

func privilegeSnapshot(t *testing.T, ctx context.Context, scratch *sql.DB) string {
	t.Helper()
	var lines []string
	appendQuery := func(query string, args ...any) {
		t.Helper()
		rows, err := scratch.QueryContext(ctx, query, args...)
		if err != nil {
			t.Fatalf("snapshot: %v", err)
		}
		defer rows.Close()
		cols, err := rows.Columns()
		if err != nil {
			t.Fatalf("snapshot columns: %v", err)
		}
		for rows.Next() {
			raw := make([]sql.NullString, len(cols))
			dest := make([]any, len(cols))
			for i := range raw {
				dest[i] = &raw[i]
			}
			if err := rows.Scan(dest...); err != nil {
				t.Fatalf("snapshot scan: %v", err)
			}
			parts := make([]string, len(cols))
			for i := range raw {
				if raw[i].Valid {
					parts[i] = strings.ReplaceAll(raw[i].String, "\n", `\n`)
				}
			}
			lines = append(lines, strings.Join(parts, "|"))
		}
		if err := rows.Err(); err != nil {
			t.Fatalf("snapshot rows: %v", err)
		}
	}

	appendQuery(`
		select 'table', c.relname, c.relrowsecurity::text, c.relforcerowsecurity::text,
		       coalesce(g.rolname, 'public'), a.privilege_type,
		       a.is_grantable::text
		from pg_class c
		join pg_namespace n on n.oid = c.relnamespace
		left join lateral aclexplode(c.relacl) a on true
		left join pg_roles g on g.oid = a.grantee
		where n.nspname = 'public' and c.relkind = 'r'`)
	appendQuery(`
		select 'sequence', c.relname, coalesce(g.rolname, 'public'), a.privilege_type, a.is_grantable::text
		from pg_class c
		join pg_namespace n on n.oid = c.relnamespace
		left join lateral aclexplode(c.relacl) a on true
		left join pg_roles g on g.oid = a.grantee
		where n.nspname = 'public' and c.relkind = 'S'`)
	appendQuery(`
		select 'function', p.proname, coalesce(array_to_string(p.proconfig, ','), ''), p.prosrc
		from pg_proc p
		join pg_namespace n on n.oid = p.pronamespace
		where n.nspname = 'public'`)
	appendQuery(`
		select 'function-acl', p.proname, coalesce(g.rolname, 'public'), a.privilege_type, a.is_grantable::text
		from pg_proc p
		join pg_namespace n on n.oid = p.pronamespace
		left join lateral aclexplode(p.proacl) a on true
		left join pg_roles g on g.oid = a.grantee
		where n.nspname = 'public'`)
	appendQuery(`
		select 'default', coalesce(n.nspname, 'global'), d.defaclobjtype,
		       coalesce(g.rolname, 'public'), a.privilege_type, a.is_grantable::text
		from pg_default_acl d
		left join pg_namespace n on n.oid = d.defaclnamespace
		left join lateral aclexplode(d.defaclacl) a on true
		left join pg_roles g on g.oid = a.grantee
		where d.defaclrole = (select oid from pg_roles where rolname = current_user)`)
	appendQuery(`
		select 'schema', coalesce(nspacl::text, '')
		from pg_namespace
		where nspname = 'public'`)
	appendQuery(`
		select 'policies', count(*)::text
		from pg_policy p
		join pg_class c on c.oid = p.polrelid
		join pg_namespace n on n.oid = c.relnamespace
		where n.nspname = 'public'`)

	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func hasTablePriv(t *testing.T, ctx context.Context, scratch *sql.DB, role, table, priv string) bool {
	t.Helper()
	var allowed bool
	query := "select has_table_privilege($1, $2, $3)"
	if err := scratch.QueryRowContext(ctx, query, role, "public."+table, priv).Scan(&allowed); err != nil {
		t.Fatalf("has_table_privilege %s %s %s: %v", role, table, priv, err)
	}
	return allowed
}

func hasSequencePriv(t *testing.T, ctx context.Context, scratch *sql.DB, role, seq, priv string) bool {
	t.Helper()
	var allowed bool
	if err := scratch.QueryRowContext(ctx, "select has_sequence_privilege($1, $2, $3)", role, "public."+seq, priv).Scan(&allowed); err != nil {
		t.Fatalf("has_sequence_privilege %s %s %s: %v", role, seq, priv, err)
	}
	return allowed
}

func hasFunctionPriv(t *testing.T, ctx context.Context, scratch *sql.DB, role, signature string) bool {
	t.Helper()
	var allowed bool
	if err := scratch.QueryRowContext(ctx, "select has_function_privilege($1, $2, 'EXECUTE')", role, "public."+signature).Scan(&allowed); err != nil {
		t.Fatalf("has_function_privilege %s %s: %v", role, signature, err)
	}
	return allowed
}

func catalogueTables() []string {
	return []string{
		"catalog_assets",
		"collections",
		"featured_products",
		"goose_db_version",
		"product_highlights",
		"products",
	}
}

func tablePrivileges() []string {
	return []string{"SELECT", "INSERT", "UPDATE", "DELETE", "TRUNCATE", "REFERENCES", "TRIGGER", "MAINTAIN"}
}

func triggerFunctions() []string {
	return []string{
		"set_updated_at()",
		"products_have_three_highlights()",
		"product_highlights_stay_three()",
	}
}

func containsAll(have, want []string) bool {
	set := map[string]bool{}
	for _, item := range have {
		set[item] = true
	}
	for _, item := range want {
		if !set[item] {
			return false
		}
	}
	return true
}
