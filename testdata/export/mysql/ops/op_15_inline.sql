CREATE TABLE `b` (
  `id` int PRIMARY KEY
);

CREATE TABLE `a` (
  `id` int PRIMARY KEY,
  `b_id` int,
  `self` int
);

CREATE TABLE `b_a` (
  `b_id` int,
  `a_b_id` int,
  PRIMARY KEY (`b_id`, `a_b_id`)
);

ALTER TABLE `b_a` ADD FOREIGN KEY (`b_id`) REFERENCES `b` (`id`);

ALTER TABLE `b_a` ADD FOREIGN KEY (`a_b_id`) REFERENCES `a` (`b_id`);


CREATE TABLE `a_a` (
  `a_id` int,
  `a_self` int,
  PRIMARY KEY (`a_id`, `a_self`)
);

ALTER TABLE `a_a` ADD FOREIGN KEY (`a_id`) REFERENCES `a` (`id`);

ALTER TABLE `a_a` ADD FOREIGN KEY (`a_self`) REFERENCES `a` (`self`);


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