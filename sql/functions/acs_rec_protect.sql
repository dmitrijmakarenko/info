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

SELECT COUNT(*) INTO cnt FROM acs.rule_record WHERE uuid_record=$2;
IF cnt = 0 THEN
	INSERT INTO acs.rule_record(uuid_record, security_rule) VALUES ($2, $1);
ELSE
	UPDATE acs.rule_record SET security_rule=$1 WHERE uuid_record=$2;
END IF;

END;
$BODY$
 LANGUAGE plpgsql;