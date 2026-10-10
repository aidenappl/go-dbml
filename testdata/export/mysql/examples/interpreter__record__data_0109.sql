CREATE TABLE `data` (
  `id` int,
  `count` integer,
  `small` smallint,
  `big` bigint
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`id`, `count`, `small`, `big`)
VALUES
  (1, 42, -100, 9999999999),
  (0, 0, 0, 0);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;