CREATE TABLE `orders` (
  `id` int DEFAULT 0,
  `code` int,
  PRIMARY KEY (`id`, `code`)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `orders` (`code`)
VALUES
  (NULL),
  (1);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;