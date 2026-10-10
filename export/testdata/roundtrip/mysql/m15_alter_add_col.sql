CREATE TABLE `invoices` (
  `id` int NOT NULL AUTO_INCREMENT,
  `number` varchar(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;
ALTER TABLE `invoices` ADD COLUMN `due_on` date DEFAULT NULL AFTER `number`;
ALTER TABLE `invoices` ADD COLUMN `paid` tinyint(1) NOT NULL DEFAULT '0';
ALTER TABLE `invoices` ADD INDEX `idx_due` (`due_on`);
ALTER TABLE `invoices` MODIFY `number` varchar(40) NOT NULL;
CREATE TABLE `invoice lines` (
  `invoice id` int NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  KEY `invoice id` (`invoice id`),
  CONSTRAINT `fk lines` FOREIGN KEY (`invoice id`) REFERENCES `invoices` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT
) ENGINE=InnoDB;
