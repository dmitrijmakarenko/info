-- name: install-functions
CREATE OR REPLACE FUNCTION acs_auth(TEXT, TEXT)
 RETURNS TEXT AS
$BODY$
DECLARE
cnt int;
token text;
BEGIN
token = '';
SELECT COUNT(*) INTO cnt FROM acs.users WHERE id=$1 AND pass=crypt($2, pass);
IF cnt > 0 THEN
	token = uuid_generate_v4();
	INSERT INTO acs.tokens(user_id, token, exp_date) VALUES ($1, token, now() + interval '1' day);
	DELETE FROM acs.tokens WHERE user_id=$1 AND exp_date < now();
END IF;

RETURN token;
END;
$BODY$
 LANGUAGE plpgsql;

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
cuuid uuid;
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

--import changed fields
SELECT data INTO json_data FROM acs.transfer_data_t WHERE tname='acs.changes_fields';
FOR item IN SELECT * FROM json_array_elements(json_data)
LOOP
	INSERT INTO acs.changes_fields SELECT * FROM json_to_record(item) AS x(db_name text,record_uuid uuid,change_uuid uuid,table_name text);
END LOOP;
--import history
SELECT data INTO json_data FROM acs.transfer_data_t WHERE tname='acs.changes_history';
FOR item IN SELECT * FROM json_array_elements(json_data)
LOOP
	INSERT INTO acs.changes_history SELECT * FROM json_to_record(item) AS x(change_uuid uuid,change_parent uuid,change_date timestamp,change_type text,change_db text,hash text);
	SELECT change.change_uuid INTO cuuid FROM (SELECT * FROM json_to_record(item) AS x(change_uuid uuid,change_parent uuid,change_date timestamp,change_type text,change_db text,hash text)) AS change;
	INSERT INTO acs.changes_history(change_parent,change_date,change_type,change_db) VALUES(cuuid,now(),'commit',current_database());
END LOOP;

DROP TABLE IF EXISTS acs.transfer_data_t;

END;
$BODY$
 LANGUAGE plpgsql;

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

CREATE OR REPLACE FUNCTION acs_get_user(text)
  RETURNS text AS
$BODY$
DECLARE
cnt int;
user_auth text;
BEGIN
user_auth = '';
SELECT user_id INTO user_auth FROM acs.tokens WHERE token=$1 AND exp_date >= now();

RETURN user_auth;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE;

CREATE OR REPLACE FUNCTION acs_install()
  RETURNS void AS
$BODY$
DECLARE
cnt int;
BEGIN

SELECT COUNT(*) INTO cnt FROM information_schema.schemata WHERE schema_name = 'acs';
IF cnt = 0 THEN
	CREATE SCHEMA acs;
END IF;

--users
CREATE TABLE IF NOT EXISTS acs.users (
	uuid_record uuid NOT NULL DEFAULT uuid_generate_v4(),
	id text NOT NULL,
	pass text NOT NULL,
	position_user text,
	realname text
);
--groups
CREATE TABLE IF NOT EXISTS acs.groups (
	uuid_record uuid NOT NULL DEFAULT uuid_generate_v4(),
	group_id text NOT NULL,
	realname text
);
--group-user
CREATE TABLE IF NOT EXISTS acs.group_user (
	uuid_record uuid NOT NULL DEFAULT uuid_generate_v4(),
	group_id uuid NOT NULL,
	user_id text NOT NULL
);
--group-struct
CREATE TABLE IF NOT EXISTS acs.groups_struct (
	uuid_record uuid NOT NULL DEFAULT uuid_generate_v4(),
	group_id uuid NOT NULL,
	parent_id uuid,
	level integer
);
--record_rule
CREATE TABLE IF NOT EXISTS acs.rule_record (
	uuid_record uuid NOT NULL,
	security_rule uuid NOT NULL
);
--rules
CREATE TABLE IF NOT EXISTS acs.rules (
	uuid_record uuid NOT NULL DEFAULT uuid_generate_v4(),
	security_rule uuid NOT NULL,
	rule_desc text
);
--rules-data
CREATE TABLE IF NOT EXISTS acs.rules_data (
	uuid_record uuid NOT NULL DEFAULT uuid_generate_v4(),
	security_rule uuid NOT NULL,
	rule_user text,
	rule_action text,
	rule_group text
);
--tokens
CREATE TABLE IF NOT EXISTS acs.tokens (
	user_id text NOT NULL,
	token text NOT NULL,
	exp_date timestamp
);
--changes_history
CREATE TABLE IF NOT EXISTS acs.changes_history
(
  change_uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  change_parent uuid,
  change_date timestamp without time zone NOT NULL,
  change_type text,
  change_db text,
  hash text
);
--changes_fields
CREATE TABLE IF NOT EXISTS acs.changes_fields
(
  db_name text NOT NULL,
  record_uuid uuid NOT NULL,
  change_uuid uuid NOT NULL,
  table_name text NOT NULL
);
--record_changes
CREATE TABLE IF NOT EXISTS acs.record_changes
(
  record_uuid uuid NOT NULL,
  time_modified timestamp without time zone NOT NULL,
  table_name text NOT NULL
);
--list tables
CREATE TABLE IF NOT EXISTS acs.vcs_tables
(
  table_name text NOT NULL,
  schema_name text NOT NULL
);

END;
$BODY$
  LANGUAGE plpgsql VOLATILE;

CREATE OR REPLACE FUNCTION acs_protect_table (TEXT)
 RETURNS void AS
$BODY$
DECLARE
user_name text;
BEGIN

--EXECUTE 'ALTER TABLE '|| $1 ||' ADD COLUMN uuid_record uuid';
--EXECUTE 'ALTER TABLE '|| $1 ||' ALTER COLUMN uuid_record SET default uuid_generate_v4()';
--EXECUTE 'UPDATE '|| $1 ||' SET uuid_record=uuid_generate_v4()';
EXECUTE 'ALTER TABLE '|| $1 ||' RENAME TO ' || $1 || '_protected';
EXECUTE 'CREATE OR REPLACE VIEW  '|| $1 ||' AS SELECT * FROM ' || $1 || '_protected';

FOR user_name IN
	SELECT usename FROM pg_user
   LOOP
	EXECUTE 'GRANT ALL PRIVILEGES ON ' || $1 || ' TO ' || user_name;
	EXECUTE 'REVOKE ALL PRIVILEGES ON ' || $1 || '_protected FROM ' || user_name;
   END LOOP;

END;
$BODY$
 LANGUAGE plpgsql;

--$1 - uuid security rule
--$2 - uuid record
--$3 - table name
CREATE OR REPLACE FUNCTION acs_rec_protect(uuid, uuid, text)
 RETURNS void AS
$BODY$
DECLARE
cnt int;
BEGIN
--check rule
SELECT COUNT(*) INTO cnt FROM acs.rules WHERE security_rule=$1;
IF cnt = 0 THEN
	RETURN;
END IF;
--check record
EXECUTE 'SELECT COUNT(*) FROM '|| $3 ||' WHERE uuid_record='|| quote_literal($2) INTO cnt;
IF cnt = 0 THEN
	RETURN;
END IF;

INSERT INTO acs.rule_record(uuid_record, security_rule) VALUES ($2, $1);

END;
$BODY$
 LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION acs_tg_audit()
  RETURNS trigger AS
$BODY$
DECLARE
BEGIN

IF TG_OP = 'INSERT' THEN
	INSERT INTO acs.record_changes(record_uuid, time_modified, table_name) VALUES(NEW.uuid_record, now(), TG_RELNAME);
RETURN NULL;
ELSIF TG_OP = 'UPDATE' THEN
	UPDATE acs.record_changes SET time_modified=now() WHERE record_uuid=NEW.uuid_record;
RETURN NULL;
ELSIF TG_OP = 'DELETE' THEN
	DELETE FROM acs.record_changes WHERE record_uuid=OLD.uuid_record;
RETURN NULL;
END IF;

END;
$BODY$
  LANGUAGE plpgsql VOLATILE;

CREATE OR REPLACE FUNCTION acs_tg_event()
  RETURNS event_trigger AS
$BODY$
DECLARE

BEGIN

IF tg_tag = 'CREATE TABLE' THEN
	RAISE notice 'command %', tg_tag;
END IF;

END;
$BODY$
  LANGUAGE plpgsql VOLATILE;

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
FOR tname IN SELECT table_name FROM acs.vcs_tables
   LOOP
	data = data || tname;
	FOR r IN EXECUTE 'SELECT '|| tname ||'.* FROM '|| tname ||' LEFT OUTER JOIN acs.record_changes ON ('|| tname ||'.uuid_record = acs.record_changes.record_uuid) WHERE acs.record_changes.time_modified >= '|| quote_literal(cdate)
	LOOP
		data = data || array_to_string(array_agg(r),',','*');
	END LOOP;
   END LOOP;

hash = md5(data);
uuid_change = uuid_generate_v4();

FOR tname IN SELECT table_name FROM acs.vcs_tables
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

CREATE OR REPLACE FUNCTION acs_vcs_init()
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

EXECUTE 'SELECT COUNT(*) FROM acs.vcs_tables WHERE table_name='|| quote_literal($1) INTO cnt;
IF cnt = 0 THEN
	EXECUTE 'INSERT INTO acs.vcs_tables(table_name, schema_name) VALUES('|| quote_literal($1) ||', '|| quote_literal('public') ||')';
END IF;

END;
$BODY$
  LANGUAGE plpgsql VOLATILE;

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
  LANGUAGE plpgsql VOLATILE;