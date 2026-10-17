import {
  collections as bundledCollections,
  featuredProductSlugs,
  products as bundledProducts,
  type Collection,
  type Product,
} from '../data/catalog.ts'
import {
  fetchCollections,
  fetchProducts,
  readApiBaseUrl,
} from './catalogClient.ts'
import { mapCatalog } from './mapCatalog.ts'

export type CatalogSource = 'bundle' | 'api'

export type CatalogSnapshot = {
  source: CatalogSource
  collections: readonly Collection[]
  products: readonly Product[]
}

/**
 * Fallback rule: the catalogue is one unit.
 * Collections and products are requested together. The storefront switches
 * to the API only when both responses parse and map. If either request
 * fails, or the payload is malformed, both lists stay on the bundled
 * catalogue. Pages never mix an API collection list with bundled products.
 * Product routes resolve by slug from that same snapshot. The slug endpoint
 * exists on the client, but pages do not call it, so one detail failure
 * cannot disagree with the list.
 * Logo, hero, and featured slugs stay in the bundle. The API does not
 * serve them. Featured cards are the snapshot products in featured-slug order.
 */
const bundleSnapshot = (): CatalogSnapshot => ({
  source: 'bundle',
  collections: bundledCollections,
  products: bundledProducts,
})

let generation = 0
let snapshot = bundleSnapshot()
let loadPromise: Promise<void> | null = null
const listeners = new Set<() => void>()

export function getCatalogSnapshot(): CatalogSnapshot {
  return snapshot
}

export function subscribeCatalog(listener: () => void): () => void {
  listeners.add(listener)
  startCatalogLoad()
  return () => {
    listeners.delete(listener)
  }
}

export function resetCatalogSession(): void {
  generation += 1
  loadPromise = null
  snapshot = bundleSnapshot()
  emit()
}

export function findCollection(
  catalog: CatalogSnapshot,
  slug: string,
): Collection | undefined {
  return catalog.collections.find((collection) => collection.slug === slug)
}

export function findProduct(
  catalog: CatalogSnapshot,
  slug: string,
): Product | undefined {
  return catalog.products.find((product) => product.slug === slug)
}

export function productsInCollection(
  catalog: CatalogSnapshot,
  slug: string,
): readonly Product[] {
  return catalog.products.filter((product) => product.collectionSlug === slug)
}

export function featuredProducts(catalog: CatalogSnapshot): readonly Product[] {
  return featuredProductSlugs.flatMap((slug) => {
    const product = findProduct(catalog, slug)
    return product ? [product] : []
  })
}

export function relatedProducts(
  catalog: CatalogSnapshot,
  slug: string,
): readonly Product[] {
  const product = findProduct(catalog, slug)
  if (!product) {
    return []
  }
  return productsInCollection(catalog, product.collectionSlug).filter(
    (item) => item.slug !== slug,
  )
}

function startCatalogLoad(): void {
  if (loadPromise) {
    return
  }
  const baseUrl = readApiBaseUrl()
  if (!baseUrl) {
    return
  }
  const gen = generation
  loadPromise = loadCatalogUnit(baseUrl)
    .then((next) => {
      if (gen !== generation || !next) {
        return
      }
      snapshot = next
      emit()
    })
    .catch(() => {
      // Keep the bundle. The rejection is consumed here.
    })
}

async function loadCatalogUnit(
  baseUrl: string,
): Promise<CatalogSnapshot | null> {
  const [collectionsResult, productsResult] = await Promise.all([
    fetchCollections(baseUrl),
    fetchProducts(baseUrl),
  ])
  if (!collectionsResult.ok || !productsResult.ok) {
    return null
  }
  const mapped = mapCatalog(collectionsResult.data, productsResult.data)
  if (!mapped) {
    return null
  }
  return { source: 'api', ...mapped }
}

function emit(): void {
  for (const listener of listeners) {
    listener()
  }
}
