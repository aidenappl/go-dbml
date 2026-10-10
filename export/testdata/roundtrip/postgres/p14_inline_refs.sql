create table authors (
  id integer primary key,
  name varchar(100) not null
);
create table books (
  id integer primary key,
  author_id integer not null references authors (id) on delete restrict,
  isbn char(13) unique,
  title text not null,
  price numeric(8,2) check (price > 0),
  lang varchar(5) default 'en-US'
);
create index on books (title);
create index books_lang_idx on books using hash (lang);
