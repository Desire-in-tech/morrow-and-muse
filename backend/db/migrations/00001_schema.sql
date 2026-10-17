-- Catalogue schema. One collection has many products, matching catalog.ts.
-- There is no product_collections join: every product has one collection.
-- +goose Up

create table collections (
    id bigint generated always as identity primary key,
    slug text not null,
    name text not null,
    description text not null,
    image_url text not null,
    image_alt text not null,
    seo_title text not null,
    meta_description text not null,
    sort_order integer not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint collections_slug_unique unique (slug),
    constraint collections_slug_format check (slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$'),
    constraint collections_name_present check (length(name) > 0),
    constraint collections_description_present check (length(description) > 0),
    constraint collections_image_url_https check (image_url ~ '^https://'),
    constraint collections_image_alt_present check (length(image_alt) > 0),
    constraint collections_seo_title_present check (length(seo_title) > 0),
    constraint collections_meta_description_present check (length(meta_description) > 0),
    constraint collections_sort_order_positive check (sort_order > 0),
    constraint collections_sort_order_unique unique (sort_order)
);

create table products (
    id text primary key,
    slug text not null,
    name text not null,
    price_cents integer not null,
    collection_id bigint not null,
    summary text not null,
    description text not null,
    card_image_url text not null,
    card_image_alt text not null,
    detail_image_url text not null,
    detail_image_alt text not null,
    sort_order integer not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint products_id_format check (id ~ '^mm-[0-9]{2}$'),
    constraint products_slug_unique unique (slug),
    constraint products_slug_format check (slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$'),
    constraint products_name_present check (length(name) > 0),
    constraint products_price_cents_nonnegative check (price_cents >= 0),
    constraint products_collection_fk foreign key (collection_id) references collections (id),
    constraint products_summary_present check (length(summary) > 0),
    constraint products_description_present check (length(description) > 0),
    constraint products_card_image_url_https check (card_image_url ~ '^https://'),
    constraint products_card_image_alt_present check (length(card_image_alt) > 0),
    constraint products_detail_image_url_https check (detail_image_url ~ '^https://'),
    constraint products_detail_image_alt_present check (length(detail_image_alt) > 0),
    constraint products_sort_order_positive check (sort_order > 0),
    constraint products_sort_order_unique unique (sort_order)
);

create index products_collection_id_idx on products (collection_id);

-- catalog.ts types highlights as a 3-tuple. Position is 1, 2, or 3.
create table product_highlights (
    product_id text not null,
    position smallint not null,
    body text not null,
    constraint product_highlights_pk primary key (product_id, position),
    constraint product_highlights_product_fk foreign key (product_id) references products (id) on delete cascade,
    constraint product_highlights_position check (position between 1 and 3),
    constraint product_highlights_body_present check (length(body) > 0)
);

create index product_highlights_product_id_idx on product_highlights (product_id);

-- Ordered subset from featuredProductSlugs. Not every product is featured.
create table featured_products (
    product_id text primary key,
    position smallint not null,
    constraint featured_products_product_fk foreign key (product_id) references products (id) on delete cascade,
    constraint featured_products_position_positive check (position > 0),
    constraint featured_products_position_unique unique (position)
);

-- logoUrl and heroImage are catalogue assets, not collection or product rows.
-- logo has a URL only. hero has a URL and alt text.
create table catalog_assets (
    key text primary key,
    url text not null,
    alt text,
    constraint catalog_assets_key_format check (key ~ '^[a-z0-9]+(-[a-z0-9]+)*$'),
    constraint catalog_assets_url_https check (url ~ '^https://'),
    constraint catalog_assets_alt_present check (alt is null or length(alt) > 0)
);

-- +goose StatementBegin
create function set_updated_at()
returns trigger
language plpgsql
as $fn$
begin
    new.updated_at = clock_timestamp();
    return new;
end;
$fn$;
-- +goose StatementEnd

create trigger collections_set_updated_at
before update on collections
for each row
execute function set_updated_at();

create trigger products_set_updated_at
before update on products
for each row
execute function set_updated_at();

-- Deferred so a product and its three highlights can be inserted in one transaction.
-- Separate functions because NEW's columns differ on the two tables.
-- +goose StatementBegin
create function products_have_three_highlights()
returns trigger
language plpgsql
as $fn$
declare
    highlight_count integer;
begin
    if not exists (select 1 from products where id = new.id) then
        return null;
    end if;
    select count(*) into highlight_count
    from product_highlights
    where product_id = new.id;
    if highlight_count <> 3 then
        raise exception 'product % must have exactly 3 highlights', new.id;
    end if;
    return null;
end;
$fn$;
-- +goose StatementEnd

-- +goose StatementBegin
create function product_highlights_stay_three()
returns trigger
language plpgsql
as $fn$
declare
    pid text;
    highlight_count integer;
begin
    pid = coalesce(new.product_id, old.product_id);
    if not exists (select 1 from products where id = pid) then
        return null;
    end if;
    select count(*) into highlight_count
    from product_highlights
    where product_id = pid;
    if highlight_count <> 3 then
        raise exception 'product % must have exactly 3 highlights', pid;
    end if;
    return null;
end;
$fn$;
-- +goose StatementEnd

create constraint trigger products_require_three_highlights
after insert or update on products
deferrable initially deferred
for each row
execute function products_have_three_highlights();

create constraint trigger product_highlights_require_three
after insert or update or delete on product_highlights
deferrable initially deferred
for each row
execute function product_highlights_stay_three();

-- +goose Down

drop trigger if exists product_highlights_require_three on product_highlights;
drop trigger if exists products_require_three_highlights on products;
drop function if exists product_highlights_stay_three();
drop function if exists products_have_three_highlights();

drop trigger if exists products_set_updated_at on products;
drop trigger if exists collections_set_updated_at on collections;
drop function if exists set_updated_at();

drop table if exists featured_products;
drop table if exists product_highlights;
drop table if exists products;
drop table if exists collections;
drop table if exists catalog_assets;
