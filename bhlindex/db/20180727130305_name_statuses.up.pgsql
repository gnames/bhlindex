CREATE TABLE name_statuses (
  name varchar(255) PRIMARY KEY,
  processed boolean DEFAULT false
);

CREATE INDEX processed_index ON name_statuses
  USING btree (processed);
