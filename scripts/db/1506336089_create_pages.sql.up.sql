CREATE TABLE pages (
  id character varying(255) NOT NULL,
  title_id character varying(255) NOT NULL,
   CONSTRAINT pages_pkey PRIMARY KEY (id)
);

CREATE INDEX title_id_index
    ON pages USING btree (title_id);

CREATE TABLE pages_tmp (
  id character varying(255) NOT NULL,
  title_id character varying(255) NOT NULL
);
