CREATE TABLE public.nodes (
    id integer NOT NULL,
    parent_id integer,
    label character varying(64) DEFAULT 'node'::character varying NOT NULL,
    weight numeric DEFAULT 1.5,
    CONSTRAINT nodes_weight_check CHECK ((weight > (0)::numeric))
);
CREATE TABLE public.edges (
    src integer NOT NULL,
    dst integer NOT NULL,
    created timestamp without time zone DEFAULT timezone('utc'::text, now())
);
ALTER TABLE ONLY public.nodes ADD CONSTRAINT nodes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.edges ADD CONSTRAINT edges_pkey PRIMARY KEY (src, dst);
ALTER TABLE ONLY public.nodes ADD CONSTRAINT nodes_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.nodes(id) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE ONLY public.edges ADD CONSTRAINT edges_src_fkey FOREIGN KEY (src) REFERENCES public.nodes(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.edges ADD CONSTRAINT edges_dst_fkey FOREIGN KEY (dst) REFERENCES public.nodes(id) ON DELETE CASCADE;
