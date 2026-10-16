import { Link } from 'react-router'
import { getCollection, type Product } from '../data/catalog.ts'
import { focusRing } from '../lib/focusRing.ts'
import { formatPrice } from '../lib/money.ts'
import { CatalogPlaceholder } from './CatalogPlaceholder.tsx'

export function ProductCard({ product }: { product: Product }) {
  const collection = getCollection(product.collectionSlug)

  return (
    <article className="min-w-0">
      <Link
        to={`/products/${product.slug}`}
        className={`group block ${focusRing}`}
      >
        <CatalogPlaceholder image={product.cardImage} decorativePlaceholder />
        <p className="mt-4 text-sm text-taupe">{collection?.name}</p>
        <h3 className="mt-1 font-display text-2xl leading-tight text-burgundy group-hover:underline">
          {product.name}
        </h3>
        <p className="mt-2 text-sm leading-relaxed text-ink">
          {product.summary}
        </p>
        <p className="mt-2 text-sm font-medium text-ink">
          {formatPrice(product.priceCents)}
        </p>
      </Link>
    </article>
  )
}
