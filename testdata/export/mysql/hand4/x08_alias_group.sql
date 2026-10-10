CREATE TABLE `very_long_users` (
  `id` int PRIMARY KEY
);

CREATE TABLE `orders` (
  `uid` int,
  `uid2` int
);

ALTER TABLE `orders` ADD FOREIGN KEY (`uid`) REFERENCES `very_long_users` (`id`);

ALTER TABLE `very_long_users` ADD FOREIGN KEY (`id`) REFERENCES `orders` (`uid2`);
