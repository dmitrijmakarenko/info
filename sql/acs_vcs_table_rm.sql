CREATE OR REPLACE FUNCTION acs_vcs_table_rm(text)
  RETURNS void AS
$BODY$
DECLARE
BEGIN

EXECUTE 'ALTER TABLE '|| $1 ||' DROP COLUMN IF EXISTS uuid_record';
EXECUTE 'DELETE FROM acs.record_changes WHERE table_name='|| quote_literal($1);
EXECUTE 'DELETE FROM acs.vcs_tables WHERE table_name='|| quote_literal($1);
EXECUTE 'DROP TRIGGER IF EXISTS t_acs_'|| $1 ||' ON ' || $1;

END;
$BODY$
  LANGUAGE plpgsql VOLATILE