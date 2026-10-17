package httpapi

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/Desire-in-tech/morrow-and-muse/backend/internal/catalog"
	"github.com/go-chi/chi/v5"
)

// slugPattern matches the catalogue slug check in the database.
var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

const readyTimeout = time.Second

// Pinger reports whether the database accepts a connection.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Options configures the HTTP handlers. Logger may be nil.
type Options struct {
	Catalog catalog.Reader
	Ready   Pinger
	Origins []string
	Logger  *slog.Logger
}

type handler struct {
	catalog catalog.Reader
	ready   Pinger
	logger  *slog.Logger
}

// NewRouter serves the process health check and the read-only catalogue.
// A malformed product slug is 400. A well-formed unknown slug is 404.
func NewRouter(opts Options) http.Handler {
	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	h := handler{catalog: opts.Catalog, ready: opts.Ready, logger: logger}

	router := chi.NewRouter()
	router.Use(cors(opts.Origins))
	router.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "not_found", "Not found.")
	})
	router.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Allow", "GET, OPTIONS")
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed.")
	})
	router.Get("/health", h.health)
	router.Get("/ready", h.readyCheck)
	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/collections", h.collections)
		r.Get("/products", h.products)
		r.Get("/products/{slug}", h.product)
	})
	return router
}

func (h handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h handler) readyCheck(w http.ResponseWriter, r *http.Request) {
	if h.ready == nil {
		writeError(w, http.StatusServiceUnavailable, "service_unavailable", "Not ready.")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), readyTimeout)
	defer cancel()
	if err := h.ready.Ping(ctx); err != nil {
		h.logger.Error("ready check failed")
		writeError(w, http.StatusServiceUnavailable, "service_unavailable", "Not ready.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (h handler) collections(w http.ResponseWriter, r *http.Request) {
	if h.catalog == nil {
		h.fail(w, "list collections", errors.New("catalog is not configured"))
		return
	}
	rows, err := h.catalog.ListCollections(r.Context())
	if err != nil {
		h.fail(w, "list collections", err)
		return
	}
	if rows == nil {
		rows = []catalog.Collection{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": rows})
}

func (h handler) products(w http.ResponseWriter, r *http.Request) {
	if h.catalog == nil {
		h.fail(w, "list products", errors.New("catalog is not configured"))
		return
	}
	rows, err := h.catalog.ListProducts(r.Context())
	if err != nil {
		h.fail(w, "list products", err)
		return
	}
	if rows == nil {
		rows = []catalog.Product{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": rows})
}

func (h handler) product(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if !slugPattern.MatchString(slug) {
		writeError(w, http.StatusBadRequest, "bad_request", "Slug is invalid.")
		return
	}
	if h.catalog == nil {
		h.fail(w, "get product", errors.New("catalog is not configured"))
		return
	}
	row, err := h.catalog.ProductBySlug(r.Context(), slug)
	if errors.Is(err, catalog.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "Product not found.")
		return
	}
	if err != nil {
		h.fail(w, "get product", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": row})
}

func (h handler) fail(w http.ResponseWriter, op string, err error) {
	state := ""
	var failure *catalog.QueryFailure
	if errors.As(err, &failure) {
		state = failure.SQLState
	}
	h.logger.Error("query failed", "op", op, "sqlstate", state)
	writeError(w, http.StatusInternalServerError, "internal", "Something went wrong.")
}
