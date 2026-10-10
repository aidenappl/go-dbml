CREATE TABLE `users` (
  `id` int PRIMARY KEY
);

CREATE TABLE `posts` (
  `id` int PRIMARY KEY,
  `user_id` int
);

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
