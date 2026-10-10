CREATE TABLE `merchants` (
  `id` int,
  `country_code` varchar(255),
  PRIMARY KEY (`id`, `country_code`)
);

CREATE TABLE `orders` (
  `id` int PRIMARY KEY,
  `merchant_id` int,
  `country` varchar(255),
  `amount` decimal
);

ALTER TABLE `orders` ADD FOREIGN KEY (`merchant_id`, `country`) REFERENCES `merchants` (`id`, `country_code`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `merchants` (`id`, `country_code`)
VALUES
  (1, 'US'),
  (1, 'UK'),
  (2, 'US');
INSERT INTO `orders` (`id`, `merchant_id`, `country`, `amount`)
VALUES
  (1, 1, 'US', 100),
  (2, 1, 'UK', 200.5),
  (3, 2, 'US', 50);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;