CREATE TABLE `user` (
  `id` int PRIMARY KEY,
  `name` string
);

CREATE TABLE `country` (
  `id` int PRIMARY KEY,
  `name` string NOT NULL
);

CREATE TABLE `product` (
  `id` int PRIMARY KEY,
  `name` string,
  `price` decimal NOT NULL
);

CREATE TABLE `merchant` (
  `id` int PRIMARY KEY,
  `user_id` int,
  `product_id` int,
  `address` string
);

ALTER TABLE `country` COMMENT = 'name is required';

ALTER TABLE `product` COMMENT = 'product must have price';

ALTER TABLE `merchant` COMMENT = 'merchants sell a lot';

ALTER TABLE `merchant` ADD FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

ALTER TABLE `merchant` ADD FOREIGN KEY (`product_id`) REFERENCES `product` (`id`);
