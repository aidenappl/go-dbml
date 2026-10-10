CREATE TABLE `bad_faith_users` (
  `id` int(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` int(11),
  `created` int(11)
);

CREATE TABLE `booking` (
  `id` int(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id_one` int(11),
  `user_id_two` int(11),
  `timezone` varchar(512) NOT NULL DEFAULT '',
  `date` varchar(255),
  `spot` varchar(255),
  `address` text,
  `time` varchar(255),
  `is_accept` int(11),
  `message_id` text,
  `channel_id` text,
  `privateChat` varchar(255),
  `cron_status_before` int(11) DEFAULT '0',
  `cron_status_during` int(11) DEFAULT '0',
  `cron_status_after` int(11) DEFAULT '0',
  `before_time_3hrs` varchar(100),
  `before_time_20mins` varchar(100),
  `after_time_8hrs` varchar(100),
  `user_one_status` tinyint(2) NOT NULL DEFAULT '0',
  `user_two_status` tinyint(2) NOT NULL DEFAULT '0',
  `feedback_notification` int(11) NOT NULL DEFAULT '0',
  `feedback_notification_two` int(11) NOT NULL DEFAULT '0',
  `Created_at` varchar(255) NOT NULL DEFAULT ''
);

CREATE TABLE `test0` (
  `id` int(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `gid` int(11)
);

CREATE TABLE `test1` (
  `id` int(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `gid` int(11)
);

CREATE TABLE `test2` (
  `id` int(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `gid` int(11)
);

CREATE TABLE `test3` (
  `id` int(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `gid` int(11)
);

CREATE INDEX `user_id` ON `bad_faith_users` (`user_id`) USING BTREE;

CREATE INDEX `venue` ON `booking` (`user_id_two`) USING BTREE;

CREATE INDEX `user_id_one` ON `booking` (`user_id_one`) USING HASH;

CREATE INDEX `gid` ON `test0` (`gid`) USING HASH;

CREATE INDEX `gid` ON `test1` (`gid`) USING HASH;

CREATE INDEX `gid` ON `test2` (`gid`) USING BTREE;

CREATE INDEX `gid` ON `test3` (`gid`) USING BTREE;
