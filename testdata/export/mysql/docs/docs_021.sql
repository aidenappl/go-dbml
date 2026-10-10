CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `name` varchar(255),
  `email` varchar(255)
);

CREATE TABLE `posts` (
  `id` int PRIMARY KEY,
  `title` varchar(255),
  `published` boolean
);

CREATE TABLE `comments` (
  `id` int PRIMARY KEY,
  `user_id` int,
  `post_id` int,
  `title` string
);

ALTER TABLE `comments` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `comments` ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `name`, `email`)
VALUES
  (1, 'Alice', 'alice@example.com'),
  (2, 'Bob', 'bob@example.com');
INSERT INTO `posts` (`id`, `title`, `published`)
VALUES
  (1, 'First Post', TRUE),
  (2, 'Second Post', FALSE);
INSERT INTO `comments` (`id`, `user_id`, `post_id`, `title`)
VALUES
  (1, 2, 1, 'First comment of first post by the second user');

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;