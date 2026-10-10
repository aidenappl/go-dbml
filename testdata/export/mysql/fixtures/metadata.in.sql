CREATE SCHEMA `sales`;

CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

CREATE TABLE `posts` (
  `id` int PRIMARY KEY,
  `user_id` int
);

CREATE TABLE `sales`.`orders` (
  `id` int PRIMARY KEY
);

ALTER TABLE `users` COMMENT = 'this will override';

ALTER TABLE `posts` COMMENT = '# Heading 1
    code block
# Heading 2
';
