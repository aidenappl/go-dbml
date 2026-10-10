CREATE TABLE `users` (
  `id` int PRIMARY KEY NOT NULL,
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp DEFAULT (now()),
  `email` varchar(255) UNIQUE,
  `name` varchar(255),
  `delete_status` boolean NOT NULL,
  `deleted_at` timestamp DEFAULT (now())
);

CREATE UNIQUE INDEX `users_index_0` ON `users` (`email`);
