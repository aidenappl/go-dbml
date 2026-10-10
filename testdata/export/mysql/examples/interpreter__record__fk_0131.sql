CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

CREATE TABLE `posts` (
  `id` int PRIMARY KEY,
  `title` varchar(255),
  `user_id` int
);

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `name`)
VALUES
  (1, 'Alice'),
  (2, 'Bob');
INSERT INTO `posts` (`id`, `title`, `user_id`)
VALUES
  (1, 'Post 1', 1),
  (2, 'Post 2', 2);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;