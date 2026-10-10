CREATE TABLE `events` (
  `id` int,
  `created_at` timestamp
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `events` (`id`, `created_at`)
VALUES
  (1, '12345'),
  (2, 'true');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;