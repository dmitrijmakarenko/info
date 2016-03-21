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
  LANGUAGE plpgsql VOLATILE