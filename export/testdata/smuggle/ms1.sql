CREATE TABLE t (
  a nvarchar(50) CHECK (a <> N'x\' OR a = N')); DROP TABLE victim; CREATE TABLE pwned (i int); -- '),
  b nvarchar(50) DEFAULT (N'\' + N')); DROP TABLE victim; CREATE TABLE pwned2 (i int); -- ')
);
GO
