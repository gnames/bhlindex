CREATE TABLE pages (
  id character varying(255) NOT NULL,
  title_id character varying(255) NOT NULL,
  page_offset integer NOT NULL DEFAULT 0,
  CONSTRAINT pages_pkey PRIMARY KEY (id)
);

CREATE INDEX title_id_index
    ON pages USING btree (title_id);
