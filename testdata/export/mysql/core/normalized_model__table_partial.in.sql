CREATE SCHEMA `a`;

CREATE TABLE `a` (
  `id` int,
  `id2` int
);

CREATE TABLE `users` (
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp DEFAULT (now()),
  `name` varchar(255),
  `email` varchar(255),
  `id` int,
  `delete_status` boolean NOT NULL,
  `deleted_at` timestamp DEFAULT (now()),
  `refColumn1` bit,
  `refColumn2` bool,
  PRIMARY KEY (`id`)
);

CREATE TABLE `testColWithTypeHaveArgs` (
  `id` int,
  `callExpColInTest1` int(4) COMMENT 'simple column with call expression as type',
  `callExpColInTest2` decimal(10,2) DEFAULT 0.1 COMMENT 'column with type that have arguments and default value being a float number'
);

CREATE TABLE `table_no_column_with_partial` (
  `email` varchar(255) UNIQUE
);

CREATE TABLE `a`.`a` (
  `id` int,
  `id2` int
);

CREATE UNIQUE INDEX `users_index_0` ON `users` (`email`);

CREATE UNIQUE INDEX `users_index_1` ON `users` (`email`, `id`);

CREATE UNIQUE INDEX `users_index_2` ON `users` (`email`);

CREATE UNIQUE INDEX `table_no_column_with_partial_index_3` ON `table_no_column_with_partial` (`email`);

ALTER TABLE `users` ADD FOREIGN KEY (`refColumn1`) REFERENCES `a` (`id2`);

ALTER TABLE `users` ADD FOREIGN KEY (`refColumn2`) REFERENCES `a`.`a` (`id`);
