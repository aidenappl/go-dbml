CREATE TABLE `products` (
  `id` int,
  `region` varchar(255),
  PRIMARY KEY (`id`, `region`)
);

CREATE TABLE `categories` (
  `id` int,
  `region` varchar(255),
  PRIMARY KEY (`id`, `region`)
);

CREATE TABLE `products_categories` (
  `products_id` int,
  `products_region` varchar,
  `categories_id` int,
  `categories_region` varchar,
  PRIMARY KEY (`products_id`, `products_region`, `categories_id`, `categories_region`)
);

ALTER TABLE `products_categories` ADD FOREIGN KEY (`products_id`, `products_region`) REFERENCES `products` (`id`, `region`);

ALTER TABLE `products_categories` ADD FOREIGN KEY (`categories_id`, `categories_region`) REFERENCES `categories` (`id`, `region`);


-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `products` (`id`, `region`)
VALUES
  (1, 'US'),
  (2, 'US');
INSERT INTO `categories` (`id`, `region`)
VALUES
  (1, 'US'),
  (3, 'EU');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;