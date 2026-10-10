CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `email` varchar(255)
);

CREATE UNIQUE INDEX `users_index_0` ON `users` (`email`);
