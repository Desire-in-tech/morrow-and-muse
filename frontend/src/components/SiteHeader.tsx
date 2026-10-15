import { House } from 'lucide-react'
import { Link, NavLink, useMatch } from 'react-router'

const focusRing =
  'focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-burgundy'

export function SiteHeader() {
  const isHome = useMatch({ path: '/', end: true })

  return (
    <header className="border-b border-sandstone">
      <div className="mx-auto flex w-full max-w-5xl flex-col gap-4 px-5 py-5 sm:flex-row sm:items-center sm:justify-between sm:px-8">
        <div className="min-w-0">
          <Link
            to="/"
            className={`font-display text-xl text-burgundy sm:text-2xl ${focusRing}`}
          >
            Morrow & Muse
          </Link>
          {isHome ? null : (
            <p className="mt-1 text-sm text-taupe">
              Thoughtfully chosen. Beautifully lived.
            </p>
          )}
        </div>
        <nav aria-label="Primary">
          <NavLink
            to="/"
            end
            className={({ isActive }) =>
              [
                'inline-flex min-h-11 items-center gap-2 rounded-full px-4 text-sm font-medium',
                focusRing,
                isActive
                  ? 'bg-burgundy text-ivory'
                  : 'text-ink hover:bg-sandstone',
              ].join(' ')
            }
          >
            <House aria-hidden="true" size={16} />
            Home
          </NavLink>
        </nav>
      </div>
    </header>
  )
}
