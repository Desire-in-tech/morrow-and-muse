import { describe, expect, it, vi } from 'vitest'
import {
  collectionsUrl,
  fetchCollections,
  fetchProduct,
  fetchProducts,
  parseCollectionList,
  parseProduct,
  productUrl,
  productsUrl,
  readApiBaseUrl,
} from './catalogClient.ts'
import { mapCatalog } from './mapCatalog.ts'
import type { ApiCollection, ApiProduct } from './catalogClient.ts'

const image = {
  src: 'https://res.cloudinary.com/sbll6gvk/image/upload/v1/mug.png',
  alt: 'A mug',
}

function apiCollection(overrides: Partial<ApiCollection> = {}): ApiCollection {
  return {
    slug: 'morning-ritual',
    name: 'Morning Ritual',
    description: 'A slower breakfast.',
    image,
    seoTitle: 'Morning Ritual | morrow & muse',
    metaDescription: 'Coffee and breakfast.',
    sortOrder: 1,
    ...overrides,
  }
}

function apiProduct(overrides: Partial<ApiProduct> = {}): ApiProduct {
  return {
    id: 'mm-01',
    slug: 'everyday-ceramic-mug',
    name: 'The Everyday Ceramic Mug',
    priceCents: 2800,
    collection: { slug: 'morning-ritual', name: 'Morning Ritual' },
    summary: 'A morning mug.',
    description: 'Make a little more of the everyday.',
    highlights: ['One.', 'Two.', 'Three.'],
    cardImage: image,
    detailImage: {
      src: 'https://res.cloudinary.com/sbll6gvk/image/upload/v1/mug-detail.png',
      alt: 'Mug detail',
    },
    sortOrder: 2,
    ...overrides,
  }
}

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

describe('catalogue client', () => {
  it('normalizes the base URL and encodes the product slug', () => {
    expect(readApiBaseUrl(undefined)).toBeNull()
    expect(readApiBaseUrl('   ')).toBeNull()
    expect(readApiBaseUrl('http://api.test///')).toBe('http://api.test')
    expect(collectionsUrl('http://api.test/')).toBe(
      'http://api.test/api/v1/collections',
    )
    expect(productsUrl('http://api.test')).toBe(
      'http://api.test/api/v1/products',
    )
    expect(productUrl('http://api.test/', 'a b/c')).toBe(
      'http://api.test/api/v1/products/a%20b%2Fc',
    )
  })

  it('parses collection and product envelopes', async () => {
    const fetchMock = vi.fn((url: string) => {
      if (url.endsWith('/collections')) {
        return Promise.resolve(
          jsonResponse({ data: [apiCollection({ sortOrder: 2 })] }),
        )
      }
      if (url.endsWith('/products/everyday-ceramic-mug')) {
        return Promise.resolve(jsonResponse({ data: apiProduct() }))
      }
      return Promise.resolve(jsonResponse({ data: [apiProduct()] }))
    })
    vi.stubGlobal('fetch', fetchMock)

    const collections = await fetchCollections('http://api.test/')
    const products = await fetchProducts('http://api.test/')
    const product = await fetchProduct(
      'http://api.test/',
      'everyday-ceramic-mug',
    )

    expect(fetchMock).toHaveBeenNthCalledWith(
      1,
      'http://api.test/api/v1/collections',
      expect.objectContaining({ method: 'GET', credentials: 'omit' }),
    )
    expect(fetchMock).toHaveBeenNthCalledWith(
      2,
      'http://api.test/api/v1/products',
      expect.objectContaining({ method: 'GET' }),
    )
    expect(fetchMock).toHaveBeenNthCalledWith(
      3,
      'http://api.test/api/v1/products/everyday-ceramic-mug',
      expect.objectContaining({ method: 'GET' }),
    )
    expect(collections.ok && collections.data[0]?.slug).toBe('morning-ritual')
    expect(products.ok && products.data[0]?.priceCents).toBe(2800)
    expect(product.ok && product.data.highlights).toEqual([
      'One.',
      'Two.',
      'Three.',
    ])
  })

  it('reports network, non-2xx, not found, and malformed bodies without throwing', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(() => Promise.reject(new TypeError('Failed to fetch'))),
    )
    const network = await fetchCollections('http://api.test')
    expect(network).toEqual({ ok: false, failure: { kind: 'network' } })

    vi.stubGlobal(
      'fetch',
      vi.fn(() =>
        Promise.resolve(
          jsonResponse(
            {
              error: {
                code: 'internal',
                message: 'postgres://user:secret@db select * from products',
              },
            },
            500,
          ),
        ),
      ),
    )
    const status = await fetchProducts('http://api.test')
    expect(status).toEqual({
      ok: false,
      failure: { kind: 'status', status: 500 },
    })
    expect(JSON.stringify(status)).not.toContain('secret')
    expect(JSON.stringify(status)).not.toContain('select')

    vi.stubGlobal(
      'fetch',
      vi.fn(() =>
        Promise.resolve(
          jsonResponse(
            { error: { code: 'not_found', message: 'Product not found.' } },
            404,
          ),
        ),
      ),
    )
    const missing = await fetchProduct('http://api.test', 'missing-piece')
    expect(missing).toEqual({ ok: false, failure: { kind: 'not_found' } })

    vi.stubGlobal(
      'fetch',
      vi.fn(() =>
        Promise.resolve(
          new Response('not-json', {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          }),
        ),
      ),
    )
    const broken = await fetchCollections('http://api.test')
    expect(broken).toEqual({ ok: false, failure: { kind: 'malformed' } })

    expect(parseCollectionList({ data: [{ slug: 'only-slug' }] })).toBeNull()
    expect(
      parseProduct({ data: { ...apiProduct(), highlights: ['one'] } }),
    ).toBeNull()
    expect(parseProduct({ data: apiProduct({ priceCents: -1 }) })).toBeNull()
  })
})

describe('catalogue mapping', () => {
  it('maps API fields onto storefront types and display order', () => {
    const mapped = mapCatalog(
      [
        apiCollection({
          slug: 'quiet-space',
          name: 'Quiet Space',
          sortOrder: 2,
        }),
        apiCollection({ sortOrder: 1 }),
      ],
      [
        apiProduct({
          slug: 'later',
          name: 'Later',
          sortOrder: 2,
          collection: { slug: 'quiet-space', name: 'Quiet Space' },
        }),
        apiProduct({ sortOrder: 1 }),
      ],
    )

    expect(mapped?.collections.map((collection) => collection.slug)).toEqual([
      'morning-ritual',
      'quiet-space',
    ])
    expect(mapped?.products.map((product) => product.slug)).toEqual([
      'everyday-ceramic-mug',
      'later',
    ])
    expect(mapped?.products[0]).toMatchObject({
      name: 'The Everyday Ceramic Mug',
      collectionSlug: 'morning-ritual',
      priceCents: 2800,
      highlights: ['One.', 'Two.', 'Three.'],
      cardImage: image,
    })
    expect(mapped?.collections[0]).not.toHaveProperty('sortOrder')
  })

  it('rejects a payload that cannot be one catalogue', () => {
    expect(
      mapCatalog(
        [apiCollection()],
        [apiProduct({ collection: { slug: 'missing', name: 'Missing' } })],
      ),
    ).toBeNull()
    expect(
      mapCatalog([apiCollection(), apiCollection()], [apiProduct()]),
    ).toBeNull()
  })
})
