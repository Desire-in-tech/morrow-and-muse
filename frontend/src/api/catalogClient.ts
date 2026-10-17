/**
 * Typed client for the read-only catalogue API.
 *
 * Response shapes follow backend/internal/catalog and the httpapi envelopes:
 * success is {"data": ...}; errors are {"error":{"code","message"}}.
 * This module does not choose the bundled fallback. The session does.
 */

export type ApiImage = {
  src: string
  alt: string
}

export type ApiCollection = {
  slug: string
  name: string
  description: string
  image: ApiImage
  seoTitle: string
  metaDescription: string
  sortOrder: number
}

export type ApiCollectionRef = {
  slug: string
  name: string
}

export type ApiProduct = {
  id: string
  slug: string
  name: string
  priceCents: number
  collection: ApiCollectionRef
  summary: string
  description: string
  highlights: [string, string, string]
  cardImage: ApiImage
  detailImage: ApiImage
  sortOrder: number
}

export type CatalogFailure =
  | { kind: 'network' }
  | { kind: 'status'; status: number }
  | { kind: 'malformed' }
  | { kind: 'not_found' }

export type CatalogResult<T> =
  { ok: true; data: T } | { ok: false; failure: CatalogFailure }

declare global {
  interface ImportMetaEnv {
    readonly VITE_API_BASE_URL?: string
  }
}

export function readApiBaseUrl(
  value: string | undefined = import.meta.env.VITE_API_BASE_URL,
): string | null {
  if (typeof value !== 'string') {
    return null
  }
  const trimmed = value.trim().replace(/\/+$/, '')
  return trimmed === '' ? null : trimmed
}

export function collectionsUrl(baseUrl: string): string {
  return `${readApiBaseUrl(baseUrl) ?? ''}/api/v1/collections`
}

export function productsUrl(baseUrl: string): string {
  return `${readApiBaseUrl(baseUrl) ?? ''}/api/v1/products`
}

export function productUrl(baseUrl: string, slug: string): string {
  return `${readApiBaseUrl(baseUrl) ?? ''}/api/v1/products/${encodeURIComponent(slug)}`
}

export function fetchCollections(
  baseUrl: string,
): Promise<CatalogResult<ApiCollection[]>> {
  return fetchData(collectionsUrl(baseUrl), parseCollectionList)
}

export function fetchProducts(
  baseUrl: string,
): Promise<CatalogResult<ApiProduct[]>> {
  return fetchData(productsUrl(baseUrl), parseProductList)
}

export function fetchProduct(
  baseUrl: string,
  slug: string,
): Promise<CatalogResult<ApiProduct>> {
  return fetchData(productUrl(baseUrl, slug), parseProduct)
}

async function fetchData<T>(
  url: string,
  parse: (body: unknown) => T | null,
): Promise<CatalogResult<T>> {
  let response: Response
  try {
    response = await fetch(url, {
      method: 'GET',
      headers: { Accept: 'application/json' },
      credentials: 'omit',
    })
  } catch {
    return { ok: false, failure: { kind: 'network' } }
  }

  let body: unknown
  try {
    body = await response.json()
  } catch {
    if (!response.ok) {
      return { ok: false, failure: { kind: 'status', status: response.status } }
    }
    return { ok: false, failure: { kind: 'malformed' } }
  }

  if (response.status === 404) {
    return { ok: false, failure: { kind: 'not_found' } }
  }
  if (!response.ok) {
    return { ok: false, failure: { kind: 'status', status: response.status } }
  }

  const data = parse(body)
  if (data === null) {
    return { ok: false, failure: { kind: 'malformed' } }
  }
  return { ok: true, data }
}

export function parseCollectionList(body: unknown): ApiCollection[] | null {
  if (!isRecord(body) || !Array.isArray(body.data)) {
    return null
  }
  if (!body.data.every(isApiCollection)) {
    return null
  }
  return body.data
}

export function parseProductList(body: unknown): ApiProduct[] | null {
  if (!isRecord(body) || !Array.isArray(body.data)) {
    return null
  }
  if (!body.data.every(isApiProduct)) {
    return null
  }
  return body.data
}

export function parseProduct(body: unknown): ApiProduct | null {
  if (!isRecord(body) || !isApiProduct(body.data)) {
    return null
  }
  return body.data
}

function isApiCollection(value: unknown): value is ApiCollection {
  if (!isRecord(value)) {
    return false
  }
  return (
    isString(value.slug) &&
    isString(value.name) &&
    isString(value.description) &&
    isApiImage(value.image) &&
    isString(value.seoTitle) &&
    isString(value.metaDescription) &&
    isInteger(value.sortOrder)
  )
}

function isApiProduct(value: unknown): value is ApiProduct {
  if (!isRecord(value)) {
    return false
  }
  return (
    isString(value.id) &&
    isString(value.slug) &&
    isString(value.name) &&
    isInteger(value.priceCents) &&
    value.priceCents >= 0 &&
    isApiCollectionRef(value.collection) &&
    isString(value.summary) &&
    isString(value.description) &&
    isHighlightTrio(value.highlights) &&
    isApiImage(value.cardImage) &&
    isApiImage(value.detailImage) &&
    isInteger(value.sortOrder)
  )
}

function isApiCollectionRef(value: unknown): value is ApiCollectionRef {
  return isRecord(value) && isString(value.slug) && isString(value.name)
}

function isApiImage(value: unknown): value is ApiImage {
  return (
    isRecord(value) &&
    isString(value.src) &&
    value.src !== '' &&
    isString(value.alt)
  )
}

function isHighlightTrio(value: unknown): value is [string, string, string] {
  return Array.isArray(value) && value.length === 3 && value.every(isString)
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function isString(value: unknown): value is string {
  return typeof value === 'string'
}

function isInteger(value: unknown): value is number {
  return typeof value === 'number' && Number.isInteger(value)
}
