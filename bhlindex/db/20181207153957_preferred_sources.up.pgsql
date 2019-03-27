CREATE TABLE preferred_sources  (
  name CHARACTER VARYING(255) COLLATE "C.UTF-8" NOT NULL,
  datasource_id INTEGER,
  datasource_title CHARACTER VARYING(255) COLLATE "C.UTF-8" NOT NULL,
  matched_name CHARACTER VARYING(255) COLLATE "C.UTF-8" NOT NULL,
  taxon_id CHARACTER VARYING(255) COLLATE "C.UTF-8" NOT NULL
);

CREATE INDEX preferred_sources_name_index
ON preferred_sources USING btree (name);