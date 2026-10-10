CREATE TABLE `users` (
  `email` varchar(255)
);

CREATE UNIQUE INDEX `idx_email` ON `users` (`email`) USING BTREE;
