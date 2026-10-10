CREATE TABLE `students` (
  `id` int PRIMARY KEY
);

CREATE TABLE `courses` (
  `id` int PRIMARY KEY
);

CREATE TABLE `students_courses` (
  `students_id` int,
  `courses_id` int,
  PRIMARY KEY (`students_id`, `courses_id`)
);

ALTER TABLE `students_courses` ADD FOREIGN KEY (`students_id`) REFERENCES `students` (`id`);

ALTER TABLE `students_courses` ADD FOREIGN KEY (`courses_id`) REFERENCES `courses` (`id`);


-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `students` (`id`)
VALUES
  (1),
  (2),
  (3);
INSERT INTO `courses` (`id`)
VALUES
  (10),
  (20);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;