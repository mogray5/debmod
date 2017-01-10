--// create plugin_files
-- Migration SQL that makes the change goes here.

CREATE SEQUENCE plugin_files_seq
   INCREMENT 1
   START 1;
ALTER TABLE plugin_files_seq OWNER TO repouser;

CREATE TABLE plugin_files
(
   FILE_ID INTEGER NOT NULL DEFAULT nextval('plugin_seq'::regclass), 
   PLUGIN_ID INTEGER NOT NULL,
   FILE_NM character varying(100) NOT NULL DEFAULT '',
   REL_PATH character varying(500) NOT NULL DEFAULT '',
   CHECKSUM character varying(500) NOT NULL DEFAULT '',
   LAST_CHANGED date NOT NULL, 
   NEW_CHECKSUM character varying(500) NOT NULL DEFAULT '',
   CONSTRAINT PK_PLUGIN_FILES PRIMARY KEY (FILE_ID) USING INDEX TABLESPACE pg_default, 
   CONSTRAINT UK_PLUGIN_FILES_1 UNIQUE (PLUGIN_ID, FILE_NM, REL_PATH) USING INDEX TABLESPACE pg_default
) 
WITH (
  OIDS = FALSE
)
;
ALTER TABLE plugin_files OWNER TO repouser;



--//@UNDO
-- SQL to undo the change goes here.

DROP TABLE plugin_files;
DROP SEQUENCE plugin_files_seq;
