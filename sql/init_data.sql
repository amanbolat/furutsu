-- Products
INSERT INTO public.product (id, name, price, description)
VALUES ('0ca41a0a-2e65-40e9-a34e-7133b4adb9b4', 'apple', 10, 'Green apples from Almaty');
INSERT INTO public.product (id, name, price, description)
VALUES ('0645a967-badb-40d2-89fb-81c7b745899c', 'pear', 15, NULL);
INSERT INTO public.product (id, name, price, description)
VALUES ('748cb518-1bd4-4f45-98ba-be016b39827e', 'banana', 3, NULL);
INSERT INTO public.product (id, name, price, description)
VALUES ('e7f83943-7044-4406-9d32-9b229725f6d0', 'orange', 5, NULL);


-- Discounts
INSERT INTO discount (name, rule, percent)
VALUES ('7 apples 10%', '{
  "0ca41a0a-2e65-40e9-a34e-7133b4adb9b4": 7
}', 10);
INSERT INTO discount (name, rule, percent)
VALUES ('set of 4 pears, 2 bananas 30%', '{
  "0645a967-badb-40d2-89fb-81c7b745899c": 4,
  "748cb518-1bd4-4f45-98ba-be016b39827e": 2
}', 10);

-- Coupons
INSERT INTO public.coupon (code, name, expire_at, rule, discount_percent)
VALUES ('orange123', 'orange 10%', '2021-11-06 15:30:31.542000', '{
  "e7f83943-7044-4406-9d32-9b229725f6d0": 0
}', 10);


-- User and its cart
INSERT INTO public.user (username, full_name, password)
VALUES ('aman', 'Amanbolat', 'pass');
WITH u AS (
    SELECT id
    FROM "user"
    WHERE username = 'aman'
)
INSERT
INTO cart (user_id)
SELECT id
FROM u;

-- User cart items
WITH tmp AS (
    SELECT cart.id as id
    FROM cart
    JOIN "user" u ON u.id = cart.user_id
    WHERE u.username = 'aman'
)
INSERT INTO cart_item (cart_id, product_id, amount)
SELECT id, 'e7f83943-7044-4406-9d32-9b229725f6d0', 20
FROM tmp;

