CREATE OR REPLACE FUNCTION acs_vcs_init()
  RETURNS void AS
$BODY$
DECLARE
tname text;
ruuid text;
BEGIN

FOR tname IN
	SELECT quote_ident(table_name)
	FROM   information_schema.tables
	WHERE  table_schema = 'public' AND table_type = 'BASE TABLE'
   LOOP
	EXECUTE 'ALTER TABLE '|| tname ||' ADD COLUMN uuid_record uuid';
	EXECUTE 'ALTER TABLE '|| tname ||' ALTER COLUMN uuid_record SET default uuid_generate_v4()';
	EXECUTE 'UPDATE '|| tname ||' SET uuid_record=uuid_generate_v4()';
   END LOOP;

INSERT INTO acs.changes_history(change_uuid, change_date, change_type, change_db) VALUES (uuid_generate_v4(), now(), 'init', current_database());

FOR tname IN
	SELECT quote_ident(table_name)
	FROM   information_schema.tables
	WHERE  table_schema = 'public' AND table_type = 'BASE TABLE'
   LOOP
	FOR ruuid IN EXECUTE 'SELECT uuid_record FROM ' || tname
	LOOP
		EXECUTE 'INSERT INTO acs.record_changes(record_uuid, time_modified, table_name) VALUES('||quote_literal(ruuid)||', now(), '||quote_literal(tname)||')';
	END LOOP;
   END LOOP;

END;
$BODY$
  LANGUAGE plpgsql VOLATILE