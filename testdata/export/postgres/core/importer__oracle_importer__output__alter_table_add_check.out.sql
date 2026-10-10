CREATE TABLE "employees" (
  "emp_id" "NUMBER(10)",
  "emp_name" "VARCHAR2(100)",
  "salary" "NUMBER(10,2)",
  "commission_pct" "NUMBER(3,2)",
  "hire_date" DATE,
  "email" "VARCHAR2(100)",
  "department_id" "NUMBER(5)",
  "status" "VARCHAR2(20)",
  "age" "NUMBER(3)",
  "performance_score" "NUMBER(3,1)",
  CONSTRAINT "chk_salary_positive" CHECK (salary > 0),
  CONSTRAINT "chk_commission_range" CHECK (commission_pct BETWEEN 0 AND 1),
  CONSTRAINT "chk_status_values" CHECK (status IN ('ACTIVE', 'INACTIVE', 'ON_LEAVE', 'TERMINATED')),
  CONSTRAINT "chk_email_format" CHECK (email LIKE '%@%'),
  CHECK (age >= 18),
  CHECK (salary >= 0),
  CONSTRAINT "chk_salary_calculation" CHECK (salary * 12 > 0),
  CONSTRAINT "chk_salary_dept_rules" CHECK (salary > 0 AND department_id IS NOT NULL),
  CONSTRAINT "chk_commission_or_salary" CHECK (commission_pct > 0 OR salary > 50000),
  CONSTRAINT "chk_email_required" CHECK (email IS NOT NULL AND LENGTH(email) > 5),
  CONSTRAINT "chk_salary_by_status" CHECK (CASE 
        WHEN status = 'ACTIVE' THEN salary > 0
        WHEN status = 'INACTIVE' THEN 1 = 1
        ELSE 1 = 1
    END = 1),
  CONSTRAINT "chk_name_length" CHECK (LENGTH(TRIM(emp_name)) > 0),
  CONSTRAINT "chk_hire_date_valid" CHECK (hire_date <= SYSDATE),
  CONSTRAINT "chk_complex_conditions" CHECK ((salary > 30000 AND age >= 21) OR 
    (salary <= 30000 AND age >= 18)),
  CONSTRAINT "chk_email_lowercase" CHECK (email = LOWER(email)),
  CONSTRAINT "chk_performance_range" CHECK (performance_score >= 0.0 AND performance_score <= 10.0),
  CONSTRAINT "chk_email_domain" CHECK (email LIKE '%@%.com' OR email LIKE '%@%.org'),
  CONSTRAINT "chk_status_not_invalid" CHECK (status NOT IN ('PENDING', 'UNKNOWN')),
  CONSTRAINT "chk_comprehensive_validation" CHECK (salary > 0 AND 
    salary < 1000000 AND
    age >= 18 AND 
    age <= 100 AND
    (commission_pct IS NULL OR (commission_pct >= 0 AND commission_pct <= 1)) AND
    status IN ('ACTIVE', 'INACTIVE', 'ON_LEAVE', 'TERMINATED')),
  CONSTRAINT "chk_age_working" CHECK (age BETWEEN 18 AND 70),
  CONSTRAINT "chk_dept_valid" CHECK (department_id IS NOT NULL AND department_id > 0)
);

CREATE TABLE "products" (
  "product_id" "NUMBER(10)",
  "product_name" "VARCHAR2(200)",
  "price" "NUMBER(10,2)",
  "discount_price" "NUMBER(10,2)",
  "stock_quantity" "NUMBER(8)",
  "category" "VARCHAR2(50)",
  "weight" "NUMBER(8,2)",
  CONSTRAINT "chk_price_comparison" CHECK (discount_price <= price),
  CHECK (stock_quantity >= 0),
  CHECK (category IS NOT NULL),
  CONSTRAINT "chk_price_margin" CHECK ((price - discount_price) >= price * 0.05),
  CONSTRAINT "chk_weight_reasonable" CHECK (weight > 0 AND weight < 10000),
  CONSTRAINT "chk_price_not_negative" CHECK (NOT (price < 0)),
  CONSTRAINT "chk_category_uppercase" CHECK (UPPER(category) IN ('ELECTRONICS', 'CLOTHING', 'FOOD', 'BOOKS')),
  CONSTRAINT "chk_price_calculations" CHECK (price * stock_quantity > 0 AND price / 100 < 100000),
  CONSTRAINT "chk_discount_rules" CHECK ((price - discount_price) * stock_quantity >= 0),
  CHECK (product_name IS NOT NULL),
  CONSTRAINT "chk_prod_name_length" CHECK (LENGTH(product_name) <= 200),
  CONSTRAINT "chk_stock_max" CHECK (stock_quantity <= 1000000)
);
