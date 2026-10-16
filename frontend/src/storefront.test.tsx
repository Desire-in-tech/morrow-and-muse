import { fireEvent, render, screen, within } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { describe, expect, it } from 'vitest'
import App from './App.tsx'

function renderAt(path: string) {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <App />
    </MemoryRouter>,
  )
}

function headingLevels() {
  return screen.getAllByRole('heading').map((heading) => {
    return Number(heading.tagName.slice(1))
  })
}

function expectNoSkippedLevels() {
  const levels = headingLevels()
  expect(levels[0]).toBe(1)
  for (let index = 1; index < levels.length; index += 1) {
    expect(levels[index]).toBeLessThanOrEqual((levels[index - 1] ?? 1) + 1)
  }
}

function description() {
  return document
    .querySelector('meta[name="description"]')
    ?.getAttribute('content')
}

describe('storefront', () => {
  it('renders the hero and a Shop link on the home route', () => {
    renderAt('/')

    expect(
      screen.getByRole('heading', {
        level: 1,
        name: 'Make room for the good things.',
      }),
    ).toBeInTheDocument()
    expect(
      screen.getByText(
        'Thoughtfully chosen objects for slower mornings, clearer spaces, and the small rituals that make a day feel like yours.',
      ),
    ).toBeInTheDocument()
    expect(
      screen.getAllByRole('link', { name: 'Shop' }).length,
    ).toBeGreaterThan(0)
    expect(
      screen.getAllByRole('link', { name: 'Explore the collection' }).length,
    ).toBeGreaterThan(0)
  })

  it('toggles the mobile menu with aria-expanded and closes on Escape', () => {
    renderAt('/')
    const menuButton = screen.getByRole('button', { name: 'Open menu' })

    expect(menuButton).toHaveAttribute('aria-expanded', 'false')
    expect(menuButton).toHaveAttribute('aria-controls', 'primary-menu')

    menuButton.focus()
    fireEvent.click(menuButton)
    expect(menuButton).toHaveAttribute('aria-expanded', 'true')
    expect(menuButton).toHaveAccessibleName('Close menu')
    expect(document.activeElement).toBe(menuButton)

    fireEvent.keyDown(document, { key: 'Escape' })
    expect(menuButton).toHaveAttribute('aria-expanded', 'false')
    expect(menuButton).toHaveAccessibleName('Open menu')
    expect(document.activeElement).toBe(menuButton)
  })

  it('closes the mobile menu when the route changes', () => {
    renderAt('/')
    const menuButton = screen.getByRole('button', { name: 'Open menu' })
    fireEvent.click(menuButton)

    const menu = document.getElementById('primary-menu')
    const shop = menu?.querySelector<HTMLAnchorElement>('a[href="/shop"]')
    expect(shop).not.toBeNull()
    shop?.focus()
    fireEvent.click(shop as HTMLAnchorElement)

    expect(menuButton).toHaveAttribute('aria-expanded', 'false')
    expect(document.activeElement).toBe(document.getElementById('main'))
    expect(
      screen.getByRole('heading', {
        level: 1,
        name: 'Shop thoughtfully chosen essentials',
      }),
    ).toBeInTheDocument()
  })

  it('does not trap Tab in the disclosure menu', () => {
    renderAt('/')
    fireEvent.click(screen.getByRole('button', { name: 'Open menu' }))
    const links = document.querySelectorAll('#primary-menu a[href]')
    const first = links[0] as HTMLElement
    const last = links[links.length - 1] as HTMLElement

    last.focus()
    fireEvent.keyDown(document, { key: 'Tab' })
    expect(document.activeElement).toBe(last)

    first.focus()
    fireEvent.keyDown(document, { key: 'Tab', shiftKey: true })
    expect(document.activeElement).toBe(first)
  })

  it('ignores Escape when focus is outside the menu and toggle', () => {
    renderAt('/')
    const menuButton = screen.getByRole('button', { name: 'Open menu' })
    fireEvent.click(menuButton)
    const heroShop = screen
      .getAllByRole('link', { name: 'Shop' })
      .find((link) => !link.closest('#primary-menu'))
    heroShop?.focus()

    fireEvent.keyDown(document, { key: 'Escape' })
    expect(menuButton).toHaveAttribute('aria-expanded', 'true')
  })

  it('returns focus to the toggle when Escape is pressed on a menu link', () => {
    renderAt('/')
    const menuButton = screen.getByRole('button', { name: 'Open menu' })
    fireEvent.click(menuButton)
    const home = document.querySelector(
      '#primary-menu a[href="/"]',
    ) as HTMLElement
    home.focus()

    fireEvent.keyDown(document, { key: 'Escape' })
    expect(menuButton).toHaveAttribute('aria-expanded', 'false')
    expect(document.activeElement).toBe(menuButton)
  })

  it('renders a collection and its products', () => {
    renderAt('/collections/morning-ritual')

    expect(
      screen.getByRole('heading', { level: 1, name: 'Morning Ritual' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', {
        level: 3,
        name: 'The Everyday Ceramic Mug',
      }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', {
        level: 3,
        name: 'The Stoneware Breakfast Bowl',
      }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', {
        level: 3,
        name: 'The Minimal Ceramic Vase',
      }),
    ).toBeInTheDocument()
  })

  it('falls back when a collection slug is unknown', () => {
    renderAt('/collections/not-a-collection')

    expect(
      screen.getByRole('heading', {
        level: 1,
        name: 'This page is taking a quieter path.',
      }),
    ).toBeInTheDocument()
  })

  it('shows a product card name, price, and card image', () => {
    renderAt('/shop')

    expect(
      screen.getByRole('heading', {
        level: 3,
        name: 'The Everyday Ceramic Mug',
      }),
    ).toBeInTheDocument()
    expect(screen.getAllByText('$28.00').length).toBeGreaterThan(0)
    const card = document.querySelector(
      'a[href="/products/everyday-ceramic-mug"] img',
    )
    expect(card).toHaveAttribute(
      'src',
      'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079402/product_card_ceramic_coffee_mug.png',
    )
    expect(card).toHaveAttribute('alt', '')
  })

  it('keeps a heading level sequence without skips', () => {
    const home = renderAt('/')
    expectNoSkippedLevels()
    expect(
      screen.getByRole('heading', {
        level: 2,
        name: 'A place for every little ritual',
      }),
    ).toBeInTheDocument()
    home.unmount()

    const shop = renderAt('/shop')
    expectNoSkippedLevels()
    expect(
      screen.getByRole('heading', { level: 2, name: 'Products' }),
    ).toBeInTheDocument()
    shop.unmount()

    const collection = renderAt('/collections/morning-ritual')
    expectNoSkippedLevels()
    expect(
      screen.getByRole('heading', { level: 1, name: 'Morning Ritual' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { level: 2, name: 'Products' }),
    ).toBeInTheDocument()
    collection.unmount()

    renderAt('/products/everyday-ceramic-mug')
    expectNoSkippedLevels()
    expect(headingLevels()[0]).toBe(1)
    expect(headingLevels()).toContain(2)
  })

  it('sets the document title and description for each route', () => {
    const home = renderAt('/')
    expect(document.title).toBe(
      'morrow & muse | Thoughtfully Chosen. Beautifully Lived.',
    )
    expect(description()).toMatch(/slower mornings/)
    home.unmount()

    const shop = renderAt('/shop')
    expect(document.title).toBe('Shop All | morrow & muse')
    expect(description()).toMatch(/all 15 considered essentials/)
    shop.unmount()

    const collection = renderAt('/collections/morning-ritual')
    expect(document.title).toBe('Morning Ritual | morrow & muse')
    expect(description()).toMatch(/coffee, breakfast/)
    collection.unmount()

    const product = renderAt('/products/everyday-ceramic-mug')
    expect(document.title).toBe('The Everyday Ceramic Mug | morrow & muse')
    expect(description()).toMatch(/first coffee/)
    product.unmount()

    const unknownCollection = renderAt('/collections/not-a-collection')
    expect(document.title).toBe('Page Not Found | morrow & muse')
    expect(description()).toBe(
      'Explore the thoughtfully chosen everyday essentials from morrow & muse.',
    )
    unknownCollection.unmount()

    const unknownProduct = renderAt('/products/missing-piece')
    expect(document.title).toBe('Page Not Found | morrow & muse')
    unknownProduct.unmount()

    renderAt('/not-a-real-page')
    expect(document.title).toBe('Page Not Found | morrow & muse')
    expect(description()).toBe(
      'Explore the thoughtfully chosen everyday essentials from morrow & muse.',
    )
  })

  it('renders a product page with the detail image and breadcrumb', () => {
    renderAt('/products/everyday-ceramic-mug')
    const main = document.getElementById('main') as HTMLElement

    expect(
      screen.getAllByRole('heading', { name: 'The Everyday Ceramic Mug' }),
    ).toHaveLength(1)
    expect(within(main).getByText('$28.00')).toBeInTheDocument()
    expect(
      within(main).getByText(/Make a little more of the everyday/),
    ).toBeInTheDocument()
    expect(main.querySelectorAll('ul li')).toHaveLength(3)
    const detail = main.querySelector('img')
    expect(detail).toHaveAttribute(
      'src',
      'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079748/product_detail_image_ceramic_coffee_mug.png',
    )
    expect(detail).toHaveAttribute(
      'alt',
      'Ceramic coffee mug photographed in a warm, neutral setting.',
    )
    const breadcrumb = screen.getByRole('navigation', { name: 'Breadcrumb' })
    expect(
      within(breadcrumb).getByRole('link', { name: 'Home' }),
    ).toHaveAttribute('href', '/')
    expect(
      within(breadcrumb).getByRole('link', { name: 'Morning Ritual' }),
    ).toHaveAttribute('href', '/collections/morning-ritual')
    expect(
      within(breadcrumb).getByText('The Everyday Ceramic Mug'),
    ).toHaveAttribute('aria-current', 'page')
    expect(
      within(main).getByRole('link', { name: 'Back to Morning Ritual' }),
    ).toHaveAttribute('href', '/collections/morning-ritual')
  })

  it('keeps dev wording out of the customer-facing pages', () => {
    const banned = /\bdemo\b|placeholder|portfolio boutique/i
    for (const path of [
      '/',
      '/shop',
      '/collections/morning-ritual',
      '/collections/workday-edit',
      '/collections/everyday-carry',
      '/collections/quiet-space',
      '/products/sculptural-table-lamp',
      '/products/missing-piece',
    ]) {
      const view = renderAt(path)
      expect(document.body.textContent ?? '').not.toMatch(banned)
      for (const img of document.querySelectorAll('img')) {
        expect(img.getAttribute('alt') ?? '').not.toMatch(banned)
      }
      view.unmount()
    }
  })
})
