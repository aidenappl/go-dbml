CREATE TABLE `products` (
  `id` INT PRIMARY KEY,
  `price` DECIMAL(10,2) CHECK (price > 0),
  `quantity` INT CHECK (quantity >= 0),
  `discount` DECIMAL(5,2) CHECK (discount >= 0 AND discount <= 100)
);

CREATE TABLE `employees` (
  `emp_id` INT PRIMARY KEY,
  `age` INT,
  `salary` DECIMAL(10,2),
  CONSTRAINT `chk_age` CHECK (age >= 18 AND age <= 65),
  CONSTRAINT `chk_salary` CHECK (salary > 0)
);

CREATE TABLE `transactions` (
  `txn_id` INT PRIMARY KEY,
  `amount` DECIMAL(10,2),
  `fee` DECIMAL(10,2),
  `total` DECIMAL(10,2),
  CHECK (total = amount + fee),
  CHECK (amount > 0 OR fee > 0)
);

CREATE TABLE `contacts` (
  `contact_id` INT PRIMARY KEY,
  `email` VARCHAR(100) CHECK (email LIKE '%@%'),
  `phone` VARCHAR(20) CHECK (phone REGEXP '^[0-9]{3}-[0-9]{3}-[0-9]{4}$')
);

CREATE TABLE `orders` (
  `order_id` INT PRIMARY KEY,
  `status` VARCHAR(20) CHECK (status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled')),
  `priority` INT CHECK (priority IN (1, 2, 3, 4, 5))
);

CREATE TABLE `users` (
  `user_id` INT PRIMARY KEY,
  `role` VARCHAR(20) CHECK (role NOT IN ('banned', 'deleted', 'suspended'))
);

CREATE TABLE `ranges` (
  `range_id` INT PRIMARY KEY,
  `min_val` INT,
  `max_val` INT,
  CHECK ((min_val >= 0) AND (max_val <= 100) AND (min_val < max_val))
);

CREATE TABLE `inventory` (
  `item_id` INT PRIMARY KEY,
  `stock` INT,
  `status` VARCHAR(20),
  CHECK (CASE
      WHEN status = 'active' THEN stock >= 0
      WHEN status = 'discontinued' THEN stock = 0
      ELSE TRUE
    END)
);

CREATE TABLE `accounts` (
  `account_id` INT PRIMARY KEY,
  `balance` DECIMAL(15,2) CHECK (balance >= 0),
  `overdraft_limit` DECIMAL(15,2) CHECK (overdraft_limit >= 0),
  CONSTRAINT `chk_balance_overdraft` CHECK (balance + overdraft_limit >= 0)
);

CREATE TABLE `shipments` (
  `shipment_id` INT PRIMARY KEY,
  `weight` DECIMAL(10,2),
  `cost` DECIMAL(10,2),
  CONSTRAINT `chk_weight` CHECK (weight > 0),
  CONSTRAINT `chk_cost` CHECK (cost >= 0)
);

CREATE TABLE `payments` (
  `payment_id` INT PRIMARY KEY,
  `amount` DECIMAL(10,2),
  CHECK (amount > 0)
);

CREATE TABLE `pricing` (
  `price_id` INT PRIMARY KEY,
  `base_price` DECIMAL(10,2),
  `discount_pct` DECIMAL(5,2),
  `final_price` DECIMAL(10,2),
  CHECK (final_price = base_price * (1 - discount_pct / 100))
);

CREATE TABLE `settings` (
  `setting_id` INT PRIMARY KEY,
  `value` VARCHAR(100),
  `default_value` VARCHAR(100),
  CHECK (COALESCE(value, default_value) IS NOT NULL)
);

CREATE TABLE `bookings` (
  `booking_id` INT PRIMARY KEY,
  `check_in` DATE,
  `check_out` DATE,
  CHECK (check_out > check_in),
  CHECK (check_in >= CURDATE())
);

CREATE TABLE `documents` (
  `doc_id` INT PRIMARY KEY,
  `title` VARCHAR(200),
  `content` TEXT,
  CHECK (CHAR_LENGTH(title) > 0),
  CHECK (CHAR_LENGTH(content) >= 10)
);

CREATE TABLE `events` (
  `event_id` INT PRIMARY KEY,
  `status` VARCHAR(20),
  `cancelled` BOOLEAN,
  CHECK (status = 'active' OR status = 'completed' OR cancelled = TRUE)
);

CREATE TABLE `grades` (
  `grade_id` INT PRIMARY KEY,
  `score` INT,
  CHECK (score BETWEEN 0 AND 100)
);

CREATE TABLE `products_special` (
  `product_id` INT PRIMARY KEY,
  `code` VARCHAR(50),
  CONSTRAINT `chk_code_format` CHECK (code REGEXP '^[A-Z]{3}-[0-9]{4}$')
);

CREATE TABLE `features` (
  `feature_id` INT PRIMARY KEY,
  `is_enabled` BOOLEAN,
  `is_premium` BOOLEAN,
  CHECK (is_enabled = TRUE OR is_premium = FALSE)
);

CREATE TABLE `configs` (
  `config_id` INT PRIMARY KEY,
  `settings` JSON,
  CHECK (JSON_VALID(settings))
);

CREATE TABLE `audit_log` (
  `log_id` INT PRIMARY KEY,
  `severity` INT,
  `message` TEXT,
  CONSTRAINT `chk_severity_range` CHECK (severity BETWEEN 1 AND 5),
  CONSTRAINT `chk_message_length` CHECK (CHAR_LENGTH(message) > 0),
  CHECK (severity > 0)
);

CREATE TABLE `limits` (
  `limit_id` INT PRIMARY KEY,
  `min_value` INT,
  `max_value` INT,
  CHECK (max_value > min_value),
  CHECK (min_value >= 0),
  CHECK (max_value <= 1000)
);

CREATE TABLE `statuses` (
  `status_id` INT PRIMARY KEY,
  `current_status` VARCHAR(20),
  `previous_status` VARCHAR(20),
  CHECK (current_status != previous_status),
  CHECK (current_status <> '')
);

CREATE TABLE `defaults` (
  `default_id` INT PRIMARY KEY,
  `value1` INT,
  `value2` INT,
  CHECK (IFNULL(value1, 0) + IFNULL(value2, 0) > 0)
);

CREATE TABLE `sequences` (
  `seq_id` INT PRIMARY KEY,
  `value` INT,
  CHECK (MOD(value, 2) = 0)
);

CREATE TABLE `temperatures` (
  `temp_id` INT PRIMARY KEY,
  `celsius` DECIMAL(5,2),
  CHECK (ABS(celsius) < 100)
);
