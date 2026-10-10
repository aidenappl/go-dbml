CREATE TABLE `articles` (
  `id` int,
  `content` text,
  `description` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `articles` (`id`, `content`, `description`)
VALUES
  (1, 'This is a very long text content that can be arbitrarily long without any length restrictions because text type does not have a length parameter', 'This is a very long description that can be arbitrarily long');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;