--// create plugin data
-- Migration SQL that makes the change goes here.

CREATE SEQUENCE plugin_seq
   INCREMENT 1
   START 1;
ALTER TABLE plugin_seq OWNER TO repouser;

CREATE TABLE plugin
(
   PLUGIN_ID INTEGER NOT NULL DEFAULT nextval('plugin_seq'::regclass), 
   PLUGIN_NM character varying(75) NOT NULL DEFAULT '',
   VCS_URL character varying(255) NOT NULL DEFAULT '', 
   VCS_CLONE_FOLDER character varying(75) NOT NULL DEFAULT '',
   VCS_CLONE_CMD character varying(255) NOT NULL DEFAULT '',
   DEST_FOLDER character varying(75) NOT NULL DEFAULT '', 
   DESCRIPTION character varying(255) NOT NULL DEFAULT '', 
   AUTHOR character varying(75) NOT NULL DEFAULT '', 
   FORUM_LINK character varying(255) NOT NULL DEFAULT '',
   PKG_NM character varying(75) NOT NULL DEFAULT '',  
   PKG_VERSION character varying(30) NOT NULL DEFAULT '',
   CONSTRAINT PK_PLUGIN PRIMARY KEY (PLUGIN_ID) USING INDEX TABLESPACE pg_default, 
   CONSTRAINT UK_PLUGIN_1 UNIQUE (PLUGIN_NM) USING INDEX TABLESPACE pg_default
) 
WITH (
  OIDS = FALSE
)
;
ALTER TABLE plugin OWNER TO repouser;



--//@UNDO
-- SQL to undo the change goes here.

DROP TABLE plugin;
DROP SEQUENCE plugin_seq;
