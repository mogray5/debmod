--// create plugin_conflicts
-- Migration SQL that makes the change goes here.

CREATE TABLE plugin_conflicts
(
   PLUGIN_ID INTEGER NOT NULL, 
   CONFLICTS_ID INTEGER NOT NULL,
   CONSTRAINT PK_PLUGIN_CONFLICTS PRIMARY KEY (PLUGIN_ID, CONFLICTS_ID) USING INDEX TABLESPACE pg_default 
) 
WITH (
  OIDS = FALSE
)
;
ALTER TABLE plugin_conflicts OWNER TO repouser;



--//@UNDO
-- SQL to undo the change goes here.

DROP TABLE plugin_conflicts;
