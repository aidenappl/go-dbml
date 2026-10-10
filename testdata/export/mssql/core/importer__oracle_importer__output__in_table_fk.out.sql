CREATE TABLE [departments] (
  [dept_id] NUMBER(10) PRIMARY KEY
)
GO

CREATE TABLE [customers] (
  [cust_id] NUMBER(10) PRIMARY KEY
)
GO

CREATE TABLE [employees] (
  [emp_id] NUMBER(10),
  [emp_name] VARCHAR2(100),
  [dept_id] NUMBER(10)
)
GO

CREATE TABLE [orders] (
  [order_id] NUMBER(10),
  [customer_id] NUMBER(10),
  [order_date] DATE
)
GO

CREATE TABLE [shipments] (
  [ship_id] NUMBER(10),
  [customer_id] NUMBER(10) NOT NULL,
  [ship_date] DATE
)
GO

CREATE TABLE [complex_parent] (
  [id] NUMBER(10),
  [code] VARCHAR2(20),
  PRIMARY KEY ([id], [code])
)
GO

CREATE TABLE [inventory] (
  [item_id] NUMBER(10) PRIMARY KEY
)
GO

CREATE TABLE [complex_child] (
  [child_id] NUMBER(10),
  [parent_id] NUMBER(10),
  [parent_code] VARCHAR2(20)
)
GO

CREATE TABLE [order_details] (
  [detail_id] NUMBER(10),
  [order_id] NUMBER(10),
  [item_id] NUMBER(10)
)
GO

CREATE TABLE [assignments] (
  [assignment_id] NUMBER(10),
  [emp_id] NUMBER(10),
  [dept_id] NUMBER(10)
)
GO

CREATE TABLE [project_members] (
  [member_id] NUMBER(10),
  [emp_id] NUMBER(10),
  [dept_id] NUMBER(10)
)
GO

CREATE TABLE [test_datatypes] (
  [col_number] NUMBER(10,2) UNIQUE,
  [col_varchar2] VARCHAR2(100) UNIQUE
)
GO

CREATE TABLE [test_datatypes_fk] (
  [id] NUMBER(10),
  [col_number] NUMBER(10,2),
  [col_varchar2] VARCHAR2(100),
  [col_char] CHAR(10),
  [col_date] DATE
)
GO

CREATE TABLE [employees_self] (
  [emp_id] NUMBER(10) PRIMARY KEY,
  [manager_id] NUMBER(10),
  [emp_name] VARCHAR2(100)
)
GO

CREATE TABLE [order_history] (
  [history_id] NUMBER(10),
  [order_id] NUMBER(10),
  [history_date] DATE
)
GO

CREATE TABLE [shipment_details] (
  [detail_id] NUMBER(10),
  [ship_id] NUMBER(10),
  [item_id] NUMBER(10)
)
GO

CREATE TABLE [project_assignments] (
  [assign_id] NUMBER(10),
  [emp_id] NUMBER(10),
  [project_id] NUMBER(10)
)
GO

ALTER TABLE [employees] ADD FOREIGN KEY ([dept_id]) REFERENCES [departments] ([dept_id])
GO

ALTER TABLE [orders] ADD CONSTRAINT [fk_order_cust] FOREIGN KEY ([customer_id]) REFERENCES [customers] ([cust_id])
GO

ALTER TABLE [shipments] ADD CONSTRAINT [fk_ship_cust] FOREIGN KEY ([customer_id]) REFERENCES [customers] ([cust_id])
GO

ALTER TABLE [complex_child] ADD CONSTRAINT [fk_child_complex] FOREIGN KEY ([parent_id], [parent_code]) REFERENCES [complex_parent] ([id], [code])
GO

ALTER TABLE [order_details] ADD CONSTRAINT [fk_detail_item] FOREIGN KEY ([item_id]) REFERENCES [inventory] ([item_id])
GO

ALTER TABLE [assignments] ADD CONSTRAINT [fk_assign_dept] FOREIGN KEY ([dept_id]) REFERENCES [departments] ([dept_id])
GO

ALTER TABLE [project_members] ADD CONSTRAINT [fk_member_dept] FOREIGN KEY ([dept_id]) REFERENCES [departments] ([dept_id])
GO

ALTER TABLE [test_datatypes_fk] ADD CONSTRAINT [fk_dt_number] FOREIGN KEY ([col_number]) REFERENCES [test_datatypes] ([col_number])
GO

ALTER TABLE [test_datatypes_fk] ADD CONSTRAINT [fk_dt_varchar2] FOREIGN KEY ([col_varchar2]) REFERENCES [test_datatypes] ([col_varchar2])
GO

ALTER TABLE [employees_self] ADD CONSTRAINT [fk_self_manager] FOREIGN KEY ([manager_id]) REFERENCES [employees_self] ([emp_id])
GO

ALTER TABLE [order_history] ADD CONSTRAINT [fk_very_long_constraint_name_for_testing_parser_limits] FOREIGN KEY ([order_id]) REFERENCES [orders] ([order_id])
GO

ALTER TABLE [shipment_details] ADD CONSTRAINT [fk_ship_detail] FOREIGN KEY ([ship_id]) REFERENCES [shipments] ([ship_id])
GO

ALTER TABLE [project_assignments] ADD CONSTRAINT [fk_assign_emp] FOREIGN KEY ([emp_id]) REFERENCES [employees] ([emp_id])
GO
