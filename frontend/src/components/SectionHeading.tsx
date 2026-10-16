export function SectionHeading({
  eyebrow,
  title,
  lede,
}: {
  eyebrow?: string
  title: string
  lede?: string
}) {
  return (
    <div className="max-w-2xl">
      {eyebrow ? (
        <p className="text-sm tracking-[0.16em] text-taupe uppercase">
          {eyebrow}
        </p>
      ) : null}
      <h2 className="mt-3 font-display text-4xl leading-tight text-burgundy sm:text-5xl">
        {title}
      </h2>
      {lede ? (
        <p className="mt-4 text-base leading-relaxed text-ink sm:text-lg">
          {lede}
        </p>
      ) : null}
    </div>
  )
}
