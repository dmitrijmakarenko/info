CREATE OR REPLACE FUNCTION acs_check_user(text, text)
  RETURNS boolean AS
$BODY$
DECLARE
cnt int;
BEGIN
SELECT COUNT(*) INTO cnt AS cnt FROM acs_tokens WHERE user_id=$1 AND token=$2;
IF cnt > 0 THEN
	RETURN TRUE;
ELSE
	RETURN FALSE;
END IF;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE