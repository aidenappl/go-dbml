CREATE TABLE [employees] (
  [employee_id] number PRIMARY KEY IDENTITY(1, 1),
  [first_name] nvarchar2(50) NOT NULL DEFAULT 'Unknown',
  [last_name] nvarchar2(50) NOT NULL,
  [middle_name] varchar2(50) DEFAULT '',
  [department] varchar2(30) UNIQUE NOT NULL DEFAULT 'General',
  [years_of_service] integer NOT NULL DEFAULT (0),
  [employee_rank] number(2) UNIQUE NOT NULL DEFAULT (1),
  [hire_date] date NOT NULL DEFAULT (SYSDATE),
  [last_review_date] timestamp DEFAULT (current_timestamp),
  [termination_date] date,
  [birth_date] date NOT NULL,
  [is_active] number(1) NOT NULL DEFAULT (1),
  [is_manager] number(1) DEFAULT (0),
  [notes] clob,
  [profile_data] nclob DEFAULT '{}'
)
GO

CREATE TABLE [projects] (
  [project_id] int PRIMARY KEY IDENTITY(1, 1),
  [project_code] varchar2(20) UNIQUE NOT NULL,
  [project_name] nvarchar2(100) NOT NULL DEFAULT 'Unnamed Project',
  [start_date] date NOT NULL DEFAULT (SYSDATE),
  [end_date] date,
  [status] varchar2(20) NOT NULL DEFAULT 'PLANNED',
  [priority] integer UNIQUE DEFAULT (999)
)
GO

CREATE TABLE [constraint_orderings] (
  [field1] varchar2(50) UNIQUE NOT NULL DEFAULT 'test1',
  [field2] varchar2(50) UNIQUE NOT NULL DEFAULT 'test2',
  [field3] varchar2(50) UNIQUE NOT NULL DEFAULT 'test3',
  [field4] varchar2(50) UNIQUE NOT NULL DEFAULT 'test4',
  [field5] varchar2(50) NOT NULL DEFAULT 'test5',
  [field6] varchar2(50) UNIQUE DEFAULT 'test6',
  [field7] int PRIMARY KEY NOT NULL DEFAULT (1)
)
GO
