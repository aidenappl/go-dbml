CREATE TABLE `articles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `body` mediumtext NOT NULL,
  `lang` char(2) NOT NULL DEFAULT 'en',
  `views` int unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_title_lang` (`title`(150),`lang`),
  KEY `idx_views` (`views` DESC),
  KEY `idx_lang_hash` (`lang`) USING HASH,
  FULLTEXT KEY `ft_title_body` (`title`,`body`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE UNIQUE INDEX `uq_id_lang` ON `articles` (`id`, `lang`);
