CREATE TABLE `database_products` (
  `product_id` NUMBER(10) PRIMARY KEY,
  `product_name` VARCHAR2(100) NOT NULL,
  `vendor_name` VARCHAR2(100),
  `release_year` NUMBER(4),
  `is_open_source` NUMBER(1) DEFAULT 0,
  `license_cost` NUMBER(10, 2),
  `market_share` NUMBER(5, 2),
  `latest_version` VARCHAR2(20),
  `description` CLOB,
  `created_at` TIMESTAMP DEFAULT (SYSTIMESTAMP),
  `updated_at` TIMESTAMP DEFAULT (SYSTIMESTAMP),
  `is_active` NUMBER(1) DEFAULT 1,
  `max_connections` NUMBER(10),
  `status` CHAR(1) DEFAULT 'A'
);

-- Disable foreign key checks for INSERT
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `database_products` (`product_id`, `product_name`, `vendor_name`, `release_year`, `is_open_source`, `license_cost`, `market_share`, `latest_version`, `description`, `is_active`, `max_connections`, `status`, `created_at`)
VALUES
  (db_products_seq.NEXTVAL, 'Oracle Database', 'Oracle Corporation', 1979, 0, 47500, 28.5, '21c', CAST('Enterprise relational database management system' AS CLOB), 1, 10000, 'A', NULL),
  (NULL, 'MySQL', 'Oracle Corporation', NULL, 1, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL),
  (NULL, 'PostgreSQL', 'PostgreSQL Global Development Group', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'SQLite', NULL, 2000, NULL, NULL, 15.5 + 2.3, NULL, NULL, NULL, 1000 * 10, NULL, NULL),
  (NULL, UPPER('microsoft sql server'), 'Microsoft' || ' Corporation', NULL, NULL, NULL, NULL, SUBSTR('2022-CU10', 1, 4), TRIM('  Relational database management system  '), NULL, NULL, NULL, NULL),
  (NULL, 'MariaDB', 'MariaDB Foundation', EXTRACT(YEAR FROM SYSDATE), NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, SYSTIMESTAMP),
  (NULL, 'MongoDB', 'MongoDB Inc.', NULL, 1, NULL, 8.2, NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'Redis', 'Redis Ltd.', NULL, 1, NULL, 6.7, NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'Cassandra', 'Apache Software Foundation', NULL, 1, NULL, 4.1, NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'IBM Db2', NULL, NULL, NULL, 1500.00 * 12, ROUND(5.678, 1), NULL, NULL, NULL, POWER(2, 14), NULL, NULL),
  (NULL, 'Amazon Aurora', 'Amazon Web Services', NULL, 0, CASE WHEN 1=1 THEN 0.00 ELSE 1000.00 END, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'CockroachDB', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'Test Database', NULL, 2020, NULL, -100.00, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'Enterprise DB', NULL, NULL, NULL, 99999.99, NULL, NULL, NULL, NULL, 999999, NULL, NULL),
  (NULL, 'Neo4j', NULL, NULL, NULL, NULL, NULL, NULL, 'Graph database for connected data' || CHR(10) || 'Supports Cypher query language', NULL, NULL, NULL, NULL),
  (NULL, 'Firebird', NULL, 2000, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, TO_TIMESTAMP('2000-07-25 10:30:00', 'YYYY-MM-DD HH24:MI:SS')),
  (NULL, 'TimescaleDB', 'Timescale Inc.', NULL, 1, NULL, NULL, NULL, NULL, 1, NULL, NULL, NULL),
  (NULL, 'MemSQL', NULL, NULL, NULL, 0, 0, NULL, NULL, NULL, 0, NULL, NULL),
  (NULL, 'InfluxDB', NULL, EXTRACT(YEAR FROM SYSDATE), NULL, NULL, NULL, 'v' || '2.7', NULL, NULL, NULL, NULL, SYSDATE),
  (NULL, 'Vertica', 'Micro Focus', NULL, NULL, NVL(NULL, 5000.00), NULL, NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'Greenplum', NULL, NULL, DECODE('yes', 'yes', 1, 'no', 0, 0), NULL, NULL, NULL, NULL, NULL, NULL, 'A', NULL),
  (NULL, 'SAP HANA', 'SAP SE', NULL, NULL, NULL, NULL, NULL, NULL, NULL, LENGTH('Maximum'), NULL, NULL),
  (NULL, INITCAP('clickhouse'), INITCAP('clickhouse inc'), NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'Snowflake', NULL, NULL, NULL, COALESCE(NULL, NULL, 2500.00), 3.5, NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'Teradata', NULL, NULL, NULL, NULL, LEAST(10.5, 8.2, 15.7), NULL, NULL, NULL, GREATEST(1000, 5000, 3000), NULL, NULL),
  (NULL, 'Informix', NULL, 1980 + MOD(45, 20), NULL, NULL, NULL, NULL, NULL, NULL, MOD(10050, 1000), NULL, NULL),
  (db_products_seq.NEXTVAL, 'DynamoDB', 'Amazon Web Services', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'SingleStore', NULL, NULL, NULL, NULL, NULL, LPAD('8', 3, '0') || '.' || RPAD('5', 2, '0'), NULL, NULL, NULL, NULL, NULL),
  (NULL, 'Azure SQL', NULL, NULL, NULL, NULL, NULL, NULL, REPLACE('Microsoft cloud database service', 'cloud', 'managed cloud'), NULL, NULL, NULL, NULL),
  (NULL, 'Db2 Warehouse', NULL, NULL, NULL, ABS(-15000), ABS(-4.5), NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, 'Actian Ingres', NULL, NULL, NULL, TRUNC(12345.6789, 2), TRUNC(7.89456, 2), NULL, NULL, NULL, NULL, NULL, NULL),
  (NULL, UPPER(SUBSTR('rockset database', 1, 7)), INITCAP(TRIM('  rockset inc  ')), TO_NUMBER(TO_CHAR(SYSDATE, 'YYYY')) - 5, NULL, NULL, NULL, 'v' || TO_CHAR(ROUND(2.75, 1)), NULL, NULL, NULL, NULL, NULL);

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;