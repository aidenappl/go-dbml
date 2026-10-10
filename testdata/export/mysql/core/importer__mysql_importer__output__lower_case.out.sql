CREATE TABLE `orders` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int UNIQUE NOT NULL,
  `status` ENUM ('created', 'running', 'done', 'failure'),
  `created_at` varchar(255)
);

CREATE TABLE `order_items` (
  `order_id` int,
  `product_id` int,
  `quantity` int DEFAULT 1
);

CREATE TABLE `products` (
  `id` int,
  `name` varchar(255),
  `merchant_id` int NOT NULL,
  `price` int COMMENT 'Products price field',
  `status` ENUM ('Out of Stock', 'In Stock'),
  `created_at` datetime DEFAULT (now()),
  PRIMARY KEY (`id`, `price`)
);

CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `full_name` varchar(255),
  `email` varchar(255) UNIQUE,
  `gender` varchar(255),
  `date_of_birth` varchar(255),
  `created_at` varchar(255),
  `country_code` int
);

CREATE TABLE `merchants` (
  `id` int PRIMARY KEY,
  `merchant_name` varchar(255),
  `country_code` int,
  `created_at` varchar(255),
  `admin_id` int
);

CREATE TABLE `countries` (
  `code` int PRIMARY KEY,
  `name` varchar(255),
  `continent_name` varchar(255)
);

CREATE TABLE `composite_service_item` (
  `composite_service_item_id` int(11) PRIMARY KEY NOT NULL,
  `service_id` int(11) NOT NULL,
  `ticket_item_id` int(11) NOT NULL,
  `discount` decimal(10,0) NOT NULL,
  `original_price` decimal(10,0) NOT NULL,
  `patient_price` decimal(10,0) NOT NULL,
  `insurance_price` decimal(10,0) NOT NULL
);

CREATE TABLE `Countries` (
  `Id` int PRIMARY KEY NOT NULL,
  `Name` varchar(32) NOT NULL
);

CREATE TABLE `States` (
  `Id` int NOT NULL,
  `CountryId` int NOT NULL,
  `Name` varchar(32) NOT NULL,
  PRIMARY KEY (`Id`, `CountryId`)
);

CREATE TABLE `Cities` (
  `Id` int NOT NULL,
  `CountryId` int NOT NULL,
  `StateId` int NOT NULL,
  `Name` varchar(32) NOT NULL,
  PRIMARY KEY (`Id`, `StateId`, `CountryId`)
);

CREATE INDEX `product_status` ON `products` (`merchant_id`, `status`);

CREATE UNIQUE INDEX `products_index_1` ON `products` (`id`) USING HASH;

CREATE INDEX `ticket_item_id` ON `composite_service_item` (`ticket_item_id`);

CREATE INDEX `service_id` ON `composite_service_item` (`service_id`);

CREATE INDEX `service_id, ticket_item_id` ON `composite_service_item` (`service_id`, `ticket_item_id`) USING HASH;

CREATE INDEX `composite_service_item_id` ON `composite_service_item` (`composite_service_item_id`) USING BTREE;

CREATE INDEX `test test test key` ON `composite_service_item` (`service_id`, `patient_price`) USING HASH;

CREATE INDEX `IX_States_CountryId` ON `States` (`CountryId`);

CREATE INDEX `IX_Cities_CountryId` ON `Cities` (`CountryId`);

CREATE INDEX `IX_Cities_StateId_CountryId` ON `Cities` (`StateId`, `CountryId`);

ALTER TABLE `products` COMMENT = 'Notes about products table';

ALTER TABLE `users` COMMENT = 'User basic information';

ALTER TABLE `States` ADD CONSTRAINT `FK_States_Countries_CountryId` FOREIGN KEY (`CountryId`) REFERENCES `Countries` (`Id`) ON DELETE CASCADE;

ALTER TABLE `Cities` ADD CONSTRAINT `FK_Cities_Countries_CountryId` FOREIGN KEY (`CountryId`) REFERENCES `Countries` (`Id`) ON DELETE CASCADE;

ALTER TABLE `Cities` ADD CONSTRAINT `FK_Cities_States_StateId_CountryId` FOREIGN KEY (`StateId`, `CountryId`) REFERENCES `States` (`Id`, `CountryId`) ON DELETE CASCADE;

ALTER TABLE `order_items` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);

ALTER TABLE `order_items` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `users` ADD FOREIGN KEY (`country_code`) REFERENCES `countries` (`code`);

ALTER TABLE `merchants` ADD FOREIGN KEY (`country_code`) REFERENCES `countries` (`code`);

ALTER TABLE `products` ADD FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`);

ALTER TABLE `merchants` ADD FOREIGN KEY (`admin_id`) REFERENCES `users` (`id`);
