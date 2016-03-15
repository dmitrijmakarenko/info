CREATE OR REPLACE FUNCTION acs_vcs_init()
  RETURNS void AS
$BODY$
DECLARE
acs_table text;
cnt int;
BEGIN

SELECT COUNT(*) INTO cnt FROM information_schema.schemata WHERE schema_name = 'acs_copy';
IF cnt = 0 THEN
	CREATE SCHEMA acs_copy;
END IF;

FOR acs_table IN
      SELECT quote_ident(table_name)
      FROM   information_schema.tables
      WHERE  table_schema = 'acs' AND table_type = 'BASE TABLE'
   LOOP
	EXECUTE 'CREATE TABLE acs_copy.' || acs_table || ' AS SELECT * FROM acs.' || acs_table;
   END LOOP;

INSERT INTO acs.changes_history(change_uuid, change_date, change_type, change_db) VALUES (uuid_generate_v4(), now(), 'init', current_database());

END;
$BODY$
  LANGUAGE plpgsql VOLATILE