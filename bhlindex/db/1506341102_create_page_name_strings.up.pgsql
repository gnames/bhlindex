CREATE TABLE page_name_strings (
  page_id varchar(255) NOT NULL,
  item_id integer NOT NULL,
  words_before varchar(255),
  name_string varchar(255) NOT NULL,
  words_after varchar(255),
  annot_nomen varchar(50),
  name_offset_start int NOT NULL,
  name_offset_end int NOT NULL,
  ends_next_page boolean DEFAULT false,
  odds float DEFAULT 0,
  kind varchar(255),
  updated_at timestamp without time zone
);

CREATE INDEX page_name_strings_page_item_index
ON page_name_strings USING btree (item_id, page_id);

CREATE INDEX page_name_strings_name_index
ON page_name_strings USING btree (name_string);
