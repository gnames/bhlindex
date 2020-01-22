CREATE TABLE name_statuses (
  name varchar(255) COLLATE "C" NOT NULL PRIMARY KEY,
  odds float NOT NULL DEFAULT 0,
  occurences int NOT NULL DEFAULT 0,
  processed boolean DEFAULT false
);

CREATE INDEX processed_index ON name_statuses
  USING btree (processed);
