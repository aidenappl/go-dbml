CREATE TABLE `users` (
  `id` integer PRIMARY KEY AUTO_INCREMENT,
  `email` varchar(255) UNIQUE NOT NULL COMMENT 'User email',
  `status` varchar(255) DEFAULT 'active',
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp
);

CREATE TABLE `orders` (
  `id` integer PRIMARY KEY AUTO_INCREMENT,
  `user_id` integer,
  `status` ENUM ('pending', 'processing', 'completed', 'cancelled') DEFAULT 'pending',
  `total` decimal(10,2) NOT NULL,
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp
);

CREATE UNIQUE INDEX `idx_users_email` ON `users` (`email`);

ALTER TABLE `users` COMMENT = 'Main users table';

ALTER TABLE `orders` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
