CREATE TABLE `ids` (
  `id` int PRIMARY KEY
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `ids` (`id`)
VALUES
  (1),
  (2),
  (3);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;