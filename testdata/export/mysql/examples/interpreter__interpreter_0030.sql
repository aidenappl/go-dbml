CREATE TABLE `users` (
  `id` integer PRIMARY KEY AUTO_INCREMENT,
  `email` varchar(255) UNIQUE NOT NULL,
  `created_at` timestamp,
  `updated_at` timestamp
);

CREATE TABLE `orders` (
  `id` integer PRIMARY KEY,
  `user_id` integer,
  `total` decimal(10,2),
  `created_at` timestamp,
  `updated_at` timestamp
);

ALTER TABLE `users` COMMENT = 'Users table';

ALTER TABLE `orders` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
