DROP INDEX product_name_idx;
DROP INDEX product_price_idx;
DROP INDEX product_textsearch_idx;

DROP TRIGGER product_search_update_trg ON product;
DROP FUNCTION product_search_update_func();

DROP TABLE product;
