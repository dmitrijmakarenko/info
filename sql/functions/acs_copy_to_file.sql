CREATE OR REPLACE FUNCTION acs_copy_to_file()
 RETURNS void AS
$BODY$
DECLARE
cdate timestamp;
json_data json;
cparent uuid;
cuuid uuid;
tname text;
BEGIN
--get change compile
SELECT change_uuid,change_parent INTO cuuid,cparent FROM acs.changes_history WHERE change_type='compile' AND change_parent IS NOT NULL ORDER BY change_date DESC LIMIT 1;
IF cuuid IS NULL THEN
	RETURN;
END IF;
--get date
SELECT change_date INTO cdate FROM acs.changes_history WHERE change_uuid=cparent;
IF cdate IS NULL THEN
	RETURN;
END IF;
--temp table
DROP TABLE IF EXISTS acs.transfer_data;
CREATE TABLE acs.transfer_data (
  tname text NOT NULL primary key,
  ttype text,
  data json
);

FOR tname IN SELECT table_name FROM acs.vcs_tables
   LOOP
	EXECUTE 'SELECT json_agg(t) FROM (SELECT '|| tname ||'.* FROM '|| tname ||' LEFT OUTER JOIN acs.record_changes ON ('|| tname ||'.uuid_record = acs.record_changes.record_uuid) WHERE acs.record_changes.time_modified >= '|| quote_literal(cdate) ||') t' INTO json_data;
	INSERT INTO acs.transfer_data(tname, ttype, data) VALUES(tname, 'data', json_data);
   END LOOP;

SELECT json_agg(t) INTO json_data FROM (SELECT * FROM acs.changes_history WHERE change_uuid=cuuid) t;
INSERT INTO acs.transfer_data(tname, ttype, data) VALUES('acs.changes_history', 'acs', json_data);
SELECT json_agg(t) INTO json_data FROM (SELECT * FROM acs.changes_fields WHERE change_uuid=cuuid) t;
INSERT INTO acs.transfer_data(tname, ttype, data) VALUES('acs.changes_fields', 'acs', json_data);

COPY acs.transfer_data TO '/tmp/transfer.data';

DROP TABLE IF EXISTS acs.transfer_data;

END;
$BODY$
 LANGUAGE plpgsql;