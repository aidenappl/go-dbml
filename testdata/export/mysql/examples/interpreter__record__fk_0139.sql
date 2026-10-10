CREATE TABLE `products` (
  `id` int,
  `region` varchar(255),
  PRIMARY KEY (`id`, `region`)
);

CREATE TABLE `orders` (
  `id` int PRIMARY KEY,
  `product_id` int,
  `region` varchar(255)
);

ALTER TABLE `orders` ADD FOREIGN KEY (`product_id`, `region`) REFERENCES `products` (`id`, `region`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `products` (`id`, `region`)
VALUES
  (1, 'US'),
  (NULL, 'EU');
INSERT INTO `orders` (`id`, `product_id`, `region`)
VALUES
  (1, 1, 'US');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;