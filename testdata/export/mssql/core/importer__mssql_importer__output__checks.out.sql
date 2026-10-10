CREATE TABLE [products] (
  [id] INT PRIMARY KEY,
  [price] DECIMAL(10,2) CHECK (price > 0),
  [quantity] INT CHECK (quantity >= 0),
  [discount] DECIMAL(5,2) CHECK (discount >= 0 AND discount <= 100)
)
GO

CREATE TABLE [employees] (
  [emp_id] INT PRIMARY KEY,
  [age] INT,
  [salary] DECIMAL(10,2),
  CONSTRAINT [chk_age] CHECK (age >= 18 AND age <= 65),
  CONSTRAINT [chk_salary] CHECK (salary > 0)
)
GO

CREATE TABLE [special_table] (
  [id] INT PRIMARY KEY,
  [price] DECIMAL CHECK ([price] > 0),
  [status] VARCHAR(20) CHECK ([status] IN ('active', 'inactive'))
)
GO

CREATE TABLE [transactions] (
  [txn_id] INT PRIMARY KEY,
  [amount] DECIMAL(10,2),
  [fee] DECIMAL(10,2),
  [total] DECIMAL(10,2),
  CHECK (total = amount + fee),
  CHECK (amount > 0 OR fee > 0)
)
GO

CREATE TABLE [contacts] (
  [contact_id] INT PRIMARY KEY,
  [email] VARCHAR(100) CHECK (email LIKE '%@%'),
  [phone] VARCHAR(20) CHECK (phone LIKE '[0-9][0-9][0-9]-[0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]')
)
GO

CREATE TABLE [orders] (
  [order_id] INT PRIMARY KEY,
  [status] VARCHAR(20) CHECK (status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled')),
  [priority] INT CHECK (priority IN (1, 2, 3, 4, 5))
)
GO

CREATE TABLE [users] (
  [user_id] INT PRIMARY KEY,
  [role] VARCHAR(20) CHECK (role NOT IN ('banned', 'deleted', 'suspended'))
)
GO

CREATE TABLE [ranges] (
  [range_id] INT PRIMARY KEY,
  [min_val] INT,
  [max_val] INT,
  CHECK ((min_val >= 0) AND (max_val <= 100) AND (min_val < max_val))
)
GO

CREATE TABLE [inventory] (
  [item_id] INT PRIMARY KEY,
  [stock] INT,
  [status] VARCHAR(20),
  CHECK (CASE
      WHEN status = 'active' THEN stock
      WHEN status = 'discontinued' THEN 0
      ELSE 1
    END >= 0)
)
GO

CREATE TABLE [accounts] (
  [account_id] INT PRIMARY KEY,
  [balance] DECIMAL(15,2) CHECK (balance >= 0),
  [overdraft_limit] DECIMAL(15,2) CHECK (overdraft_limit >= 0),
  CONSTRAINT [chk_balance_overdraft] CHECK (balance + overdraft_limit >= 0)
)
GO

CREATE TABLE [shipments] (
  [shipment_id] INT PRIMARY KEY,
  [weight] DECIMAL(10,2),
  [cost] DECIMAL(10,2),
  CONSTRAINT [chk_weight] CHECK (weight > 0),
  CONSTRAINT [chk_cost] CHECK (cost >= 0)
)
GO

CREATE TABLE [payments] (
  [payment_id] INT PRIMARY KEY,
  [amount] DECIMAL(10,2),
  CHECK (amount > 0)
)
GO

CREATE TABLE [pricing] (
  [price_id] INT PRIMARY KEY,
  [base_price] DECIMAL(10,2),
  [discount_pct] DECIMAL(5,2),
  [final_price] DECIMAL(10,2),
  CHECK (final_price = base_price * (1 - discount_pct / 100))
)
GO

CREATE TABLE [settings] (
  [setting_id] INT PRIMARY KEY,
  [value] VARCHAR(100),
  [default_value] VARCHAR(100),
  CHECK (COALESCE(value, default_value) IS NOT NULL)
)
GO

CREATE TABLE [bookings] (
  [booking_id] INT PRIMARY KEY,
  [check_in] DATE,
  [check_out] DATE,
  CHECK (check_out > check_in),
  CHECK (check_in >= GETDATE())
)
GO

CREATE TABLE [documents] (
  [doc_id] INT PRIMARY KEY,
  [title] VARCHAR(200),
  [content] VARCHAR(MAX),
  CHECK (LEN(title) > 0),
  CHECK (LEN(content) >= 10)
)
GO

CREATE TABLE [events] (
  [event_id] INT PRIMARY KEY,
  [status] VARCHAR(20),
  [cancelled] BIT,
  CHECK (status = 'active' OR status = 'completed' OR cancelled = 1)
)
GO

CREATE TABLE [grades] (
  [grade_id] INT PRIMARY KEY,
  [score] INT,
  CHECK (score BETWEEN 0 AND 100)
)
GO

CREATE TABLE [products_special] (
  [product_id] INT PRIMARY KEY,
  [code] VARCHAR(50),
  CONSTRAINT [chk_code_format] CHECK (code LIKE '[A-Z][A-Z][A-Z]-[0-9][0-9][0-9][0-9]')
)
GO

CREATE TABLE [features] (
  [feature_id] INT PRIMARY KEY,
  [is_enabled] BIT,
  [is_premium] BIT,
  CHECK (is_enabled = 1 OR is_premium = 0)
)
GO

CREATE TABLE [configs] (
  [config_id] INT PRIMARY KEY,
  [value1] INT,
  [value2] INT,
  CHECK (ISNULL(value1, 0) + ISNULL(value2, 0) > 0)
)
GO

CREATE TABLE [audit_log] (
  [log_id] INT PRIMARY KEY,
  [severity] INT,
  [message] VARCHAR(MAX),
  CONSTRAINT [chk_severity_range] CHECK (severity BETWEEN 1 AND 5),
  CONSTRAINT [chk_message_length] CHECK (LEN(message) > 0),
  CHECK (severity > 0)
)
GO

CREATE TABLE [limits] (
  [limit_id] INT PRIMARY KEY,
  [min_value] INT,
  [max_value] INT,
  CHECK (max_value > min_value),
  CHECK (min_value >= 0),
  CHECK (max_value <= 1000)
)
GO

CREATE TABLE [statuses] (
  [status_id] INT PRIMARY KEY,
  [current_status] VARCHAR(20),
  [previous_status] VARCHAR(20),
  CHECK (current_status != previous_status),
  CHECK (current_status <> '')
)
GO

CREATE TABLE [order_items] (
  [order_item_id] INT PRIMARY KEY,
  [order_id] INT,
  [quantity] INT,
  [unit_price] DECIMAL(10,2),
  CHECK ([quantity] > 0),
  CONSTRAINT [chk_price_positive] CHECK ([unit_price] > 0)
)
GO

CREATE TABLE [sequences] (
  [seq_id] INT PRIMARY KEY,
  [value] INT,
  CHECK (value % 2 = 0)
)
GO

CREATE TABLE [temperatures] (
  [temp_id] INT PRIMARY KEY,
  [celsius] DECIMAL(5,2),
  CHECK (ABS(celsius) < 100)
)
GO

CREATE TABLE [emails] (
  [email_id] INT PRIMARY KEY,
  [email_address] VARCHAR(100),
  CHECK (CHARINDEX('@', email_address) > 0)
)
GO

CREATE TABLE [codes] (
  [code_id] INT PRIMARY KEY,
  [code] VARCHAR(50),
  CHECK (code = UPPER(code))
)
GO

CREATE TABLE [order_discounts] (
  [discount_id] INT PRIMARY KEY,
  [order_total] DECIMAL(10,2),
  [discount_pct] DECIMAL(5,2),
  [discount_amt] DECIMAL(10,2),
  [final_total] DECIMAL(10,2),
  CONSTRAINT [chk_discount_calculation] CHECK ((discount_amt = order_total * discount_pct / 100) AND
    (final_total = order_total - discount_amt) AND
    (final_total >= 0))
)
GO
