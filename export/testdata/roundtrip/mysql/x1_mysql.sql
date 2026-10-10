/*!40101 SET NAMES utf8mb4 */;
DROP TABLE IF EXISTS `acct`;
CREATE TABLE `acct` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT 'O\'Brien''s \\ "x"',
  `kind` enum('a''b','c\\d','e,f') NOT NULL DEFAULT 'a''b',
  `bal` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT 'it\'s \"money\"\nline2',
  `full` varchar(200) GENERATED ALWAYS AS (concat(`name`,' ',`kind`)) VIRTUAL,
  `created` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `flags` set('x','y') DEFAULT NULL,
  `parent_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `ck_bal` CHECK (`bal` >= 0 and `name` <> 'it''s'),
  CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `acct` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Accounts – “quoted” \'t\'';
