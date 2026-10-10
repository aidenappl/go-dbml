--
-- PostgreSQL database dump
--
SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SELECT pg_catalog.set_config('search_path', '', false);
CREATE SCHEMA app;
CREATE TYPE app.role AS ENUM (
    'admin',
    'member'
);
CREATE FUNCTION app.touch() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$;
SET default_tablespace = '';
CREATE TABLE app.users (
    id bigint NOT NULL,
    email character varying(320) NOT NULL,
    role app.role DEFAULT 'member'::app.role NOT NULL,
    meta jsonb DEFAULT '{}'::jsonb,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    tags text[] DEFAULT ARRAY[]::text[]
);
COMMENT ON TABLE app.users IS 'Application users';
COMMENT ON COLUMN app.users.email IS 'login, it''s unique';
CREATE SEQUENCE app.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE app.users_id_seq OWNED BY app.users.id;
CREATE TABLE app.posts (
    id integer NOT NULL,
    user_id bigint,
    title text NOT NULL,
    body text,
    score numeric(10,2) DEFAULT 0.0,
    published boolean DEFAULT false
);
ALTER TABLE ONLY app.users ALTER COLUMN id SET DEFAULT nextval('app.users_id_seq'::regclass);
ALTER TABLE ONLY app.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
ALTER TABLE ONLY app.users
    ADD CONSTRAINT users_email_key UNIQUE (email);
ALTER TABLE ONLY app.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);
CREATE INDEX posts_user_id_idx ON app.posts USING btree (user_id);
CREATE UNIQUE INDEX posts_title_lower ON app.posts USING btree (lower(title));
ALTER TABLE ONLY app.posts
    ADD CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES app.users(id) ON DELETE CASCADE;
