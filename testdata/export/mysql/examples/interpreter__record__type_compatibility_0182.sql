CREATE TABLE `data` (
  `id` int,
  `name` varchar(255),
  `title` text,
  `content` text
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`id`, `name`, `title`, `content`)
VALUES
  (1, 'Alice', 'Bob', 'Line 1
Line 2'),
  (2, '', '', 'Tab	Separated'),
  (3, '123', 'true', 'Quote: "test"');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;