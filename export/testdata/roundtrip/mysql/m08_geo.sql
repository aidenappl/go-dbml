CREATE TABLE `countries` (
  `iso2` char(2) NOT NULL,
  `name` varchar(80) NOT NULL,
  PRIMARY KEY (`iso2`)
) ENGINE=InnoDB;
CREATE TABLE `cities` (
  `id` mediumint unsigned NOT NULL AUTO_INCREMENT,
  `country_iso2` char(2) NOT NULL,
  `name` varchar(120) NOT NULL,
  `population` int unsigned DEFAULT NULL,
  `location` point NOT NULL /*!80003 SRID 4326 */,
  PRIMARY KEY (`id`),
  KEY `country_iso2` (`country_iso2`),
  SPATIAL KEY `location` (`location`),
  FULLTEXT KEY `ft_name` (`name`),
  CONSTRAINT `cities_ibfk_1` FOREIGN KEY (`country_iso2`) REFERENCES `countries` (`iso2`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB;
