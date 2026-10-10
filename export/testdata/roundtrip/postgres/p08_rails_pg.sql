SET statement_timeout = 0;
SET client_encoding = 'UTF8';

CREATE TABLE public.ar_internal_metadata (
    key character varying NOT NULL,
    value character varying,
    created_at timestamp(6) without time zone NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);

CREATE TABLE public.comments (
    id bigint NOT NULL,
    article_id bigint NOT NULL,
    commenter character varying,
    body text,
    status character varying DEFAULT 'public'::character varying,
    created_at timestamp(6) without time zone NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);

CREATE SEQUENCE public.comments_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER SEQUENCE public.comments_id_seq OWNED BY public.comments.id;

CREATE TABLE public.articles (
    id bigint NOT NULL,
    title character varying,
    body text,
    created_at timestamp(6) without time zone NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);

ALTER TABLE ONLY public.articles ALTER COLUMN id SET DEFAULT nextval('public.articles_id_seq'::regclass);
ALTER TABLE ONLY public.comments ALTER COLUMN id SET DEFAULT nextval('public.comments_id_seq'::regclass);

ALTER TABLE ONLY public.ar_internal_metadata ADD CONSTRAINT ar_internal_metadata_pkey PRIMARY KEY (key);
ALTER TABLE ONLY public.articles ADD CONSTRAINT articles_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.comments ADD CONSTRAINT comments_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.schema_migrations ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);
CREATE INDEX index_comments_on_article_id ON public.comments USING btree (article_id);
ALTER TABLE ONLY public.comments ADD CONSTRAINT fk_rails_3bf61a60d3 FOREIGN KEY (article_id) REFERENCES public.articles(id);
