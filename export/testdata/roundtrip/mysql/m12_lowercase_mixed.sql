create table if not exists `app_settings` (
  `key` varchar(100) not null primary key,
  `value` text null,
  `updated_by` int null default null,
  `is_secret` bit(1) not null default b'0'
) engine=innodb default charset=utf8;
create table `feature_flags` (
  `id` int unsigned auto_increment,
  `name` varchar(100) not null,
  `rollout_pct` tinyint unsigned not null default 0,
  `config` longtext character set utf8mb4 collate utf8mb4_bin default null check (json_valid(`config`)),
  primary key (`id`),
  unique key (`name`)
) engine=innodb;
insert into `feature_flags` (`name`,`rollout_pct`) values ('new_ui', 10), ('beta', 0);
