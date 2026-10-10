CREATE TABLE `orders` (
  `id` SERIAL PRIMARY KEY AUTO_INCREMENT,
  `user_id` int UNIQUE NOT NULL,
  `status` ENUM ('created', 'running', 'done', 'failure'),
  `created_at` varchar(255),
  `modified_at` timestamp(2)
);

CREATE TABLE `order_items` (
  `order_id` int,
  `product_id` int,
  `quantity` int DEFAULT 1,
  `created_at` timestamp,
  `modified_at` time
);

CREATE TABLE `products` (
  `id` int,
  `name` varchar(255),
  `merchant_id` int NOT NULL,
  `price` int,
  `status` ENUM ('Out of Stock', 'In Stock'),
  `created_at` datetime DEFAULT (now()),
  `modified_at` timestamp,
  PRIMARY KEY (`id`, `name`)
);

CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `full_name` varchar(255),
  `email` varchar(255) UNIQUE,
  `gender` varchar(255),
  `date_of_birth` varchar(255),
  `created_at` varchar(255),
  `modified_at` time(2),
  `country_code` int
);

CREATE TABLE `merchants` (
  `id` int PRIMARY KEY,
  `merchant_name` varchar(255),
  `country_code` int,
  `created_at` varchar(255),
  `modified_at` time,
  `admin_id` int
);

CREATE TABLE `countries` (
  `code` int PRIMARY KEY,
  `name` varchar(255),
  `continent_name` varchar(255)
);

CREATE TABLE `foo` (
  `bar` text[],
  `bar2` int[1],
  `bar3` int[2][3],
  `bar4` int[],
  `bar5` int[2],
  `bar6` text[8],
  `bar7` text[100],
  `bar8` time(2)[],
  `bar9` time(1)[1],
  `bar10` time(1)[],
  `bar11` time[5],
  `bar12` timestamp(2)[10][2][5],
  `bar13` character varying[],
  `bar14` character varying(25)[][2][],
  `bar15` character varying[76]
);

CREATE TABLE `accounts` (
  `id` integer NOT NULL AUTO_INCREMENT,
  `domain` character varying(100),
  `name` character varying(100),
  `slug` character varying(10),
  PRIMARY KEY (`id`)
);

CREATE TABLE `tests` (
  `id` bigint NOT NULL,
  `type` character varying DEFAULT ('testing'::charactervarying)
);

CREATE TABLE `users2` (
  `username` character varying NOT NULL,
  `hashed_password` character varying NOT NULL,
  `full_name` character varying NOT NULL,
  `email` character varying NOT NULL,
  `password_changed_at` timestamp NOT NULL DEFAULT ('0001-01-01 00:00:00+00'::timestampwithtimezone),
  `created_at` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `table1` (
  `field` text NOT NULL,
  `field2` text NOT NULL,
  PRIMARY KEY (`field`, `field2`)
);

CREATE TABLE `table2` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `field` text,
  `field2` text
);

CREATE TABLE `test_table` (
  `name` varchar(255) NOT NULL,
  `id` varchar(255),
  `linked_name` varchar(255) NOT NULL,
  PRIMARY KEY (`name`, `id`)
);

CREATE INDEX `product_status` ON `products` (`merchant_id`, `status`);

CREATE UNIQUE INDEX `products_index_1` ON `products` (`id`) USING HASH;

ALTER TABLE `products` COMMENT = 'Products table comment';

ALTER TABLE `users` COMMENT = 'Store user data';

ALTER TABLE `order_items` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);

ALTER TABLE `order_items` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `users` ADD FOREIGN KEY (`country_code`) REFERENCES `countries` (`code`);

ALTER TABLE `merchants` ADD FOREIGN KEY (`country_code`) REFERENCES `countries` (`code`);

ALTER TABLE `products` ADD FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`);

ALTER TABLE `merchants` ADD FOREIGN KEY (`admin_id`) REFERENCES `users` (`id`);

ALTER TABLE `table2` ADD FOREIGN KEY (`field`, `field2`) REFERENCES `table1` (`field`, `field2`);

ALTER TABLE `test_table` ADD CONSTRAINT `test_table_fk` FOREIGN KEY (`linked_name`, `id`) REFERENCES `test_table` (`name`, `id`);
