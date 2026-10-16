import { useState } from 'react'
import { resolveCatalogImage, type CatalogImage } from '../data/catalog.ts'

/**
 * Renders a catalogue photograph, or the on-brand block if `src` is missing
 * or the request fails. Product cards pass `decorativePlaceholder` so the
 * image alt is empty: the surrounding link already names the piece.
 * Detail, hero, and collection images keep their descriptive alt.
 * The fallback block does not paint a product name.
 */
export function CatalogPlaceholder({
  image,
  decorativePlaceholder = false,
  loading = 'lazy',
  width = 800,
  height = 1000,
  className = 'aspect-[4/5] w-full object-cover',
}: {
  image: CatalogImage
  decorativePlaceholder?: boolean
  loading?: 'eager' | 'lazy'
  width?: number
  height?: number
  className?: string
}) {
  const resolved = resolveCatalogImage(image)
  const [failed, setFailed] = useState(false)
  const frame = className.includes('aspect-')
    ? className
    : `aspect-[4/5] ${className}`

  if (!resolved.src || failed) {
    return (
      <div
        role={decorativePlaceholder ? undefined : 'img'}
        aria-label={decorativePlaceholder ? undefined : resolved.alt}
        aria-hidden={decorativePlaceholder ? true : undefined}
        className={`relative overflow-hidden bg-sandstone ${frame}`}
      >
        <svg
          aria-hidden="true"
          viewBox="0 0 400 500"
          className="h-full w-full"
          focusable="false"
        >
          <rect width="400" height="500" fill="#E8DED0" />
          <rect x="0" y="48" width="400" height="8" fill="#A8B09A" />
          <circle cx="150" cy="190" r="78" fill="#722F3E" />
          <rect x="170" y="250" width="150" height="190" fill="#F7F2EA" />
          <path
            d="M40 420 C 140 300, 220 360, 360 300"
            fill="none"
            stroke="#292521"
            strokeWidth="2"
          />
        </svg>
      </div>
    )
  }

  return (
    <img
      src={resolved.src}
      alt={decorativePlaceholder ? '' : resolved.alt}
      loading={loading}
      decoding="async"
      width={width}
      height={height}
      className={frame}
      onError={() => {
        setFailed(true)
      }}
    />
  )
}
