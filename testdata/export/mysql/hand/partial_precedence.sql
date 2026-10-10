CREATE TABLE `users` (
  `name` varchar(255),
  `created_at` datetime NOT NULL,
  `updated_at` datetime,
  `id` bigint COMMENT 'local',
  CONSTRAINT `users.chk_upd` CHECK (updated_at >= created_at)
);

CREATE INDEX `users_index_0` ON `users` (`created_at`);

ALTER TABLE `users` COMMENT = 'base note';
