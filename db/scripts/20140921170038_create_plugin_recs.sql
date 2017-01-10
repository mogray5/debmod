--// create plugin_recs
-- Migration SQL that makes the change goes here.

CREATE TABLE plugin_recs
(
   PLUGIN_ID INTEGER NOT NULL, 
   REC_ID INTEGER NOT NULL,
   CONSTRAINT PK_PLUGIN_RECS PRIMARY KEY (PLUGIN_ID, REC_ID) USING INDEX TABLESPACE pg_default 
) 
WITH (
  OIDS = FALSE
)
;
ALTER TABLE plugin_recs OWNER TO repouser;



--//@UNDO
-- SQL to undo the change goes here.

DROP TABLE plugin_recs;

