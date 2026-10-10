CREATE TABLE `parent` (
  `id` int
);

CREATE TABLE `child` (
  `parent_id` int
);

ALTER TABLE `child` ADD FOREIGN KEY (`parent_id`) REFERENCES `parent` (`id`);
