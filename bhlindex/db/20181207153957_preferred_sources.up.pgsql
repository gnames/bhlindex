CREATE TABLE preferred_sources  (
  name CHARACTER VARYING(255) COLLATE "C" NOT NULL,
  taxon_id CHARACTER VARYING(255) COLLATE "C" NOT NULL,
  match_type CHARACTER VARYING(100) NOT NULL,
  edit_distance INTEGER DEFAULT 0,
  stem_edit_distance INTEGER DEFAULT 0,
  matched_name CHARACTER VARYING(255) COLLATE "C" NOT NULL,
  matched_canonical CHARACTER VARYING(255) COLLATE "C" NOT NULL,
  current_name CHARACTER VARYING(255) COLLATE "C",
  classification CHARACTER VARYING COLLATE "C",
  datasource_id INTEGER NOT NULL,
  datasource_title CHARACTER VARYING(255) COLLATE "C" NOT NULL
);

CREATE INDEX preferred_sources_name_index
ON preferred_sources USING btree (name);