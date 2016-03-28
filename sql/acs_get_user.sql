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