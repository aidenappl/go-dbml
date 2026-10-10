CREATE TABLE `orders` (
  `id` int PRIMARY KEY,
  `total` float,
  `placed` date,
  `status` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `orders` (`id`, `total`, `placed`, `status`)
VALUES
  (1, 10.5, '2024-03-01', 'paid'),
  (2, 20, '2024-03-02', NULL);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;