CREATE TABLE pages (
  page_id character varying(255) NOT NULL,
  item_id integer,
  page_offset integer,
  CONSTRAINT pages_pkey PRIMARY KEY (item_id, page_id)
);