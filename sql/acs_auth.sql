CREATE OR REPLACE FUNCTION acs_auth(TEXT, TEXT)
 RETURNS void AS
$BODY$
DECLARE
BEGIN

INSERT INTO acs.tokens(user_id, token, exp_date) VALUES ($1, $2, now() + interval '1' day);

IF NOT EXISTS (SELECT * FROM   pg_catalog.pg_user WHERE  usename = $1) THEN
      EXECUTE 'CREATE USER '|| $1 ||' WITH PASSWORD ' || quote_literal(12345);
END IF;

DELETE FROM acs.tokens WHERE user_id=$1 AND exp_date < now();

END;
$BODY$
 LANGUAGE plpgsql;