CREATE TABLE `follows` (
  `following_user_id` integer,
  `followed_user_id` integer,
  `created_at` timestamp
);

CREATE TABLE `users` (
  `id` integer PRIMARY KEY,
  `username` varchar(255),
  `role` varchar(255),
  `created_at` timestamp
);

CREATE TABLE `posts` (
  `id` integer PRIMARY KEY,
  `title` varchar(255),
  `body` text COMMENT '# Heading 1
    code block
# Heading 2
  * 1
  * 2
  ',
  `user_id` integer,
  `status` varchar(255),
  `created_at` timestamp
);

ALTER TABLE `follows` COMMENT = '# Heading 1
    code block
# Heading 2
  * 1
  * 1
';

ALTER TABLE `users` COMMENT = '# Heading 1
    code block
# Heading 2
  * 1
  * 2
';

ALTER TABLE `posts` COMMENT = '# Heading 1
    code block
# Heading 2
  * 1
    * 1.1
    * 1.2
  * 2
';

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `follows` ADD FOREIGN KEY (`following_user_id`) REFERENCES `users` (`id`);

ALTER TABLE `follows` ADD FOREIGN KEY (`followed_user_id`) REFERENCES `users` (`id`);
