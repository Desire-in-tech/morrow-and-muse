import { describe, expect, it } from 'vitest'
import {
  collections,
  featuredProductSlugs,
  heroImage,
  logoUrl,
  products,
} from './data/catalog.ts'

const expectedUrls = [
  [
    'logo',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079260/Logo.png',
  ],
  [
    'hero',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079279/Homepage_hero.png',
  ],
  [
    'morning-ritual',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792078635/Home_living.png',
  ],
  [
    'workday-edit',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792076358/Desk_stationary.png',
  ],
  [
    'everyday-carry',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792076237/Everyday_essentials.png',
  ],
  [
    'quiet-space',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792078654/Home_accessories.png',
  ],
  [
    'everyday-ceramic-mug-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079402/product_card_ceramic_coffee_mug.png',
  ],
  [
    'everyday-ceramic-mug-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079748/product_detail_image_ceramic_coffee_mug.png',
  ],
  [
    'stoneware-breakfast-bowl-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079570/product_card_stoneware_breakfast_bowel.png',
  ],
  [
    'stoneware-breakfast-bowl-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079727/product_detail_image_breakfast_bowl.png',
  ],
  [
    'minimal-ceramic-vase-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079435/product_card_minimal_ceramic_vase.png',
  ],
  [
    'minimal-ceramic-vase-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079824/product_detail_image_minimal_ceramic_vase.png',
  ],
  [
    'linen-covered-journal-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079412/product_card_linen_covered_journal.png',
  ],
  [
    'linen-covered-journal-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079765/product_detail_image_linen_covered_journal.png',
  ],
  [
    'wooden-desk-organizer-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079587/product_card_wooden_desk_organizer.png',
  ],
  [
    'wooden-desk-organizer-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079928/product_detail_image_wooden_desk_organizer.png',
  ],
  [
    'black-pen-set-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079381/product_card_black_pen_set.png',
  ],
  [
    'black-pen-set-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079705/product_detail_image_black_pen_set.png',
  ],
  [
    'neutral-notebook-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079523/product_card_neutral_notebook.png',
  ],
  [
    'neutral-notebook-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079877/product_detail_image_neutral_notebook.png',
  ],
  [
    'reusable-insulated-bottle-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079534/product_card_reusable_insulated_bottle.png',
  ],
  [
    'reusable-insulated-bottle-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079891/product_detail_image_reusable_insulated_bottle.png',
  ],
  [
    'canvas-tote-bag-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079390/product_card_canvas_tote_bag.png',
  ],
  [
    'canvas-tote-bag-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079736/product_detail_image_canvas_tote_bag.png',
  ],
  [
    'minimal-travel-pouch-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079484/product_card_minimal_travel_pouch.png',
  ],
  [
    'minimal-travel-pouch-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079849/product_detail_image_minmal_travel_pouch.png',
  ],
  [
    'everyday-accessory-tray-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079305/product_card_accessory_tray.png',
  ],
  [
    'everyday-accessory-tray-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079691/product_detail_image_accessory_tray.png',
  ],
  [
    'sculptural-table-lamp-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079542/product_card_sculptual_table_lamp.png',
  ],
  [
    'sculptural-table-lamp-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079904/product_detail_image_sculptural_table_lamp.png',
  ],
  [
    'linen-cushion-cover-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079423/product_card_linen_cushion_cover.png',
  ],
  [
    'linen-cushion-cover-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079811/product_detail_image_linen_sushion_cover.png',
  ],
  [
    'minimal-decorative-tray-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079447/product_card_minimal_decorative_tray.png',
  ],
  [
    'minimal-decorative-tray-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079837/product_detail_image_minimal_decorative_tray.png',
  ],
  [
    'small-ceramic-planter-card',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079555/product_card_small_ceramic_plant.png',
  ],
  [
    'small-ceramic-planter-detail',
    'https://res.cloudinary.com/sbll6gvk/image/upload/v1792079915/product_detail_image_small_ceramic_plant.png',
  ],
] as const

describe('catalogue', () => {
  it('lists 15 unique products in four collections', () => {
    expect(products).toHaveLength(15)
    expect(new Set(products.map((product) => product.id)).size).toBe(15)
    expect(new Set(products.map((product) => product.slug)).size).toBe(15)
    expect(collections.map((collection) => collection.slug)).toEqual([
      'morning-ritual',
      'workday-edit',
      'everyday-carry',
      'quiet-space',
    ])
    expect(
      products
        .filter((product) => product.collectionSlug === 'morning-ritual')
        .map((product) => product.slug),
    ).toEqual([
      'everyday-ceramic-mug',
      'stoneware-breakfast-bowl',
      'minimal-ceramic-vase',
    ])
    expect(
      products.filter((product) => product.collectionSlug === 'workday-edit'),
    ).toHaveLength(4)
    expect(
      products.filter((product) => product.collectionSlug === 'everyday-carry'),
    ).toHaveLength(4)
    expect(
      products.filter((product) => product.collectionSlug === 'quiet-space'),
    ).toHaveLength(4)
    expect(featuredProductSlugs).toEqual([
      'everyday-ceramic-mug',
      'linen-covered-journal',
      'canvas-tote-bag',
      'sculptural-table-lamp',
    ])

    for (const product of products) {
      expect(product.priceCents).toBeGreaterThan(0)
      expect(product.summary.length).toBeGreaterThan(0)
      expect(product.description.length).toBeGreaterThan(0)
      expect(product.highlights).toHaveLength(3)
      expect(product.cardImage.src).toMatch(/^https:\/\/res\.cloudinary\.com\//)
      expect(product.detailImage.src).toMatch(
        /^https:\/\/res\.cloudinary\.com\//,
      )
      expect(
        collections.some(
          (collection) => collection.slug === product.collectionSlug,
        ),
      ).toBe(true)
    }
  })

  it('matches every Cloudinary URL to the approved list', () => {
    const actual = new Map<string, string>([
      ['logo', logoUrl],
      ['hero', heroImage.src ?? ''],
      ...collections.map(
        (collection) => [collection.slug, collection.image.src ?? ''] as const,
      ),
      ...products.flatMap((product) => [
        [`${product.slug}-card`, product.cardImage.src ?? ''] as const,
        [`${product.slug}-detail`, product.detailImage.src ?? ''] as const,
      ]),
    ])

    expect([...actual.keys()].sort()).toEqual(
      expectedUrls.map(([key]) => key).sort(),
    )
    for (const [key, url] of expectedUrls) {
      expect(actual.get(key)).toBe(url)
    }
  })
})
