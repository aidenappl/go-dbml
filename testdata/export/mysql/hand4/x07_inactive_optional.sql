CREATE TABLE `users` (
  `id` integer PRIMARY KEY
);

CREATE TABLE `posts` (
  `user_id` integer NOT NULL,
  `editor_id` int,
  `owner_id` int UNIQUE,
  `other_id` int
);

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `users` ADD FOREIGN KEY (`id`) REFERENCES `posts` (`editor_id`);

ALTER TABLE `posts` ADD FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`);

CREATE TABLE `posts_users` (
  `posts_other_id` int,
  `users_id` integer,
  PRIMARY KEY (`posts_other_id`, `users_id`)
);

ALTER TABLE `posts_users` ADD FOREIGN KEY (`posts_other_id`) REFERENCES `posts` (`other_id`);

ALTER TABLE `posts_users` ADD FOREIGN KEY (`users_id`) REFERENCES `users` (`id`);

