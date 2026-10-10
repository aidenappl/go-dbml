CREATE TABLE `merchants` (
  `id` int PRIMARY KEY,
  `country_code` varchar(255)
);

CREATE TABLE `orders` (
  `merchant_id` int,
  `country` varchar(255)
);

ALTER TABLE `orders` ADD FOREIGN KEY (`merchant_id`, `country`) REFERENCES `merchants` (`id`, `country_code`);
