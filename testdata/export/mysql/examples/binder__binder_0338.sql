CREATE SCHEMA `shop`;

CREATE TABLE `shop`.`products` (
  `id` int PRIMARY KEY,
  `category_id` int
);

CREATE TABLE `shop`.`orders` (
  `product_id` int,
  `category` int
);

ALTER TABLE `shop`.`orders` ADD FOREIGN KEY (`product_id`, `category`) REFERENCES `shop`.`products` (`id`, `category_id`);
