CREATE TABLE `users` (
  `id` int PRIMARY KEY
);

CREATE TABLE `posts` (
  `user_id` int
);

CREATE TABLE `comments` (
  `user_id` int
);

CREATE TABLE `likes` (
  `user_id` int
);

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `comments` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `likes` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
