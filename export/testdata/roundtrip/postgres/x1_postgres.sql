SET standard_conforming_strings = on;
CREATE TYPE public."Mood" AS ENUM ('it''s ok', 'sad\', 'happy');
CREATE TABLE public."Acct" (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "Name" text DEFAULT 'O''Brien\' NOT NULL,
    mood public."Mood" DEFAULT 'it''s ok'::public."Mood",
    note text DEFAULT E'tab\there\'s',
    dollar text DEFAULT $q$it's $$ fine$q$,
    total numeric(10,2) GENERATED ALWAYS AS (("Name" IS NOT NULL)::int * 2) STORED,
    parent bigint REFERENCES public."Acct"(id) ON DELETE SET NULL,
    CONSTRAINT "ck name" CHECK (("Name" <> 'x''y') AND (length("Name") > 0))
);
COMMENT ON TABLE public."Acct" IS 'Accounts: it''s "quoted"';
COMMENT ON COLUMN public."Acct"."Name" IS E'multi\nline \\ backslash';
