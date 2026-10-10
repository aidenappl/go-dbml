CREATE TABLE `t` (
  `c` int COMMENT 'it''s escaped 
 newline'
);

ALTER TABLE `t` COMMENT = 'table-level note
';
