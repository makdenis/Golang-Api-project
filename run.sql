
drop table if exists users cascade;
drop table if exists forums cascade;
drop table if exists threads cascade;
drop table if exists votes cascade;
drop table if exists posts2 cascade;

CREATE TABLE IF NOT EXISTS users
(
  nickname VARCHAR(64) NOT NULL UNIQUE primary key,
  email    TEXT NOT NULL UNIQUE,

  about    TEXT DEFAULT '',
  fullname VARCHAR(96) DEFAULT ''
);



CREATE TABLE IF NOT EXISTS forums
(
  id      BIGSERIAL primary key,

  slug    TEXT not null unique,

  title   TEXT,

  threads INTEGER DEFAULT 0,
  posts   INTEGER DEFAULT 0,

  author  VARCHAR references users(nickname)
);

CREATE TABLE if not exists threads
(
  id         BIGSERIAL PRIMARY KEY,
  slug       TEXT unique,

  created    TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,

  message    TEXT default '',
  title      TEXT default '',

  author     VARCHAR REFERENCES users (nickname),
  forum      TEXT REFERENCES forums(slug),

  votes      INTEGER DEFAULT 0
);

create table if not exists posts2
(
  id        bigserial  primary key,

  created   TIMESTAMPTZ NOT NULL,--TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,

  is_edited boolean default FALSE,

  parent    bigint DEFAULT 0,
  tree_path integer[],

  message   text not null,

  author    varchar not null references users(nickname),
  forum     TEXT references forums(slug),
  thread    bigint references threads(id)
);



CREATE TABLE if not exists votes
(
  id        bigserial   NOT NULL PRIMARY KEY,
  username  VARCHAR     ,
  thread    INTEGER     ,
  voice     INTEGER,

  UNIQUE(username, thread,voice)
);

