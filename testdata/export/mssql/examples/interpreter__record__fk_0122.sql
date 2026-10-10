CREATE TABLE [products] (
  [id] int,
  [region] nvarchar(255),
  PRIMARY KEY ([id], [region])
)
GO

CREATE TABLE [categories] (
  [id] int,
  [region] nvarchar(255),
  PRIMARY KEY ([id], [region])
)
GO

CREATE TABLE [products_categories] (
  [products_id] int,
  [products_region] varchar,
  [categories_id] int,
  [categories_region] varchar,
  PRIMARY KEY ([products_id], [products_region], [categories_id], [categories_region])
);
GO

ALTER TABLE [products_categories] ADD FOREIGN KEY ([products_id], [products_region]) REFERENCES [products] ([id], [region]);
GO

ALTER TABLE [products_categories] ADD FOREIGN KEY ([categories_id], [categories_region]) REFERENCES [categories] ([id], [region]);
GO


-- Disable constraint checks for tables with data
ALTER TABLE [products] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [categories] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [products] ([id], [region])
VALUES
  (1, 'US'),
  (2, 'US');
GO
INSERT INTO [categories] ([id], [region])
VALUES
  (1, 'US'),
  (3, 'EU');
GO

-- Re-enable constraint checks
ALTER TABLE [products] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [categories] WITH CHECK CHECK CONSTRAINT ALL;
GO