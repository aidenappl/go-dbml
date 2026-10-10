CREATE UNLOGGED TABLE public.cache_entries (
    key text NOT NULL PRIMARY KEY,
    value bytea,
    expires_at timestamptz
);
CREATE TABLE public.metrics (
    ts timestamptz NOT NULL,
    name text NOT NULL,
    value double precision NOT NULL
) PARTITION BY RANGE (ts);
CREATE TABLE public.metrics_2026 PARTITION OF public.metrics FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');
ALTER TABLE public.cache_entries ADD COLUMN hits integer DEFAULT 0 NOT NULL;
