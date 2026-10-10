CREATE TABLE `raw_orders` (
  `id` int PRIMARY KEY,
  `user_id` int,
  `amount` decimal,
  `created_at` timestamp
);

CREATE TABLE `stg_orders` (
  `id` int PRIMARY KEY,
  `user_id` int,
  `amount` decimal
);

CREATE TABLE `fct_orders` (
  `id` int PRIMARY KEY,
  `user_id` int,
  `revenue` decimal
);

CREATE TABLE `mart_orders` (
  `id` int PRIMARY KEY,
  `total` decimal
);
