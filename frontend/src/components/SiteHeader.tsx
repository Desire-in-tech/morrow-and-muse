import { Menu, X } from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import { Link, NavLink, useLocation } from 'react-router'
import { useCatalog } from '../api/useCatalog.ts'
import { logoUrl } from '../data/catalog.ts'
import { focusRing } from '../lib/focusRing.ts'
import { Container } from './Container.tsx'

const menuId = 'primary-menu'

export function SiteHeader() {
  const { collections } = useCatalog()
  const navLinks = [
    { to: '/', label: 'Home', end: true },
    { to: '/shop', label: 'Shop', end: true },
    ...collections.map((collection) => ({
      to: `/collections/${collection.slug}`,
      label: collection.name,
      end: true,
    })),
  ]
  const { pathname } = useLocation()
  const [openPath, setOpenPath] = useState<string | null>(null)
  const open = openPath === pathname
  const buttonRef = useRef<HTMLButtonElement>(null)
  const menuRef = useRef<HTMLElement>(null)
  const menuLabel = open ? 'Close menu' : 'Open menu'

  useEffect(() => {
    function onKeyDown(event: KeyboardEvent) {
      if (event.key !== 'Escape' || !open) {
        return
      }

      const button = buttonRef.current
      const menu = menuRef.current
      const active = document.activeElement
      const insideToggle =
        button === active || Boolean(button?.contains(active))
      const insideMenu = Boolean(menu?.contains(active))
      if (!insideToggle && !insideMenu) {
        return
      }

      event.preventDefault()
      setOpenPath(null)
      button?.focus()
    }

    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [open])

  useEffect(() => {
    if (openPath === null || openPath === pathname) {
      return
    }

    const menu = menuRef.current
    const active = document.activeElement
    if (menu && active instanceof Node && menu.contains(active)) {
      document.getElementById('main')?.focus()
    }
  }, [pathname, openPath])

  return (
    <header className="border-b border-sandstone">
      <Container className="py-5">
        <div className="flex items-center justify-between gap-4">
          <div className="min-w-0">
            <Link
              to="/"
              className={`relative block h-8 w-[7.25rem] overflow-hidden ${focusRing}`}
            >
              <img
                src={logoUrl}
                alt="morrow & muse"
                width={1200}
                height={1200}
                decoding="async"
                className="absolute top-1/2 left-1/2 h-64 w-64 max-w-none -translate-x-1/2 -translate-y-1/2"
              />
            </Link>
            <p className="mt-2 max-w-[14rem] text-[0.6875rem] leading-relaxed tracking-[0.12em] text-taupe">
              THOUGHTFULLY CHOSEN. BEAUTIFULLY LIVED.
            </p>
          </div>
          <button
            ref={buttonRef}
            type="button"
            className={`inline-flex min-h-11 min-w-11 items-center justify-center gap-2 px-3 text-sm text-ink md:hidden ${focusRing}`}
            aria-expanded={open}
            aria-controls={menuId}
            onClick={() => {
              setOpenPath(open ? null : pathname)
            }}
          >
            {open ? (
              <X aria-hidden="true" size={18} />
            ) : (
              <Menu aria-hidden="true" size={18} />
            )}
            {menuLabel}
          </button>
        </div>
        <nav
          id={menuId}
          ref={menuRef}
          aria-label="Primary"
          className={`${open ? 'mt-4 flex' : 'hidden'} flex-col gap-1 md:mt-4 md:flex md:flex-row md:flex-wrap md:items-center md:gap-x-5`}
        >
          {navLinks.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.end}
              className={({ isActive }) =>
                [
                  'inline-flex min-h-11 items-center text-sm',
                  focusRing,
                  isActive
                    ? 'text-burgundy underline decoration-burgundy underline-offset-8'
                    : 'text-ink hover:text-burgundy',
                ].join(' ')
              }
            >
              {item.label}
            </NavLink>
          ))}
        </nav>
      </Container>
    </header>
  )
}
