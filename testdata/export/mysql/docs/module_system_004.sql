CREATE SCHEMA `auth`;

CREATE TABLE `auth`.`users` (
  `id` int PRIMARY KEY,
  `email` varchar(255)
);

CREATE TABLE `auth`.`roles` (
  `id` int PRIMARY KEY,
  `name` varchar(255)
);

CREATE TABLE `auth`.`sessions` (
  `id` int PRIMARY KEY
);
