CREATE TABLE `alert` (
  `i` int
);

CREATE TABLE `condition` (
  `i` int
);

CREATE TABLE `dt_t` (
  `dynamic` int,
  `downstream` int,
  `target_lag` int
);

CREATE TABLE `wh_id` (
  `MAX_CONCURRENCY_LEVEL` int,
  `WAREHOUSE_TYPE` int
);

CREATE TABLE `CAST` (
  `KEY` int,
  `TRY_CAST` int
);

CREATE TABLE `DIRECTION` (
  `LENGTH` INT,
  `LANGUAGE` INT
);

CREATE TABLE `ROLES` (
  `USERS` INT,
  `WAREHOUSES` int,
  `INTEGRATIONS` int,
  `DATABASES` int,
  `SCHEMAS` int
);

CREATE TABLE `TABLES` (
  `VIEWS` INT,
  `STAGES` int,
  `STREAMS` int,
  `TASKS` int,
  `ALERTS` int
);

CREATE TABLE `ACTION` (
  `MODE` int
);

CREATE TABLE `interval` (
  `internal` int
);
