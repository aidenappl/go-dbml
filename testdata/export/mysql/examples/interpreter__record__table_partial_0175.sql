CREATE TABLE `T` (
  `a` varchar(255) PRIMARY KEY,
  `b` text,
  `c` int
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `T` (`a`, `b`, `c`)
VALUES
  ('value_a', 'value_b', 3);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;