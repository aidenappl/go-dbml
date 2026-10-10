CREATE TABLE t (
  a varchar(50) CHECK (a <> 'x\')); DROP TABLE victim; CREATE TABLE pwned (i int); -- '),
  b varchar(50) DEFAULT (concat('y\'', ')); DROP TABLE victim; CREATE TABLE pwned2 (i int); -- ')),
  c varchar(50),
  CONSTRAINT ck CHECK (c <> 'z\''),
  KEY ix ((concat(c, 'w\''))),
  d int COMMENT 'ok\'); DROP TABLE victim; -- '
);
