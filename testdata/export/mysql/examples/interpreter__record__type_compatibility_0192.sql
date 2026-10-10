CREATE TABLE `events` (
  `id` int,
  `created_at` timestamp,
  `event_date` date
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `events` (`id`, `created_at`, `event_date`)
VALUES
  (1, '2024-01-15T10:30:00+07:00', '2024-01-15');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;