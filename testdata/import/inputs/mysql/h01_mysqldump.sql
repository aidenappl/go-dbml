-- MySQL dump 10.13  Distrib 8.0.36, for Linux (x86_64)
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
DROP TABLE IF EXISTS `accounts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
CREATE TABLE `accounts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(320) COLLATE utf8mb4_unicode_ci NOT NULL,
  `display_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'shown in UI',
  `status` enum('active','suspended','deleted') NOT NULL DEFAULT 'active',
  `balance` decimal(12,2) NOT NULL DEFAULT '0.00',
  `flags` set('a','b') DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_admin` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `accounts_email_unique` (`email`),
  KEY `accounts_status_index` (`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User accounts';
/*!40101 SET character_set_client = @saved_cs_client */;

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
  `id` char(36) NOT NULL,
  `account_id` bigint unsigned NOT NULL,
  `ip` varbinary(16) DEFAULT NULL,
  `expires_at` datetime(6) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `sessions_account_id_foreign` (`account_id`),
  CONSTRAINT `sessions_account_id_foreign` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `accounts` WRITE;
/*!40000 ALTER TABLE `accounts` DISABLE KEYS */;
INSERT INTO `accounts` VALUES (1,'a@b.c','Ann','active',10.50,NULL,'2024-01-01 00:00:00','2024-01-01 00:00:00',1),(2,'d@e.f',NULL,'deleted',0.00,'a,b','2024-02-01 10:00:00','2024-02-01 10:00:00',0);
/*!40000 ALTER TABLE `accounts` ENABLE KEYS */;
UNLOCK TABLES;
