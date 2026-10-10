CREATE TABLE t (
  a text CHECK (a <> 'a\' OR a = ')); DROP TABLE victim; CREATE TABLE pwned (i int); -- '),
  b text DEFAULT ('q\' || ')); DROP TABLE victim; CREATE TABLE pwned2 (i int); -- '),
  c text CHECK (c <> E'\')); DROP TABLE victim; CREATE TABLE pwned3 (i int); -- '),
  d text CHECK (d <> $$')); DROP TABLE victim; CREATE TABLE pwned4 (i int); -- $$)
);
CREATE INDEX ix ON t ((a || 'p\' || '); DROP TABLE victim; -- '));
