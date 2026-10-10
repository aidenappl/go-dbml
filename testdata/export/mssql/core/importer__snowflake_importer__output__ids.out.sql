CREATE TABLE [alert] (
  [i] int
)
GO

CREATE TABLE [condition] (
  [i] int
)
GO

CREATE TABLE [dt_t] (
  [dynamic] int,
  [downstream] int,
  [target_lag] int
)
GO

CREATE TABLE [wh_id] (
  [MAX_CONCURRENCY_LEVEL] int,
  [WAREHOUSE_TYPE] int
)
GO

CREATE TABLE [CAST] (
  [KEY] int,
  [TRY_CAST] int
)
GO

CREATE TABLE [DIRECTION] (
  [LENGTH] INT,
  [LANGUAGE] INT
)
GO

CREATE TABLE [ROLES] (
  [USERS] INT,
  [WAREHOUSES] int,
  [INTEGRATIONS] int,
  [DATABASES] int,
  [SCHEMAS] int
)
GO

CREATE TABLE [TABLES] (
  [VIEWS] INT,
  [STAGES] int,
  [STREAMS] int,
  [TASKS] int,
  [ALERTS] int
)
GO

CREATE TABLE [ACTION] (
  [MODE] int
)
GO

CREATE TABLE [interval] (
  [internal] int
)
GO
