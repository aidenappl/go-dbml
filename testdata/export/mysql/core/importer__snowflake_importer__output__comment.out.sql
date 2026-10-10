CREATE SCHEMA `TEST_SCHEMA`;

CREATE TABLE `TEST_SCHEMA`.`ONE` (
  `employee_id` INT NOT NULL COMMENT 'Unique identifier for each employee',
  `first_name` VARCHAR(50) NOT NULL COMMENT 'Employee''s first name',
  `last_name` VARCHAR(50) NOT NULL COMMENT 'Employee''s last name',
  `email` VARCHAR(100) COMMENT 'Employee''s work email address',
  `hire_date` DATE COMMENT 'Date when the employee was hired',
  `department` VARCHAR(50) COMMENT 'Department where the employee works',
  `salary` DECIMAL(10,2) COMMENT 'Employee''s annual salary'
);

CREATE TABLE `TEST_SCHEMA`.`TWO` (
  `employee_id` INT NOT NULL COMMENT 'Unique identifier for each employee',
  `first_name` VARCHAR(50) NOT NULL COMMENT 'Employee''s first name',
  `last_name` VARCHAR(50) NOT NULL COMMENT 'Employee''s last name',
  `email` VARCHAR(100) COMMENT 'Employee''s work email address',
  `hire_date` DATE COMMENT 'Date when the employee was hired',
  `department` VARCHAR(50) COMMENT 'Department where the employee works',
  `salary` DECIMAL(10,2) COMMENT 'Employee''s annual salary'
);

ALTER TABLE `TEST_SCHEMA`.`TWO` COMMENT = 'Another comment for employees table';
