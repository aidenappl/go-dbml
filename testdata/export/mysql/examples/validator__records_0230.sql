CREATE TABLE `data` (
  `col1` varchar(255),
  `col2` boolean,
  `col3` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`col1`, `col2`, `col3`)
VALUES
  (NULL, TRUE, 'false'),
  (NULL, TRUE, 'false'),
  (NULL, TRUE, 'false'),
  (NULL, TRUE, 'false');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;