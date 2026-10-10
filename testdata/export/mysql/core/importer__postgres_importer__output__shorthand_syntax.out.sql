CREATE TABLE `booking_reference` (
  `reference_id` NVARCHAR(10) NOT NULL,
  `cust_id` NUMBER(10) NOT NULL,
  `status` NVARCHAR(1) NOT NULL,
  PRIMARY KEY (`reference_id`, `cust_id`)
);

CREATE TABLE `br_flight` (
  `reference_id` NVARCHAR(10) NOT NULL,
  `cust_id` NUMBER(10) NOT NULL,
  `flight_id` NVARCHAR(10) NOT NULL,
  PRIMARY KEY (`reference_id`, `flight_id`)
);

CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `full_name` varchar(255),
  `email` varchar(255) UNIQUE,
  `gender` varchar(255),
  `date_of_birth` varchar(255),
  `created_at` varchar(255),
  `modified_at` time(2),
  `country_code` int
);

CREATE TABLE `countries` (
  `code` int PRIMARY KEY,
  `name` varchar(255),
  `continent_name` varchar(255)
);

ALTER TABLE `users` ADD CONSTRAINT `fk_country_code` FOREIGN KEY (`country_code`) REFERENCES `countries` (`code`);

ALTER TABLE `br_flight` ADD CONSTRAINT `fk_composite` FOREIGN KEY (`reference_id`, `cust_id`) REFERENCES `booking_reference` (`reference_id`, `cust_id`);
