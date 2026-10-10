CREATE TABLE IF NOT EXISTS organizations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL CHECK (char_length(name) > 0),
  plan varchar(20) NOT NULL DEFAULT 'free',
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE memberships (
  org_id uuid NOT NULL REFERENCES organizations ON DELETE CASCADE,
  user_id integer NOT NULL,
  role text,
  PRIMARY KEY (org_id, user_id),
  CONSTRAINT role_ok CHECK (role IN ('owner', 'member'))
);
CREATE TABLE invoices (
  id bigserial PRIMARY KEY,
  org_id uuid REFERENCES organizations (id) ON UPDATE CASCADE ON DELETE SET NULL,
  amount_cents bigint NOT NULL DEFAULT 0,
  currency char(3) DEFAULT 'USD',
  due date,
  paid_at timestamp(6) without time zone,
  data json,
  ratio double precision,
  UNIQUE (org_id, due)
);
CREATE INDEX CONCURRENTLY IF NOT EXISTS invoices_due ON invoices (due DESC NULLS LAST) WHERE paid_at IS NULL;
ALTER TABLE invoices ADD CONSTRAINT positive CHECK (amount_cents >= 0);
INSERT INTO organizations (id, name) VALUES ('00000000-0000-0000-0000-000000000001', 'Acme'), ('00000000-0000-0000-0000-000000000002', E'O\'Brien');
INSERT INTO invoices (org_id, amount_cents, due, ratio) VALUES (NULL, -500, '2024-03-01', 0.5);
