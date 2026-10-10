CREATE SCHEMA `demographic`;

CREATE TABLE `user` (
  `name` text,
  `id` int,
  `status` ENUM ('active', 'churned', 'inactive') DEFAULT 'active',
  `gender` ENUM ('male', 'female') DEFAULT 'male',
  `age_type` ENUM ('toddler', 'children', 'teenager', 'young_adult', 'elder') DEFAULT 'young_adult'
);
