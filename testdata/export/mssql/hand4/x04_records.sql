CREATE TABLE [t] (
  [id] int PRIMARY KEY,
  [r] real,
  [f] float,
  [b] boolean,
  [bt] bit,
  [s] varchar(3),
  [d] date,
  [ts] timestamp,
  [j] json,
  [by] bytea,
  [en] nvarchar(255) CHECK ([en] IN ('a', 'b')),
  [x] decimal(4,1)
)
GO

CREATE TABLE [u] (
  [id] int PRIMARY KEY,
  [name] text
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [t] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [u] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [t] ([id], [r], [f], [b], [bt], [s], [d], [ts], [j], [by], [en], [x])
VALUES
  (1, 1.5, -2, 1, 1, 'it''s', '2024-01-02', '2024-01-02T03:04:05', CAST('{"a":1}' AS json), 0xff, 'a', 123.4),
  (2, NULL, 1000, 0, 0, 'a\b', NULL, NULL, now(), NULL, 'b', 1),
  (3, NULL, NULL, 1, 0, 'toolong', 'bad', 'bad', CAST('3' AS json), 0xx, 'zz', 99999);
GO

-- Re-enable constraint checks
ALTER TABLE [t] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [u] WITH CHECK CHECK CONSTRAINT ALL;
GO