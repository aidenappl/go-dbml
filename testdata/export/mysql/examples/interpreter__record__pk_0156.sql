CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

CREATE TABLE `countries` (
  `code` varchar(2) PRIMARY KEY,
  `name` varchar(255)
);

CREATE TABLE `accounts` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `name`)
VALUES
  (1, 'Alice'),
  (2, 'Bob');
INSERT INTO `countries` (`code`, `name`)
VALUES
  ('US', 'United States'),
  ('UK', 'United Kingdom');
INSERT INTO `accounts` (`id`, `name`)
VALUES
  (0, 'Default Account'),
  (1, 'User Account'),
  (2, 'Business Account');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;