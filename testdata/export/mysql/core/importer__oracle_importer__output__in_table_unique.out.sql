CREATE TABLE `test_datatypes` (
  `col_number` NUMBER(10,2),
  `col_varchar2` VARCHAR2(100),
  `col_char` CHAR(10),
  `col_date` DATE,
  `col_timestamp` TIMESTAMP,
  `col_clob` CLOB,
  `col_raw` RAW(100),
  `col_float` FLOAT,
  `col_integer` INTEGER,
  `col_nvarchar2` NVARCHAR2(100),
  `col_nchar` NCHAR(10)
);

CREATE TABLE `order_details` (
  `order_id` NUMBER(10),
  `item_id` NUMBER(10),
  `quantity` NUMBER(10),
  `order_ref` VARCHAR2(20)
);

CREATE TABLE `reports` (
  `report_id` NUMBER(10),
  `report_name` VARCHAR2(100),
  `report_date` DATE,
  `report_code` VARCHAR2(20)
);

CREATE TABLE `sales` (
  `sale_id` NUMBER(10),
  `sale_date` DATE,
  `sale_code` VARCHAR2(20)
);

CREATE TABLE `accounts` (
  `account_id` NUMBER(10),
  `account_number` VARCHAR2(20),
  `balance` NUMBER(12,2)
);

CREATE UNIQUE INDEX `uk_number` ON `test_datatypes` (`col_number`);

CREATE UNIQUE INDEX `uk_varchar2` ON `test_datatypes` (`col_varchar2`);

CREATE UNIQUE INDEX `uk_char` ON `test_datatypes` (`col_char`);

CREATE UNIQUE INDEX `uk_date` ON `test_datatypes` (`col_date`);

CREATE UNIQUE INDEX `uk_timestamp` ON `test_datatypes` (`col_timestamp`);

CREATE UNIQUE INDEX `uk_clob` ON `test_datatypes` (`col_clob`);

CREATE UNIQUE INDEX `uk_raw` ON `test_datatypes` (`col_raw`);

CREATE UNIQUE INDEX `uk_float` ON `test_datatypes` (`col_float`);

CREATE UNIQUE INDEX `uk_integer` ON `test_datatypes` (`col_integer`);

CREATE UNIQUE INDEX `uk_nvarchar2` ON `test_datatypes` (`col_nvarchar2`);

CREATE UNIQUE INDEX `uk_nchar` ON `test_datatypes` (`col_nchar`);

CREATE UNIQUE INDEX `uk_order_ref` ON `order_details` (`order_id`, `item_id`);

CREATE UNIQUE INDEX `uk_very_long_constraint_name_for_testing_parser_limits` ON `reports` (`report_code`);

CREATE UNIQUE INDEX `uk_sale_code` ON `sales` (`sale_code`);

CREATE UNIQUE INDEX `uk_account_number` ON `accounts` (`account_number`);
