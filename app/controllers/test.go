package controllers

import (
	"github.com/robfig/revel"
	"encoding/hex"
	"database/sql"
	"os"
	"strings"
)

type TestCntl struct {
	*revel.Controller
}

func (c TestCntl) Reset() revel.Result {
	var err error
	ret := make(map[string]string)

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

	return c.RenderJson(ret)
}

func (c TestCntl) Compile() revel.Result {
	ret := make(map[string]string)
	revel.INFO.Println("[compile] start")

	filepath := "gocode/src/info/dump/copy"

	f, err := os.Create(filepath)
	defer f.Close()
	revel.INFO.Println("[compile] file created")

	tables := []string{"ttt", "fruits"}
	for _, table := range tables {
		_, err = f.WriteString(table + "\n")
		rows, err := DB.Query("SELECT * FROM " + table)
		if err != nil {
			ret["error"] = err.Error()
			_ = os.Remove(filepath)
			return c.RenderJson(ret)
		}
		columnNames, _ := rows.Columns()
		for rows.Next() {
			var retRow []string
			columnPointers := make([]interface{}, len(columnNames))
			for i := 0; i < len(columnNames); i++ {
				columnPointers[i] = new(sql.RawBytes)
			}
			if err := rows.Scan(columnPointers...); err != nil {
				ret["error"] = err.Error()
				_ = os.Remove(filepath)
				return c.RenderJson(ret)
			}
			for i := 0; i < len(columnNames); i++ {
				if rb, ok := columnPointers[i].(*sql.RawBytes); ok {
					retRow = append(retRow, string(*rb))
				} else {
					ret["error"] = "error row"
					_ = os.Remove(filepath)
					return c.RenderJson(ret)
				}
			}
			s := strings.Join(retRow, " ")
			_, err = f.WriteString(s + "\n")
		}
	}

	//compute checksum
	b, err := ComputeMd5(filepath)
	if err != nil {
		_ = os.Remove(filepath)
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	s := hex.EncodeToString(b)

	//add to history
	_, err = DB.Exec("INSERT INTO acs.changes_history(change_uuid, change_date, change_type, change_db, hash) VALUES (uuid_generate_v4(), now(), 'compile', current_database(), $1)", s)
	if err != nil {
		_ = os.Remove(filepath)
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	return c.RenderJson(ret)
}