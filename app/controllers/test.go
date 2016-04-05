package controllers

import (
	"github.com/robfig/revel"
	"database/sql"
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