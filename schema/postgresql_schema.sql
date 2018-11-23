CREATE TABLE product (
 id text PRIMARY KEY,
 name text NOT NULL,
 display_image text,
 thumbnail text,
 price numeric(15,6) NOT NULL,
 description text,
 short_description text,
 qty_in_stock int NOT NULL DEFAULT 0
);
