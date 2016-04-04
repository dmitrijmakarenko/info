CREATE OR REPLACE FUNCTION acs_vcs_compile()
 RETURNS void AS
$BODY$
DECLARE
tname text;
ruuid text;
uuid_change uuid;
cdate timestamp;
cuuid uuid;
r record;
data text;
hash text;
BEGIN

SELECT change_date,change_uuid INTO cdate,cuuid FROM acs.changes_history ORDER BY change_date DESC LIMIT 1;
IF cdate IS NULL THEN
	RETURN;
END IF;

data = '';
FOR tname IN SELECT table_name FROM acs.tables
   LOOP
	data = data || tname;
	FOR r IN EXECUTE 'SELECT '|| tname ||'.* FROM '|| tname ||' LEFT OUTER JOIN acs.record_changes ON ('|| tname ||'.uuid_record = acs.record_changes.record_uuid) WHERE acs.record_changes.time_modified >= '|| quote_literal(cdate)
	LOOP
		data = data || array_to_string(array_agg(r),',','*');
	END LOOP;
   END LOOP;

hash = md5(data);
uuid_change = uuid_generate_v4();

FOR tname IN SELECT table_name FROM acs.tables
   LOOP
	FOR ruuid IN EXECUTE 'SELECT '|| tname ||'.uuid_record FROM '|| tname ||' LEFT OUTER JOIN acs.record_changes ON ('|| tname ||'.uuid_record = acs.record_changes.record_uuid) WHERE acs.record_changes.time_modified >= '|| quote_literal(cdate)
	LOOP
		--RAISE notice 'uuid record %', ruuid;
		EXECUTE 'INSERT INTO acs.changes_fields(db_name,record_uuid,change_uuid,table_name) VALUES(current_database(), '|| quote_literal(ruuid) ||', '|| quote_literal(uuid_change) ||', '|| quote_literal(tname) ||')';
	END LOOP;
   END LOOP;

INSERT INTO acs.changes_history(change_uuid,change_parent,change_date,change_type,change_db,hash) VALUES (uuid_change,cuuid,now(),'compile',current_database(),hash);

END;
$BODY$
 LANGUAGE plpgsql;