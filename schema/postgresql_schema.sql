CREATE TABLE product (
  id text PRIMARY KEY,
  name text NOT NULL,
  display_image text,
  thumbnail text,
  price numeric(15,6) NOT NULL,
  description text,
  short_description text,
  qty_in_stock int NOT NULL DEFAULT 0,
  textsearchable_index_col tsvector NOT NULL
);

CREATE INDEX product_name_idx ON product (name);
CREATE INDEX product_price_idx ON product (price);
CREATE INDEX product_textsearch_idx ON product USING GIN (textsearchable_index_col);

CREATE FUNCTION product_search_update_func() RETURNS trigger AS $$
begin
  new.textsearchable_index_col :=
  setweight(to_tsvector('pg_catalog.english', coalesce(new.name,'')), 'A') ||
  setweight(to_tsvector('pg_catalog.english', coalesce(new.id,'')), 'B') ||
  setweight(to_tsvector('pg_catalog.english', coalesce(new.short_description,'')), 'C') ||
  setweight(to_tsvector('pg_catalog.english', coalesce(new.description,'')), 'D');
  return new;
end
$$ LANGUAGE plpgsql;

CREATE TRIGGER product_search_update_trg BEFORE INSERT OR UPDATE
  ON product FOR EACH ROW EXECUTE PROCEDURE product_search_update_func();
