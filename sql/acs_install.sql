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
CREATE TABLE IF NOT EXISTS acs.users
(
	id text NOT NULL,
	pass text NOT NULL,
	position_user text,
	realname text
);
--user rule
CREATE TABLE IF NOT EXISTS acs.user_rule
(
	user_id text NOT NULL,
	table_name text,
	security_label uuid,
	temp_label text,
	temp_time interval
);
--groups
CREATE TABLE IF NOT EXISTS acs.groups
(
	group_id text NOT NULL,
	realname text
);
--group-user
CREATE TABLE IF NOT EXISTS acs.group_user
(
	group_id uuid NOT NULL,
	user_id text NOT NULL
);
--group-struct
CREATE TABLE IF NOT EXISTS acs.groups_struct
(
	group_id uuid NOT NULL,
	parent_id uuid,
	level integer
);
--record_rule
CREATE TABLE IF NOT EXISTS acs.rule_record
(
	uuid_record uuid NOT NULL,
	security_rule uuid NOT NULL
);
--rules
CREATE TABLE IF NOT EXISTS acs.rules
(
	security_rule uuid NOT NULL,
	rule_desc text
);
--rules-data
CREATE TABLE IF NOT EXISTS acs.rules_data
(
	security_rule uuid NOT NULL,
	rule_user text,
	rule_action text,
	rule_group text,
	temp_label text,
	temp_time timestamp
);
--tokens
CREATE TABLE IF NOT EXISTS acs.tokens
(
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
CREATE TABLE IF NOT EXISTS acs.tables
(
	table_name text NOT NULL,
	schema_name text NOT NULL
);

END;
$BODY$
  LANGUAGE plpgsql VOLATILE;