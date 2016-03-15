CREATE OR REPLACE FUNCTION acs_protect_table (TEXT)
 RETURNS void AS
$BODY$
DECLARE
BEGIN

EXECUTE 'ALTER TABLE '|| $1 ||' ADD COLUMN rule uuid';
EXECUTE 'ALTER TABLE '|| $1 ||' RENAME TO ' || $1 || '_protected';
EXECUTE 'CREATE OR REPLACE VIEW  '|| $1 ||' AS SELECT * FROM ' || $1 || '_protected';

END;
$BODY$
 LANGUAGE plpgsql;