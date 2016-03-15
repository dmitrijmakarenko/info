CREATE OR REPLACE FUNCTION acs_protect_table (TEXT)
 RETURNS void AS
$BODY$
DECLARE
user_name text;
BEGIN

EXECUTE 'ALTER TABLE '|| $1 ||' ADD COLUMN rule uuid';
EXECUTE 'ALTER TABLE '|| $1 ||' RENAME TO ' || $1 || '_protected';
EXECUTE 'CREATE OR REPLACE VIEW  '|| $1 ||' AS SELECT * FROM ' || $1 || '_protected';

FOR user_name IN
	SELECT usename FROM pg_user
   LOOP
	EXECUTE 'GRANT ALL PRIVILEGES ON ' || $1 || ' TO ' || user_name;
	EXECUTE 'REVOKE ALL PRIVILEGES ON ' || $1 || '_protected FROM ' || user_name;
   END LOOP;

END;
$BODY$
 LANGUAGE plpgsql;