CREATE TABLE `t` (
  `first_name` varchar(255),
  `last_name` varchar(255)
);

CREATE INDEX `t_index_0` ON `t` ((lower(first_name)), `last_name`);
