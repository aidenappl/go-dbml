CREATE SCHEMA `app`;

CREATE SCHEMA `other`;

CREATE TABLE `app`.`users` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `st` ENUM ('active', 'on hold') NOT NULL,
  `st2` ENUM ('a', 'b'),
  `cur` public.money_t,
  `cur2` other.money_t,
  `name` varchar(255) DEFAULT 'x' COMMENT 'Bob''s name'
);

CREATE TABLE `posts` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `user_id` int,
  `t` timestamp with time zone DEFAULT (now())
);

CREATE TABLE `other`.`tags` (
  `id` smallint PRIMARY KEY AUTO_INCREMENT,
  `name` My Type
);

ALTER TABLE `app`.`users` COMMENT = 'It''s users';

ALTER TABLE `posts` COMMENT = 'posts';

ALTER TABLE `posts` ADD FOREIGN KEY (`user_id`) REFERENCES `app`.`users` (`id`);

CREATE TABLE `app`.`users_tags` (
  `users_id` int,
  `tags_id` smallint,
  PRIMARY KEY (`users_id`, `tags_id`)
);

ALTER TABLE `app`.`users_tags` ADD FOREIGN KEY (`users_id`) REFERENCES `app`.`users` (`id`);

ALTER TABLE `app`.`users_tags` ADD FOREIGN KEY (`tags_id`) REFERENCES `other`.`tags` (`id`);


CREATE TABLE `posts_tags` (
  `posts_id` bigint,
  `tags_id` smallint,
  PRIMARY KEY (`posts_id`, `tags_id`)
);

ALTER TABLE `posts_tags` ADD FOREIGN KEY (`posts_id`) REFERENCES `posts` (`id`);

ALTER TABLE `posts_tags` ADD FOREIGN KEY (`tags_id`) REFERENCES `other`.`tags` (`id`);


ALTER TABLE `other`.`tags` ADD CONSTRAINT `r1` FOREIGN KEY (`id`) REFERENCES `posts` (`user_id`) ON DELETE CASCADE ON UPDATE SET NULL;
