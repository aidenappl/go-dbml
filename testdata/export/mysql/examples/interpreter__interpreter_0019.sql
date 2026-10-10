CREATE TABLE `a` (
  `id` int
);

CREATE TABLE `b` (
  `a_id` int
);

ALTER TABLE `a` ADD FOREIGN KEY (`id`) REFERENCES `b` (`a_id`);
