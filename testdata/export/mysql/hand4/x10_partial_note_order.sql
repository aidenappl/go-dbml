CREATE TABLE `t` (
  `x` int,
  `b` int PRIMARY KEY,
  `a` text
);

CREATE TABLE `t2` (
  `b` int PRIMARY KEY,
  `a` int
);

CREATE UNIQUE INDEX `t_index_0` ON `t` (`b`);

CREATE INDEX `t_index_1` ON `t` (`a`);

CREATE INDEX `t_index_2` ON `t` (`x`);

CREATE INDEX `t2_index_3` ON `t2` (`a`);

CREATE UNIQUE INDEX `t2_index_4` ON `t2` (`b`);

ALTER TABLE `t` COMMENT = 'p2 note';

ALTER TABLE `t2` COMMENT = 'own';
