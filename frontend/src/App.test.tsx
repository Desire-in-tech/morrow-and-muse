import { render, screen } from '@testing-library/react'
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

describe('shell', () => {
  it('renders the brand and tagline on the home route', () => {
    renderAt('/')

    expect(
      screen.getByRole('heading', { level: 1, name: 'Morrow & Muse' }),
    ).toBeInTheDocument()
    expect(
      screen.getByText('Thoughtfully chosen. Beautifully lived.'),
    ).toBeInTheDocument()
  })

  it('renders a not-found page for an unknown route', () => {
    renderAt('/not-a-real-page')

    expect(
      screen.getByRole('heading', { level: 1, name: 'Page not found' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('link', { name: 'Back to home' }),
    ).toBeInTheDocument()
  })
})
