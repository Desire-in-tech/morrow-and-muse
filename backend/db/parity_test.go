package db_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	"github.com/Desire-in-tech/morrow-and-muse/backend/db"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestSeedMatchesCatalogExport(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		if os.Getenv("CI") == "true" {
			t.Fatal("TEST_DATABASE_URL is unset; database integration tests must run in CI")
		}
		t.Skip("TEST_DATABASE_URL is unset; skipping database integration tests")
	}

	ctx := context.Background()
	if err := db.Up(ctx, databaseURL); err != nil {
		t.Fatalf("migrate up: %v", err)
	}

	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	exportBytes, err := os.ReadFile("testdata/catalog_export.json")
	if err != nil {
		t.Fatalf("read export: %v", err)
	}
	var want catalogExport
	if err := json.Unmarshal(exportBytes, &want); err != nil {
		t.Fatalf("parse export: %v", err)
	}

	assertCatalog(t, ctx, sqlDB, want)
}

type catalogImage struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}

type catalogExport struct {
	LogoURL              string              `json:"logoUrl"`
	HeroImage            catalogImage        `json:"heroImage"`
	FeaturedProductSlugs []string            `json:"featuredProductSlugs"`
	Collections          []catalogCollection `json:"collections"`
	Products             []catalogProduct    `json:"products"`
}

type catalogCollection struct {
	Slug            string       `json:"slug"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	Image           catalogImage `json:"image"`
	SEOTitle        string       `json:"seoTitle"`
	MetaDescription string       `json:"metaDescription"`
	SortOrder       int          `json:"sortOrder"`
}

type catalogProduct struct {
	ID             string       `json:"id"`
	Slug           string       `json:"slug"`
	Name           string       `json:"name"`
	PriceCents     int          `json:"priceCents"`
	CollectionSlug string       `json:"collectionSlug"`
	Summary        string       `json:"summary"`
	Description    string       `json:"description"`
	Highlights     []string     `json:"highlights"`
	CardImage      catalogImage `json:"cardImage"`
	DetailImage    catalogImage `json:"detailImage"`
	SortOrder      int          `json:"sortOrder"`
}

func assertCatalog(t *testing.T, ctx context.Context, sqlDB *sql.DB, want catalogExport) {
	t.Helper()

	collections := queryCollections(t, ctx, sqlDB)
	if len(collections) != len(want.Collections) {
		t.Fatalf("collections = %d, export = %d", len(collections), len(want.Collections))
	}
	for i := range want.Collections {
		if collections[i] != want.Collections[i] {
			t.Fatalf("collection %d\n got %#v\nwant %#v", i, collections[i], want.Collections[i])
		}
	}

	products := queryProducts(t, ctx, sqlDB)
	if len(products) != len(want.Products) {
		t.Fatalf("products = %d, export = %d", len(products), len(want.Products))
	}
	for i := range want.Products {
		if !sameProduct(products[i], want.Products[i]) {
			t.Fatalf("product %s\n got %#v\nwant %#v", want.Products[i].Slug, products[i], want.Products[i])
		}
	}

	featured := queryFeatured(t, ctx, sqlDB)
	if len(featured) != len(want.FeaturedProductSlugs) {
		t.Fatalf("featured = %v, export = %v", featured, want.FeaturedProductSlugs)
	}
	for i := range want.FeaturedProductSlugs {
		if featured[i] != want.FeaturedProductSlugs[i] {
			t.Fatalf("featured = %v, export = %v", featured, want.FeaturedProductSlugs)
		}
	}

	logo, hero := queryAssets(t, ctx, sqlDB)
	if logo != want.LogoURL {
		t.Fatalf("logo = %q, export = %q", logo, want.LogoURL)
	}
	if hero != want.HeroImage {
		t.Fatalf("hero = %#v, export = %#v", hero, want.HeroImage)
	}

	gotURLs := map[string]int{}
	add := func(url string) { gotURLs[url]++ }
	add(logo)
	add(hero.Src)
	for _, collection := range collections {
		add(collection.Image.Src)
	}
	for _, product := range products {
		add(product.CardImage.Src)
		add(product.DetailImage.Src)
	}
	wantURLs := map[string]int{}
	addWant := func(url string) { wantURLs[url]++ }
	addWant(want.LogoURL)
	addWant(want.HeroImage.Src)
	for _, collection := range want.Collections {
		addWant(collection.Image.Src)
	}
	for _, product := range want.Products {
		addWant(product.CardImage.Src)
		addWant(product.DetailImage.Src)
	}
	if len(gotURLs) != 36 || len(wantURLs) != 36 {
		t.Fatalf("unique urls db=%d export=%d, want 36", len(gotURLs), len(wantURLs))
	}
	for url, count := range wantURLs {
		if gotURLs[url] != count {
			t.Fatalf("url %s count db=%d export=%d", url, gotURLs[url], count)
		}
	}

	var missingTimestamps int
	if err := sqlDB.QueryRowContext(ctx, `
		select count(*) from (
			select created_at, updated_at from collections
			union all
			select created_at, updated_at from products
		) rows
		where created_at is null or updated_at is null
	`).Scan(&missingTimestamps); err != nil {
		t.Fatalf("timestamps: %v", err)
	}
	if missingTimestamps != 0 {
		t.Fatalf("rows missing timestamps: %d", missingTimestamps)
	}
}

func sameProduct(got, want catalogProduct) bool {
	if got.ID != want.ID || got.Slug != want.Slug || got.Name != want.Name ||
		got.PriceCents != want.PriceCents || got.CollectionSlug != want.CollectionSlug ||
		got.Summary != want.Summary || got.Description != want.Description ||
		got.CardImage != want.CardImage || got.DetailImage != want.DetailImage ||
		got.SortOrder != want.SortOrder || len(got.Highlights) != len(want.Highlights) {
		return false
	}
	for i := range want.Highlights {
		if got.Highlights[i] != want.Highlights[i] {
			return false
		}
	}
	return true
}

func queryCollections(t *testing.T, ctx context.Context, sqlDB *sql.DB) []catalogCollection {
	t.Helper()
	rows, err := sqlDB.QueryContext(ctx, `
		select slug, name, description, image_url, image_alt, seo_title, meta_description, sort_order
		from collections
		order by sort_order
	`)
	if err != nil {
		t.Fatalf("collections: %v", err)
	}
	defer rows.Close()
	var out []catalogCollection
	for rows.Next() {
		var row catalogCollection
		if err := rows.Scan(
			&row.Slug, &row.Name, &row.Description, &row.Image.Src, &row.Image.Alt,
			&row.SEOTitle, &row.MetaDescription, &row.SortOrder,
		); err != nil {
			t.Fatalf("scan collection: %v", err)
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("collections: %v", err)
	}
	return out
}

func queryProducts(t *testing.T, ctx context.Context, sqlDB *sql.DB) []catalogProduct {
	t.Helper()
	rows, err := sqlDB.QueryContext(ctx, `
		select p.id, p.slug, p.name, p.price_cents, c.slug, p.summary, p.description,
			p.card_image_url, p.card_image_alt, p.detail_image_url, p.detail_image_alt, p.sort_order
		from products p
		join collections c on c.id = p.collection_id
		order by p.sort_order
	`)
	if err != nil {
		t.Fatalf("products: %v", err)
	}
	defer rows.Close()
	var out []catalogProduct
	for rows.Next() {
		var row catalogProduct
		if err := rows.Scan(
			&row.ID, &row.Slug, &row.Name, &row.PriceCents, &row.CollectionSlug,
			&row.Summary, &row.Description, &row.CardImage.Src, &row.CardImage.Alt,
			&row.DetailImage.Src, &row.DetailImage.Alt, &row.SortOrder,
		); err != nil {
			t.Fatalf("scan product: %v", err)
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("products: %v", err)
	}

	highlights, err := sqlDB.QueryContext(ctx, `
		select product_id, position, body
		from product_highlights
		order by product_id, position
	`)
	if err != nil {
		t.Fatalf("highlights: %v", err)
	}
	defer highlights.Close()
	byProduct := map[string][]string{}
	for highlights.Next() {
		var productID, body string
		var position int
		if err := highlights.Scan(&productID, &position, &body); err != nil {
			t.Fatalf("scan highlight: %v", err)
		}
		byProduct[productID] = append(byProduct[productID], body)
	}
	if err := highlights.Err(); err != nil {
		t.Fatalf("highlights: %v", err)
	}
	for i := range out {
		out[i].Highlights = byProduct[out[i].ID]
	}
	return out
}

func queryFeatured(t *testing.T, ctx context.Context, sqlDB *sql.DB) []string {
	t.Helper()
	rows, err := sqlDB.QueryContext(ctx, `
		select p.slug
		from featured_products f
		join products p on p.id = f.product_id
		order by f.position
	`)
	if err != nil {
		t.Fatalf("featured: %v", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			t.Fatalf("scan featured: %v", err)
		}
		out = append(out, slug)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("featured: %v", err)
	}
	return out
}

func queryAssets(t *testing.T, ctx context.Context, sqlDB *sql.DB) (string, catalogImage) {
	t.Helper()
	rows, err := sqlDB.QueryContext(ctx, `select key, url, alt from catalog_assets order by key`)
	if err != nil {
		t.Fatalf("assets: %v", err)
	}
	defer rows.Close()
	var logo string
	var hero catalogImage
	count := 0
	for rows.Next() {
		var key, url string
		var alt sql.NullString
		if err := rows.Scan(&key, &url, &alt); err != nil {
			t.Fatalf("scan asset: %v", err)
		}
		count++
		switch key {
		case "logo":
			if alt.Valid {
				t.Fatal("logo alt is not in catalog.ts and must stay null")
			}
			logo = url
		case "hero":
			if !alt.Valid {
				t.Fatal("hero alt is null")
			}
			hero = catalogImage{Src: url, Alt: alt.String}
		default:
			t.Fatalf("unexpected asset %s", key)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("assets: %v", err)
	}
	if count != 2 {
		t.Fatalf("assets = %d, want 2", count)
	}
	return logo, hero
}
