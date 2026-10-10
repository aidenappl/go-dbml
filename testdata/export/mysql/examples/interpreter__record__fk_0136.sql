CREATE TABLE `nodes` (
  `id` int PRIMARY KEY,
  `name` varchar(255),
  `parent_id` int
);

ALTER TABLE `nodes` ADD FOREIGN KEY (`parent_id`) REFERENCES `nodes` (`id`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `nodes` (`id`, `name`, `parent_id`)
VALUES
  (1, 'Root', NULL),
  (2, 'Child 1', 1),
  (3, 'Invalid Child', 999);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;