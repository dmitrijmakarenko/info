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