CREATE TABLE `a` (
  `id` int,
  `bid` int
);

CREATE TABLE `b` (
  `id` int
);

ALTER TABLE `b` ADD CONSTRAINT `named_ref` FOREIGN KEY (`id`) REFERENCES `a` (`bid`) ON UPDATE CASCADE;

CREATE TABLE `a_b` (
  `a_id` int,
  `b_id` int,
  PRIMARY KEY (`a_id`, `b_id`)
);

ALTER TABLE `a_b` ADD FOREIGN KEY (`a_id`) REFERENCES `a` (`id`);

ALTER TABLE `a_b` ADD FOREIGN KEY (`b_id`) REFERENCES `b` (`id`);

