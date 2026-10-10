CREATE SCHEMA `v2`;

CREATE TABLE `v2`.`jobs` (
  `status` ENUM ('created', 'running', 'done', 'failure') DEFAULT 'created'
);
