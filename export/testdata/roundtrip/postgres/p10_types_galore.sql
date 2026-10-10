CREATE TABLE measurements (
    id integer NOT NULL,
    sensor_id integer NOT NULL,
    reading double precision NOT NULL,
    reading_int bigint,
    tsr tstzrange,
    taken_at timestamp without time zone NOT NULL,
    taken_on date GENERATED ALWAYS AS ((taken_at)::date) STORED,
    meta hstore,
    flags bit varying(8),
    note character(10),
    ip inet,
    mac macaddr,
    dur interval,
    raw bytea,
    PRIMARY KEY (id)
);
