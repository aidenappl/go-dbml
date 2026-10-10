CREATE TABLE `a` (
  `id` int,
  `b_id` int
);

CREATE TABLE `b` (
  `id` int,
  `a_id` int
);

ALTER TABLE `a` ADD FOREIGN KEY (`b_id`) REFERENCES `b` (`id`);

ALTER TABLE `b` ADD FOREIGN KEY (`a_id`) REFERENCES `a` (`id`);
