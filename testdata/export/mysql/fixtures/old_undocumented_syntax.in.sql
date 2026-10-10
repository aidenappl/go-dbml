CREATE TABLE `customer` (
  `customer_id` SMALLINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `store_id` TINYINT NOT NULL,
  `first_name` VARCHAR(45) NOT NULL,
  `last_name` VARCHAR(45) NOT NULL DEFAULT false,
  `email` VARCHAR(50),
  `address_id` SMALLINT NOT NULL,
  `active` BOOLEAN NOT NULL DEFAULT true,
  `create_date` DATETIME NOT NULL,
  `last_update` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE `cities` (
  `id` integer PRIMARY KEY,
  `name` e DEFAULT 'hello',
  `country_id` integer
);

CREATE TABLE `country` (
  `id` integer,
  `cities` string[]
);

CREATE TABLE `citites` (
  `id` integer,
  `name` string
);

CREATE INDEX `citites_index_0` ON `citites` (`id`, `name`);

ALTER TABLE `cities` COMMENT = 'sasasa';
