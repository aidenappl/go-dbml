CREATE TABLE `tags` (
  `id` int PRIMARY KEY
);

CREATE TABLE `posts` (
  `id` int PRIMARY KEY
);

CREATE TABLE `tags_posts` (
  `tags_id` int,
  `posts_id` int,
  PRIMARY KEY (`tags_id`, `posts_id`)
);

ALTER TABLE `tags_posts` ADD FOREIGN KEY (`tags_id`) REFERENCES `tags` (`id`);

ALTER TABLE `tags_posts` ADD FOREIGN KEY (`posts_id`) REFERENCES `posts` (`id`);


-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `tags` (`id`)
VALUES
  (NULL),
  (1);
INSERT INTO `posts` (`id`)
VALUES
  (NULL),
  (10);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;