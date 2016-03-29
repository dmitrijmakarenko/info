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