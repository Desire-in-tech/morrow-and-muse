package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/Desire-in-tech/morrow-and-muse/backend/db"
	"github.com/Desire-in-tech/morrow-and-muse/backend/internal/catalog"
	"github.com/Desire-in-tech/morrow-and-muse/backend/internal/dbpool"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestAPIMatchesCatalogExport(t *testing.T) {
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
	pool, err := dbpool.New(databaseURL)
	if err != nil {
		t.Fatal("open pool")
	}
	t.Cleanup(pool.Close)

	export := readExport(t)
	router := NewRouter(Options{
		Catalog: catalog.NewPostgres(pool),
		Ready:   pool,
		Origins: []string{"http://localhost:5173"},
	})

	health := request(router, http.MethodGet, "/health")
	if health.Code != http.StatusOK {
		t.Fatalf("health = %d", health.Code)
	}
	ready := request(router, http.MethodGet, "/ready")
	if ready.Code != http.StatusOK {
		t.Fatalf("ready = %d body %s", ready.Code, ready.Body.Bytes())
	}

	var collectionsBody struct {
		Data []catalog.Collection `json:"data"`
	}
	decode(t, request(router, http.MethodGet, "/api/v1/collections").Body.Bytes(), &collectionsBody)
	if !reflect.DeepEqual(collectionsBody.Data, export.collections()) {
		t.Fatalf("collections\n got %#v\nwant %#v", collectionsBody.Data, export.collections())
	}

	var productsBody struct {
		Data []catalog.Product `json:"data"`
	}
	productsResponse := request(router, http.MethodGet, "/api/v1/products")
	if productsResponse.Code != http.StatusOK {
		t.Fatalf("products = %d %s", productsResponse.Code, productsResponse.Body.Bytes())
	}
	decode(t, productsResponse.Body.Bytes(), &productsBody)
	wantProducts := export.products()
	if !reflect.DeepEqual(productsBody.Data, wantProducts) {
		t.Fatalf("products mismatch\n got %s\nwant %s", mustJSON(t, productsBody.Data), mustJSON(t, wantProducts))
	}

	slug := wantProducts[0].Slug
	var one struct {
		Data catalog.Product `json:"data"`
	}
	single := request(router, http.MethodGet, "/api/v1/products/"+slug)
	if single.Code != http.StatusOK {
		t.Fatalf("product = %d %s", single.Code, single.Body.Bytes())
	}
	decode(t, single.Body.Bytes(), &one)
	if !reflect.DeepEqual(one.Data, wantProducts[0]) {
		t.Fatalf("product %s\n got %#v\nwant %#v", slug, one.Data, wantProducts[0])
	}

	missing := request(router, http.MethodGet, "/api/v1/products/missing-piece")
	if missing.Code != http.StatusNotFound {
		t.Fatalf("missing = %d", missing.Code)
	}
	malformed := request(router, http.MethodGet, "/api/v1/products/Not_A_Slug")
	if malformed.Code != http.StatusBadRequest {
		t.Fatalf("malformed = %d", malformed.Code)
	}

	gotURLs := map[string]int{}
	add := func(url string) { gotURLs[url]++ }
	for _, collection := range collectionsBody.Data {
		add(collection.Image.Src)
	}
	for _, product := range productsBody.Data {
		add(product.CardImage.Src)
		add(product.DetailImage.Src)
	}
	assets := assetURLs(t, ctx, pool)
	add(assets["logo"])
	add(assets["hero"])

	wantURLs := export.urls()
	if len(wantURLs) != 36 || len(gotURLs) != 36 {
		t.Fatalf("unique urls api+assets=%d export=%d", len(gotURLs), len(wantURLs))
	}
	for url, count := range wantURLs {
		if gotURLs[url] != count {
			t.Fatalf("url count %s api=%d export=%d", url, gotURLs[url], count)
		}
	}
}

func request(handler http.Handler, method string, path string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(method, path, nil))
	return response
}

func assetURLs(t *testing.T, ctx context.Context, pool *pgxpool.Pool) map[string]string {
	t.Helper()
	rows, err := pool.Query(ctx, `select key, url from catalog_assets`)
	if err != nil {
		t.Fatalf("assets: %v", err)
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var key, url string
		if err := rows.Scan(&key, &url); err != nil {
			t.Fatalf("scan asset: %v", err)
		}
		out[key] = url
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("assets: %v", err)
	}
	return out
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(body)
}

type catalogExport struct {
	LogoURL     string       `json:"logoUrl"`
	HeroImage   catalogImage `json:"heroImage"`
	Collections []struct {
		Slug            string       `json:"slug"`
		Name            string       `json:"name"`
		Description     string       `json:"description"`
		Image           catalogImage `json:"image"`
		SEOTitle        string       `json:"seoTitle"`
		MetaDescription string       `json:"metaDescription"`
		SortOrder       int          `json:"sortOrder"`
	} `json:"collections"`
	Products []struct {
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
	} `json:"products"`
}

type catalogImage struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}

func readExport(t *testing.T) catalogExport {
	t.Helper()
	body, err := os.ReadFile("../../db/testdata/catalog_export.json")
	if err != nil {
		t.Fatalf("read export: %v", err)
	}
	var export catalogExport
	if err := json.Unmarshal(body, &export); err != nil {
		t.Fatalf("parse export: %v", err)
	}
	return export
}

func (e catalogExport) collections() []catalog.Collection {
	out := make([]catalog.Collection, len(e.Collections))
	for i, row := range e.Collections {
		out[i] = catalog.Collection{
			Slug:            row.Slug,
			Name:            row.Name,
			Description:     row.Description,
			Image:           catalog.Image{Src: row.Image.Src, Alt: row.Image.Alt},
			SEOTitle:        row.SEOTitle,
			MetaDescription: row.MetaDescription,
			SortOrder:       row.SortOrder,
		}
	}
	return out
}

func (e catalogExport) products() []catalog.Product {
	names := map[string]string{}
	for _, collection := range e.Collections {
		names[collection.Slug] = collection.Name
	}
	out := make([]catalog.Product, len(e.Products))
	for i, row := range e.Products {
		out[i] = catalog.Product{
			ID:          row.ID,
			Slug:        row.Slug,
			Name:        row.Name,
			PriceCents:  row.PriceCents,
			Collection:  catalog.CollectionRef{Slug: row.CollectionSlug, Name: names[row.CollectionSlug]},
			Summary:     row.Summary,
			Description: row.Description,
			Highlights:  append([]string(nil), row.Highlights...),
			CardImage:   catalog.Image{Src: row.CardImage.Src, Alt: row.CardImage.Alt},
			DetailImage: catalog.Image{Src: row.DetailImage.Src, Alt: row.DetailImage.Alt},
			SortOrder:   row.SortOrder,
		}
	}
	return out
}

func (e catalogExport) urls() map[string]int {
	out := map[string]int{}
	add := func(url string) { out[url]++ }
	add(e.LogoURL)
	add(e.HeroImage.Src)
	for _, collection := range e.Collections {
		add(collection.Image.Src)
	}
	for _, product := range e.Products {
		add(product.CardImage.Src)
		add(product.DetailImage.Src)
	}
	return out
}
