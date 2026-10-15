import { Outlet } from 'react-router'
import { SiteFooter } from './SiteFooter.tsx'
import { SiteHeader } from './SiteHeader.tsx'

export function Shell() {
  return (
    <div className="flex min-h-screen flex-col bg-ivory font-sans text-ink">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-10 focus:bg-burgundy focus:px-4 focus:py-3 focus:text-ivory"
      >
        Skip to content
      </a>
      <SiteHeader />
      <main id="main" tabIndex={-1} className="flex-1">
        <Outlet />
      </main>
      <SiteFooter />
    </div>
  )
}
