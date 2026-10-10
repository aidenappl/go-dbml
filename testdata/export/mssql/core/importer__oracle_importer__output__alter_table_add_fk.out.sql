CREATE TABLE [departments] (
  [dept_id] NUMBER(10) PRIMARY KEY,
  [dept_name] VARCHAR2(100),
  [location_id] NUMBER(10),
  [manager_id] NUMBER(10)
)
GO

CREATE TABLE [locations] (
  [location_id] NUMBER(10) PRIMARY KEY,
  [city] VARCHAR2(100),
  [country_code] CHAR(2)
)
GO

CREATE TABLE [countries] (
  [country_code] CHAR(2) PRIMARY KEY,
  [country_name] VARCHAR2(100)
)
GO

CREATE TABLE [employees] (
  [emp_id] NUMBER(10) PRIMARY KEY,
  [emp_name] VARCHAR2(100),
  [department_id] NUMBER(10),
  [manager_id] NUMBER(10),
  [hire_date] DATE,
  [location_id] NUMBER(10)
)
GO

CREATE TABLE [projects] (
  [project_id] NUMBER(10) PRIMARY KEY,
  [project_name] VARCHAR2(200),
  [dept_id] NUMBER(10),
  [lead_emp_id] NUMBER(10),
  [backup_emp_id] NUMBER(10)
)
GO

CREATE TABLE [assignments] (
  [assignment_id] NUMBER(10) PRIMARY KEY,
  [emp_id] NUMBER(10),
  [project_id] NUMBER(10),
  [start_date] DATE
)
GO

CREATE TABLE [project_tasks] (
  [project_id] NUMBER(10),
  [task_id] NUMBER(10),
  [task_name] VARCHAR2(200),
  [assigned_emp_id] NUMBER(10),
  PRIMARY KEY ([project_id], [task_id])
)
GO

CREATE TABLE [task_hours] (
  [hour_id] NUMBER(10) PRIMARY KEY,
  [project_id] NUMBER(10),
  [task_id] NUMBER(10),
  [hours_worked] NUMBER(5,2)
)
GO

ALTER TABLE [employees] ADD CONSTRAINT [fk_emp_dept] FOREIGN KEY ([department_id]) REFERENCES [departments] ([dept_id])
GO

ALTER TABLE [departments] ADD CONSTRAINT [fk_dept_location] FOREIGN KEY ([location_id]) REFERENCES [locations] ([location_id])
GO

ALTER TABLE [locations] ADD CONSTRAINT [fk_loc_country] FOREIGN KEY ([country_code]) REFERENCES [countries] ([country_code])
GO

ALTER TABLE [employees] ADD CONSTRAINT [fk_emp_manager] FOREIGN KEY ([manager_id]) REFERENCES [employees] ([emp_id])
GO

ALTER TABLE [departments] ADD CONSTRAINT [fk_dept_manager] FOREIGN KEY ([manager_id]) REFERENCES [employees] ([emp_id])
GO

ALTER TABLE [projects] ADD FOREIGN KEY ([dept_id]) REFERENCES [departments] ([dept_id])
GO

ALTER TABLE [assignments] ADD FOREIGN KEY ([emp_id]) REFERENCES [employees] ([emp_id])
GO

ALTER TABLE [assignments] ADD FOREIGN KEY ([project_id]) REFERENCES [projects] ([project_id])
GO

ALTER TABLE [task_hours] ADD CONSTRAINT [fk_task_hours_project_task] FOREIGN KEY ([project_id], [task_id]) REFERENCES [project_tasks] ([project_id], [task_id])
GO

ALTER TABLE [projects] ADD CONSTRAINT [fk_proj_lead_emp] FOREIGN KEY ([lead_emp_id]) REFERENCES [employees] ([emp_id])
GO

ALTER TABLE [projects] ADD CONSTRAINT [fk_proj_backup_emp] FOREIGN KEY ([backup_emp_id]) REFERENCES [employees] ([emp_id])
GO

ALTER TABLE [project_tasks] ADD CONSTRAINT [fk_task_assigned_emp] FOREIGN KEY ([assigned_emp_id]) REFERENCES [employees] ([emp_id])
GO
