CREATE TABLE `test_table` (
  `id` INT PRIMARY KEY,
  `name` VARCHAR(100),
  `email` VARCHAR(100),
  `age` INT,
  `active` BOOLEAN DEFAULT true
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `test_table` (`id`, `name`, `email`, `age`, `active`)
VALUES
  (1, 'Alice', 'alice@example.com', 25, TRUE),
  (2, 'Bob', NULL, 30, FALSE),
  (3, NULL, NULL, NULL, NULL),
  (4, 'Charlie', NULL, NULL, TRUE);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;