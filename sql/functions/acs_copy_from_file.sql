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
  ttype text,
  data json
);

COPY acs.transfer_data_t FROM '/tmp/transfer.data';

FOR table_name,json_data IN SELECT tname,data FROM acs.transfer_data_t WHERE ttype='data'
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

--import from acs
SELECT data INTO json_data FROM acs.transfer_data_t WHERE table_name='acs.changes_history';
FOR item IN SELECT * FROM json_array_elements(json_data)
LOOP
	INSERT INTO acs.changes_history SELECT * FROM json_to_record(item) AS x(change_uuid uuid,change_parent uuid,change_date timestamp,change_type text,change_db text,hash text);
END LOOP;

SELECT data INTO json_data FROM acs.transfer_data_t WHERE table_name='acs.changes_fields';
FOR item IN SELECT * FROM json_array_elements(json_data)
LOOP
	INSERT INTO acs.changes_fields SELECT * FROM json_to_record(item) AS x(db_name text,record_uuid uuid,change_uuid uuid,table_name text);
END LOOP;

DROP TABLE IF EXISTS acs.transfer_data_t;

END;
$BODY$
 LANGUAGE plpgsql;