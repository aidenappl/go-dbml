CREATE TABLE `T` (
  `a` int,
  `b` int,
  `x` int PRIMARY KEY,
  `shared` varchar(255),
  `c` int,
  `y` int
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `T` (`a`, `b`, `x`, `shared`, `c`, `y`)
VALUES
  (1, 2, 3, 'shared_value', 4, 5);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;