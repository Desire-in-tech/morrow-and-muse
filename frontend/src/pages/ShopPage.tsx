import { Container } from '../components/Container.tsx'
import { ProductCard } from '../components/ProductCard.tsx'
import { useCatalog } from '../api/useCatalog.ts'
import { useDocumentMeta } from '../lib/useDocumentMeta.ts'

export function ShopPage() {
  const { products } = useCatalog()
  useDocumentMeta(
    'Shop All | morrow & muse',
    'Explore all 15 considered essentials from morrow & muse, from morning rituals and stationery to everyday carry and quiet home accents.',
  )

  return (
    <section className="py-16 sm:py-24">
      <Container>
        <h1 className="font-display text-5xl leading-tight text-burgundy">
          Shop thoughtfully chosen essentials
        </h1>
        <p className="mt-4 max-w-xl text-base leading-relaxed text-ink">
          Explore the Morrow & Muse collection: understated pieces for morning
          rituals, focused workdays, everyday carry and quieter spaces at home.
        </p>
        <h2 className="sr-only">Products</h2>
        <div className="mt-12 grid grid-cols-1 gap-x-6 gap-y-12 sm:grid-cols-2 lg:grid-cols-3">
          {products.map((product) => (
            <ProductCard key={product.slug} product={product} />
          ))}
        </div>
      </Container>
    </section>
  )
}
