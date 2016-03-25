CREATE OR REPLACE FUNCTION acs_copy_from_file()
 RETURNS void AS
$BODY$
DECLARE
json_data json;
item json;
r record;
table_name text;
BEGIN

DROP TABLE IF EXISTS acs.transfer_data_t;

CREATE TABLE acs.transfer_data_t (
  id serial primary key,
  tname text NOT NULL,
  data json
);

COPY acs.transfer_data_t FROM '/tmp/transfer_data.copy';

FOR table_name,json_data IN SELECT tname,data FROM acs.transfer_data_t
   LOOP
	--RAISE notice 'tname %', table_name;
	FOR item IN SELECT * FROM json_array_elements(json_data)
	LOOP
		RAISE NOTICE 'item %', item;
		EXECUTE 'INSERT INTO '|| table_name ||' SELECT * FROM json_to_record('|| quote_literal(item) ||')';
	END LOOP;
   END LOOP;

--DROP TABLE IF EXISTS acs.transfer_data_t;

END;
$BODY$
 LANGUAGE plpgsql;