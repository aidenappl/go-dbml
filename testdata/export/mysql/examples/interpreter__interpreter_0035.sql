CREATE TABLE `users` (
  `id` int
);

CREATE TABLE `posts` (
  `user_id` int
);

ALTER TABLE `posts` ADD CONSTRAINT `post_author` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION;
