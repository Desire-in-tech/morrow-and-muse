package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// QueryFailure is a database error safe to log. Error text never includes SQL or a DSN.
type QueryFailure struct {
	SQLState string
}

func (f *QueryFailure) Error() string {
	return "catalog query failed"
}

const queryTimeout = 3 * time.Second

// Postgres reads the catalogue through a pgx pool.
type Postgres struct {
	pool *pgxpool.Pool
}

// NewPostgres returns a reader. The caller owns the pool and closes it.
func NewPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{pool: pool}
}

func (s *Postgres) ListCollections(ctx context.Context) ([]Collection, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctx, `
		select slug, name, description, image_url, image_alt, seo_title, meta_description, sort_order
		from collections
		order by sort_order
	`)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	var out []Collection
	for rows.Next() {
		var row Collection
		if err := rows.Scan(
			&row.Slug, &row.Name, &row.Description, &row.Image.Src, &row.Image.Alt,
			&row.SEOTitle, &row.MetaDescription, &row.SortOrder,
		); err != nil {
			return nil, queryError(err)
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, queryError(err)
	}
	if out == nil {
		out = []Collection{}
	}
	return out, nil
}

func (s *Postgres) ListProducts(ctx context.Context) ([]Product, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctx, productQuery+` order by p.sort_order`)
	if err != nil {
		return nil, queryError(err)
	}
	defer rows.Close()

	var out []Product
	for rows.Next() {
		row, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, queryError(err)
	}
	if out == nil {
		out = []Product{}
	}
	return out, nil
}

func (s *Postgres) ProductBySlug(ctx context.Context, slug string) (Product, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctx, productQuery+` where p.slug = $1`, slug)
	if err != nil {
		return Product{}, queryError(err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return Product{}, queryError(err)
		}
		return Product{}, ErrNotFound
	}
	row, err := scanProduct(rows)
	if err != nil {
		return Product{}, err
	}
	if rows.Next() {
		return Product{}, errors.New("catalog query failed")
	}
	return row, nil
}

const productQuery = `
	select p.id, p.slug, p.name, p.price_cents,
		c.slug, c.name,
		p.summary, p.description,
		p.card_image_url, p.card_image_alt,
		p.detail_image_url, p.detail_image_alt,
		p.sort_order,
		(
			select json_agg(h.body order by h.position)
			from product_highlights h
			where h.product_id = p.id
		)
	from products p
	join collections c on c.id = p.collection_id
`

type scanner interface {
	Scan(dest ...any) error
}

func scanProduct(row scanner) (Product, error) {
	var (
		product    Product
		highlights []byte
	)
	if err := row.Scan(
		&product.ID, &product.Slug, &product.Name, &product.PriceCents,
		&product.Collection.Slug, &product.Collection.Name,
		&product.Summary, &product.Description,
		&product.CardImage.Src, &product.CardImage.Alt,
		&product.DetailImage.Src, &product.DetailImage.Alt,
		&product.SortOrder,
		&highlights,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Product{}, ErrNotFound
		}
		return Product{}, queryError(err)
	}
	if err := json.Unmarshal(highlights, &product.Highlights); err != nil || len(product.Highlights) != 3 {
		return Product{}, errors.New("catalog query failed")
	}
	return product, nil
}

func queryError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	failure := &QueryFailure{}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		failure.SQLState = pgErr.Code
	}
	return failure
}
