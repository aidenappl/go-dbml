CREATE SCHEMA `schemaB`;

CREATE SCHEMA `ecommerce`;

CREATE SCHEMA `schemaA`;

CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `name` varchar(255),
  `pjs` ENUM ('created2', 'running2', 'done2', 'failure2'),
  `pjs2` ENUM ('created2', 'running2', 'done2', 'failure2'),
  `pg` ENUM ('male', 'female'),
  `pg2` ENUM ('male2', 'female2')
);

CREATE TABLE `products` (
  `id` int PRIMARY KEY,
  `name` varchar(255) COMMENT 'Product name of public schema'
);

CREATE TABLE `ecommerce`.`users` (
  `id` int PRIMARY KEY,
  `name` varchar(255),
  `ejs` ENUM ('created2', 'running2', 'done2', 'failure2'),
  `ejs2` ENUM ('created2', 'running2', 'done2', 'failure2'),
  `eg` ENUM ('male', 'female'),
  `eg2` ENUM ('male2', 'female2')
);

CREATE TABLE `schemaA`.`products` (
  `id` int PRIMARY KEY,
  `name` varchar(255) COMMENT 'Product name of schemaA',
  `lid` int
);

CREATE TABLE `schemaA`.`locations` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

CREATE UNIQUE INDEX `products_index_0` ON `products` (`id`) USING HASH;

CREATE INDEX `product_status` ON `schemaA`.`products` (`id`, `name`);

ALTER TABLE `schemaA`.`locations` COMMENT = 'This is a note in table "locations"';

ALTER TABLE `schemaA`.`products` ADD CONSTRAINT `fk_1` FOREIGN KEY (`lid`) REFERENCES `schemaA`.`locations` (`id`);

ALTER TABLE `ecommerce`.`users` ADD FOREIGN KEY (`id`) REFERENCES `users` (`id`);

ALTER TABLE `ecommerce`.`users` ADD CONSTRAINT `name_optional` FOREIGN KEY (`id`) REFERENCES `users` (`name`);

ALTER TABLE `schemaA`.`products` ADD FOREIGN KEY (`name`) REFERENCES `ecommerce`.`users` (`id`);

ALTER TABLE `schemaA`.`locations` ADD FOREIGN KEY (`name`) REFERENCES `users` (`id`);
