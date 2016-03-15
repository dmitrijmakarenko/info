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
	record_uuid uuid NOT NULL,
	id text NOT NULL,
	position_user text,
	realname text
);
--groups
CREATE TABLE IF NOT EXISTS acs.groups (
	record_uuid uuid NOT NULL,
	group_id text NOT NULL,
	realname text
);

END;
$BODY$
  LANGUAGE plpgsql VOLATILE