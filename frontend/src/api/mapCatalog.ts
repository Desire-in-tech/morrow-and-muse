import type { Collection, Product } from '../data/catalog.ts'
import type { ApiCollection, ApiProduct } from './catalogClient.ts'

export type MappedCatalog = {
  collections: readonly Collection[]
  products: readonly Product[]
}

/**
 * Maps a successful API payload onto the storefront types.
 * Returns null when the two lists cannot be used together: duplicate slugs,
 * or a product whose collection slug is missing from the collection list.
 * Sort order is the API sortOrder, which is the display order.
 */
export function mapCatalog(
  apiCollections: readonly ApiCollection[],
  apiProducts: readonly ApiProduct[],
): MappedCatalog | null {
  const collections = [...apiCollections].sort(bySortOrder)
  const products = [...apiProducts].sort(bySortOrder)
  if (
    !uniqueSlugs(collections.map((collection) => collection.slug)) ||
    !uniqueSlugs(products.map((product) => product.slug))
  ) {
    return null
  }

  const collectionSlugs = new Set(
    collections.map((collection) => collection.slug),
  )
  if (
    products.some((product) => !collectionSlugs.has(product.collection.slug))
  ) {
    return null
  }

  return {
    collections: collections.map(mapCollection),
    products: products.map(mapProduct),
  }
}

function mapCollection(collection: ApiCollection): Collection {
  return {
    slug: collection.slug,
    name: collection.name,
    description: collection.description,
    image: { src: collection.image.src, alt: collection.image.alt },
    seoTitle: collection.seoTitle,
    metaDescription: collection.metaDescription,
  }
}

function mapProduct(product: ApiProduct): Product {
  const [first, second, third] = product.highlights
  return {
    id: product.id,
    slug: product.slug,
    name: product.name,
    priceCents: product.priceCents,
    collectionSlug: product.collection.slug,
    summary: product.summary,
    description: product.description,
    highlights: [first, second, third],
    cardImage: { src: product.cardImage.src, alt: product.cardImage.alt },
    detailImage: {
      src: product.detailImage.src,
      alt: product.detailImage.alt,
    },
  }
}

function bySortOrder(
  left: { sortOrder: number },
  right: { sortOrder: number },
) {
  return left.sortOrder - right.sortOrder
}

function uniqueSlugs(slugs: readonly string[]): boolean {
  return new Set(slugs).size === slugs.length
}
