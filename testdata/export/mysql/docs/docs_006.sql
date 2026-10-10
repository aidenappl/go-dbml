CREATE TABLE `users` (
  `id` integer,
  `wealth` integer,
  `debt` integer,
  CONSTRAINT `chk_positive_money` CHECK (debt + wealth >= 0)
);
