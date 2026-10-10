CREATE TABLE `items` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `items` (`name`)
VALUES
  ('Alice'),
  ('Bob');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;