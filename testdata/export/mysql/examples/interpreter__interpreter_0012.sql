CREATE TABLE `users` (
  `email` varchar(255)
);

CREATE UNIQUE INDEX `users_index_0` ON `users` (`email`);
