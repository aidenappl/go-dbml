CREATE TABLE `customers` (
  `cust_id` NUMBER(10),
  `cust_name` VARCHAR2(100),
  `email` VARCHAR2(100),
  `phone` VARCHAR2(20),
  `status` VARCHAR2(20),
  `reg_date` DATE,
  `region` VARCHAR2(50),
  `credit_limit` NUMBER(10,2)
);

CREATE TABLE `inventory` (
  `item_id` NUMBER(10),
  `item_name` VARCHAR2(200),
  `price` NUMBER(10,2),
  `category` VARCHAR2(50),
  `sku` VARCHAR2(50),
  `stock_level` NUMBER(10),
  `warehouse_id` NUMBER(10),
  `supplier_id` NUMBER(10)
);

CREATE TABLE `shipments` (
  `ship_id` NUMBER(10),
  `customer_id` NUMBER(10),
  `ship_date` DATE,
  `total_amount` NUMBER(12,2),
  `tracking_number` VARCHAR2(50),
  `carrier` VARCHAR2(100),
  `priority` NUMBER(2)
);

CREATE TABLE `departments` (
  `dept_id` NUMBER(10),
  `dept_name` VARCHAR2(100),
  `manager_id` NUMBER(10),
  `location` VARCHAR2(50)
);

CREATE TABLE `suppliers` (
  `supp_id` NUMBER(10),
  `supp_name` VARCHAR2(100),
  `contact` VARCHAR2(100),
  `city` VARCHAR2(50)
);

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

CREATE INDEX `idx_cust_id` ON `customers` (`cust_id`);

CREATE UNIQUE INDEX `idx_cust_email` ON `customers` (`email`);

CREATE INDEX `idx_cust_status` ON `customers` (`status`) USING BITMAP;

CREATE INDEX `idx_cust_name_upper` ON `customers` ((UPPER(cust_name)));

CREATE UNIQUE INDEX `idx_cust_email_lower` ON `customers` ((LOWER(email)));

CREATE INDEX `idx_very_long_index_name_for_testing_parser_limits` ON `customers` (`cust_id`);

CREATE INDEX `idx_cust_id_bitmap` ON `customers` (`cust_id`) USING BITMAP;

CREATE INDEX `idx_cust_complex` ON `customers` ((TO_CHAR(reg_date, 'YYYYMM') || cust_name));

CREATE INDEX `idx_inv_item_price` ON `inventory` (`item_name`, `price`);

CREATE INDEX `idx_inv_category_stock` ON `inventory` (`category`, `stock_level`) USING BITMAP;

CREATE INDEX `idx_inv_name_price` ON `inventory` ((UPPER(item_name)), (price * 1.1));

CREATE INDEX `idx_inv_hash_like` ON `inventory` ((ORA_HASH(item_name)));

CREATE UNIQUE INDEX `idx_ship_cust_date` ON `shipments` (`customer_id`, `ship_date`);

CREATE INDEX `idx_ship_nullable` ON `shipments` (`customer_id`);

CREATE INDEX `idx_dt_mixed` ON `test_datatypes` (`col_number`, `col_varchar2`, `col_date`);
