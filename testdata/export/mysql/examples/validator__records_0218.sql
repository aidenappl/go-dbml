CREATE TABLE `data` (
  `int_col` int,
  `float_col` decimal(10,2),
  `bool_col` boolean,
  `str_col` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `data` (`int_col`, `float_col`, `bool_col`, `str_col`)
VALUES
  (1, 3.14, TRUE, 'hello'),
  (2, -2.5, FALSE, 'world');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;