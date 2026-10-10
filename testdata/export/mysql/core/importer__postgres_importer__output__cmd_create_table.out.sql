CREATE TABLE `a` (
  `id` integer UNIQUE PRIMARY KEY DEFAULT 1,
  `name` varchar(255) CHECK ("id" > -1)
);

CREATE TABLE `b` (
  `id` uuid,
  `name` varchar(255) NOT NULL CHECK (name <> ''),
  `email` varchar(55),
  PRIMARY KEY (`email`)
);

CREATE TABLE `c` (
  `index` CHARACTER VARYING(255) UNIQUE,
  `phone` SERIAL AUTO_INCREMENT,
  `time` timestamp
);

CREATE TABLE `circles` (
  `c` circle
);

CREATE TABLE `foo` (
  `bar` text[],
  `bar2` int[1],
  `bar3` int[2][3],
  `bar4` int[],
  `bar5` int[2],
  `bar6` text[8],
  `bar7` text[100],
  `bar8` time(2)[],
  `bar9` time(1)[1],
  `bar10` time(1)[],
  `bar11` time[5],
  `bar12` timestamp(2)[10][2][5],
  `bar13` character varying[],
  `bar14` character varying(25)[][2][],
  `bar15` character varying[76]
);

CREATE UNIQUE INDEX `b_index_0` ON `b` (`name`);

ALTER TABLE `b` ADD FOREIGN KEY (`id`) REFERENCES `a` (`id`) ON DELETE NO ACTION;

ALTER TABLE `c` ADD CONSTRAINT `c_time_constraint` FOREIGN KEY (`time`) REFERENCES `a` (`name`);
