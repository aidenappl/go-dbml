CREATE TABLE `messages` (
  `id` int,
  `text` varchar(20),
  `short` varchar(10)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `messages` (`id`, `text`, `short`)
VALUES
  (1, 'Hello', 'Test'),
  (2, '😀😁😂😃😄', '😀😁😂');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;