CREATE TABLE pages (
  page_id character varying(255) NOT NULL,
  title_id integer NOT NULL,
  page_offset integer NOT NULL DEFAULT 0,
  CONSTRAINT pages_pkey PRIMARY KEY (title_id, page_id)
);