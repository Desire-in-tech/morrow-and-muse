// Package catalog is the read-only storefront catalogue.
package catalog

import (
	"context"
	"errors"
)

// ErrNotFound is returned when a slug is well formed and absent.
var ErrNotFound = errors.New("catalog: not found")

// Image is a catalogue photograph.
type Image struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}

// Collection is one of the four storefront collections, in sort order.
type Collection struct {
	Slug            string `json:"slug"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Image           Image  `json:"image"`
	SEOTitle        string `json:"seoTitle"`
	MetaDescription string `json:"metaDescription"`
	SortOrder       int    `json:"sortOrder"`
}

// CollectionRef is the product's collection. Name is the display title.
type CollectionRef struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// Product is one storefront piece.
// Name is the display title, matching catalog.ts and products.name.
type Product struct {
	ID          string        `json:"id"`
	Slug        string        `json:"slug"`
	Name        string        `json:"name"`
	PriceCents  int           `json:"priceCents"`
	Collection  CollectionRef `json:"collection"`
	Summary     string        `json:"summary"`
	Description string        `json:"description"`
	Highlights  []string      `json:"highlights"`
	CardImage   Image         `json:"cardImage"`
	DetailImage Image         `json:"detailImage"`
	SortOrder   int           `json:"sortOrder"`
}

// Reader loads the catalogue. Implementations must return rows in sort order
// and exactly three highlights per product.
type Reader interface {
	ListCollections(ctx context.Context) ([]Collection, error)
	ListProducts(ctx context.Context) ([]Product, error)
	ProductBySlug(ctx context.Context, slug string) (Product, error)
}
