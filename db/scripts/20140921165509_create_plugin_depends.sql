--// create plugin_depends
-- Migration SQL that makes the change goes here.

CREATE TABLE plugin_depends
(
   PLUGIN_ID INTEGER NOT NULL, 
   DEPENDS_ID INTEGER NOT NULL,
   CONSTRAINT PK_PLUGIN_DEPENDS PRIMARY KEY (PLUGIN_ID, DEPENDS_ID) USING INDEX TABLESPACE pg_default 
) 
WITH (
  OIDS = FALSE
)
;
ALTER TABLE plugin_depends OWNER TO repouser;



--//@UNDO
-- SQL to undo the change goes here.

DROP TABLE plugin_depends;


