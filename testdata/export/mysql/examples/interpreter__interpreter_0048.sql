CREATE TABLE `t` (
  `first_name` varchar(255),
  `last_name` varchar(255)
);

CREATE INDEX `idx_fullname` ON `t` (`first_name`, `last_name`);
