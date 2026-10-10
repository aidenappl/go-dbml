CREATE SCHEMA `s1`;

CREATE TABLE `s1`.`t1` (
  `a` int,
  `b` int,
  `c` text,
  PRIMARY KEY (`a`, `b`)
);

CREATE TABLE `s1`.`t3` (
  `a` int PRIMARY KEY
);

CREATE TABLE `t2` (
  `a` int
);

CREATE UNIQUE INDEX ``s1`.t1_index_0` ON `s1`.`t1` (`a`) USING HASH;

CREATE INDEX `expr_ix` ON `s1`.`t1` ((lower(c)), `b`) USING BTREE;

CREATE INDEX ``s1`.t1_index_2` ON `s1`.`t1` ((a + b));

CREATE INDEX ``s1`.t1_index_3` ON `s1`.`t1` (`c`);

CREATE INDEX ``s1`.t3_index_4` ON `s1`.`t3` (`a`);

CREATE INDEX `t2_index_0` ON `t2` (`a`);

CREATE UNIQUE INDEX `t2_index_1` ON `t2` (`a`);
