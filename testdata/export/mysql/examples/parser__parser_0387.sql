CREATE TABLE `users` (
  `id` integer PRIMARY KEY AUTO_INCREMENT,
  `email` varchar(255) UNIQUE NOT NULL,
  `status` user_status DEFAULT 'active',
  PRIMARY KEY (`id`, `email`)
);

CREATE UNIQUE INDEX `idx_email` ON `users` (`email`);

ALTER TABLE `users` COMMENT = 'Main users table';
