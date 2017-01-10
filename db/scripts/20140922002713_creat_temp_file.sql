--// creat temp file
-- Migration SQL that makes the change goes here.

CREATE TABLE temp_files
(
   PLUGIN_ID INTEGER NOT NULL,
   FILE_NM character varying(100) NOT NULL DEFAULT '',
   REL_PATH character varying(500) NOT NULL DEFAULT '',
   CHECKSUM character varying(500) NOT NULL DEFAULT '',
   LAST_CHANGED date NOT NULL, 
   NEW_CHECKSUM character varying(500) NOT NULL DEFAULT '',
   CONSTRAINT PK_TEMP_FILES PRIMARY KEY (PLUGIN_ID, FILE_NM, REL_PATH) USING INDEX TABLESPACE pg_default
) 
WITH (
  OIDS = FALSE
)
;
ALTER TABLE temp_files OWNER TO repouser;



--//@UNDO
-- SQL to undo the change goes here.

DROP TABLE temp_files;
