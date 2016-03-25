CREATE OR REPLACE FUNCTION acs_copy_from_file()
 RETURNS void AS
$BODY$
DECLARE
json_data json;
item json;
column_name text;
column_type text;
table_name text;
structure text;
BEGIN
--temp table
DROP TABLE IF EXISTS acs.transfer_data_t;
CREATE TABLE acs.transfer_data_t (
  tname text NOT NULL primary key,
  data json
);

COPY acs.transfer_data_t FROM '/tmp/transfer.data';

FOR table_name,json_data IN SELECT tname,data FROM acs.transfer_data_t
   LOOP
	--stucture table
	structure = '';
	FOR column_name, column_type IN EXECUTE 'SELECT c.column_name, c.data_type 
	FROM information_schema.tables t JOIN information_schema.columns c ON t.table_name = c.table_name 
	WHERE t.table_schema = '|| quote_literal('public') ||' AND t.table_catalog = current_database() AND t.table_name = ' || quote_literal(table_name)
	LOOP
		structure = structure || column_name || ' ' || column_type || ',';
	END LOOP;
	structure = substr(structure, 0, char_length(structure));
	RAISE notice 'structure %', structure;
	--insert values
	FOR item IN SELECT * FROM json_array_elements(json_data)
	LOOP
		EXECUTE 'INSERT INTO '|| table_name ||' SELECT * FROM json_to_record('|| quote_literal(item) ||') as x(' || structure || ')';
	END LOOP;
   END LOOP;

DROP TABLE IF EXISTS acs.transfer_data_t;

END;
$BODY$
 LANGUAGE plpgsql;