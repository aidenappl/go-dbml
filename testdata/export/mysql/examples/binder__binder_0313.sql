CREATE TABLE `users` (
  `first_name` varchar(255),
  `last_name` varchar(255),
  `email` varchar(255)
);

CREATE INDEX `users_index_0` ON `users` (`first_name`, `last_name`);

CREATE UNIQUE INDEX `users_index_1` ON `users` (`email`);
