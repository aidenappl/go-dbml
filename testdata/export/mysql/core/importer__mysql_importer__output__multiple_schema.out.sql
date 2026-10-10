CREATE SCHEMA `ecommerce`;

CREATE SCHEMA `schemaA`;

CREATE SCHEMA `schemaB`;

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
  `name` varchar(255) COMMENT 'Product name'
);

CREATE TABLE `orders` (
  `id1` int,
  `id2` int,
  `pid` int,
  `name` varchar(255),
  `name2` varchar(255),
  PRIMARY KEY (`id1`, `id2`)
);

CREATE TABLE `orders2` (
  `id1` int,
  `id2` int,
  `pid` int,
  PRIMARY KEY (`id1`, `id2`)
);

CREATE TABLE `orders4` (
  `id1` int,
  `id2` int,
  PRIMARY KEY (`id1`, `id2`)
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
  `name` varchar(255) COMMENT 'Sample field comment on multiples schema',
  `lid` int
);

CREATE TABLE `schemaA`.`locations` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

CREATE TABLE `schemaB`.`orders3` (
  `id1` int PRIMARY KEY,
  `id2` int
);

CREATE UNIQUE INDEX `I_UNIQUE_1` ON `orders` (`name`, `name2`) USING BTREE;

CREATE INDEX `INDEX_2` ON `orders` (`id1`, `id2`) USING HASH;

ALTER TABLE `orders` COMMENT = 'this is table orders';

ALTER TABLE `schemaA`.`products` ADD CONSTRAINT `FK_1` FOREIGN KEY (`lid`) REFERENCES `schemaA`.`locations` (`id`) ON DELETE CASCADE;

ALTER TABLE `orders` ADD FOREIGN KEY (`pid`) REFERENCES `schemaA`.`products` (`id`);

ALTER TABLE `orders2` ADD FOREIGN KEY (`pid`) REFERENCES `schemaA`.`products` (`id`);

ALTER TABLE `orders2` ADD CONSTRAINT `CFK_1` FOREIGN KEY (`id1`, `id2`) REFERENCES `orders` (`id1`, `id2`) ON DELETE CASCADE ON UPDATE SET NULL;

ALTER TABLE `ecommerce`.`users` ADD FOREIGN KEY (`id`) REFERENCES `users` (`id`);

ALTER TABLE `ecommerce`.`users` ADD CONSTRAINT `name_optional` FOREIGN KEY (`id`) REFERENCES `users` (`name`);

ALTER TABLE `schemaA`.`products` ADD FOREIGN KEY (`name`) REFERENCES `ecommerce`.`users` (`id`);

ALTER TABLE `schemaA`.`locations` ADD FOREIGN KEY (`name`) REFERENCES `users` (`id`);
