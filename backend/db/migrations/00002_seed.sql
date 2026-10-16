-- Generated from frontend/src/data/catalog.ts by scripts/export-catalog.ts.
-- Do not edit by hand. Regenerate, then commit the result:
--   node --experimental-strip-types scripts/export-catalog.ts
-- +goose Up

insert into collections (
    slug,
    name,
    description,
    image_url,
    image_alt,
    seo_title,
    meta_description,
    sort_order
) values
    (
        $mmseed$morning-ritual$mmseed$,
        $mmseed$Morning Ritual$mmseed$,
        $mmseed$Everyday pieces for the first coffee, a slower breakfast and the small routines that help the day begin well.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792078635/Home_living.png$mmseed$,
        $mmseed$A breakfast table set with a mug, a bowl of fruit, and a small vase.$mmseed$,
        $mmseed$Morning Ritual | morrow & muse$mmseed$,
        $mmseed$Discover considered pieces for coffee, breakfast and everyday table rituals with morrow & muse.$mmseed$,
        1
    ),
    (
        $mmseed$workday-edit$mmseed$,
        $mmseed$Workday Edit$mmseed$,
        $mmseed$Thoughtful stationery and desk essentials for collecting ideas, finding focus and bringing a little order to the working day.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792076358/Desk_stationary.png$mmseed$,
        $mmseed$A desk arranged with a notebook, pens, and a wooden organizer.$mmseed$,
        $mmseed$Workday Edit | morrow & muse$mmseed$,
        $mmseed$Explore understated journals, notebooks, writing tools and desk organisation for a more considered workspace.$mmseed$,
        2
    ),
    (
        $mmseed$everyday-carry$mmseed$,
        $mmseed$Everyday Carry$mmseed$,
        $mmseed$Practical companions for commutes, errands and the everyday essentials you take wherever the day leads.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792076237/Everyday_essentials.png$mmseed$,
        $mmseed$A tote bag, bottle, pouch, and small tray arranged together.$mmseed$,
        $mmseed$Everyday Carry | morrow & muse$mmseed$,
        $mmseed$Discover everyday bags, a reusable bottle, a travel pouch and small accessories selected for life on the move.$mmseed$,
        3
    ),
    (
        $mmseed$quiet-space$mmseed$,
        $mmseed$Quiet Space$mmseed$,
        $mmseed$Understated lighting, textiles and decorative accents for making your home feel more personal and considered.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792078654/Home_accessories.png$mmseed$,
        $mmseed$A living area with a table lamp, cushion, tray, and small planter.$mmseed$,
        $mmseed$Quiet Space | morrow & muse$mmseed$,
        $mmseed$Explore understated home accents, including a sculptural table lamp, linen cushion cover, decorative tray and ceramic planter.$mmseed$,
        4
    );

insert into products (
    id,
    slug,
    name,
    price_cents,
    collection_id,
    summary,
    description,
    card_image_url,
    card_image_alt,
    detail_image_url,
    detail_image_alt,
    sort_order
) values
    (
        $mmseed$mm-01$mmseed$,
        $mmseed$everyday-ceramic-mug$mmseed$,
        $mmseed$The Everyday Ceramic Mug$mmseed$,
        2800,
        (select id from collections where slug = $mmseed$morning-ritual$mmseed$),
        $mmseed$A considered companion for your first coffee, afternoon tea or any small pause worth taking.$mmseed$,
        $mmseed$Make a little more of the everyday. The Everyday Ceramic Mug brings a simple, understated presence to the moments that bookend your day, from a quiet first coffee to an unhurried cup of tea. Its clean silhouette sits comfortably within a thoughtfully arranged kitchen or on a well-loved desk.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079402/product_card_ceramic_coffee_mug.png$mmseed$,
        $mmseed$Ceramic coffee mug photographed in a warm, neutral setting.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079748/product_detail_image_ceramic_coffee_mug.png$mmseed$,
        $mmseed$Ceramic coffee mug photographed in a warm, neutral setting.$mmseed$,
        1
    ),
    (
        $mmseed$mm-02$mmseed$,
        $mmseed$stoneware-breakfast-bowl$mmseed$,
        $mmseed$The Stoneware Breakfast Bowl$mmseed$,
        3400,
        (select id from collections where slug = $mmseed$morning-ritual$mmseed$),
        $mmseed$A simple, versatile bowl for breakfast, fruit and everyday table moments.$mmseed$,
        $mmseed$From a slow breakfast to a handful of seasonal fruit, the Stoneware Breakfast Bowl gives everyday serving a calmer, more considered feel. Its understated design makes it an easy addition to a relaxed table, whether used on its own or alongside your favourite morning pieces.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079570/product_card_stoneware_breakfast_bowel.png$mmseed$,
        $mmseed$Stoneware breakfast bowl displayed in a softly lit neutral setting.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079727/product_detail_image_breakfast_bowl.png$mmseed$,
        $mmseed$Stoneware breakfast bowl displayed in a softly lit neutral setting.$mmseed$,
        2
    ),
    (
        $mmseed$mm-03$mmseed$,
        $mmseed$minimal-ceramic-vase$mmseed$,
        $mmseed$The Minimal Ceramic Vase$mmseed$,
        4200,
        (select id from collections where slug = $mmseed$morning-ritual$mmseed$),
        $mmseed$An understated decorative accent for a single stem or a quiet corner.$mmseed$,
        $mmseed$Sometimes a room needs only one thoughtful detail. The Minimal Ceramic Vase brings a simple sculptural note to a shelf, bedside table or dining surface. Style it with a favourite stem or let its silhouette stand on its own as part of a calm, collected space.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079435/product_card_minimal_ceramic_vase.png$mmseed$,
        $mmseed$Minimal ceramic vase styled as a decorative accent.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079824/product_detail_image_minimal_ceramic_vase.png$mmseed$,
        $mmseed$Minimal ceramic vase styled as a decorative accent.$mmseed$,
        3
    ),
    (
        $mmseed$mm-04$mmseed$,
        $mmseed$linen-covered-journal$mmseed$,
        $mmseed$The Linen-Covered Journal$mmseed$,
        3200,
        (select id from collections where slug = $mmseed$workday-edit$mmseed$),
        $mmseed$A considered place for thoughts, plans, lists and ideas worth returning to.$mmseed$,
        $mmseed$Give your thoughts somewhere to land. The Linen-Covered Journal is made for morning pages, meeting notes, passing ideas and the plans you want to keep close. Its tactile, understated appearance brings a quieter feeling to the everyday habit of putting pen to paper.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079412/product_card_linen_covered_journal.png$mmseed$,
        $mmseed$Linen-covered journal arranged on a neutral desk.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079765/product_detail_image_linen_covered_journal.png$mmseed$,
        $mmseed$Linen-covered journal arranged on a neutral desk.$mmseed$,
        4
    ),
    (
        $mmseed$mm-05$mmseed$,
        $mmseed$wooden-desk-organizer$mmseed$,
        $mmseed$The Wooden Desk Organizer$mmseed$,
        4800,
        (select id from collections where slug = $mmseed$workday-edit$mmseed$),
        $mmseed$A warm, practical home for the small essentials that make up your working day.$mmseed$,
        $mmseed$A clearer desk can make room for clearer thinking. The Wooden Desk Organizer gives everyday workspace essentials a more intentional place, helping pens and frequently used items feel considered rather than scattered. Keep it within reach on a home-office desk or study surface.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079587/product_card_wooden_desk_organizer.png$mmseed$,
        $mmseed$Wooden desk organizer arranged with everyday stationery.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079928/product_detail_image_wooden_desk_organizer.png$mmseed$,
        $mmseed$Wooden desk organizer arranged with everyday stationery.$mmseed$,
        5
    ),
    (
        $mmseed$mm-06$mmseed$,
        $mmseed$black-pen-set$mmseed$,
        $mmseed$The Black Pen Set$mmseed$,
        2400,
        (select id from collections where slug = $mmseed$workday-edit$mmseed$),
        $mmseed$Understated writing essentials for notes, lists and ideas in progress.$mmseed$,
        $mmseed$Good writing tools belong where ideas happen. The Black Pen Set brings a clean, coordinated look to daily notes, work lists and pages of unfinished thoughts. Keep it beside your journal or notebook so the next idea is always ready to be written down.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079381/product_card_black_pen_set.png$mmseed$,
        $mmseed$Black pen set presented with coordinated stationery.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079705/product_detail_image_black_pen_set.png$mmseed$,
        $mmseed$Black pen set presented with coordinated stationery.$mmseed$,
        6
    ),
    (
        $mmseed$mm-07$mmseed$,
        $mmseed$neutral-notebook$mmseed$,
        $mmseed$The Neutral Notebook$mmseed$,
        2200,
        (select id from collections where slug = $mmseed$workday-edit$mmseed$),
        $mmseed$An everyday notebook for working thoughts, quick lists and new ideas.$mmseed$,
        $mmseed$Keep the moving parts of your day in one simple place. The Neutral Notebook is intended for the lists, notes, sketches and reminders that help ideas take shape. Its quiet presentation fits neatly into a home workspace, a study corner or an everyday bag.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079523/product_card_neutral_notebook.png$mmseed$,
        $mmseed$Neutral notebook styled on a tidy workspace.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079877/product_detail_image_neutral_notebook.png$mmseed$,
        $mmseed$Neutral notebook styled on a tidy workspace.$mmseed$,
        7
    ),
    (
        $mmseed$mm-08$mmseed$,
        $mmseed$reusable-insulated-bottle$mmseed$,
        $mmseed$The Reusable Insulated Bottle$mmseed$,
        3600,
        (select id from collections where slug = $mmseed$everyday-carry$mmseed$),
        $mmseed$A reusable bottle for workdays, errands and time spent out and about.$mmseed$,
        $mmseed$A useful everyday essential should fit into the rhythm of your day. The Reusable Insulated Bottle is designed to accompany your commute, work bag, errands and time outdoors. Its clean, understated appearance makes it a natural fit with the other pieces in the Everyday Carry collection.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079534/product_card_reusable_insulated_bottle.png$mmseed$,
        $mmseed$Reusable insulated bottle photographed in a neutral setting.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079891/product_detail_image_reusable_insulated_bottle.png$mmseed$,
        $mmseed$Reusable insulated bottle photographed in a neutral setting.$mmseed$,
        8
    ),
    (
        $mmseed$mm-09$mmseed$,
        $mmseed$canvas-tote-bag$mmseed$,
        $mmseed$The Canvas Tote Bag$mmseed$,
        3000,
        (select id from collections where slug = $mmseed$everyday-carry$mmseed$),
        $mmseed$A straightforward carry-all for the everyday things you want close at hand.$mmseed$,
        $mmseed$From a quick trip out to a day of errands, the Canvas Tote Bag offers an uncomplicated way to carry everyday essentials. Its relaxed, versatile look works with a range of daily routines, making it a natural companion for the pieces you reach for most often.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079390/product_card_canvas_tote_bag.png$mmseed$,
        $mmseed$Canvas tote bag styled as an everyday carry accessory.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079736/product_detail_image_canvas_tote_bag.png$mmseed$,
        $mmseed$Canvas tote bag styled as an everyday carry accessory.$mmseed$,
        9
    ),
    (
        $mmseed$mm-10$mmseed$,
        $mmseed$minimal-travel-pouch$mmseed$,
        $mmseed$The Minimal Travel Pouch$mmseed$,
        2600,
        (select id from collections where slug = $mmseed$everyday-carry$mmseed$),
        $mmseed$A considered way to keep small essentials together when you are on the move.$mmseed$,
        $mmseed$The little things are easier to find when they have a place of their own. The Minimal Travel Pouch is suited to organising the small essentials you take between home, work and time away. Keep it inside your everyday bag or pack it when your routine takes you somewhere new.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079484/product_card_minimal_travel_pouch.png$mmseed$,
        $mmseed$Minimal travel pouch displayed with everyday carry essentials.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079849/product_detail_image_minmal_travel_pouch.png$mmseed$,
        $mmseed$Minimal travel pouch displayed with everyday carry essentials.$mmseed$,
        10
    ),
    (
        $mmseed$mm-11$mmseed$,
        $mmseed$everyday-accessory-tray$mmseed$,
        $mmseed$The Everyday Accessory Tray$mmseed$,
        2400,
        (select id from collections where slug = $mmseed$everyday-carry$mmseed$),
        $mmseed$A small landing place for the everyday pieces you do not want to misplace.$mmseed$,
        $mmseed$Give your daily essentials a place to land. The Everyday Accessory Tray helps gather the small items you reach for when leaving home or settling in for the evening. Place it on an entry table, bedside surface or desk to bring a little more order to the everyday.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079305/product_card_accessory_tray.png$mmseed$,
        $mmseed$Accessory tray arranged on a tidy surface.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079691/product_detail_image_accessory_tray.png$mmseed$,
        $mmseed$Accessory tray arranged on a tidy surface.$mmseed$,
        11
    ),
    (
        $mmseed$mm-12$mmseed$,
        $mmseed$sculptural-table-lamp$mmseed$,
        $mmseed$The Sculptural Table Lamp$mmseed$,
        14800,
        (select id from collections where slug = $mmseed$quiet-space$mmseed$),
        $mmseed$A sculptural lighting accent for evenings, reading corners and slower moments.$mmseed$,
        $mmseed$Lighting helps set the feeling of a room. The Sculptural Table Lamp brings a considered focal point to a bedside table, reading corner or side table. Its expressive form can sit comfortably among other carefully chosen objects or hold its own as a quiet statement piece.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079542/product_card_sculptual_table_lamp.png$mmseed$,
        $mmseed$Sculptural table lamp displayed in a warm, quiet interior.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079904/product_detail_image_sculptural_table_lamp.png$mmseed$,
        $mmseed$Sculptural table lamp displayed in a warm, quiet interior.$mmseed$,
        12
    ),
    (
        $mmseed$mm-13$mmseed$,
        $mmseed$linen-cushion-cover$mmseed$,
        $mmseed$The Linen Cushion Cover$mmseed$,
        4600,
        (select id from collections where slug = $mmseed$quiet-space$mmseed$),
        $mmseed$A simple textile accent for adding softness and a more considered finish to a room.$mmseed$,
        $mmseed$Small changes can shift the feeling of a space. The Linen Cushion Cover brings a relaxed textile layer to a sofa, armchair or bed, complementing a palette of quiet neutrals and natural-looking textures. Use it as a subtle finishing touch or combine it with other understated furnishings.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079423/product_card_linen_cushion_cover.png$mmseed$,
        $mmseed$Linen cushion cover styled on neutral furnishings.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079811/product_detail_image_linen_sushion_cover.png$mmseed$,
        $mmseed$Linen cushion cover styled on neutral furnishings.$mmseed$,
        13
    ),
    (
        $mmseed$mm-14$mmseed$,
        $mmseed$minimal-decorative-tray$mmseed$,
        $mmseed$The Minimal Decorative Tray$mmseed$,
        3800,
        (select id from collections where slug = $mmseed$quiet-space$mmseed$),
        $mmseed$A quiet styling piece that brings smaller objects together into one considered arrangement.$mmseed$,
        $mmseed$A few favourite objects can feel more intentional when arranged together. The Minimal Decorative Tray provides a simple base for styling a coffee table, console or shelf. Use it to bring visual order to a small grouping or let it contribute its own understated shape to the room.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079447/product_card_minimal_decorative_tray.png$mmseed$,
        $mmseed$Minimal decorative tray styled with carefully arranged home objects.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079837/product_detail_image_minimal_decorative_tray.png$mmseed$,
        $mmseed$Minimal decorative tray styled with carefully arranged home objects.$mmseed$,
        14
    ),
    (
        $mmseed$mm-15$mmseed$,
        $mmseed$small-ceramic-planter$mmseed$,
        $mmseed$The Small Ceramic Planter$mmseed$,
        3200,
        (select id from collections where slug = $mmseed$quiet-space$mmseed$),
        $mmseed$A compact ceramic accent for bringing a little greenery into a considered space.$mmseed$,
        $mmseed$A small touch of greenery can make a space feel more personal. The Small Ceramic Planter is designed to sit naturally on a desk, shelf or windowsill, adding a simple decorative element to everyday surroundings. Pair it with other quiet accents or allow it to bring a little life to a smaller corner.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079555/product_card_small_ceramic_plant.png$mmseed$,
        $mmseed$Small ceramic planter styled as a home accent.$mmseed$,
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079915/product_detail_image_small_ceramic_plant.png$mmseed$,
        $mmseed$Small ceramic planter styled as a home accent.$mmseed$,
        15
    );

insert into product_highlights (product_id, position, body) values
    ($mmseed$mm-01$mmseed$, 1, $mmseed$Made for everyday coffee and tea rituals.$mmseed$),
    ($mmseed$mm-01$mmseed$, 2, $mmseed$An understated shape that complements a considered table.$mmseed$),
    ($mmseed$mm-01$mmseed$, 3, $mmseed$Easy to pair with the rest of the Morning Ritual collection.$mmseed$),
    ($mmseed$mm-02$mmseed$, 1, $mmseed$Suited to breakfast and everyday serving.$mmseed$),
    ($mmseed$mm-02$mmseed$, 2, $mmseed$A restrained look that works across different table settings.$mmseed$),
    ($mmseed$mm-02$mmseed$, 3, $mmseed$Designed to complement the Morning Ritual collection.$mmseed$),
    ($mmseed$mm-03$mmseed$, 1, $mmseed$Works as a standalone decorative accent.$mmseed$),
    ($mmseed$mm-03$mmseed$, 2, $mmseed$Pairs naturally with neutral interiors.$mmseed$),
    ($mmseed$mm-03$mmseed$, 3, $mmseed$Easy to style on shelves, tables and smaller surfaces.$mmseed$),
    ($mmseed$mm-04$mmseed$, 1, $mmseed$For journaling, lists, notes and personal planning.$mmseed$),
    ($mmseed$mm-04$mmseed$, 2, $mmseed$An understated cover that sits naturally on a desk or bedside table.$mmseed$),
    ($mmseed$mm-04$mmseed$, 3, $mmseed$A thoughtful addition to a daily writing ritual.$mmseed$),
    ($mmseed$mm-05$mmseed$, 1, $mmseed$Helps bring order to everyday desk essentials.$mmseed$),
    ($mmseed$mm-05$mmseed$, 2, $mmseed$Adds a warm visual accent to a workspace.$mmseed$),
    ($mmseed$mm-05$mmseed$, 3, $mmseed$Designed to sit naturally alongside notebooks and writing tools.$mmseed$),
    ($mmseed$mm-06$mmseed$, 1, $mmseed$For everyday writing and note-taking.$mmseed$),
    ($mmseed$mm-06$mmseed$, 2, $mmseed$A restrained look that pairs with the Workday Edit collection.$mmseed$),
    ($mmseed$mm-06$mmseed$, 3, $mmseed$A useful desk companion or considered gift.$mmseed$),
    ($mmseed$mm-07$mmseed$, 1, $mmseed$Suited to lists, notes and everyday planning.$mmseed$),
    ($mmseed$mm-07$mmseed$, 2, $mmseed$An understated addition to a desk or bag.$mmseed$),
    ($mmseed$mm-07$mmseed$, 3, $mmseed$Easy to pair with the Black Pen Set and Linen-Covered Journal.$mmseed$),
    ($mmseed$mm-08$mmseed$, 1, $mmseed$Designed for repeated everyday use.$mmseed$),
    ($mmseed$mm-08$mmseed$, 2, $mmseed$Suited to workdays and time on the move.$mmseed$),
    ($mmseed$mm-08$mmseed$, 3, $mmseed$A coordinated addition to an everyday carry setup.$mmseed$),
    ($mmseed$mm-09$mmseed$, 1, $mmseed$Suited to errands and everyday outings.$mmseed$),
    ($mmseed$mm-09$mmseed$, 2, $mmseed$An easy addition to a practical daily carry.$mmseed$),
    ($mmseed$mm-09$mmseed$, 3, $mmseed$Coordinates with the Everyday Carry collection.$mmseed$),
    ($mmseed$mm-10$mmseed$, 1, $mmseed$Helps organise small everyday essentials.$mmseed$),
    ($mmseed$mm-10$mmseed$, 2, $mmseed$Suited to bags, commutes and travel.$mmseed$),
    ($mmseed$mm-10$mmseed$, 3, $mmseed$A simple companion to the Canvas Tote Bag.$mmseed$),
    ($mmseed$mm-11$mmseed$, 1, $mmseed$For organising small daily essentials.$mmseed$),
    ($mmseed$mm-11$mmseed$, 2, $mmseed$Suited to entry tables, desks and bedside surfaces.$mmseed$),
    ($mmseed$mm-11$mmseed$, 3, $mmseed$Pairs naturally with the Quiet Space collection.$mmseed$),
    ($mmseed$mm-12$mmseed$, 1, $mmseed$A decorative lighting accent for considered interiors.$mmseed$),
    ($mmseed$mm-12$mmseed$, 2, $mmseed$Suited to bedside tables and reading corners.$mmseed$),
    ($mmseed$mm-12$mmseed$, 3, $mmseed$Pairs with neutral furnishings and understated decor.$mmseed$),
    ($mmseed$mm-13$mmseed$, 1, $mmseed$A textile accent for sofas, chairs and beds.$mmseed$),
    ($mmseed$mm-13$mmseed$, 2, $mmseed$Suited to calm, neutral interiors.$mmseed$),
    ($mmseed$mm-13$mmseed$, 3, $mmseed$Easy to incorporate into a considered home palette.$mmseed$),
    ($mmseed$mm-14$mmseed$, 1, $mmseed$For styling and arranging decorative objects.$mmseed$),
    ($mmseed$mm-14$mmseed$, 2, $mmseed$Suited to coffee tables, shelves and consoles.$mmseed$),
    ($mmseed$mm-14$mmseed$, 3, $mmseed$Complements the Quiet Space collection.$mmseed$),
    ($mmseed$mm-15$mmseed$, 1, $mmseed$Suited to small decorative planting arrangements.$mmseed$),
    ($mmseed$mm-15$mmseed$, 2, $mmseed$An understated accent for desks, shelves and windowsills.$mmseed$),
    ($mmseed$mm-15$mmseed$, 3, $mmseed$Works naturally with the Quiet Space collection.$mmseed$);

insert into featured_products (product_id, position) values
    ((select id from products where slug = $mmseed$everyday-ceramic-mug$mmseed$), 1),
    ((select id from products where slug = $mmseed$linen-covered-journal$mmseed$), 2),
    ((select id from products where slug = $mmseed$canvas-tote-bag$mmseed$), 3),
    ((select id from products where slug = $mmseed$sculptural-table-lamp$mmseed$), 4);

insert into catalog_assets (key, url, alt) values
    ('logo', $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079260/Logo.png$mmseed$, null),
    (
        'hero',
        $mmseed$https://res.cloudinary.com/sbll6gvk/image/upload/v1792079279/Homepage_hero.png$mmseed$,
        $mmseed$A quiet bedroom corner with a side table, lamp, and layered neutral bedding.$mmseed$
    );

-- +goose Down

delete from featured_products;
delete from product_highlights;
delete from products;
delete from collections;
delete from catalog_assets;
