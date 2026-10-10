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
  `name` varchar(255)
);

CREATE TABLE `b` (
  `id` int,
  `c_id` int
);

CREATE TABLE `c` (
  `id` int,
  `b_id` int
);

CREATE TABLE `unary_expression` (
  `id` int(-1) DEFAULT -2,
  `id2` int(--1) DEFAULT -2,
  `id3` int(+-+---+0.1) DEFAULT 7.2225
);

CREATE TABLE `injected_unary_expression` (
  `id` int,
  `id2` int(--1) DEFAULT -2,
  `id3` int(+-+---+0.1) DEFAULT -7.2225
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
  `name` varchar(255)
);

CREATE TABLE `schemaA`.`locations` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

ALTER TABLE `schemaA`.`products` ADD FOREIGN KEY (`name`) REFERENCES `ecommerce`.`users` (`id`);

ALTER TABLE `schemaA`.`locations` ADD FOREIGN KEY (`name`) REFERENCES `users` (`id`);

ALTER TABLE `ecommerce`.`users` ADD FOREIGN KEY (`id`) REFERENCES `users` (`id`);

ALTER TABLE `ecommerce`.`users` ADD CONSTRAINT `name_optional` FOREIGN KEY (`id`) REFERENCES `users` (`name`);

ALTER TABLE `c` ADD CONSTRAINT `short_ref` FOREIGN KEY (`id`) REFERENCES `b` (`id`);

ALTER TABLE `b` ADD CONSTRAINT `long_ref` FOREIGN KEY (`c_id`) REFERENCES `c` (`id`);
