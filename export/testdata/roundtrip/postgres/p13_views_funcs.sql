SET search_path = public, pg_catalog;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE TABLE public.accounts (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    balance numeric DEFAULT 0 NOT NULL,
    CONSTRAINT accounts_pkey PRIMARY KEY (id)
);
CREATE FUNCTION public.touch() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$;
CREATE VIEW public.rich_accounts AS
 SELECT accounts.id, accounts.balance FROM public.accounts WHERE (accounts.balance > (1000)::numeric);
CREATE TABLE public.transfers (
    id bigint NOT NULL,
    src uuid NOT NULL,
    dst uuid NOT NULL,
    amount numeric(14,2) NOT NULL
);
ALTER TABLE ONLY public.transfers ADD CONSTRAINT transfers_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.transfers ADD CONSTRAINT transfers_src_fkey FOREIGN KEY (src) REFERENCES public.accounts(id);
ALTER TABLE ONLY public.transfers ADD CONSTRAINT transfers_dst_fkey FOREIGN KEY (dst) REFERENCES public.accounts(id);
