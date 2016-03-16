CREATE OR REPLACE FUNCTION acs_vcs_init()
  RETURNS void AS
$BODY$
DECLARE
tname text;
BEGIN

FOR tname IN
	SELECT quote_ident(table_name)
	FROM   information_schema.tables
	WHERE  table_schema = 'public' AND table_type = 'BASE TABLE'
   LOOP
	EXECUTE 'ALTER TABLE '|| tname ||' ADD COLUMN time_modified timestamp';
	EXECUTE 'ALTER TABLE '|| tname ||' ALTER COLUMN time_modified SET default current_timestamp';
	EXECUTE 'UPDATE '|| tname ||' SET time_modified=current_timestamp';
   END LOOP;

INSERT INTO acs.changes_history(change_uuid, change_date, change_type, change_db) VALUES (uuid_generate_v4(), now(), 'init', current_database());

END;
$BODY$
  LANGUAGE plpgsql VOLATILE