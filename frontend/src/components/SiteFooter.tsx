import { Link } from 'react-router'
import { useCatalog } from '../api/useCatalog.ts'
import { focusRing } from '../lib/focusRing.ts'
import { Container } from './Container.tsx'

export function SiteFooter() {
  const { collections } = useCatalog()
  const footerLinks = [
    { to: '/', label: 'Home' },
    { to: '/shop', label: 'Shop' },
    ...collections.map((collection) => ({
      to: `/collections/${collection.slug}`,
      label: collection.name,
    })),
  ]

  return (
    <footer className="border-t border-sage">
      <Container className="flex flex-col gap-8 py-10">
        <div className="flex flex-col gap-2">
          <p className="font-display text-2xl text-burgundy">morrow & muse</p>
          <p className="text-sm text-ink">
            Thoughtfully chosen. Beautifully lived.
          </p>
          <p className="max-w-md text-sm leading-relaxed text-ink">
            Objects for slower mornings, clearer spaces, and everyday rituals.
          </p>
        </div>
        <nav
          aria-label="Footer"
          className="flex flex-col gap-1 sm:flex-row sm:flex-wrap sm:gap-x-6"
        >
          {footerLinks.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              className={`inline-flex min-h-11 items-center text-sm text-ink hover:text-burgundy ${focusRing}`}
            >
              {item.label}
            </Link>
          ))}
        </nav>
      </Container>
    </footer>
  )
}
