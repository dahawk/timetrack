CREATE EXTENSION "uuid-ossp"
  SCHEMA public
  VERSION "1.0";

CREATE TABLE public."user"
(
  id uuid NOT NULL DEFAULT uuid_generate_v4(),
  name text NOT NULL,
  username text NOT NULL,
  password text NOT NULL,
  admin boolean DEFAULT false,
  CONSTRAINT user_pkey PRIMARY KEY (id)
);
ALTER TABLE public."user"
  OWNER TO timetrack;

CREATE TABLE public.entry_data
(
  id bigserial NOT NULL,
  created timestamp without time zone,
  begin timestamp without time zone,
  "end" timestamp without time zone,
  type text,
  create_type text,
  CONSTRAINT entry_data_pkey PRIMARY KEY (id)
);
  ALTER TABLE public.entry_data
    OWNER TO timetrack;

CREATE TABLE public.entry
(
  id bigserial NOT NULL,
  entry_data bigint,
  modified timestamp without time zone,
  user_id uuid,
  entry_id uuid default uuid_generate_v4(),
  active boolean,
  CONSTRAINT entry_pkey PRIMARY KEY (id),
  CONSTRAINT entry_entry_data_fkey FOREIGN KEY (entry_data)
      REFERENCES public.entry_data (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT entry_user_fk FOREIGN KEY (user_id)
      REFERENCES public."user" (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);
ALTER TABLE public.entry
  OWNER TO timetrack;

#intial admin user, since i had to create one every time i set up a new instance, here is a prepared admin.
#intial password is AdminInitia1Pa5sword
#if you use this script change the password under any circumstances!!
insert into user (name, username, password, admin) ('admin','admin','$2a$12$1JRT5Ti4g3D6m27GVwy.3.6ozHJyn9xi6xJOuaXD4FqVC30m/6B7C',true);
