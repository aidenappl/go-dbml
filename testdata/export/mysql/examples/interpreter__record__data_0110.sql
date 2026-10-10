CREATE TABLE `data` (
  `price` decimal(10,2),
  `rate` float,
  `amount` numeric
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`price`, `rate`, `amount`)
VALUES
  (99.99, 3.14159, 0.001),
  (50.5, 0.5, 100);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;