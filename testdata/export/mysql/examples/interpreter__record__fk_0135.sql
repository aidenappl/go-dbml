CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

CREATE TABLE `categories` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

CREATE TABLE `posts` (
  `id` int PRIMARY KEY,
  `title` varchar(255),
  `user_id` int,
  `category_id` int
);

ALTER TABLE `posts` ADD FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`);

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `name`)
VALUES
  (1, 'Alice');
INSERT INTO `categories` (`id`, `name`)
VALUES
  (1, 'Tech');
INSERT INTO `posts` (`id`, `title`, `user_id`, `category_id`)
VALUES
  (1, 'Valid Post', 1, 1),
  (2, 'Invalid Category', 1, 999),
  (3, 'Invalid User', 999, 1);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;