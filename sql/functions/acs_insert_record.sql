--$1 - uuid record
--$2 - user
CREATE OR REPLACE FUNCTION acs_insert_record(uuid, text)
 RETURNS void AS
$BODY$
DECLARE
tmp_rec bool;
rule_uuid uuid;
BEGIN

rule_uuid = uuid_generate_v4();

tmp_rec = true;

INSERT INTO acs.rule_record(uuid_record, security_rule) VALUES ($1, rule_uuid);

IF tmp_rec THEN
	INSERT INTO acs.rules_data(security_rule, rule_user, rule_action, temp_label, temp_time) 
	VALUES (rule_uuid, $2, 'r', 't', now()+interval '1' day);
ELSE
	INSERT INTO acs.rules_data(security_rule, rule_user, rule_action) 
	VALUES (rule_uuid, $2, 'r');
END IF;

END;
$BODY$
 LANGUAGE plpgsql;