CREATE TABLE `ledger_entries` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `account_id` bigint unsigned NOT NULL,
  `amount_minor` bigint NOT NULL,
  `currency` char(3) NOT NULL DEFAULT 'USD',
  `memo` varchar(255) NOT NULL DEFAULT '',
  `rate` decimal(20,10) NOT NULL DEFAULT '1.0000000000',
  `posted` tinyint(1) NOT NULL DEFAULT '0',
  `posted_on` date DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `ledger_entries_account_id_posted_on_index` (`account_id`,`posted_on`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE `accounts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `type` set('asset','liability','equity','income','expense') NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;
ALTER TABLE `ledger_entries`
  ADD CONSTRAINT `ledger_entries_account_id_foreign` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE;
