import { act, render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { describe, expect, it, vi } from 'vitest'
import App from '../App.tsx'
import { readApiBaseUrl } from './catalogClient.ts'
import type { ApiCollection, ApiProduct } from './catalogClient.ts'

const image = {
  src: 'https://res.cloudinary.com/sbll6gvk/image/upload/v1/api-mug.png',
  alt: 'API mug',
}

const apiCollection: ApiCollection = {
  slug: 'morning-ritual',
  name: 'API Morning Ritual',
  description: 'API description for the morning.',
  image,
  seoTitle: 'API Morning | morrow & muse',
  metaDescription: 'API meta for coffee.',
  sortOrder: 1,
}

const otherCollection: ApiCollection = {
  ...apiCollection,
  slug: 'quiet-space',
  name: 'API Quiet Space',
  seoTitle: 'API Quiet | morrow & muse',
  sortOrder: 2,
}

const apiProduct: ApiProduct = {
  id: 'mm-01',
  slug: 'everyday-ceramic-mug',
  name: 'API Ceramic Mug',
  priceCents: 2800,
  collection: { slug: 'morning-ritual', name: 'API Morning Ritual' },
  summary: 'API summary for the first coffee.',
  description: 'API description of the mug.',
  highlights: ['API one.', 'API two.', 'API three.'],
  cardImage: image,
  detailImage: image,
  sortOrder: 1,
}

const otherProduct: ApiProduct = {
  ...apiProduct,
  id: 'mm-15',
  slug: 'small-ceramic-planter',
  name: 'API Planter',
  collection: { slug: 'quiet-space', name: 'API Quiet Space' },
  sortOrder: 2,
}

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

async function settleFetches(fetchMock: ReturnType<typeof vi.fn>) {
  await vi.waitFor(() => {
    expect(fetchMock).toHaveBeenCalled()
  })
  await act(async () => {
    await Promise.resolve()
    await Promise.resolve()
  })
}

function renderAt(path: string) {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <App />
    </MemoryRouter>,
  )
}

function stubApi(handler: (url: string) => Promise<Response> | Response) {
  vi.stubEnv('VITE_API_BASE_URL', 'http://api.test/')
  const fetchMock = vi.fn((url: string) => Promise.resolve(handler(url)))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

describe('catalogue session', () => {
  it('uses the bundle and does not fetch when the base URL is unset or blank', () => {
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)

    const unset = renderAt('/shop')
    expect(
      screen.getByRole('heading', {
        level: 3,
        name: 'The Everyday Ceramic Mug',
      }),
    ).toBeInTheDocument()
    expect(fetchMock).not.toHaveBeenCalled()
    unset.unmount()

    vi.stubEnv('VITE_API_BASE_URL', '   ')
    expect(readApiBaseUrl()).toBeNull()
    renderAt('/shop')
    expect(
      screen.getByRole('heading', {
        level: 3,
        name: 'The Everyday Ceramic Mug',
      }),
    ).toBeInTheDocument()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('displays API data when both catalogue requests succeed', async () => {
    const fetchMock = stubApi((url) => {
      if (url.endsWith('/collections')) {
        return jsonResponse({ data: [otherCollection, apiCollection] })
      }
      return jsonResponse({ data: [otherProduct, apiProduct] })
    })

    const shop = renderAt('/shop')
    expect(
      await screen.findByRole('heading', { level: 3, name: 'API Ceramic Mug' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { level: 3, name: 'API Planter' }),
    ).toBeInTheDocument()
    expect(
      screen.queryByRole('heading', { name: 'The Everyday Ceramic Mug' }),
    ).not.toBeInTheDocument()
    expect(fetchMock).toHaveBeenCalledTimes(2)
    shop.unmount()

    renderAt('/collections/morning-ritual')
    expect(
      await screen.findByRole('heading', {
        level: 1,
        name: 'API Morning Ritual',
      }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { level: 3, name: 'API Ceramic Mug' }),
    ).toBeInTheDocument()
    expect(
      screen.queryByRole('heading', { name: 'API Planter' }),
    ).not.toBeInTheDocument()
    expect(document.title).toBe('API Morning | morrow & muse')

    renderAt('/products/everyday-ceramic-mug')
    expect(
      await screen.findByRole('heading', { level: 1, name: 'API Ceramic Mug' }),
    ).toBeInTheDocument()
    expect(screen.getByText('API description of the mug.')).toBeInTheDocument()
    expect(screen.getByText('API one.')).toBeInTheDocument()
    expect(document.title).toBe('API Ceramic Mug | morrow & muse')
    expect(fetchMock).toHaveBeenCalledTimes(2)
  })

  it('keeps the bundle when the network fails', async () => {
    vi.stubEnv('VITE_API_BASE_URL', 'http://api.test')
    const fetchMock = vi.fn(() =>
      Promise.reject(new TypeError('Failed to fetch')),
    )
    vi.stubGlobal('fetch', fetchMock)

    renderAt('/shop')
    await settleFetches(fetchMock)
    expect(
      screen.getByRole('heading', { name: 'The Everyday Ceramic Mug' }),
    ).toBeInTheDocument()
    expect(document.body.textContent ?? '').not.toContain('Failed to fetch')
  })

  it('keeps the whole bundle when only one catalogue request fails', async () => {
    const secret = 'postgres://user:secret@db select * from products'
    vi.stubEnv('VITE_API_BASE_URL', 'http://api.test')
    const fetchMock = vi.fn((url: string) => {
      if (String(url).endsWith('/collections')) {
        return Promise.resolve(jsonResponse({ data: [apiCollection] }))
      }
      return Promise.resolve(
        jsonResponse({ error: { code: 'internal', message: secret } }, 500),
      )
    })
    vi.stubGlobal('fetch', fetchMock)

    renderAt('/shop')
    await settleFetches(fetchMock)
    expect(
      screen.getByRole('heading', { name: 'The Everyday Ceramic Mug' }),
    ).toBeInTheDocument()
    expect(screen.queryByText('API Morning Ritual')).not.toBeInTheDocument()
    expect(document.body.textContent ?? '').not.toContain('secret')
    expect(document.body.textContent ?? '').not.toContain('select')
  })

  it('keeps the bundle when the JSON is not JSON', async () => {
    vi.stubEnv('VITE_API_BASE_URL', 'http://api.test')
    const fetchMock = vi.fn(() =>
      Promise.resolve(
        new Response('not-json', {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        }),
      ),
    )
    vi.stubGlobal('fetch', fetchMock)

    renderAt('/collections/morning-ritual')
    await settleFetches(fetchMock)
    expect(
      screen.getByRole('heading', { level: 1, name: 'Morning Ritual' }),
    ).toBeInTheDocument()
    expect(document.body.textContent ?? '').not.toContain('not-json')
  })

  it('keeps the bundle when the JSON shape is wrong', async () => {
    vi.stubEnv('VITE_API_BASE_URL', 'http://api.test')
    const fetchMock = vi.fn(() =>
      Promise.resolve(jsonResponse({ data: { slug: 'nope' } })),
    )
    vi.stubGlobal('fetch', fetchMock)

    renderAt('/shop')
    await settleFetches(fetchMock)
    expect(
      screen.getByRole('heading', { name: 'The Everyday Ceramic Mug' }),
    ).toBeInTheDocument()
  })

  it('keeps the not-found page for an unknown slug when the API catalogue loads', async () => {
    const fetchMock = stubApi((url) => {
      if (url.endsWith('/collections')) {
        return jsonResponse({ data: [apiCollection] })
      }
      return jsonResponse({ data: [apiProduct] })
    })

    renderAt('/products/missing-piece')
    await settleFetches(fetchMock)
    expect(
      screen.getByRole('heading', {
        level: 1,
        name: 'This page is taking a quieter path.',
      }),
    ).toBeInTheDocument()
    expect(
      screen.queryByRole('heading', { name: 'API Ceramic Mug' }),
    ).not.toBeInTheDocument()
  })
})
