CREATE TABLE `users` (
  `id` int PRIMARY KEY
);

CREATE TABLE `posts` (
  `user_id` int NOT NULL
);

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
