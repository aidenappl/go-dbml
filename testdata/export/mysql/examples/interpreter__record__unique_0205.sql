CREATE TABLE `event_registrations` (
  `event_id` int,
  `attendee_id` int,
  `registration_date` timestamp
);

CREATE UNIQUE INDEX `event_registrations_index_0` ON `event_registrations` (`event_id`, `attendee_id`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `event_registrations` (`event_id`, `attendee_id`, `registration_date`)
VALUES
  (1, 100, '2024-01-01'),
  (1, 101, '2024-01-02'),
  (2, 100, '2024-01-03');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;