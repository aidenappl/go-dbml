CREATE SCHEMA `myschema`;

CREATE TABLE `sample_table` (
  `status` VARCHAR(20) CHECK ([status] IN ('active', 'inactive', 'pending')),
  `email` VARCHAR(100) CHECK ([email] LIKE '%@%.%'),
  `user_type` VARCHAR(20) CHECK ([user_type] NOT IN ('banned', 'suspended'))
);

CREATE TABLE `myschema`.`sample_table` (
  `status` VARCHAR(20),
  `email` VARCHAR(100),
  `user_type` VARCHAR(20),
  CONSTRAINT `chk_status` CHECK (status IN ('active', 'inactive', 'pending')),
  CONSTRAINT `chk_email` CHECK (email LIKE '%@%.%'),
  CONSTRAINT `chk_user_type` CHECK (user_type NOT IN ('banned', 'suspended'))
);
