CREATE TYPE public.order_status AS ENUM (
    'pending',
    'paid',
    'shipped',
    'refunded'
);

CREATE TABLE public.orders (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    status public.order_status DEFAULT 'pending'::public.order_status NOT NULL,
    total numeric(12,2) DEFAULT 0 NOT NULL,
    metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
    tags text[] DEFAULT '{}'::text[],
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE public.orders IS 'Customer orders';
COMMENT ON COLUMN public.orders.total IS 'Order total in dollars';
COMMENT ON COLUMN public.orders.metadata IS 'Free-form metadata, don''t index';

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);
CREATE INDEX orders_status_created_at_idx ON public.orders USING btree (status, created_at DESC);
CREATE INDEX orders_metadata_idx ON public.orders USING gin (metadata);
