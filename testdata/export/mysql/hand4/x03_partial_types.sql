CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `created_by` s.s.user_ref,
  `kind` ENUM ('x', 'y'),
  `amount` decimal(10,2) CHECK (amount > 0) DEFAULT 1.5,
  CONSTRAINT `users.small` CHECK (amount < 100)
);

CREATE TABLE `posts` (
  `id` int PRIMARY KEY,
  `created_by` s.s.user_ref,
  `kind` ENUM ('x', 'y'),
  `amount` decimal(10,2) CHECK (amount > 0) DEFAULT 1.5,
  `note2` text COMMENT 'n',
  CONSTRAINT `posts.small` CHECK (amount < 100)
);

CREATE INDEX `by_creator` ON `users` (`created_by`);

CREATE INDEX `by_creator` ON `posts` (`created_by`);

ALTER TABLE `users` COMMENT = 'audit';

ALTER TABLE `posts` COMMENT = 'audit';

ALTER TABLE `users` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `posts` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);
