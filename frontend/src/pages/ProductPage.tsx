import { Link, useParams } from 'react-router'
import { CatalogPlaceholder } from '../components/CatalogPlaceholder.tsx'
import { Container } from '../components/Container.tsx'
import { ProductCard } from '../components/ProductCard.tsx'
import {
  findCollection,
  findProduct,
  relatedProducts,
} from '../api/catalogSession.ts'
import { useCatalog } from '../api/useCatalog.ts'
import { focusRing } from '../lib/focusRing.ts'
import { formatPrice } from '../lib/money.ts'
import { useDocumentMeta } from '../lib/useDocumentMeta.ts'
import { NotFoundPage } from './NotFoundPage.tsx'

export function ProductPage() {
  const { slug } = useParams()
  const catalog = useCatalog()
  const product = slug ? findProduct(catalog, slug) : undefined
  const collection = product
    ? findCollection(catalog, product.collectionSlug)
    : undefined

  useDocumentMeta(
    product
      ? `${product.name} | morrow & muse`
      : 'Page Not Found | morrow & muse',
    product?.summary ??
      'Explore the thoughtfully chosen everyday essentials from morrow & muse.',
  )

  if (!product || !collection) {
    return <NotFoundPage />
  }

  const related = relatedProducts(catalog, product.slug)

  return (
    <section className="py-16 sm:py-24">
      <Container>
        <nav aria-label="Breadcrumb">
          <ol className="flex flex-wrap items-center gap-x-2 gap-y-1 text-sm text-taupe">
            <li>
              <Link to="/" className={`hover:text-burgundy ${focusRing}`}>
                Home
              </Link>
            </li>
            <li aria-hidden="true">/</li>
            <li>
              <Link
                to={`/collections/${collection.slug}`}
                className={`hover:text-burgundy ${focusRing}`}
              >
                {collection.name}
              </Link>
            </li>
            <li aria-hidden="true">/</li>
            <li aria-current="page" className="text-ink">
              {product.name}
            </li>
          </ol>
        </nav>
        <div className="mt-8 grid items-start gap-10 lg:grid-cols-2">
          <CatalogPlaceholder image={product.detailImage} />
          <div className="min-w-0">
            <h1 className="font-display text-5xl leading-tight text-burgundy">
              {product.name}
            </h1>
            <p className="mt-4 text-base font-medium text-ink">
              {formatPrice(product.priceCents)}
            </p>
            <p className="mt-6 max-w-md text-base leading-relaxed text-ink">
              {product.description}
            </p>
            <h2 className="mt-8 font-display text-2xl text-burgundy">
              Highlights
            </h2>
            <ul className="mt-4 list-disc space-y-2 pl-5 text-sm leading-relaxed text-ink">
              {product.highlights.map((highlight) => (
                <li key={highlight}>{highlight}</li>
              ))}
            </ul>
            <Link
              to={`/collections/${collection.slug}`}
              className={`mt-8 inline-flex min-h-11 items-center text-burgundy underline decoration-sandstone underline-offset-4 hover:decoration-burgundy ${focusRing}`}
            >
              Back to {collection.name}
            </Link>
          </div>
        </div>
        {related.length > 0 ? (
          <div className="mt-16">
            <h2 className="font-display text-3xl text-burgundy">
              More from {collection.name}
            </h2>
            <div className="mt-8 grid grid-cols-1 gap-x-6 gap-y-12 sm:grid-cols-2 lg:grid-cols-3">
              {related.map((item) => (
                <ProductCard key={item.slug} product={item} />
              ))}
            </div>
          </div>
        ) : null}
      </Container>
    </section>
  )
}
