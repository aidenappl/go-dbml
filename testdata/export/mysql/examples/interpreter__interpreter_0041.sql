CREATE TABLE `users` (
  `id` int
);

CREATE TABLE `posts` (
  `author_id` int
);

ALTER TABLE `posts` ADD FOREIGN KEY (`author_id`) REFERENCES `users` (`id`);
