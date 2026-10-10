CREATE SCHEMA [HP_TEST]
GO

CREATE SCHEMA [TEST]
GO

CREATE SCHEMA [TEST_2024_06_28_10_01_ZQMIMY]
GO

CREATE TABLE [HP_TEST].[CUSTOMERS] (
  [CUSTOMER_ID] NUMBER(38, 0) PRIMARY KEY NOT NULL IDENTITY(1, 1),
  [FIRST_NAME] VARCHAR(50) NOT NULL DEFAULT 'John',
  [LAST_NAME] VARCHAR(50) NOT NULL DEFAULT 'Doe',
  [EMAIL] VARCHAR(100),
  [SECOND_EMAIL] VARCHAR(100),
  [UNIQUE_NAME] VARCHAR(100) UNIQUE NOT NULL DEFAULT (CONCAT(FIRST_NAME, LAST_NAME, EMAIL)),
  [IS_VIP] BOOLEAN DEFAULT (false),
  [IS_SOMETHING] BOOLEAN,
  [CREATED_AT] TIMESTAMP_NTZ(9) DEFAULT (CURRENT_TIMESTAMP())
)
GO

CREATE TABLE [HP_TEST].[ORDERS] (
  [ORDER_ID] NUMBER(38, 0) PRIMARY KEY NOT NULL IDENTITY(1, 1),
  [CUSTOMER_ID] NUMBER(38, 0),
  [ORDER_DATE] TIMESTAMP_NTZ(9) DEFAULT (CURRENT_TIMESTAMP()),
  [TOTAL_AMOUNT] NUMBER(10, 2) DEFAULT (0)
)
GO

CREATE TABLE [HP_TEST].[ORDER_ITEMS] (
  [ORDER_ITEM_ID] NUMBER(38, 0) UNIQUE PRIMARY KEY NOT NULL IDENTITY(1, 1),
  [ORDER_ID] NUMBER(38, 0),
  [PRODUCT_ID] NUMBER(38, 0),
  [QUANTITY] NUMBER(38, 0) DEFAULT (1),
  [PRICE] NUMBER(10, 2)
)
GO

CREATE TABLE [HP_TEST].[PRODUCTS] (
  [PRODUCT_ID] NUMBER(38, 0) PRIMARY KEY NOT NULL IDENTITY(1, 1),
  [PRODUCT_NAME] VARCHAR(100),
  [DESCRIPTION] VARCHAR(16777216),
  [PRICE] NUMBER(10, 2),
  [STOCK_QUANTITY] NUMBER(38, 0)
)
GO

CREATE TABLE [HP_TEST].[COMPOSITE_PK_TABLE] (
  [COLUMN1] NUMBER(38,0) NOT NULL,
  [COLUMN2] VARCHAR(16777216) NOT NULL,
  [COLUMN3] DATE,
  PRIMARY KEY ([COLUMN1], [COLUMN2])
)
GO

CREATE TABLE [HP_TEST].[CLONE_ORDERS] (
  [ORDER_ID] NUMBER(38, 0) PRIMARY KEY NOT NULL IDENTITY(1, 1),
  [CUSTOMER_ID] NUMBER(38, 0),
  [ORDER_DATE] TIMESTAMP_NTZ(9) DEFAULT (CURRENT_TIMESTAMP()),
  [TOTAL_AMOUNT] NUMBER(10, 2) DEFAULT (0)
)
GO

CREATE TABLE [HP_TEST].[employees] (
  [employee_id] INT PRIMARY KEY IDENTITY(1, 1),
  [first_name] VARCHAR(50) NOT NULL,
  [last_name] VARCHAR(50) NOT NULL,
  [email] VARCHAR(100),
  [hire_date] DATE DEFAULT (CURRENT_DATE())
)
GO

CREATE TABLE [TEST].[_holistics_persisted_test_20240703094645] (
  [a] NUMBER(1, 0),
  [b] NUMBER(1, 0),
  [c] NUMBER(1, 0)
)
GO

CREATE TABLE [TEST_2024_06_28_10_01_ZQMIMY].[ONE] (
  [FOO] VARCHAR(16777216) PRIMARY KEY,
  [BAR] VARCHAR(16777216)
)
GO

CREATE TABLE [TEST_2024_06_28_10_01_ZQMIMY].[TWO] (
  [ABC] VARCHAR(16777216),
  [XYZ] VARCHAR(16777216),
  [QWE] VARCHAR(16777216)
)
GO

CREATE TABLE [TESTKEYWORD] (
  [IFNULL] int,
  [NVL] int,
  [GET] int,
  [LEFT] int,
  [RIGHT] int,
  [DATE_PART] int,
  [TO_DATE] int,
  [DATE] int,
  [SPLIT] int,
  [NULLIF] int,
  [EQUAL_NULL] int
)
GO

CREATE TABLE [t] (
  [i] int
)
GO

CREATE TABLE [t1] (
  [v] varchar(16777216)
)
GO

CREATE TABLE [t4] (
  [v] varchar(16777216)
)
GO

CREATE TABLE [t5] (
  [v] varchar(16777216) NOT NULL
)
GO

CREATE TABLE [t6] (
  [v] varchar(32) UNIQUE NOT NULL
)
GO

CREATE TABLE [t7] (
  [i] integer DEFAULT (1)
)
GO

CREATE TABLE [t8] (
  [i] integer
)
GO

CREATE TABLE [t9] (
  [i] integer
)
GO

CREATE TABLE [t10] (
  [i] integer
)
GO

CREATE TABLE [t_default_1] (
  [i] integer NOT NULL
)
GO

CREATE TABLE [t_default_2] (
  [i] integer UNIQUE NOT NULL
)
GO

CREATE TABLE [t_default_3] (
  [i] integer UNIQUE NOT NULL
)
GO

CREATE TABLE [t_default_4] (
  [i] integer PRIMARY KEY NOT NULL
)
GO

CREATE TABLE [t_default_5] (
  [i] integer UNIQUE PRIMARY KEY NOT NULL
)
GO

CREATE TABLE [t_default_6] (
  [i] integer UNIQUE PRIMARY KEY NOT NULL
)
GO

CREATE TABLE [t_default_7] (
  [i] integer NOT NULL DEFAULT (1)
)
GO

CREATE TABLE [t_default_8] (
  [i] integer NOT NULL DEFAULT (1)
)
GO

CREATE TABLE [t_default_9] (
  [i] integer NOT NULL IDENTITY(1, 1)
)
GO

CREATE TABLE [t_default_10] (
  [i] integer NOT NULL IDENTITY(1, 1)
)
GO

CREATE TABLE [t_default_11] (
  [v] varchar(32) NOT NULL
)
GO

CREATE TABLE [t_default_12] (
  [v] varchar(32) NOT NULL
)
GO

CREATE TABLE [t_default_13] (
  [v] varchar(32) UNIQUE NOT NULL
)
GO

CREATE TABLE [t_default_14] (
  [v] varchar(32) UNIQUE NOT NULL
)
GO

CREATE TABLE [t_default_15] (
  [v] varchar(32) DEFAULT 'hey'
)
GO

CREATE TABLE [t_default_16] (
  [v] varchar(32) UNIQUE DEFAULT 'hey'
)
GO

CREATE TABLE [t_default_17] (
  [v] varchar(32) UNIQUE DEFAULT 'hey'
)
GO

CREATE TABLE [t_default_18] (
  [v] varchar(32) PRIMARY KEY DEFAULT 'hey'
)
GO

CREATE TABLE [t_masking_1] (
  [i] integer
)
GO

CREATE TABLE [t_masking_2] (
  [i] integer
)
GO

CREATE TABLE [t_policy_1] (
  [i] integer UNIQUE NOT NULL DEFAULT (1)
)
GO

CREATE TABLE [t_policy_2] (
  [v] nvarchar(255) UNIQUE NOT NULL
)
GO

CREATE TABLE [public] (
  [public] int
)
GO

CREATE TABLE [T] (
  [Source] string NOT NULL,
  [Query_Id] string,
  [State] string NOT NULL,
  [Procedure_Name] string
)
GO

CREATE TABLE [t_not_exist_1] (
  [v] varchar(16777216)
)
GO

CREATE TABLE [t_not_exist_2] (
  [v] varchar(16777216)
)
GO

CREATE TABLE [TESTSEED1] (
  [IDENT] int DEFAULT (SEQID.NEXTVAL),
  [mycol] string
)
GO

CREATE TABLE [TESTSEED2] (
  [IDENT] int DEFAULT (SCHEM.SEQID.NEXTVAL),
  [mycol] string
)
GO

CREATE TABLE [TESTSEED3] (
  [ident] int IDENTITY(1, 1)
)
GO

CREATE TABLE [TESTSEED4] (
  [ident] int IDENTITY(1, 1)
)
GO

CREATE TABLE [TESTSEED5] (
  [ident] int IDENTITY(1, 1)
)
GO

CREATE TABLE [TESTSEED6] (
  [ident] int IDENTITY(1, 1)
)
GO

CREATE TABLE [t_date_time_1] (
  [v] datetime(9)
)
GO

CREATE TABLE [T_DATE_TIME_2] (
  [TIMESTAMP] DATETIME,
  [VALUE] STRING,
  [NAME] STRING
)
GO

CREATE TABLE [T_OUT_1] (
  [C1] STRING
)
GO

CREATE TABLE [T_OUT_2] (
  [C2] STRING PRIMARY KEY
)
GO

CREATE TABLE [T_OUT_3] (
  [C3] STRING
)
GO

CREATE TABLE [T_OUT_4] (
  [C1] STRING
)
GO

CREATE TABLE [T_OUT_5] (
  [C2] STRING PRIMARY KEY
)
GO

CREATE TABLE [T_OUT_6] (
  [C3] STRING
)
GO

CREATE TABLE [T_CONSTRAINT_1] (
  [C1] STRING
)
GO

CREATE TABLE [T_CONSTRAINT_2] (
  [C1] STRING
)
GO

CREATE TABLE [T_CONSTRAINT_3] (
  [C3] STRING
)
GO

CREATE TABLE [T_CONSTRAINT_4] (
  [C3] STRING
)
GO

CREATE TABLE [T_CONSTRAINT_5] (
  [C3] STRING
)
GO

CREATE TABLE [T_CONSTRAINT_6] (
  [C3] STRING
)
GO

CREATE TABLE [T_CONSTRAINT_7] (
  [C1] STRING UNIQUE
)
GO

CREATE TABLE [T_CONSTRAINT_8] (
  [C3] STRING
)
GO

CREATE TABLE [T_CONSTRAINT_9] (
  [C3] STRING
)
GO

CREATE TABLE [T_CONSTRAINT_10] (
  [C3] STRING
)
GO

CREATE TABLE [TestK] (
  [NVL2] string,
  [FIRST_VALUE] string,
  [RESTRICT] int,
  [NVL] int,
  [RESPECT] int
)
GO

CREATE TABLE [T_TAG_1] (
  [C1] string
)
GO

CREATE TABLE [tpk] (
  [i] int PRIMARY KEY
)
GO

CREATE TABLE [tc1] (
  [i] int
)
GO

CREATE TABLE [tc2] (
  [i] int
)
GO

CREATE TABLE [tc3] (
  [c] char(4),
  [c2] character(2)
)
GO

CREATE TABLE [tz1] (
  [i] TIMESTAMPLTZ
)
GO

CREATE TABLE [tz2] (
  [i] TIMESTAMPNTZ
)
GO

CREATE TABLE [tz3] (
  [i] TIMESTAMPTZ
)
GO

CREATE TABLE [TESTSEED7] (
  [ident] int IDENTITY(1, 1)
)
GO

CREATE TABLE [TESTSEED8] (
  [ident] int IDENTITY(1, 1)
)
GO

CREATE TABLE [DIRECTION] (
  [LENGTH] INT,
  [LANGUAGE] INT
)
GO

CREATE TABLE [t_comment] (
  [i] int
)
GO

CREATE TABLE [T_NULL_NOT_NULL] (
  [C_NULL] INT NOT NULL
)
GO

CREATE TABLE [t_cluster] (
  [f1] nvarchar(255),
  [f2] number
)
GO

CREATE TABLE [t_virtual_1] (
  [f1] number,
  [f2] number,
  [f3] number
)
GO

CREATE TABLE [t_virtual_2] (
  [f1] number,
  [f2] number,
  [f3] number
)
GO

CREATE TABLE [floor] (
  [any_value] int,
  [getdate] int
)
GO

CREATE UNIQUE INDEX [CUSTOMERS_index_0] ON [HP_TEST].[CUSTOMERS] ("EMAIL", "SECOND_EMAIL")
GO

CREATE UNIQUE INDEX [ONE_index_0] ON [TEST_2024_06_28_10_01_ZQMIMY].[ONE] ("BAR")
GO

CREATE UNIQUE INDEX [ONE_index_1] ON [TEST_2024_06_28_10_01_ZQMIMY].[ONE] ("FOO", "BAR")
GO

CREATE UNIQUE INDEX [ONE_index_2] ON [TEST_2024_06_28_10_01_ZQMIMY].[ONE] ("BAR", "FOO")
GO

CREATE UNIQUE INDEX [TWO_index_3] ON [TEST_2024_06_28_10_01_ZQMIMY].[TWO] ("QWE", "ABC")
GO

CREATE UNIQUE INDEX [T_OUT_1_index_0] ON [T_OUT_1] ("C1")
GO

CREATE UNIQUE INDEX [T_OUT_4_index_1] ON [T_OUT_4] ("C1")
GO

CREATE UNIQUE INDEX [T_CONSTRAINT_1_index_2] ON [T_CONSTRAINT_1] ("C1")
GO

CREATE UNIQUE INDEX [T_CONSTRAINT_2_index_3] ON [T_CONSTRAINT_2] ("C1")
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'hello world',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 't4',
@level2type = N'Column', @level2name = 'v';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'hello world',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 't5',
@level2type = N'Column', @level2name = 'v';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'hello world',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 't6',
@level2type = N'Column', @level2name = 'v';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'hello world',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 't7',
@level2type = N'Column', @level2name = 'i';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'hello world',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 't_policy_1',
@level2type = N'Column', @level2name = 'i';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'hello world',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 't_policy_2',
@level2type = N'Column', @level2name = 'v';
GO

ALTER TABLE [HP_TEST].[ORDERS] ADD FOREIGN KEY ([CUSTOMER_ID]) REFERENCES [HP_TEST].[CUSTOMERS] ([CUSTOMER_ID])
GO

ALTER TABLE [HP_TEST].[ORDER_ITEMS] ADD FOREIGN KEY ([ORDER_ID]) REFERENCES [HP_TEST].[ORDERS] ([ORDER_ID])
GO

ALTER TABLE [HP_TEST].[ORDER_ITEMS] ADD FOREIGN KEY ([PRODUCT_ID]) REFERENCES [HP_TEST].[PRODUCTS] ([PRODUCT_ID])
GO

ALTER TABLE [TEST_2024_06_28_10_01_ZQMIMY].[TWO] ADD FOREIGN KEY ([XYZ], [ABC]) REFERENCES [TEST_2024_06_28_10_01_ZQMIMY].[ONE] ([BAR], [FOO])
GO

ALTER TABLE [TEST_2024_06_28_10_01_ZQMIMY].[TWO] ADD FOREIGN KEY ([QWE]) REFERENCES [TEST_2024_06_28_10_01_ZQMIMY].[ONE] ([FOO])
GO

ALTER TABLE [T_OUT_3] ADD FOREIGN KEY ([C3]) REFERENCES [T_OUT_2] ([C2])
GO

ALTER TABLE [T_OUT_6] ADD FOREIGN KEY ([C3]) REFERENCES [T_OUT_5] ([C2])
GO

ALTER TABLE [T_CONSTRAINT_3] ADD FOREIGN KEY ([C3]) REFERENCES [T_OUT_5] ([C2])
GO

ALTER TABLE [T_CONSTRAINT_4] ADD FOREIGN KEY ([C3]) REFERENCES [T_OUT_5] ([C2])
GO

ALTER TABLE [T_CONSTRAINT_5] ADD FOREIGN KEY ([C3]) REFERENCES [T_OUT_5] ([C2])
GO

ALTER TABLE [T_CONSTRAINT_6] ADD FOREIGN KEY ([C3]) REFERENCES [T_OUT_5] ([C2])
GO

ALTER TABLE [T_CONSTRAINT_8] ADD FOREIGN KEY ([C3]) REFERENCES [T_OUT_5] ([C2])
GO

ALTER TABLE [T_CONSTRAINT_9] ADD FOREIGN KEY ([C3]) REFERENCES [T_OUT_5] ([C2])
GO
