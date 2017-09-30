CREATE TABLE page_name_strings (
  page_id varchar(255) NOT NULL,
  name_string varchar(255) NOT NULL,
  name_offset_start int NOT NULL,
  name_offset_end int NOT NULL,
  ends_next_page boolean DEFAULT false,
  updated_at timestamp without time zone
);
