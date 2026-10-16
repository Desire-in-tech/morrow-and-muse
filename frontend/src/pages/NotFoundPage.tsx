import { Link } from 'react-router'
import { Container } from '../components/Container.tsx'
import { focusRing } from '../lib/focusRing.ts'
import { useDocumentMeta } from '../lib/useDocumentMeta.ts'

export function NotFoundPage() {
  useDocumentMeta(
    'Page Not Found | morrow & muse',
    'Explore the thoughtfully chosen everyday essentials from morrow & muse.',
  )

  return (
    <section className="py-16 sm:py-24">
      <Container>
        <h1 className="font-display text-4xl leading-tight text-burgundy sm:text-5xl">
          This page is taking a quieter path.
        </h1>
        <p className="mt-4 max-w-md text-lg text-ink">
          We couldn&apos;t find the page you were looking for. The pieces
          you&apos;re after may be elsewhere in the collection.
        </p>
        <div className="mt-8 flex flex-wrap gap-4">
          <Link
            to="/shop"
            className={`inline-flex min-h-11 items-center bg-burgundy px-6 text-sm text-ivory hover:bg-ink ${focusRing}`}
          >
            Explore the collection
          </Link>
          <Link
            to="/"
            className={`inline-flex min-h-11 items-center text-burgundy underline decoration-sandstone underline-offset-4 hover:decoration-burgundy ${focusRing}`}
          >
            Visit the homepage
          </Link>
        </div>
      </Container>
    </section>
  )
}
