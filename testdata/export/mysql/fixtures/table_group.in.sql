CREATE TABLE `users` (
  `id` int,
  `full_name` varchar(255),
  `created_at` timestamp,
  `country_code` int
);

CREATE TABLE `merchants` (
  `id` int,
  `merchant_name` varchar(255),
  `country_code` int,
  `created at` varchar(255),
  `admin_id` int
);

ALTER TABLE `merchants` ADD FOREIGN KEY (`admin_id`) REFERENCES `users` (`id`);
