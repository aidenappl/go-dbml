CREATE TABLE `users` (
  `id` integer PRIMARY KEY,
  `username` varchar(255),
  `role` varchar(255),
  `created_at` timestamp
);

CREATE TABLE `posts` (
  `id` integer PRIMARY KEY,
  `title` varchar(255),
  `body` text COMMENT 'Content of the post',
  `user_id` integer NOT NULL,
  `status` varchar(255),
  `created_at` timestamp
);

ALTER TABLE `posts` ADD CONSTRAINT `user_posts` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `users` (`id`, `username`, `role`)
VALUES
  (0, 'Alice', 'admin'),
  (1, 'Bob', 'moderator'),
  (2, 'Candice', 'moderator'),
  (3, 'David', 'member');
INSERT INTO `posts` (`id`, `title`, `user_id`)
VALUES
  (0, 'Welcome to the forum!', 0),
  (1, 'Guidelines', 1),
  (2, 'Hello all!', 3);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;