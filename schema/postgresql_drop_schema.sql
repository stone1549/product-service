DROP INDEX product_created_at_idx;
DROP INDEX product_updated_at_idx;
DROP INDEX product_name_idx;
DROP INDEX product_price_idx;
DROP INDEX product_textsearch_idx;

DROP TRIGGER product_search_update_trg ON product;
DROP FUNCTION product_search_update_func();
DROP TRIGGER product_set_updated_at_trg ON product;
DROP FUNCTION set_updated_at();
DROP TABLE product;
