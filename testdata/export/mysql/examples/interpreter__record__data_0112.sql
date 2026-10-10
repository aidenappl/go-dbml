CREATE TABLE `data` (
  `name` varchar(255),
  `description` text,
  `code` char(10)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`name`, `description`, `code`)
VALUES
  ('Alice', 'A short description', 'ABC123'),
  ('Bob', 'Another description', '');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;