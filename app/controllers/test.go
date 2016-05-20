package controllers

import (
	"github.com/robfig/revel"
	"database/sql"
	"strconv"
	"time"
)

type TestCntl struct {
	*revel.Controller
}

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

	//insert data
	/*_, err = dbGeneral.Exec("INSERT INTO persons VALUES(1, 'Иванов', 'Иван', 'ул.Парфенова', 'г.Тула')")
	_, err = dbGeneral.Exec("INSERT INTO persons VALUES(2, 'Петров', 'Петр', 'ул.Ленина', 'г.Томск')")
	_, err = dbGeneral.Exec("INSERT INTO clients VALUES(1, 'Автозапчасти', 'г.Москва', '8 (906) 531-85-86')")
	_, err = dbGeneral.Exec("INSERT INTO clients VALUES(2, 'Продукты', 'г.Тула', '8 (906) 531-85-86')")*/

	/*_, err = dbStat1.Exec("INSERT INTO persons VALUES(1, 'Иванов', 'Иван', 'ул.Парфенова', 'г.Тула')")
	_, err = dbStat1.Exec("INSERT INTO persons VALUES(2, 'Петров', 'Петр', 'ул.Ленина', 'г.Томск')")
	_, err = dbStat1.Exec("INSERT INTO clients VALUES(1, 'Автозапчасти', 'г.Москва', '8 (906) 531-85-86')")
	_, err = dbStat1.Exec("INSERT INTO clients VALUES(2, 'Продукты', 'г.Тула', '8 (906) 531-85-86')")*/

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

func GetTest(table string, token string) DataRecords {
	var ret DataRecords
	var p DataGetParams
	var hasRights bool
	var err error

	p.Table = table;
	p.Token = token;

	//check user
	var user string
	revel.INFO.Println("time started")
	start := time.Now()
	err = DB.QueryRow("SELECT COALESCE(acs_get_user($1), '')", p.Token).Scan(&user)
	if err != nil {
		revel.ERROR.Println(err.Error())
		return ret
	}

	if user == "" {
		revel.ERROR.Println("authorization error")
		return ret
	}

	//get rights
	rows, err := DB.Query("SELECT DISTINCT security_rule FROM "+TABLE_RULES_DATA+" WHERE rule_user=$1 AND rule_action='r'", user)
	if err != nil {
		revel.ERROR.Println(err.Error())
		return ret
	}
	var rules []interface{}
	for rows.Next() {
		var rule string
		err := rows.Scan(&rule)
		if err != nil {
			revel.ERROR.Println(err.Error())
			return ret
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
	var stmt *sql.Stmt
	if (hasRights) {
		stmt, err = DB.Prepare("SELECT "+p.Table+".*  FROM "+p.Table+" LEFT OUTER JOIN acs.rule_record AS rules ON ("+p.Table+".uuid_record = rules.uuid_record) WHERE security_rule IS NULL OR security_rule IN " + rulesList)
	} else {
		stmt, err = DB.Prepare("SELECT "+p.Table+".*  FROM "+p.Table+" LEFT OUTER JOIN acs.rule_record AS rules ON ("+p.Table+".uuid_record = rules.uuid_record) WHERE security_rule IS NULL")
	}
	if err != nil {
		revel.ERROR.Println(err.Error())
		return ret
	}
	if (hasRights) {
		rows, err = stmt.Query(rules...)
	} else {
		rows, err = stmt.Query()
	}
	if err != nil {
		revel.ERROR.Println(err.Error())
		return ret
	}
	elapsed := time.Since(start)
	revel.INFO.Println("elapsed time", elapsed.Seconds());

	columnNames, err := rows.Columns()
	ret.Columns = columnNames
	var retRows [][]string
	for rows.Next() {
		var retRow []string
		columnPointers := make([]interface{}, len(columnNames))
		for i := 0; i < len(columnNames); i++ {
			columnPointers[i] = new(sql.RawBytes)
		}
		if err := rows.Scan(columnPointers...); err != nil {
			revel.ERROR.Println(err.Error())
			return ret
		}
		for i := 0; i < len(columnNames); i++ {
			if rb, ok := columnPointers[i].(*sql.RawBytes); ok {
				retRow = append(retRow, string(*rb))
			} else {
				revel.ERROR.Println("erorr get row")
				return ret
			}
		}
		retRows = append(retRows, retRow)
	}
	ret.Rows = retRows

	return ret
}

func (c TestCntl) SelectDataNormal() revel.Result {
	ret := make(map[string]string)
	var err error
	revel.INFO.Println("Start test select normal")

	_, err = DB.Exec("DROP TABLE IF EXISTS test_data")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	_, err = DB.Exec("CREATE TABLE test_data(val_id int, field1 text, field2 text, field3 text, field4 text, field5 text)")
	for i := 0; i < 500; i++ {
		_, err = DB.Exec("INSERT INTO test_data VALUES($1,$2,$3,$4,$5,$6)", i, "d1_"+strconv.Itoa(i), "d2_"+strconv.Itoa(i), "d3_"+strconv.Itoa(i), "d4_"+strconv.Itoa(i), "d5_"+strconv.Itoa(i))
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
	}
	revel.INFO.Println("time started")
	start := time.Now()
	_, err = DB.Query("SELECT * FROM test_data")
	elapsed := time.Since(start)
	revel.INFO.Println("elapsed time", elapsed.Seconds());


	return c.RenderJson(ret)
}

func (c TestCntl) SelectDataSecure() revel.Result {
	ret := make(map[string]string)
	var err error
	revel.INFO.Println("Start test select secure");

	_, err = DB.Exec("DROP TABLE IF EXISTS test_data")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	_, err = DB.Exec("CREATE TABLE test_data(val_id int, field1 text, field2 text, field3 text, field4 text, field5 text)")
	for i := 0; i < 500; i++ {
		_, err = DB.Exec("INSERT INTO test_data VALUES($1,$2,$3,$4,$5,$6)", i, "d1_"+strconv.Itoa(i), "d2_"+strconv.Itoa(i), "d3_"+strconv.Itoa(i), "d4_"+strconv.Itoa(i), "d5_"+strconv.Itoa(i))
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
	}
	_, err = DB.Query("SELECT acs_vcs_table_add('test_data')")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	_ = GetTest("test_data", "214364c4-7eba-41db-8257-9c75fcbe243d")


	return c.RenderJson(ret)
}