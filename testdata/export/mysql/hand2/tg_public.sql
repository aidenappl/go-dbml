CREATE TABLE `users` (
  `id` int
);

CREATE TABLE `users2` (
  `id` int
);

ALTER TABLE `users` ADD FOREIGN KEY (`id`) REFERENCES `users2` (`id`);
