CREATE TABLE name_strings (
  id serial NOT NULL,
  name character varying(255) COLLATE "C.UTF-8" NOT NULL,
  status integer NOT NULL DEFAULT 0,
  match_type varchar(100),
  in_curated_source boolean DEFAULT false,
  CONSTRAINT name_strings_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX name_id_index
    ON name_strings USING btree (name);
