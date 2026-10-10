CREATE SCHEMA [demographic]
GO

CREATE TABLE [user] (
  [name] text,
  [id] int,
  [status] nvarchar(255) CHECK ([status] IN ('active', 'churned', 'inactive')) DEFAULT 'active',
  [gender] nvarchar(255) CHECK ([gender] IN ('male', 'female')) DEFAULT 'male',
  [age_type] nvarchar(255) CHECK ([age_type] IN ('toddler', 'children', 'teenager', 'young_adult', 'elder')) DEFAULT 'young_adult'
)
GO
