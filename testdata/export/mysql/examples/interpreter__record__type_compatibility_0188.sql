CREATE TABLE `data` (
  `id` int,
  `count` int,
  `active` boolean,
  `name` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`id`, `count`, `active`, `name`)
VALUES
  (NULL, NULL, NULL, '');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;