BEGIN;

DROP TABLE payment;
DROP TABLE coupon;
DROP TABLE order_item;
DROP TABLE "order";
DROP TABLE discount;
DROP TABLE cart_item;
DROP TABLE cart;
DROP TABLE product;
DROP TABLE session;
DROP TABLE "user";
DROP FUNCTION sanitize_timestamps;
DROP FUNCTION trigger_sanitize_id;
DROP EXTENSION "uuid-ossp";
DROP EXTENSION pgcrypto;

COMMIT;