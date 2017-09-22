CREATE TABLE titles (
    id integer NOT NULL,
    path character varying(255) NOT NULL,
    internet_archive_id character varying(255) NOT NULL,
    gnrd_url character varying(255) NOT NULL,
    status integer DEFAULT 0
);

ALTER TABLE ONLY titles
    ADD CONSTRAINT titles_pkey PRIMARY KEY (id);

CREATE UNIQUE INDEX internet_archive_id_index
    ON titles USING btree (internet_archive_id);

CREATE INDEX status_index
    ON titles USING btree (status);
