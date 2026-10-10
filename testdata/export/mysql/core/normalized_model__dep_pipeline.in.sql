CREATE TABLE `raw_users` (
  `id` int PRIMARY KEY,
  `email` varchar(255),
  `name` varchar(255),
  `created_at` timestamp,
  `raw_data` json
);

CREATE TABLE `raw_orders` (
  `id` int PRIMARY KEY,
  `user_id` int,
  `ordered_at` timestamp,
  `status` varchar(255),
  `raw_data` json
);

CREATE TABLE `raw_order_items` (
  `id` int PRIMARY KEY,
  `order_id` int,
  `product_id` int,
  `quantity` int,
  `price` decimal
);

CREATE TABLE `stg_users` (
  `user_id` int PRIMARY KEY,
  `email` varchar(255),
  `name` varchar(255),
  `created_at` timestamp
);

CREATE TABLE `stg_orders` (
  `order_id` int PRIMARY KEY,
  `user_id` int,
  `ordered_at` timestamp,
  `status` varchar(255)
);

CREATE TABLE `stg_order_items` (
  `order_item_id` int PRIMARY KEY,
  `order_id` int,
  `product_id` int,
  `quantity` int,
  `unit_price` decimal,
  `line_total` decimal
);

CREATE TABLE `dim_users` (
  `user_id` int PRIMARY KEY,
  `email` varchar(255),
  `name` varchar(255),
  `first_order_at` timestamp,
  `total_orders` int,
  `total_spent` decimal,
  `created_at` timestamp
);

CREATE TABLE `fct_orders` (
  `order_id` int PRIMARY KEY,
  `user_id` int,
  `ordered_at` timestamp,
  `status` varchar(255),
  `total_items` int,
  `order_total` decimal
);

ALTER TABLE `raw_orders` ADD FOREIGN KEY (`user_id`) REFERENCES `raw_users` (`id`);

ALTER TABLE `raw_order_items` ADD FOREIGN KEY (`order_id`) REFERENCES `raw_orders` (`id`);

ALTER TABLE `stg_orders` ADD FOREIGN KEY (`user_id`) REFERENCES `stg_users` (`user_id`);

ALTER TABLE `stg_order_items` ADD FOREIGN KEY (`order_id`) REFERENCES `stg_orders` (`order_id`);

ALTER TABLE `fct_orders` ADD FOREIGN KEY (`user_id`) REFERENCES `dim_users` (`user_id`);
