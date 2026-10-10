CREATE TABLE `authorization_records` (
  `id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `project_name` VARCHAR(100) NOT NULL,
  `updated_date` DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `created_date` DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE `DemoTable737` (
  `StudentId` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `StudentName` varchar(100),
  `StudentAdmissiondate` datetime
);
