CREATE TABLE `users` (
  `email` varchar(255)
);

CREATE INDEX `idx_email` ON `users` (`email`);
