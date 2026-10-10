CREATE TABLE `users` (
  `id` int PRIMARY KEY
);

CREATE TABLE `posts` (
  `created_at` timestamp DEFAULT (now()),
  `title` varchar(255),
  `owner_id` int,
  `id` bigint,
  CONSTRAINT `posts.pos` CHECK (id > 0)
);

CREATE INDEX `posts_index_0` ON `posts` (`created_at`);

ALTER TABLE `posts` COMMENT = 'p';

ALTER TABLE `posts` ADD FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `posts` (`id`, `title`, `owner_id`, `created_at`)
VALUES
  (1, 'a', 1, '2024-01-01T00:00:00'),
  (1, 'b', 2, NULL);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;