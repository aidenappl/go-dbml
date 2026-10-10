CREATE TABLE `t` (
  `email` varchar(255)
);

CREATE UNIQUE INDEX `idx_email` ON `t` (`email`) USING BTREE;
