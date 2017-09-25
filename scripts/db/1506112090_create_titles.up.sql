CREATE TABLE titles (
    id serial NOT NULL,
    path character varying(255) NOT NULL,
    internet_archive_id character varying(255) NOT NULL,
    gnrd_url character varying(255),
    status integer NOT NULL DEFAULT 0,
    language character varying(100),
    english_detected boolean NOT NULL DEFAULT false,
    updated_at timestamp without time zone,
    CONSTRAINT titles_pkey PRIMARY KEY (id)
);

CREATE INDEX internet_archive_id_index
    ON titles USING btree (internet_archive_id);

CREATE INDEX status_index
    ON titles USING btree (status);

CREATE INDEX updated_at_index
    ON titles USING btree (updated_at);
