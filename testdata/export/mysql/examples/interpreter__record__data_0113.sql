CREATE TABLE `events` (
  `created_at` timestamp,
  `event_date` date,
  `event_time` time
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `events` (`created_at`, `event_date`, `event_time`)
VALUES
  ('2024-01-15T10:30:00', '2024-01-15', '10:30:00'),
  ('2024-12-31T23:59:59', '2024-12-31', '23:59:59');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;