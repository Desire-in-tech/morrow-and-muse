/**
 * Reads frontend/src/data/catalog.ts and writes the goose seed migration
 * plus the JSON fixture the database parity test compares against.
 *
 * The catalogue file is the only source of the strings. This script does
 * not modify it.
 *
 *   node --experimental-strip-types scripts/export-catalog.ts
 *   node --experimental-strip-types scripts/export-catalog.ts --check
 */
import { readFileSync, writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import {
  collections,
  featuredProductSlugs,
  heroImage,
  logoUrl,
  products,
} from '../frontend/src/data/catalog.ts'

const quoteTag = '$mmseed$'
const root = dirname(dirname(fileURLToPath(import.meta.url)))
const seedPath = join(root, 'backend/db/migrations/00002_seed.sql')
const jsonPath = join(root, 'backend/db/testdata/catalog_export.json')

type Image = { src?: string; alt: string }

function requiredSrc(image: Image, label: string): string {
  if (image.src === undefined || image.src === '') {
    throw new Error(`catalog.ts image is missing src: ${label}`)
  }
  return image.src
}

function quote(value: string): string {
  if (value.includes(quoteTag)) {
    throw new Error('catalogue text contains the seed dollar-quote tag')
  }
  return `${quoteTag}${value}${quoteTag}`
}

const exportDocument = {
  logoUrl,
  heroImage: {
    src: requiredSrc(heroImage, 'hero'),
    alt: heroImage.alt,
  },
  featuredProductSlugs: [...featuredProductSlugs],
  collections: collections.map((collection, index) => ({
    slug: collection.slug,
    name: collection.name,
    description: collection.description,
    image: {
      src: requiredSrc(collection.image, collection.slug),
      alt: collection.image.alt,
    },
    seoTitle: collection.seoTitle,
    metaDescription: collection.metaDescription,
    sortOrder: index + 1,
  })),
  products: products.map((product, index) => ({
    id: product.id,
    slug: product.slug,
    name: product.name,
    priceCents: product.priceCents,
    collectionSlug: product.collectionSlug,
    summary: product.summary,
    description: product.description,
    highlights: [...product.highlights],
    cardImage: {
      src: requiredSrc(product.cardImage, `${product.slug} card`),
      alt: product.cardImage.alt,
    },
    detailImage: {
      src: requiredSrc(product.detailImage, `${product.slug} detail`),
      alt: product.detailImage.alt,
    },
    sortOrder: index + 1,
  })),
}

const jsonText = `${JSON.stringify(exportDocument, null, 2)}\n`

const collectionRows = exportDocument.collections
  .map(
    (collection) => `    (
        ${quote(collection.slug)},
        ${quote(collection.name)},
        ${quote(collection.description)},
        ${quote(collection.image.src)},
        ${quote(collection.image.alt)},
        ${quote(collection.seoTitle)},
        ${quote(collection.metaDescription)},
        ${collection.sortOrder}
    )`,
  )
  .join(',\n')

const productRows = exportDocument.products
  .map(
    (product) => `    (
        ${quote(product.id)},
        ${quote(product.slug)},
        ${quote(product.name)},
        ${product.priceCents},
        (select id from collections where slug = ${quote(product.collectionSlug)}),
        ${quote(product.summary)},
        ${quote(product.description)},
        ${quote(product.cardImage.src)},
        ${quote(product.cardImage.alt)},
        ${quote(product.detailImage.src)},
        ${quote(product.detailImage.alt)},
        ${product.sortOrder}
    )`,
  )
  .join(',\n')

const highlightRows = exportDocument.products
  .flatMap((product) =>
    product.highlights.map(
      (body, index) =>
        `    (${quote(product.id)}, ${index + 1}, ${quote(body)})`,
    ),
  )
  .join(',\n')

const featuredRows = exportDocument.featuredProductSlugs
  .map(
    (slug, index) =>
      `    ((select id from products where slug = ${quote(slug)}), ${index + 1})`,
  )
  .join(',\n')

const seedText = `-- Generated from frontend/src/data/catalog.ts by scripts/export-catalog.ts.
-- Do not edit by hand. Regenerate, then commit the result:
--   node --experimental-strip-types scripts/export-catalog.ts
-- +goose Up

insert into collections (
    slug,
    name,
    description,
    image_url,
    image_alt,
    seo_title,
    meta_description,
    sort_order
) values
${collectionRows};

insert into products (
    id,
    slug,
    name,
    price_cents,
    collection_id,
    summary,
    description,
    card_image_url,
    card_image_alt,
    detail_image_url,
    detail_image_alt,
    sort_order
) values
${productRows};

insert into product_highlights (product_id, position, body) values
${highlightRows};

insert into featured_products (product_id, position) values
${featuredRows};

insert into catalog_assets (key, url, alt) values
    ('logo', ${quote(exportDocument.logoUrl)}, null),
    (
        'hero',
        ${quote(exportDocument.heroImage.src)},
        ${quote(exportDocument.heroImage.alt)}
    );

-- +goose Down

delete from featured_products;
delete from product_highlights;
delete from products;
delete from collections;
delete from catalog_assets;
`

const check = process.argv.includes('--check')
if (check) {
  const mismatches = [
    [seedPath, seedText],
    [jsonPath, jsonText],
  ].filter(([path, text]) => {
    try {
      return readFileSync(path, 'utf8') !== text
    } catch {
      return true
    }
  })
  if (mismatches.length > 0) {
    const names = mismatches.map(([path]) => path).join('\n')
    console.error(
      `catalogue seed is stale. Regenerate with: node --experimental-strip-types scripts/export-catalog.ts\n${names}`,
    )
    process.exit(1)
  }
  console.log('catalogue seed matches catalog.ts')
} else {
  writeFileSync(seedPath, seedText)
  writeFileSync(jsonPath, jsonText)
  console.log(`wrote ${seedPath}`)
  console.log(`wrote ${jsonPath}`)
}
