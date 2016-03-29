CREATE OR REPLACE FUNCTION acs_copy_to_file()
 RETURNS void AS
$BODY$
DECLARE
cdate timestamp;
json_data json;
tname text;
BEGIN
--get last date
SELECT change_date INTO cdate FROM acs.changes_history ORDER BY change_date DESC LIMIT 1;
IF cdate IS NULL THEN
	RETURN;
END IF;
--temp table
DROP TABLE IF EXISTS acs.transfer_data;
CREATE TABLE acs.transfer_data (
  tname text NOT NULL primary key,
  data json
);

FOR tname IN SELECT table_name FROM acs.vcs_tables
   LOOP
	EXECUTE 'SELECT json_agg(t) FROM (SELECT '|| tname ||'.* FROM '|| tname ||' LEFT OUTER JOIN acs.record_changes ON ('|| tname ||'.uuid_record = acs.record_changes.record_uuid) WHERE acs.record_changes.time_modified >= '|| quote_literal(cdate) ||') t' INTO json_data;
	INSERT INTO acs.transfer_data(tname, data) VALUES(tname,json_data);
   END LOOP;

COPY acs.transfer_data TO '/tmp/transfer.data';

DROP TABLE IF EXISTS acs.transfer_data;

END;
$BODY$
 LANGUAGE plpgsql;