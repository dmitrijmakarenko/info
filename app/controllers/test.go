package controllers

import (
	"github.com/robfig/revel"
	"database/sql"
	"strconv"
	"time"
	"os"
	"errors"
)

type TestCntl struct {
	*revel.Controller
}

const CNTRECORDS = 100

var dbGeneral *sql.DB
var dbStat1 *sql.DB

func Connect() {
	var err error
	//connect to general
	connstring := "host=" + DB_HOST + " port=" + DB_PORT + " user=" + DB_USER+ " dbname=general password=" + DB_PASSWORD + " sslmode=disable"
	dbGeneral, err = sql.Open(DB_DRIVER, connstring)
	if err != nil {
		revel.ERROR.Println(err.Error())
	}
	//connect to station1
	connstring = "host=" + DB_HOST + " port=" + DB_PORT + " user=" + DB_USER+ " dbname=station1 password=" + DB_PASSWORD + " sslmode=disable"
	dbStat1, err = sql.Open(DB_DRIVER, connstring)
	if err != nil {
		revel.ERROR.Println(err.Error())
	}
}

func (c TestCntl) Reset() revel.Result {
	var err error
	ret := make(map[string]string)

	Connect()

	//drop db
	_, err = DB.Exec("DROP DATABASE IF EXISTS general")
	_, err = DB.Exec("DROP DATABASE IF EXISTS station1")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	//create db
	_, err = DB.Exec("CREATE DATABASE general WITH OWNER = postgres ENCODING = 'UTF8' TABLESPACE = pg_default LC_COLLATE = 'ru_RU.UTF-8' LC_CTYPE = 'ru_RU.UTF-8' CONNECTION LIMIT = -1")
	_, err = DB.Exec("CREATE DATABASE station1 WITH OWNER = postgres ENCODING = 'UTF8' TABLESPACE = pg_default LC_COLLATE = 'ru_RU.UTF-8' LC_CTYPE = 'ru_RU.UTF-8' CONNECTION LIMIT = -1")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	//create data tables
	_, err = dbGeneral.Exec("CREATE TABLE persons(personId int, lastname text, firstname text, address text, city text)")
	_, err = dbGeneral.Exec("CREATE TABLE clients(clientId int, name text, city text, telnum text)")
	_, err = dbGeneral.Exec("CREATE TABLE sales(product text, count int, price int)")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	_, err = dbStat1.Exec("CREATE TABLE persons(personId int, lastname text, firstname text, address text, city text)")
	_, err = dbStat1.Exec("CREATE TABLE clients(clientId int, name text, firstname text, city text, telnum text)")
	_, err = dbStat1.Exec("CREATE TABLE sales(product text, count int, price int)")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	return c.RenderJson(ret)
}

func (c TestCntl) Install() revel.Result {
	ret := make(map[string]string)
	err := InstallFunc(DB)
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	return c.RenderJson(ret)
}

func (c TestCntl) Init() revel.Result {
	var err error
	ret := make(map[string]string)

	//install tools
	_, err = dbGeneral.Exec("CREATE EXTENSION \"uuid-ossp\"")
	_, err = dbStat1.Exec("CREATE EXTENSION \"uuid-ossp\"")

	InstallFunc(dbGeneral)
	InstallFunc(dbStat1)

	//install acs
	_, err = dbGeneral.Query("SELECT acs_install()")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	_, err = dbStat1.Query("SELECT acs_install()")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	//init
	_, err = dbGeneral.Query("SELECT acs_vcs_init()")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	_, err = dbStat1.Query("SELECT acs_vcs_init()")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	return c.RenderJson(ret)
}

func (c TestCntl) Work() revel.Result {
	ret := make(map[string]string)
	var err error

	_, err = dbStat1.Exec("INSERT INTO persons VALUES(1, 'Иванов', 'Иван', 'ул.Парфенова', 'г.Тула')")
	_, err = dbStat1.Exec("INSERT INTO persons VALUES(2, 'Петров', 'Петр', 'ул.Ленина', 'г.Томск')")
	_, err = dbStat1.Exec("INSERT INTO clients VALUES(1, 'Автозапчасти', 'г.Москва', '8 (906) 531-85-86')")
	_, err = dbStat1.Exec("INSERT INTO clients VALUES(2, 'Продукты', 'г.Тула', '8 (906) 531-85-86')")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	return c.RenderJson(ret)
}

func (c TestCntl) Compile() revel.Result {
	ret := make(map[string]string)
	var err error

	//compile
	_, err = dbStat1.Query("SELECT acs_vcs_compile()")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	return c.RenderJson(ret)
}

func (c TestCntl) CopyToFile() revel.Result {
	ret := make(map[string]string)
	var err error

	_, err = dbStat1.Query("SELECT acs_copy_to_file()")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	return c.RenderJson(ret)
}

func (c TestCntl) CopyFromFile() revel.Result {
	ret := make(map[string]string)
	var err error

	_, err = dbGeneral.Query("SELECT acs_copy_from_file()")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	return c.RenderJson(ret)
}

func GetTest(table string, token string) (error, time.Duration) {
	var p DataGetParams
	var hasRights bool
	var err error

	p.Table = table;
	p.Token = token;

	//check user
	var user string
	revel.INFO.Println("[TEST] start query...")
	err = DB.QueryRow("SELECT COALESCE(acs_get_user($1), '')", p.Token).Scan(&user)
	if err != nil {
		return err, 0
	}

	if user == "" {
		revel.ERROR.Println("[TEST] authorization error...")
		return errors.New("authorization error"), 0
	}

	//get rights
	rows, err := DB.Query("SELECT DISTINCT security_rule FROM "+TABLE_RULES_DATA+" WHERE rule_user=$1 AND rule_action='r'", user)
	if err != nil {
		return err, 0
	}
	var rules []interface{}
	for rows.Next() {
		var rule string
		err := rows.Scan(&rule)
		if err != nil {
			return err, 0
		} else {
			rules = append(rules, rule)
		}
	}
	if len(rules) > 0 {
		hasRights = true
	} else {
		hasRights = false
	}
	var rulesList string
	rulesList += "("
	for i, _ := range rules {
		rulesList += "$" + strconv.Itoa(i+1)
		if (i != len(rules) - 1) {
			rulesList += ", "
		}
	}
	rulesList += ")"

	//sql query
	start := time.Now()
	var stmt *sql.Stmt
	if (hasRights) {
		stmt, err = DB.Prepare("SELECT "+p.Table+".*  FROM "+p.Table+" LEFT OUTER JOIN acs.rule_record AS rules ON ("+p.Table+".uuid_record = rules.uuid_record) WHERE security_rule IS NULL OR security_rule IN " + rulesList)
	} else {
		stmt, err = DB.Prepare("SELECT "+p.Table+".*  FROM "+p.Table+" LEFT OUTER JOIN acs.rule_record AS rules ON ("+p.Table+".uuid_record = rules.uuid_record) WHERE security_rule IS NULL")
	}
	if err != nil {
		return err, 0
	}
	if (hasRights) {
		rows, err = stmt.Query(rules...)
	} else {
		rows, err = stmt.Query()
	}
	if err != nil {
		return err, 0
	}
	elapsed := time.Since(start)
	revel.INFO.Println("[TEST] elapsed time", elapsed.Seconds())

	return nil, elapsed
}

func UpdateTest(table string, token string) (error, time.Duration) {
	var p DataGetParams
	var hasRights bool
	var err error

	p.Table = table;
	p.Token = token;

	//check user
	var user string
	revel.INFO.Println("[TEST] start query...")
	err = DB.QueryRow("SELECT COALESCE(acs_get_user($1), '')", p.Token).Scan(&user)
	if err != nil {
		return err, 0
	}

	if user == "" {
		return errors.New("authorization error"), 0
	}

	//get rights
	rows, err := DB.Query("SELECT DISTINCT security_rule FROM "+TABLE_RULES_DATA+" WHERE rule_user=$1 AND rule_action='r'", user)
	if err != nil {
		return err, 0
	}
	var rules []interface{}
	for rows.Next() {
		var rule string
		err := rows.Scan(&rule)
		if err != nil {
			return err, 0
		} else {
			rules = append(rules, rule)
		}
	}
	if len(rules) > 0 {
		hasRights = true
	} else {
		hasRights = false
	}
	var rulesList string
	rulesList += "("
	for i, _ := range rules {
		rulesList += "$" + strconv.Itoa(i+1)
		if (i != len(rules) - 1) {
			rulesList += ", "
		}
	}
	rulesList += ")"
	revel.INFO.Println("[TEST] rules", rules);

	//sql query
	start := time.Now()
	var stmt *sql.Stmt
	if (hasRights) {
		stmt, err = DB.Prepare("UPDATE "+p.Table+" SET val_id = -1 FROM (SELECT t.* FROM "+p.Table+" t LEFT OUTER JOIN acs.rule_record AS rules ON (t.uuid_record = rules.uuid_record) WHERE security_rule IS NULL OR security_rule IN " + rulesList + ") AS d WHERE d.uuid_record = test_data.uuid_record")
	} else {
		stmt, err = DB.Prepare("UPDATE "+p.Table+" SET val_id = -1 FROM (SELECT t.* FROM "+p.Table+" t LEFT OUTER JOIN acs.rule_record AS rules ON (t.uuid_record = rules.uuid_record) WHERE security_rule IS NULL) AS d WHERE d.uuid_record = test_data.uuid_record")
	}
	if err != nil {
		return err, 0
	}
	if (hasRights) {
		_, err = stmt.Exec(rules...)
	} else {
		_, err = stmt.Exec()
	}
	if err != nil {
		return err, 0
	}
	elapsed := time.Since(start)
	revel.INFO.Println("[TEST] elapsed time", elapsed.Seconds());

	return nil, elapsed
}

func DeleteTest(table string, token string) (error, time.Duration) {
	var p DataGetParams
	var hasRights bool
	var err error

	p.Table = table;
	p.Token = token;

	//check user
	var user string
	revel.INFO.Println("[TEST] start query...")
	err = DB.QueryRow("SELECT COALESCE(acs_get_user($1), '')", p.Token).Scan(&user)
	if err != nil {
		return err, 0
	}

	if user == "" {
		return errors.New("authorization error"), 0
	}

	//get rights
	rows, err := DB.Query("SELECT DISTINCT security_rule FROM "+TABLE_RULES_DATA+" WHERE rule_user=$1 AND rule_action='r'", user)
	if err != nil {
		return err, 0
	}
	var rules []interface{}
	for rows.Next() {
		var rule string
		err := rows.Scan(&rule)
		if err != nil {
			return err, 0
		} else {
			rules = append(rules, rule)
		}
	}
	if len(rules) > 0 {
		hasRights = true
	} else {
		hasRights = false
	}
	var rulesList string
	rulesList += "("
	for i, _ := range rules {
		rulesList += "$" + strconv.Itoa(i+1)
		if (i != len(rules) - 1) {
			rulesList += ", "
		}
	}
	rulesList += ")"
	revel.INFO.Println("[TEST] rules", rules);

	//sql query
	start := time.Now()
	var stmt *sql.Stmt
	if (hasRights) {
		stmt, err = DB.Prepare("DELETE FROM "+p.Table+" USING "+p.Table+" AS t LEFT OUTER JOIN acs.rule_record AS rules ON t.uuid_record = rules.uuid_record WHERE security_rule IS NULL OR security_rule IN " + rulesList + " AND t.val_id = 1")
	} else {
		stmt, err = DB.Prepare("DELETE FROM "+p.Table+" USING "+p.Table+" AS t LEFT OUTER JOIN acs.rule_record AS rules ON t.uuid_record = rules.uuid_record WHERE security_rule IS NULL AND t.val_id = 1")
	}
	if err != nil {
		return err, 0
	}
	if (hasRights) {
		_, err = stmt.Exec(rules...)
	} else {
		_, err = stmt.Exec()
	}
	if err != nil {
		return err, 0
	}
	elapsed := time.Since(start)
	revel.INFO.Println("[TEST] elapsed time", elapsed.Seconds());

	return nil, elapsed
}

func SelectDataStandartProcess(cntRecords int) (error, time.Duration) {
	var err error
	_, err = DB.Exec("DROP TABLE IF EXISTS test_data")
	if err != nil {
		return err, 0
	}
	_, err = DB.Exec("CREATE TABLE test_data(val_id int, field1 text, field2 text, field3 text, field4 text, field5 text)")
	for i := 0; i < cntRecords; i++ {
		_, err = DB.Exec("INSERT INTO test_data VALUES($1,$2,$3,$4,$5,$6)", i, "d1_"+strconv.Itoa(i), "d2_"+strconv.Itoa(i), "d3_"+strconv.Itoa(i), "d4_"+strconv.Itoa(i), "d5_"+strconv.Itoa(i))
		if err != nil {
			return err, 0
		}
	}
	revel.INFO.Println("[TEST] start query...")
	start := time.Now()
	_, err = DB.Query("SELECT * FROM test_data")
	elapsed := time.Since(start)
	revel.INFO.Println("[TEST] elapsed time", elapsed.Seconds());
	return err, elapsed
}

func SelectDataSecureProcess(cntRecords int) (error, time.Duration) {
	var err error
	_, err = DB.Exec("DROP TABLE IF EXISTS test_data")
	if err != nil {
		return err, 0
	}
	_, err = DB.Exec("CREATE TABLE test_data(val_id int, field1 text, field2 text, field3 text, field4 text, field5 text)")
	for i := 0; i < cntRecords; i++ {
		_, err = DB.Exec("INSERT INTO test_data VALUES($1,$2,$3,$4,$5,$6)", i, "d1_"+strconv.Itoa(i), "d2_"+strconv.Itoa(i), "d3_"+strconv.Itoa(i), "d4_"+strconv.Itoa(i), "d5_"+strconv.Itoa(i))
		if err != nil {
			return err, 0
		}
	}
	_, err = DB.Query("SELECT acs_vcs_table_add('test_data')")
	if err != nil {
		return err, 0
	}
	err, elapsed := GetTest("test_data", "25b9acf4-3694-4739-830e-6df25e80fe33")
	if err != nil {
		return err, 0
	}

	return err, elapsed
}

func UpdateDataStandartProcess(cntRecords int) (error, time.Duration) {
	var err error
	_, err = DB.Exec("DROP TABLE IF EXISTS test_data")
	if err != nil {
		return err, 0
	}
	_, err = DB.Exec("CREATE TABLE test_data(val_id int, field1 text, field2 text, field3 text, field4 text, field5 text)")
	for i := 0; i < cntRecords; i++ {
		_, err = DB.Exec("INSERT INTO test_data VALUES($1,$2,$3,$4,$5,$6)", i, "d1_"+strconv.Itoa(i), "d2_"+strconv.Itoa(i), "d3_"+strconv.Itoa(i), "d4_"+strconv.Itoa(i), "d5_"+strconv.Itoa(i))
		if err != nil {
			return err, 0
		}
	}
	revel.INFO.Println("[TEST] start query...")
	start := time.Now()
	_, err = DB.Query("UPDATE test_data SET val_id = -1 WHERE val_id < $1", cntRecords/2)
	elapsed := time.Since(start)
	revel.INFO.Println("[TEST] elapsed time", elapsed.Seconds());
	return err, elapsed
}

func UpdateDataSecureProcess(cntRecords int) (error, time.Duration) {
	var err error
	_, err = DB.Exec("DROP TABLE IF EXISTS test_data")
	if err != nil {
		return err, 0
	}
	_, err = DB.Exec("CREATE TABLE test_data(val_id int, field1 text, field2 text, field3 text, field4 text, field5 text)")
	for i := 0; i < cntRecords; i++ {
		_, err = DB.Exec("INSERT INTO test_data VALUES($1,$2,$3,$4,$5,$6)", i, "d1_"+strconv.Itoa(i), "d2_"+strconv.Itoa(i), "d3_"+strconv.Itoa(i), "d4_"+strconv.Itoa(i), "d5_"+strconv.Itoa(i))
		if err != nil {
			return err, 0
		}
	}
	_, err = DB.Query("SELECT acs_vcs_table_add('test_data')")
	if err != nil {
		return err, 0
	}
	err, elapsed := UpdateTest("test_data", "5c3fab69-f2c1-4f8b-8cea-2cd5e74b195f")
	if err != nil {
		return err, 0
	}
	return err, elapsed
}

func DeleteDataStandartProcess(cntRecords int) (error, time.Duration) {
	var err error
	_, err = DB.Exec("DROP TABLE IF EXISTS test_data")
	if err != nil {
		return err, 0
	}
	_, err = DB.Exec("CREATE TABLE test_data(val_id int, field1 text, field2 text, field3 text, field4 text, field5 text)")
	for i := 0; i < cntRecords; i++ {
		_, err = DB.Exec("INSERT INTO test_data VALUES($1,$2,$3,$4,$5,$6)", i, "d1_"+strconv.Itoa(i), "d2_"+strconv.Itoa(i), "d3_"+strconv.Itoa(i), "d4_"+strconv.Itoa(i), "d5_"+strconv.Itoa(i))
		if err != nil {
			return err, 0
		}
	}
	revel.INFO.Println("[TEST] start query...")
	start := time.Now()
	_, err = DB.Query("DELETE FROM test_data WHERE val_id < $1", cntRecords/2)
	elapsed := time.Since(start)
	revel.INFO.Println("[TEST] elapsed time", elapsed.Seconds());
	return err, elapsed
}

func DeleteDataSecureProcess(cntRecords int) (error, time.Duration) {
	var err error
	_, err = DB.Exec("DROP TABLE IF EXISTS test_data")
	if err != nil {
		return err, 0
	}
	_, err = DB.Exec("CREATE TABLE test_data(val_id int, field1 text, field2 text, field3 text, field4 text, field5 text)")
	for i := 0; i < cntRecords; i++ {
		_, err = DB.Exec("INSERT INTO test_data VALUES($1,$2,$3,$4,$5,$6)", i, "d1_"+strconv.Itoa(i), "d2_"+strconv.Itoa(i), "d3_"+strconv.Itoa(i), "d4_"+strconv.Itoa(i), "d5_"+strconv.Itoa(i))
		if err != nil {
			return err, 0
		}
	}
	_, err = DB.Query("SELECT acs_vcs_table_add('test_data')")
	if err != nil {
		return err, 0
	}
	err, elapsed := DeleteTest("test_data", "5c3fab69-f2c1-4f8b-8cea-2cd5e74b195f")
	if err != nil {
		return err, 0
	}
	return err, elapsed
}

func (c TestCntl) SelectDataNormal() revel.Result {
	ret := make(map[string]string)
	var err error
	var resTest string
	revel.INFO.Println("[TEST] select unsecure")

	file, err := os.Create("gocode/src/info/tests/selectStandart.txt")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	defer file.Close()

	for i := 100; i <= 500; i += 100 {
		err, elapsed := SelectDataStandartProcess(i)
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
		resTest = strconv.Itoa(i) + " " + strconv.FormatFloat(elapsed.Seconds(), 'f', 10, 64)
		file.WriteString(resTest)
		file.WriteString("\n")
	}
	revel.INFO.Println("[TEST] test done.")

	return c.RenderJson(ret)
}

func (c TestCntl) SelectDataSecure() revel.Result {
	ret := make(map[string]string)
	var err error
	var resTest string
	revel.INFO.Println("[TEST] select secure")

	file, err := os.Create("gocode/src/info/tests/selectSecure.txt")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	defer file.Close()

	for i := 100; i <= 500; i += 100 {
		err, elapsed := SelectDataSecureProcess(i)
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
		resTest = strconv.Itoa(i) + " " + strconv.FormatFloat(elapsed.Seconds(), 'f', 10, 64)
		file.WriteString(resTest)
		file.WriteString("\n")
	}
	revel.INFO.Println("[TEST] test done.")

	return c.RenderJson(ret)
}

func (c TestCntl) UpdateDataNormal() revel.Result {
	ret := make(map[string]string)
	var err error
	var resTest string
	revel.INFO.Println("[TEST] update unsecure")

	file, err := os.Create("gocode/src/info/tests/updateStandart.txt")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	defer file.Close()

	for i := 100; i <= 500; i += 100 {
		err, elapsed := UpdateDataStandartProcess(i)
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
		resTest = strconv.Itoa(i) + " " + strconv.FormatFloat(elapsed.Seconds(), 'f', 10, 64)
		file.WriteString(resTest)
		file.WriteString("\n")
	}
	revel.INFO.Println("[TEST] test done.")

	return c.RenderJson(ret)
}

func (c TestCntl) UpdateDataSecure() revel.Result {
	ret := make(map[string]string)
	var err error
	var resTest string
	revel.INFO.Println("[TEST] update data secure")

	file, err := os.Create("gocode/src/info/tests/updateSecure.txt")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	defer file.Close()

	for i := 100; i <= 500; i += 100 {
		err, elapsed := UpdateDataSecureProcess(i)
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
		resTest = strconv.Itoa(i) + " " + strconv.FormatFloat(elapsed.Seconds(), 'f', 10, 64)
		file.WriteString(resTest)
		file.WriteString("\n")
	}
	revel.INFO.Println("[TEST] test done.")

	return c.RenderJson(ret)
}

func (c TestCntl) DeleteDataNormal() revel.Result {
	ret := make(map[string]string)
	var err error
	var resTest string
	revel.INFO.Println("[TEST] delete unsecure")

	file, err := os.Create("gocode/src/info/tests/deleteStandart.txt")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	defer file.Close()

	for i := 100; i <= 500; i += 100 {
		err, elapsed := DeleteDataStandartProcess(i)
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
		resTest = strconv.Itoa(i) + " " + strconv.FormatFloat(elapsed.Seconds(), 'f', 10, 64)
		file.WriteString(resTest)
		file.WriteString("\n")
	}
	revel.INFO.Println("[TEST] test done.")

	return c.RenderJson(ret)
}

func (c TestCntl) DeleteDataSecure() revel.Result {
	ret := make(map[string]string)
	var err error
	var resTest string
	revel.INFO.Println("[TEST] delete data secure")

	file, err := os.Create("gocode/src/info/tests/deleteSecure.txt")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	defer file.Close()

	for i := 100; i <= 500; i += 100 {
		err, elapsed := DeleteDataSecureProcess(i)
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
		resTest = strconv.Itoa(i) + " " + strconv.FormatFloat(elapsed.Seconds(), 'f', 10, 64)
		file.WriteString(resTest)
		file.WriteString("\n")
	}
	revel.INFO.Println("[TEST] test done.")

	return c.RenderJson(ret)
}