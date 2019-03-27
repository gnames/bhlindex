CREATE TABLE name_statuses (
  name varchar(255) COLLATE "C.UTF-8" NOT NULL PRIMARY KEY,
  odds float,
  occurances int,
  processed boolean DEFAULT false
);

CREATE INDEX processed_index ON name_statuses
  USING btree (processed);
