import { Link } from 'react-router'
import { CatalogPlaceholder } from '../components/CatalogPlaceholder.tsx'
import { CollectionCard } from '../components/CollectionCard.tsx'
import { Container } from '../components/Container.tsx'
import { ProductCard } from '../components/ProductCard.tsx'
import { SectionHeading } from '../components/SectionHeading.tsx'
import { collections, getFeaturedProducts, heroImage } from '../data/catalog.ts'
import { focusRing } from '../lib/focusRing.ts'
import { useDocumentMeta } from '../lib/useDocumentMeta.ts'

const values = [
  {
    title: 'Choose with intention',
    copy: 'A thoughtful collection begins with the details: useful forms, considered presentation and pieces that feel at home in everyday life.',
  },
  {
    title: 'Make space to breathe',
    copy: 'A calmer space is not about having more. It is about bringing together the things that serve a purpose and make a place feel like yours.',
  },
  {
    title: 'Find beauty in the familiar',
    copy: 'The everyday deserves attention, too. A morning cup, an open notebook or a favourite corner can become a ritual worth looking forward to.',
  },
]

export function HomePage() {
  const featured = getFeaturedProducts()
  useDocumentMeta(
    'morrow & muse | Thoughtfully Chosen. Beautifully Lived.',
    'Thoughtfully chosen home, stationery and everyday essentials for slower mornings, clearer spaces and everyday rituals.',
  )

  return (
    <>
      <section className="border-b border-sandstone">
        <Container className="grid items-end gap-10 py-16 sm:py-24 lg:grid-cols-[minmax(0,1.1fr)_minmax(0,0.9fr)]">
          <div className="min-w-0">
            <p className="text-sm tracking-[0.16em] text-taupe">
              THE ART OF EVERYDAY LIVING
            </p>
            <h1 className="mt-4 font-display text-5xl leading-tight text-burgundy sm:text-6xl">
              Make room for the good things.
            </h1>
            <p className="mt-6 max-w-md text-base leading-relaxed text-ink">
              Thoughtfully chosen objects for slower mornings, clearer spaces,
              and the small rituals that make a day feel like yours.
            </p>
            <div className="mt-8 flex flex-wrap gap-4">
              <a
                href="#collections"
                className={`inline-flex min-h-11 items-center bg-burgundy px-6 text-sm text-ivory hover:bg-ink ${focusRing}`}
              >
                Explore the collection
              </a>
              <Link
                to="/shop"
                className={`inline-flex min-h-11 items-center border border-burgundy px-6 text-sm text-burgundy hover:bg-sandstone ${focusRing}`}
              >
                Shop all
              </Link>
            </div>
          </div>
          <CatalogPlaceholder
            image={heroImage}
            loading="eager"
            width={1800}
            height={1200}
            className="aspect-[3/2] w-full object-cover"
          />
        </Container>
      </section>

      <section id="collections" className="py-16 sm:py-24">
        <Container>
          <SectionHeading
            title="A place for every little ritual"
            lede="From the first cup of the morning to the details that make a room feel like home, discover pieces chosen to bring more intention to everyday living."
          />
          <div className="mt-10 grid grid-cols-1 gap-10 sm:grid-cols-2">
            {collections.map((collection) => (
              <CollectionCard key={collection.slug} collection={collection} />
            ))}
          </div>
        </Container>
      </section>

      <section className="border-t border-sandstone py-16 sm:py-24">
        <Container>
          <SectionHeading
            title="Considered favourites"
            lede="Everyday pieces with a quiet point of view, chosen to work beautifully with the way you live."
          />
          <div className="mt-10 grid grid-cols-1 gap-x-6 gap-y-12 sm:grid-cols-2 lg:grid-cols-4">
            {featured.map((product) => (
              <ProductCard key={product.slug} product={product} />
            ))}
          </div>
          <Link
            to="/shop"
            className={`mt-12 inline-flex min-h-11 items-center text-sm text-burgundy underline decoration-sandstone underline-offset-4 hover:decoration-burgundy ${focusRing}`}
          >
            Discover all pieces
          </Link>
        </Container>
      </section>

      <section className="bg-sandstone py-16 sm:py-24">
        <Container>
          <SectionHeading
            title="The beauty is in the everyday."
            lede="We believe a considered life is shaped by small, meaningful choices. The objects we use, the spaces we make, and the rituals we return to can bring a little more ease to an ordinary day."
          />
          <ul className="mt-10 grid grid-cols-1 gap-8 sm:grid-cols-3">
            {values.map((value) => (
              <li key={value.title} className="min-w-0">
                <h3 className="font-display text-2xl text-burgundy">
                  {value.title}
                </h3>
                <p className="mt-3 text-sm leading-relaxed text-ink">
                  {value.copy}
                </p>
              </li>
            ))}
          </ul>
        </Container>
      </section>

      <section className="py-16 sm:py-20">
        <Container>
          <SectionHeading
            title="A slower kind of everyday"
            lede="Small rituals. Thoughtful objects. Spaces that feel like your own. Morrow & Muse is an invitation to make a little more room for the things that matter in daily life."
          />
          <a
            href="#collections"
            className={`mt-8 inline-flex min-h-11 items-center text-sm text-burgundy underline decoration-sandstone underline-offset-4 hover:decoration-burgundy ${focusRing}`}
          >
            Explore the collection
          </a>
        </Container>
      </section>
    </>
  )
}
