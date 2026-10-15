import { Link } from 'react-router'

export function NotFoundPage() {
  return (
    <section className="mx-auto w-full max-w-5xl px-5 py-16 sm:px-8 sm:py-24">
      <h1 className="font-display text-4xl leading-tight text-burgundy sm:text-5xl">
        Page not found
      </h1>
      <p className="mt-4 max-w-md text-lg text-taupe">
        That address is not part of the boutique.
      </p>
      <Link
        to="/"
        className="mt-8 inline-flex min-h-11 items-center text-burgundy underline decoration-sandstone underline-offset-4 hover:decoration-burgundy focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-burgundy"
      >
        Back to home
      </Link>
    </section>
  )
}
