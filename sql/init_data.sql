-- Products
INSERT INTO public.product (id, name, price, description)
VALUES ('0ca41a0a-2e65-40e9-a34e-7133b4adb9b4', 'apple', 200, 'An apple is an edible fruit produced by an apple tree (Malus domestica). Apple trees are cultivated worldwide and are the most widely grown species in the genus Malus. The tree originated in Central Asia, where its wild ancestor, Malus sieversii, is still found today.');
INSERT INTO public.product (id, name, price, description)
VALUES ('0645a967-badb-40d2-89fb-81c7b745899c', 'pear', 350, 'The pear tree and shrub are a species of genus Pyrus, in the family Rosaceae, bearing the pomaceous fruit of the same name. Several species of pears are valued for their edible fruit and juices, while others are cultivated as trees.');
INSERT INTO public.product (id, name, price, description)
VALUES ('748cb518-1bd4-4f45-98ba-be016b39827e', 'banana', 400, 'A banana is an elongated, edible fruit – botanically a berry – produced by several kinds of large herbaceous flowering plants in the genus Musa. In some countries, bananas used for cooking may be called "plantains", distinguishing them from dessert bananas.');
INSERT INTO public.product (id, name, price, description)
VALUES ('e7f83943-7044-4406-9d32-9b229725f6d0', 'orange', 500, 'The orange is the fruit of various citrus species in the family Rutaceae; it primarily refers to Citrus × sinensis, which is also called sweet orange, to distinguish it from the related Citrus × aurantium, referred to as bitter orange.');


-- Discounts
INSERT INTO discount (name, rule, percent)
VALUES ('7 apples 10%', '{
  "0ca41a0a-2e65-40e9-a34e-7133b4adb9b4": 7
}', 10);
INSERT INTO discount (name, rule, percent)
VALUES ('set of 4 pears, 2 bananas 30%', '{
  "0645a967-badb-40d2-89fb-81c7b745899c": 4,
  "748cb518-1bd4-4f45-98ba-be016b39827e": 2
}', 30);

-- Coupons
INSERT INTO public.coupon (code, name, expire_at, rule, discount_percent)
VALUES ('orange555', 'Orange 30% discount coupon', '2021-11-06 15:30:31.542000', '{
  "e7f83943-7044-4406-9d32-9b229725f6d0": 0
}', 30);

INSERT INTO public.coupon (code, name, expire_at, rule, discount_percent)
VALUES ('123123', 'Orange 30% discount coupon', '2021-11-06 15:30:31.542000', '{
  "e7f83943-7044-4406-9d32-9b229725f6d0": 0
}', 30);


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
SELECT id, '0645a967-badb-40d2-89fb-81c7b745899c', 20
FROM tmp;

