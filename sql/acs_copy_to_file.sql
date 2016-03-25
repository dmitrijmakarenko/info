CREATE OR REPLACE FUNCTION acs_copy_to_file()
 RETURNS void AS
$BODY$
DECLARE
json_data json;
tname text;
BEGIN

DROP TABLE IF EXISTS acs.transfer_data;

CREATE TABLE acs.transfer_data (
  id serial primary key,
  tname text NOT NULL,
  data json
);

FOR tname IN SELECT table_name FROM acs.vcs_tables
   LOOP
	EXECUTE 'SELECT json_agg(t) FROM (SELECT * FROM '|| tname ||') t' INTO json_data;
	INSERT INTO acs.transfer_data(tname, data) VALUES(tname,json_data);
	--RAISE notice 'json data %', json_data;
   END LOOP;

COPY acs.transfer_data TO '/tmp/transfer_data.copy';

--DROP TABLE IF EXISTS acs.transfer_data;

END;
$BODY$
 LANGUAGE plpgsql;