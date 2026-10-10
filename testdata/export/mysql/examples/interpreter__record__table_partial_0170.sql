CREATE TABLE `posts` (
  `id` int PRIMARY KEY,
  `created_at` timestamp,
  `title` varchar(255),
  `content` text,
  `deleted_at` timestamp,
  `is_deleted` boolean,
  `author_id` int
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `posts` (`id`, `title`, `content`, `author_id`)
VALUES
  (1, 'First Post', 'Hello World', 100),
  (2, 'Second Post', 'Welcome', 101);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;