package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Desire-in-tech/morrow-and-muse/backend/internal/catalog"
)

func TestHealthDoesNotTouchTheDatabase(t *testing.T) {
	ready := &countingReady{}
	response := httptest.NewRecorder()
	NewRouter(Options{Ready: ready, Catalog: &fakeCatalog{}}).ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/health", nil),
	)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
	if response.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("content type = %q", response.Header().Get("Content-Type"))
	}
	var payload map[string]string
	decode(t, response.Body.Bytes(), &payload)
	if len(payload) != 1 || payload["status"] != "ok" {
		t.Fatalf("body = %#v", payload)
	}
	if ready.calls != 0 {
		t.Fatalf("ready calls = %d", ready.calls)
	}
}

func TestReady(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		response := serve(Options{Ready: &countingReady{}}, http.MethodGet, "/ready", "")
		if response.Code != http.StatusOK {
			t.Fatalf("status = %d body %s", response.Code, response.Body.Bytes())
		}
		var payload map[string]string
		decode(t, response.Body.Bytes(), &payload)
		if payload["status"] != "ready" {
			t.Fatalf("body = %#v", payload)
		}
	})

	t.Run("down", func(t *testing.T) {
		var logs bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&logs, nil))
		secret := "postgres://user:super-secret-password@db.internal/morrow"
		response := serve(Options{
			Ready:  &countingReady{err: errors.New(secret + " password authentication failed")},
			Logger: logger,
		}, http.MethodGet, "/ready", "")
		if response.Code != http.StatusServiceUnavailable {
			t.Fatalf("status = %d", response.Code)
		}
		assertError(t, response, "service_unavailable", "Not ready.")
		if strings.Contains(response.Body.String(), "super-secret-password") || strings.Contains(logs.String(), "super-secret-password") || strings.Contains(logs.String(), "postgres://") {
			t.Fatalf("leaked credentials\nbody %s\nlog %s", response.Body.String(), logs.String())
		}
	})

	health := serve(Options{Ready: &countingReady{err: errors.New("down")}}, http.MethodGet, "/health", "")
	if health.Code != http.StatusOK {
		t.Fatalf("health status = %d", health.Code)
	}
}

func TestCollectionsAndProducts(t *testing.T) {
	store := &fakeCatalog{
		collections: []catalog.Collection{
			{Slug: "workday-edit", Name: "Workday Edit", Description: "Desk", Image: catalog.Image{Src: "https://example.com/b", Alt: "desk"}, SEOTitle: "Workday", MetaDescription: "meta", SortOrder: 2},
			{Slug: "morning-ritual", Name: "Morning Ritual", Description: "Table", Image: catalog.Image{Src: "https://example.com/a", Alt: "table"}, SEOTitle: "Morning", MetaDescription: "meta", SortOrder: 1},
		},
		products: []catalog.Product{
			sampleProduct("second", 2),
			sampleProduct("everyday-ceramic-mug", 1),
		},
	}
	collections := serve(Options{Catalog: store}, http.MethodGet, "/api/v1/collections", "")
	if collections.Code != http.StatusOK {
		t.Fatalf("collections status = %d", collections.Code)
	}
	var listed struct {
		Data []catalog.Collection `json:"data"`
	}
	decode(t, collections.Body.Bytes(), &listed)
	if len(listed.Data) != 2 || listed.Data[0].Slug != "workday-edit" || listed.Data[1].Slug != "morning-ritual" {
		t.Fatalf("order = %#v", listed.Data)
	}

	products := serve(Options{Catalog: store}, http.MethodGet, "/api/v1/products", "")
	var productList struct {
		Data []catalog.Product `json:"data"`
	}
	decode(t, products.Body.Bytes(), &productList)
	if products.Code != http.StatusOK || len(productList.Data) != 2 {
		t.Fatalf("products status %d body %s", products.Code, products.Body.Bytes())
	}
	if productList.Data[0].Slug != "second" || productList.Data[1].Slug != "everyday-ceramic-mug" {
		t.Fatalf("order = %s then %s", productList.Data[0].Slug, productList.Data[1].Slug)
	}
	if len(productList.Data[1].Highlights) != 3 || productList.Data[1].PriceCents != 2800 {
		t.Fatalf("product = %#v", productList.Data[1])
	}
	if !strings.Contains(products.Body.String(), `"highlights":["one","two","three"]`) {
		t.Fatalf("highlights were not a JSON array: %s", products.Body.String())
	}

	one := serve(Options{Catalog: store}, http.MethodGet, "/api/v1/products/everyday-ceramic-mug", "")
	var single struct {
		Data catalog.Product `json:"data"`
	}
	decode(t, one.Body.Bytes(), &single)
	if one.Code != http.StatusOK || single.Data.Slug != "everyday-ceramic-mug" || single.Data.Collection.Slug != "morning-ritual" {
		t.Fatalf("status %d body %s", one.Code, one.Body.Bytes())
	}
}

func TestProductSlugCases(t *testing.T) {
	store := &fakeCatalog{products: []catalog.Product{sampleProduct("everyday-ceramic-mug", 1)}}

	missing := serve(Options{Catalog: store}, http.MethodGet, "/api/v1/products/missing-piece", "")
	if missing.Code != http.StatusNotFound {
		t.Fatalf("status = %d", missing.Code)
	}
	assertError(t, missing, "not_found", "Product not found.")

	before := store.calls
	malformed := serve(Options{Catalog: store}, http.MethodGet, "/api/v1/products/Not_A_Slug", "")
	if malformed.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body %s", malformed.Code, malformed.Body.Bytes())
	}
	assertError(t, malformed, "bad_request", "Slug is invalid.")
	if store.calls != before {
		t.Fatal("malformed slug queried the repository")
	}
}

func TestDatabaseFailureDoesNotLeak(t *testing.T) {
	var logs bytes.Buffer
	secret := "postgres://user:super-secret-password@db.internal/morrow select * from products"
	response := serve(Options{
		Catalog: &fakeCatalog{listErr: errors.New(secret)},
		Logger:  slog.New(slog.NewJSONHandler(&logs, nil)),
	}, http.MethodGet, "/api/v1/products", "")
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", response.Code)
	}
	assertError(t, response, "internal", "Something went wrong.")
	combined := response.Body.String() + logs.String()
	for _, leaked := range []string{"super-secret-password", "postgres://", "select *"} {
		if strings.Contains(combined, leaked) {
			t.Fatalf("leaked %q\n%s", leaked, combined)
		}
	}
}

func TestUnknownRouteAndMethod(t *testing.T) {
	unknown := serve(Options{}, http.MethodGet, "/api/v1/missing", "")
	if unknown.Code != http.StatusNotFound {
		t.Fatalf("status = %d", unknown.Code)
	}
	assertError(t, unknown, "not_found", "Not found.")

	method := serve(Options{}, http.MethodPost, "/health", "")
	if method.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d", method.Code)
	}
	assertError(t, method, "method_not_allowed", "Method not allowed.")
	if allow := method.Header().Get("Allow"); !strings.Contains(allow, "GET") {
		t.Fatalf("Allow = %q", allow)
	}
}

func TestCORS(t *testing.T) {
	opts := Options{Origins: []string{"http://localhost:5173"}, Ready: &countingReady{}}

	allowed := serve(opts, http.MethodGet, "/health", "http://localhost:5173")
	if allowed.Code != http.StatusOK {
		t.Fatalf("status = %d", allowed.Code)
	}
	if allowed.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("allow origin = %q", allowed.Header().Get("Access-Control-Allow-Origin"))
	}
	if allowed.Header().Get("Vary") != "Origin" {
		t.Fatalf("vary = %q", allowed.Header().Get("Vary"))
	}
	if allowed.Header().Get("Access-Control-Allow-Credentials") != "" {
		t.Fatal("credentials were allowed")
	}
	if allowed.Header().Get("Access-Control-Allow-Methods") != "GET, OPTIONS" {
		t.Fatalf("methods = %q", allowed.Header().Get("Access-Control-Allow-Methods"))
	}

	blocked := serve(opts, http.MethodGet, "/health", "https://shop.example")
	if blocked.Code != http.StatusOK {
		t.Fatalf("status = %d", blocked.Code)
	}
	if blocked.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("disallowed origin was reflected: %q", blocked.Header().Get("Access-Control-Allow-Origin"))
	}
	if blocked.Header().Get("Vary") != "Origin" {
		t.Fatalf("vary = %q", blocked.Header().Get("Vary"))
	}

	preflight := serve(opts, http.MethodOptions, "/api/v1/products", "http://localhost:5173")
	if preflight.Code != http.StatusNoContent {
		t.Fatalf("preflight status = %d body %s", preflight.Code, preflight.Body.Bytes())
	}
	if preflight.Body.Len() != 0 {
		t.Fatalf("preflight body = %q", preflight.Body.String())
	}
	if preflight.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatal("preflight missing allow origin")
	}
	if preflight.Header().Get("Access-Control-Allow-Headers") != "Accept, Content-Type" {
		t.Fatalf("headers = %q", preflight.Header().Get("Access-Control-Allow-Headers"))
	}
	if preflight.Header().Get("Access-Control-Allow-Credentials") != "" {
		t.Fatal("preflight allowed credentials")
	}

	denied := serve(opts, http.MethodOptions, "/api/v1/products", "https://shop.example")
	if denied.Code != http.StatusForbidden {
		t.Fatalf("status = %d", denied.Code)
	}
	assertError(t, denied, "forbidden", "Origin is not allowed.")
	if denied.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("disallowed preflight reflected an origin")
	}
}

func sampleProduct(slug string, sort int) catalog.Product {
	return catalog.Product{
		ID:          "mm-01",
		Slug:        slug,
		Name:        "Everyday Ceramic Mug",
		PriceCents:  2800,
		Collection:  catalog.CollectionRef{Slug: "morning-ritual", Name: "Morning Ritual"},
		Summary:     "A morning mug",
		Description: "A description",
		Highlights:  []string{"one", "two", "three"},
		CardImage:   catalog.Image{Src: "https://example.com/card", Alt: "card"},
		DetailImage: catalog.Image{Src: "https://example.com/detail", Alt: "detail"},
		SortOrder:   sort,
	}
}

func serve(opts Options, method string, path string, origin string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil)
	if origin != "" {
		request.Header.Set("Origin", origin)
	}
	response := httptest.NewRecorder()
	NewRouter(opts).ServeHTTP(response, request)
	return response
}

func decode(t *testing.T, body []byte, dest any) {
	t.Helper()
	if err := json.Unmarshal(body, dest); err != nil {
		t.Fatalf("unmarshal %s: %v", body, err)
	}
}

func assertError(t *testing.T, response *httptest.ResponseRecorder, code string, message string) {
	t.Helper()
	if response.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("content type = %q", response.Header().Get("Content-Type"))
	}
	var payload errorEnvelope
	decode(t, response.Body.Bytes(), &payload)
	if payload.Error.Code != code || payload.Error.Message != message {
		t.Fatalf("error = %#v", payload.Error)
	}
}

type countingReady struct {
	calls int
	err   error
}

func (c *countingReady) Ping(context.Context) error {
	c.calls++
	return c.err
}

type fakeCatalog struct {
	collections []catalog.Collection
	products    []catalog.Product
	listErr     error
	calls       int
}

func (f *fakeCatalog) ListCollections(context.Context) ([]catalog.Collection, error) {
	f.calls++
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.collections, nil
}

func (f *fakeCatalog) ListProducts(context.Context) ([]catalog.Product, error) {
	f.calls++
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.products, nil
}

func (f *fakeCatalog) ProductBySlug(_ context.Context, slug string) (catalog.Product, error) {
	f.calls++
	if f.listErr != nil {
		return catalog.Product{}, f.listErr
	}
	for _, product := range f.products {
		if product.Slug == slug {
			return product, nil
		}
	}
	return catalog.Product{}, catalog.ErrNotFound
}
