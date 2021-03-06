﻿CREATE OR REPLACE FUNCTION acs_vcs_init()
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
	EXECUTE 'SELECT acs_vcs_table_add('|| quote_literal(tname) ||')';
   END LOOP;

INSERT INTO acs.changes_history(change_uuid, change_date, change_type, change_db) VALUES (uuid_generate_v4(), now(), 'init', current_database());

DROP EVENT TRIGGER IF EXISTS acs_tg_event;
CREATE EVENT TRIGGER acs_tg_event ON ddl_command_end
   EXECUTE PROCEDURE acs_tg_event();

END;
$BODY$
  LANGUAGE plpgsql VOLATILE;