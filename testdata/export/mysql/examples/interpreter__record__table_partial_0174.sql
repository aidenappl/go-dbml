CREATE TABLE `T` (
  `a` int PRIMARY KEY,
  `b` int
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `T` (`a`, `b`)
VALUES
  (1, 2);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;