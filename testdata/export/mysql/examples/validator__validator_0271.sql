CREATE TABLE `users` (
  `email` varchar(255)
);

CREATE INDEX `users_index_0` ON `users` (`email`) USING HASH;
