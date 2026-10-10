CREATE TABLE `data` (
  `active` boolean,
  `verified` bool
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`active`, `verified`)
VALUES
  (TRUE, FALSE),
  (FALSE, TRUE);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;