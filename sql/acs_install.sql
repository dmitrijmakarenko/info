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
	record_id SERIAL PRIMARY KEY,
	id text NOT NULL,
	pass text NOT NULL,
	position_user text,
	realname text
);
--groups
CREATE TABLE IF NOT EXISTS acs.groups (
	record_id SERIAL PRIMARY KEY,
	group_id text NOT NULL,
	realname text
);
--group-user
CREATE TABLE IF NOT EXISTS acs.group_user (
	record_id SERIAL PRIMARY KEY,
	group_id uuid NOT NULL,
	user_id text NOT NULL
);
--group-struct
CREATE TABLE IF NOT EXISTS acs.groups_struct (
	record_id SERIAL PRIMARY KEY,
	group_id uuid NOT NULL,
	parent_id uuid,
	level integer
);
--rules
CREATE TABLE IF NOT EXISTS acs.rules (
	record_id SERIAL PRIMARY KEY,
	rule_id uuid NOT NULL,
	rule_desc text
);
--rules-data
CREATE TABLE IF NOT EXISTS acs.rules_data (
	record_id SERIAL PRIMARY KEY,
	rule_id uuid NOT NULL,
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
CREATE TABLE acs.changes_history
(
  change_uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  change_date timestamp without time zone NOT NULL,
  change_type text,
  change_db text,
  hash text
);

END;
$BODY$
  LANGUAGE plpgsql VOLATILE