CREATE TABLE `warehouses` (
  `id` smallint unsigned NOT NULL AUTO_INCREMENT,
  `code` char(4) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  `name` varchar(100) NOT NULL,
  `lat` double DEFAULT NULL,
  `lng` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `stock_levels` (
  `warehouse_id` smallint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL,
  `on_hand` int NOT NULL DEFAULT '0',
  `reserved` int NOT NULL DEFAULT '0',
  `available` int GENERATED ALWAYS AS ((`on_hand` - `reserved`)) VIRTUAL,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`warehouse_id`,`product_id`),
  KEY `idx_product` (`product_id`),
  CONSTRAINT `stock_levels_chk_1` CHECK ((`on_hand` >= 0)),
  CONSTRAINT `fk_stock_wh` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
