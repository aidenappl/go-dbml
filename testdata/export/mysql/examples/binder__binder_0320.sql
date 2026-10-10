CREATE SCHEMA `types`;

CREATE TABLE `users` (
  `status` ENUM ('active', 'inactive') DEFAULT 'active'
);
