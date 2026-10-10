CREATE TABLE `t` (
  `email` varchar(255)
);

CREATE INDEX `idx_email_lower` ON `t` ((lower(email)));
