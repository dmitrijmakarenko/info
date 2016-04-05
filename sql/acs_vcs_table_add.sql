CREATE OR REPLACE FUNCTION acs_vcs_table_add(text)
  RETURNS void AS
$BODY$
DECLARE
ruuid text;
cnt int;
BEGIN

EXECUTE 'ALTER TABLE '|| $1 ||' ADD COLUMN uuid_record uuid';
EXECUTE 'ALTER TABLE '|| $1 ||' ALTER COLUMN uuid_record SET default uuid_generate_v4()';
EXECUTE 'UPDATE '|| $1 ||' SET uuid_record=uuid_generate_v4()';

FOR ruuid IN EXECUTE 'SELECT uuid_record FROM ' || $1
	LOOP
		EXECUTE 'INSERT INTO acs.record_changes(record_uuid, time_modified, table_name) VALUES('||quote_literal(ruuid)||', now(), '||quote_literal($1)||')';
	END LOOP;

EXECUTE 'CREATE TRIGGER t_acs_'|| $1 ||'
AFTER INSERT OR UPDATE OR DELETE ON '|| $1 ||' FOR EACH ROW
EXECUTE PROCEDURE acs_tg_audit()';

EXECUTE 'SELECT COUNT(*) FROM acs.tables WHERE table_name='|| quote_literal($1) INTO cnt;
IF cnt = 0 THEN
	EXECUTE 'INSERT INTO acs.tables(table_name, schema_name) VALUES('|| quote_literal($1) ||', '|| quote_literal('public') ||')';
END IF;

END;
$BODY$
  LANGUAGE plpgsql VOLATILE;