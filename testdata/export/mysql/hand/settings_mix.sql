CREATE TABLE `t1` (
  `id` integer UNIQUE PRIMARY KEY NOT NULL AUTO_INCREMENT DEFAULT 1,
  `name` varchar(255) DEFAULT 'x',
  `ratio` float DEFAULT -1.5,
  `flag` bool DEFAULT false,
  `ts` timestamp DEFAULT (now()),
  `arr` int[],
  `json` jsonb,
  PRIMARY KEY (`name`)
);

CREATE UNIQUE INDEX `idx_id` ON `t1` (`id`) USING HASH;

CREATE INDEX `t1_index_1` ON `t1` ((lower(name)), `name`);

ALTER TABLE `t1` COMMENT = 'n';
