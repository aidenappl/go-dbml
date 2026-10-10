CREATE TABLE `users` (
  `id` int
);

CREATE TABLE `posts` (
  `user_id` int
);

ALTER TABLE `posts` ADD CONSTRAINT `user_posts` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
