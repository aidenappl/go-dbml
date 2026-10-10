CREATE TABLE `users` (
  `id` int,
  `quantity` int,
  `name` varchar(5),
  `price` decimal(5,2),
  `active` boolean,
  `status` ENUM ('active', 'inactive')
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `quantity`, `name`, `price`, `active`, `status`)
VALUES
  (invalid, 10, 'VeryLongName', 12345.67, FALSE, 'invalid_status');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;