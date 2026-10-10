CREATE TABLE `categories` (
  `id` int PRIMARY KEY,
  `parent_id` int,
  `name` varchar(255)
);

ALTER TABLE `categories` ADD FOREIGN KEY (`parent_id`) REFERENCES `categories` (`id`);
