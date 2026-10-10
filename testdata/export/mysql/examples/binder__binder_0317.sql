CREATE SCHEMA `types`;

CREATE TABLE `users` (
  `id` int,
  `status` ENUM ('active', 'inactive')
);
