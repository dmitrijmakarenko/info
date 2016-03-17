CREATE OR REPLACE FUNCTION acs_protect_table (TEXT)
 RETURNS void AS
$BODY$
DECLARE
user_name text;
BEGIN

--EXECUTE 'ALTER TABLE '|| $1 ||' ADD COLUMN uuid_record uuid';
--EXECUTE 'ALTER TABLE '|| $1 ||' ALTER COLUMN uuid_record SET default uuid_generate_v4()';
--EXECUTE 'UPDATE '|| $1 ||' SET uuid_record=uuid_generate_v4()';
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