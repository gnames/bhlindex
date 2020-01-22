CREATE TABLE name_strings (
  id serial NOT NULL,
  name CHARACTER VARYING(255) COLLATE "C" NOT NULL,
  taxon_id CHARACTER VARYING(255) COLLATE "C",
  match_type CHARACTER VARYING(100),
  edit_distance INTEGER DEFAULT 0,
  stem_edit_distance INTEGER DEFAULT 0,
  matched_name CHARACTER VARYING(255) COLLATE "C",
  matched_canonical CHARACTER VARYING(255) COLLATE "C",
  current_name CHARACTER VARYING(255) COLLATE "C",
  classification CHARACTER VARYING COLLATE "C",
  datasource_id INTEGER,
  datasource_title CHARACTER VARYING(255) COLLATE "C",
  datasources_number INTEGER,
  curation CHARACTER VARYING(255),
  occurences int,
  odds float,
  retries  INTEGER NOT NULL DEFAULT 0,
  error CHARACTER VARYING(255),
  updated_at timestamp without time zone,
  CONSTRAINT name_strings_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX name_strings_name_index
ON name_strings USING btree (name);