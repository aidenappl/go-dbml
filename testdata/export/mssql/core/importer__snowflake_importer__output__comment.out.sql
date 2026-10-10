CREATE SCHEMA [TEST_SCHEMA]
GO

CREATE TABLE [TEST_SCHEMA].[ONE] (
  [employee_id] INT NOT NULL,
  [first_name] VARCHAR(50) NOT NULL,
  [last_name] VARCHAR(50) NOT NULL,
  [email] VARCHAR(100),
  [hire_date] DATE,
  [department] VARCHAR(50),
  [salary] DECIMAL(10,2)
)
GO

CREATE TABLE [TEST_SCHEMA].[TWO] (
  [employee_id] INT NOT NULL,
  [first_name] VARCHAR(50) NOT NULL,
  [last_name] VARCHAR(50) NOT NULL,
  [email] VARCHAR(100),
  [hire_date] DATE,
  [department] VARCHAR(50),
  [salary] DECIMAL(10,2)
)
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Unique identifier for each employee',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'ONE',
@level2type = N'Column', @level2name = 'employee_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee''s first name',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'ONE',
@level2type = N'Column', @level2name = 'first_name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee''s last name',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'ONE',
@level2type = N'Column', @level2name = 'last_name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee''s work email address',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'ONE',
@level2type = N'Column', @level2name = 'email';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date when the employee was hired',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'ONE',
@level2type = N'Column', @level2name = 'hire_date';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Department where the employee works',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'ONE',
@level2type = N'Column', @level2name = 'department';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee''s annual salary',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'ONE',
@level2type = N'Column', @level2name = 'salary';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Another comment for employees table',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'TWO';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Unique identifier for each employee',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'TWO',
@level2type = N'Column', @level2name = 'employee_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee''s first name',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'TWO',
@level2type = N'Column', @level2name = 'first_name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee''s last name',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'TWO',
@level2type = N'Column', @level2name = 'last_name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee''s work email address',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'TWO',
@level2type = N'Column', @level2name = 'email';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date when the employee was hired',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'TWO',
@level2type = N'Column', @level2name = 'hire_date';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Department where the employee works',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'TWO',
@level2type = N'Column', @level2name = 'department';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee''s annual salary',
@level0type = N'Schema', @level0name = 'TEST_SCHEMA',
@level1type = N'Table',  @level1name = 'TWO',
@level2type = N'Column', @level2name = 'salary';
GO
