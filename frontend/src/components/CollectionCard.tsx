import { Link } from 'react-router'
import type { Collection } from '../data/catalog.ts'
import { focusRing } from '../lib/focusRing.ts'
import { CatalogPlaceholder } from './CatalogPlaceholder.tsx'

export function CollectionCard({ collection }: { collection: Collection }) {
  return (
    <article className="group relative min-w-0">
      <CatalogPlaceholder
        image={collection.image}
        width={1200}
        height={800}
        className="aspect-[3/2] w-full object-cover"
      />
      <h3 className="mt-4 font-display text-3xl text-burgundy">
        <Link
          to={`/collections/${collection.slug}`}
          className={`after:absolute after:inset-0 hover:underline ${focusRing}`}
        >
          {collection.name}
        </Link>
      </h3>
      <p className="mt-2 text-sm leading-relaxed text-ink">
        {collection.description}
      </p>
    </article>
  )
}
