CREATE TABLE `devices` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'surrogate key',
  `serial` varchar(32) NOT NULL COMMENT 'manufacturer serial, e.g. "SN-001"',
  `kind` enum('sensor','gateway') DEFAULT 'sensor' COMMENT 'device kind',
  `fw_version` varchar(16) DEFAULT '1.0.0',
  `battery` float(5,2) DEFAULT NULL,
  `last_seen` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `serial` (`serial`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='IoT devices: it''s tracked';
CREATE TABLE `readings` (
  `device_id` int NOT NULL,
  `ts` datetime(6) NOT NULL,
  `value` decimal(18,6) NOT NULL,
  PRIMARY KEY (`device_id`,`ts`),
  CONSTRAINT `readings_ibfk_1` FOREIGN KEY (`device_id`) REFERENCES `devices` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
/*!50100 PARTITION BY RANGE (year(`ts`))
(PARTITION p2024 VALUES LESS THAN (2025) ENGINE = InnoDB,
 PARTITION pmax VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */;
