CREATE TABLE users (
  id INT NOT NULL PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  manager_id INT FOREIGN KEY REFERENCES users(id),
  full_name AS (first + ' ' + last) PERSISTED,
  first NVARCHAR(50) DEFAULT N'x',
  last NVARCHAR(50) DEFAULT '',
  score FLOAT DEFAULT -1.5,
  ts ROWVERSION
);
CREATE TABLE sales.items (
  sku CHAR(8) NOT NULL,
  user_id INT NOT NULL,
  qty SMALLINT CHECK (qty > 0),
  CONSTRAINT pk_items PRIMARY KEY (sku, user_id),
  CONSTRAINT fk_items_user FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE NO ACTION,
  INDEX ix_qty UNIQUE (qty DESC)
);
CREATE UNIQUE INDEX ux_email ON users (email);
INSERT INTO users (id, email, first) VALUES (1, 'a@b.c', N'Ann'), (2, NULL, 'Bob');
INSERT INTO sales.items VALUES ('A1', 1, 3);
ALTER TABLE users ADD nickname NVARCHAR(30) NULL, CONSTRAINT ck_nick CHECK (nickname <> '');
