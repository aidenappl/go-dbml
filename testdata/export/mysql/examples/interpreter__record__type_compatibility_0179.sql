CREATE TABLE `data` (
  `id` int,
  `value` decimal,
  `count` int,
  `temperature` float
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`id`, `value`, `count`, `temperature`)
VALUES
  (1, 10000000000, 100, 314),
  (2, 0.0000314, 2000, -0.0025),
  (3, 200000000, NULL, NULL);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;