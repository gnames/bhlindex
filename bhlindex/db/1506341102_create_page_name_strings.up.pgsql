CREATE TABLE page_name_strings (
  page_id varchar(255) NOT NULL,
  title_id integer NOT NULL,
  name_string varchar(255) NOT NULL,
  name_offset_start int NOT NULL,
  name_offset_end int NOT NULL,
  ends_next_page boolean DEFAULT false,
  odds float DEFAULT 0,
  kind varchar(255),
  updated_at timestamp without time zone
);

CREATE INDEX page_name_strings_page_title_index
ON page_name_strings USING btree (page_id, title_id);