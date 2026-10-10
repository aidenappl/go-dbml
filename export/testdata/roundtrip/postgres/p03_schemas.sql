CREATE SCHEMA billing;
CREATE SCHEMA auth;

CREATE TABLE auth.accounts (
    id bigserial PRIMARY KEY,
    email citext NOT NULL UNIQUE,
    password_digest text NOT NULL
);

CREATE TABLE billing.invoices (
    id bigserial PRIMARY KEY,
    account_id bigint NOT NULL REFERENCES auth.accounts(id),
    amount_cents integer NOT NULL CHECK (amount_cents >= 0),
    issued_on date NOT NULL DEFAULT CURRENT_DATE,
    paid boolean NOT NULL DEFAULT false
);

CREATE TABLE billing.invoice_lines (
    invoice_id bigint NOT NULL,
    position smallint NOT NULL,
    description text NOT NULL,
    amount_cents integer NOT NULL,
    PRIMARY KEY (invoice_id, position),
    FOREIGN KEY (invoice_id) REFERENCES billing.invoices (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX invoices_account_issued ON billing.invoices (account_id, issued_on);
