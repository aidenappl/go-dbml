CREATE TABLE `merchants` (
  `id` int,
  `country_code` varchar(255),
  PRIMARY KEY (`id`, `country_code`)
);

CREATE TABLE `orders` (
  `id` int PRIMARY KEY,
  `merchant_id` int,
  `country` varchar(255)
);

ALTER TABLE `orders` ADD FOREIGN KEY (`merchant_id`, `country`) REFERENCES `merchants` (`id`, `country_code`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `merchants` (`id`, `country_code`)
VALUES
  (1, 'US'),
  (2, 'UK');
INSERT INTO `orders` (`id`, `merchant_id`, `country`)
VALUES
  (1, 1, 'US'),
  (2, 1, 'UK');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;