CREATE TABLE `users` (
  `id` int PRIMARY KEY
);

CREATE TABLE `posts` (
  `user_id` int
);

CREATE TABLE `orders` (
  `user_id` int
);

ALTER TABLE `posts` ADD CONSTRAINT `r1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `orders` ADD CONSTRAINT `r2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
