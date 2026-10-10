CREATE TABLE `b` (
  `id` int PRIMARY KEY
);

CREATE TABLE `a` (
  `id` int PRIMARY KEY,
  `b_id` int,
  `self` int
);

ALTER TABLE `b` ADD FOREIGN KEY (`id`) REFERENCES `a` (`b_id`);

ALTER TABLE `a` ADD FOREIGN KEY (`id`) REFERENCES `a` (`self`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `b` (`id`)
VALUES
  (1);
INSERT INTO `a` (`id`, `b_id`, `self`)
VALUES
  (1, 1, 1),
  (2, 9, NULL);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;