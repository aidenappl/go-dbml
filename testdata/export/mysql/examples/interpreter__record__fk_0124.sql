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

CREATE TABLE `posts` (
  `id` int PRIMARY KEY,
  `user_id` int,
  `title` varchar(255)
);

CREATE TABLE `cities` (
  `id` int PRIMARY KEY,
  `country_code` varchar(2),
  `name` varchar(255)
);

CREATE TABLE `transactions` (
  `id` int PRIMARY KEY,
  `account_id` int,
  `amount` decimal
);

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `cities` ADD FOREIGN KEY (`country_code`) REFERENCES `countries` (`code`);

ALTER TABLE `transactions` ADD FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`);

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
INSERT INTO `posts` (`id`, `user_id`, `title`)
VALUES
  (1, 1, 'Alice Post'),
  (2, 1, 'Another Post'),
  (3, 2, 'Bob Post');
INSERT INTO `cities` (`id`, `country_code`, `name`)
VALUES
  (1, 'US', 'New York'),
  (2, 'UK', 'London');
INSERT INTO `transactions` (`id`, `account_id`, `amount`)
VALUES
  (1, 0, 100),
  (2, 1, 50),
  (3, 2, 25);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;