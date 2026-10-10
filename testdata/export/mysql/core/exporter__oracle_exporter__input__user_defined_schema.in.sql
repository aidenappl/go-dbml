CREATE SCHEMA `C##dbx`;

CREATE TABLE `C##dbx`.`users` (
  `user_id` int PRIMARY KEY AUTO_INCREMENT,
  `username` varchar2(40) UNIQUE NOT NULL,
  `password` varchar2(255) NOT NULL
);

CREATE TABLE `C##dbx`.`category` (
  `cat_id` int PRIMARY KEY AUTO_INCREMENT,
  `cat_name` nvarchar2(50) NOT NULL,
  `super_cat_id` int,
  `created_date` timestamp DEFAULT (current_timestamp)
);

CREATE TABLE `C##dbx`.`product` (
  `product_id` int PRIMARY KEY AUTO_INCREMENT,
  `product_name` nvarchar2(255) NOT NULL,
  `current_price` float NOT NULL,
  `expired_date` timestamp NOT NULL,
  `cat_id` int NOT NULL
);

CREATE TABLE `C##dbx`.`wishlist` (
  `user_id` int NOT NULL,
  `product_id` int NOT NULL,
  PRIMARY KEY (`user_id`, `product_id`)
);

CREATE INDEX `FULLTEXT_INDEX_CATEGORY` ON `C##dbx`.`category` (`cat_name`);

CREATE INDEX `FULLTEXT_INDEX_PRODUCT` ON `C##dbx`.`product` (`product_name`);

ALTER TABLE `C##dbx`.`category` ADD CONSTRAINT `FK_CATEGORY_SUPER` FOREIGN KEY (`super_cat_id`) REFERENCES `C##dbx`.`category` (`cat_id`);

ALTER TABLE `C##dbx`.`product` ADD CONSTRAINT `FK_PRODUCT_CATEGORY` FOREIGN KEY (`cat_id`) REFERENCES `C##dbx`.`category` (`cat_id`);

ALTER TABLE `C##dbx`.`wishlist` ADD CONSTRAINT `FK_WISHLIST_PRODUCT` FOREIGN KEY (`product_id`) REFERENCES `C##dbx`.`product` (`product_id`);

ALTER TABLE `C##dbx`.`wishlist` ADD CONSTRAINT `FK_WISHLIST_USER` FOREIGN KEY (`user_id`) REFERENCES `C##dbx`.`users` (`user_id`);
