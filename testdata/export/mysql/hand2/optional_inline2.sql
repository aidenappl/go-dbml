CREATE TABLE `x` (
  `id` int PRIMARY KEY
);

CREATE TABLE `y` (
  `x_id` int,
  `x2` int
);

CREATE TABLE `z` (
  `x3` int
);

ALTER TABLE `y` ADD FOREIGN KEY (`x_id`) REFERENCES `x` (`id`);

ALTER TABLE `y` ADD FOREIGN KEY (`x2`) REFERENCES `x` (`id`);

CREATE TABLE `x_z` (
  `x_id` int,
  `z_x3` int,
  PRIMARY KEY (`x_id`, `z_x3`)
);

ALTER TABLE `x_z` ADD FOREIGN KEY (`x_id`) REFERENCES `x` (`id`);

ALTER TABLE `x_z` ADD FOREIGN KEY (`z_x3`) REFERENCES `z` (`x3`);

