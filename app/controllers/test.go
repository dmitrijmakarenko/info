package controllers

import (
	"github.com/robfig/revel"
	"database/sql"
	"github.com/gchaincl/dotsql"
	"os"
	"io/ioutil"
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
	_, err := DB.Query("SELECT acs_install()")
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

	//CompileFunc()

	dot, err := dotsql.LoadFromFile("gocode/src/info/sql/func_acs_f.sql")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	_, err = dot.Exec(dbGeneral, "install-functions")
	_, err = dot.Exec(dbStat1, "install-functions")
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

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

func CompileFunc() error {
	f, err := os.Create("gocode/src/info/sql/func_acs.sql")
	if err != nil {
		return err
	}
	_, err = f.WriteString("-- name: install-functions\n")

	files, _ := ioutil.ReadDir("gocode/src/info/sql/functions/")
	for _, file := range files {
		b, err := ioutil.ReadFile("gocode/src/info/sql/functions/" + file.Name())
		_, err = f.WriteString(string(b))
		_, err = f.WriteString("\n\n")
		if err != nil {
			return err
		}
	}
	return nil
}