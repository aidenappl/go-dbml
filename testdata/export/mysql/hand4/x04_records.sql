CREATE TABLE `t` (
  `id` int PRIMARY KEY,
  `r` real,
  `f` float,
  `b` boolean,
  `bt` bit,
  `s` varchar(3),
  `d` date,
  `ts` timestamp,
  `j` json,
  `by` bytea,
  `en` ENUM ('a', 'b'),
  `x` decimal(4,1)
);

CREATE TABLE `u` (
  `id` int PRIMARY KEY,
  `name` text
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `t` (`id`, `r`, `f`, `b`, `bt`, `s`, `d`, `ts`, `j`, `by`, `en`, `x`)
VALUES
  (1, 1.5, -2, TRUE, TRUE, 'it''s', '2024-01-02', '2024-01-02T03:04:05', CAST('{"a":1}' AS json), 'ff', 'a', 123.4),
  (2, NULL, 1000, FALSE, FALSE, 'a\\b', NULL, NULL, now(), NULL, 'b', 1),
  (3, NULL, NULL, TRUE, FALSE, 'toolong', 'bad', 'bad', CAST('3' AS json), 'x', 'zz', 99999);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;