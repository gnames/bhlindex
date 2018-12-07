CREATE TABLE name_strings (
  id serial NOT NULL,
  name CHARACTER VARYING(255) COLLATE "C.UTF-8" NOT NULL,
  match_type CHARACTER VARYING(100),
  edit_distance INTEGER DEFAULT 0,
  stem_edit_distance INTEGER DEFAULT 0,
  matched_name CHARACTER VARYING(255) COLLATE "C.UTF-8",
  matched_canonical CHARACTER VARYING(255) COLLATE "C.UTF-8",
  current_name CHARACTER VARYING(255) COLLATE "C.UTF-8",
  classification CHARACTER VARYING COLLATE "C.UTF-8",
  datasource_id INTEGER,
  datasources_number INTEGER,
  curation CHARACTER VARYING(255),
  retries  INTEGER NOT NULL DEFAULT 0,
  error CHARACTER VARYING(255),
  updated_at timestamp without time zone,
  CONSTRAINT name_strings_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX name_strings_name_index
ON name_strings USING btree (name);

-- CREATE INDEX name_strings_status_index
--     ON name_strings USING btree (status);