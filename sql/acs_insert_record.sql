--$1 - uuid record
--$2 - user
--$3 - table
CREATE OR REPLACE FUNCTION acs_insert_record(uuid, text, text)
 RETURNS void AS
$BODY$
DECLARE
use_rule bool;
tmp_rec bool;
rule_uuid uuid;
time_temp interval;
table_rule uuid;
BEGIN
rule_uuid = uuid_generate_v4();
tmp_rec = false;
use_rule = false;
--check temp rules
SELECT temp_time INTO time_temp FROM acs.user_rules WHERE user_id=$2 AND temp_use=true AND (table_name=$3 OR table_all=true);
IF time_temp IS NOT NULL THEN
	tmp_rec = true;
	use_rule = true;
END IF;
--table rules
IF NOT tmp_rec THEN
	SELECT security_rule INTO table_rule FROM acs.user_rules WHERE user_id=$2 AND temp_use=false AND (table_name=$3 OR table_all=true);
	IF table_rule IS NOT NULL THEN
		rule_uuid = table_rule;
		use_rule = true;
	END IF;
END IF;

IF use_rule THEN
	INSERT INTO acs.rule_record(uuid_record, security_rule) VALUES ($1, rule_uuid);
	IF tmp_rec THEN
		INSERT INTO acs.rules_data(security_rule, rule_user, rule_action, temp_label, temp_time)
		VALUES (rule_uuid, $2, 'r', 't', now()+time_temp);
	ELSE
		--INSERT INTO acs.rules_data(security_rule, rule_user, rule_action)
		--VALUES (rule_uuid, $2, 'r');
	END IF;
END IF;

END;
$BODY$
 LANGUAGE plpgsql;