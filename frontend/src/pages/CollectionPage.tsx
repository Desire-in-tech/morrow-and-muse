import { useParams } from 'react-router'
import { Container } from '../components/Container.tsx'
import { ProductCard } from '../components/ProductCard.tsx'
import { getCollection, getProductsByCollection } from '../data/catalog.ts'
import { useDocumentMeta } from '../lib/useDocumentMeta.ts'
import { NotFoundPage } from './NotFoundPage.tsx'

export function CollectionPage() {
  const { slug } = useParams()
  const collection = slug ? getCollection(slug) : undefined

  useDocumentMeta(
    collection?.seoTitle ?? 'Page Not Found | morrow & muse',
    collection?.metaDescription ??
      'Explore the thoughtfully chosen everyday essentials from morrow & muse.',
  )

  if (!collection) {
    return <NotFoundPage />
  }

  const items = getProductsByCollection(collection.slug)

  return (
    <section className="py-16 sm:py-24">
      <Container>
        <h1 className="font-display text-5xl leading-tight text-burgundy">
          {collection.name}
        </h1>
        <p className="mt-4 max-w-xl text-base leading-relaxed text-ink">
          {collection.description}
        </p>
        {items.length === 0 ? (
          <>
            <h2 className="mt-12 font-display text-3xl text-burgundy">
              Nothing here just yet.
            </h2>
            <p className="mt-4 max-w-md text-base leading-relaxed text-ink">
              There are no pieces to show in this collection at the moment.
              Explore the full collection to find something that feels right.
            </p>
          </>
        ) : (
          <>
            <h2 className="sr-only">Products</h2>
            <div className="mt-12 grid grid-cols-1 gap-x-6 gap-y-12 sm:grid-cols-2 lg:grid-cols-3">
              {items.map((product) => (
                <ProductCard key={product.slug} product={product} />
              ))}
            </div>
          </>
        )}
      </Container>
    </section>
  )
}
