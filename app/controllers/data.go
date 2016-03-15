package controllers

import (
	"github.com/robfig/revel"
	"database/sql"
	"encoding/json"
	"strconv"
)

type DataCntl struct {
	*revel.Controller
}

type DataRecords struct {
	Id string `json:"id"`
	Columns []string `json:"columns"`
	Rows [][]string `json:"rows"`
	Error string `json:"error"`
}

type DataParams struct {
	Table string `json:"table"`
	User string `json:"user"`
	Token string `json:"token"`
}

func (c DataCntl) Get(params string) revel.Result {
	var ret DataRecords
	var p DataParams
	var hasRights bool

	err := json.Unmarshal([]byte(params), &p)
	if err != nil {
		ret.Error = "invalid params"
		return c.RenderJson(ret)
	}
	revel.INFO.Println("[get data] user", p.User)
	revel.INFO.Println("[get data] table", p.Table)

	//check user
	var auth bool
	err = DB.QueryRow("SELECT check_user($1,$2) AS auth", p.User, p.Token).Scan(&auth)
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}

	if !auth {
		ret.Error = "authorization error"
		return c.RenderJson(ret)
	}

	//get rights
	rows, err := DB.Query("SELECT DISTINCT rule_id FROM "+TABLE_RULES_DATA+" WHERE rule_user=$1 AND rule_action='select'", p.User)
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	var rules []interface{}
	for rows.Next() {
		var rule string
		err := rows.Scan(&rule)
		if err != nil {
			ret.Error = err.Error()
			return c.RenderJson(ret)
		} else {
			rules = append(rules, rule)
		}
	}
	if len(rules) > 0 {
		hasRights = true
	} else {
		hasRights = false
	}
	revel.INFO.Println("[get data] rules", rules)
	var rulesList string
	rulesList += "("
	for i, _ := range rules {
		rulesList += "$" + strconv.Itoa(i+1)
		if (i != len(rules) - 1) {
			rulesList += ", "
		}
	}
	rulesList += ")"
	revel.INFO.Println("[get data] rulesList", rulesList)

	connstring := "host=" + DB_HOST + " port=" + DB_PORT + " user=" + p.User + " dbname=" + DB_NAME + " password=" + DB_PASSWORD + " sslmode=disable"
	dbUser, err := sql.Open(DB_DRIVER, connstring)

	//sql query
	var stmt *sql.Stmt
	if (hasRights) {
		stmt, err = dbUser.Prepare("SELECT * FROM " + p.Table + " WHERE rule IS NULL OR rule IN " + rulesList)
	} else {
		stmt, err = dbUser.Prepare("SELECT * FROM " + p.Table + " WHERE rule IS NULL")
	}
	if err != nil {
		revel.ERROR.Println("[get data] stmt", err)
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	if (hasRights) {
		rows, err = stmt.Query(rules...)
	} else {
		rows, err = stmt.Query()
	}
	if err != nil {
		revel.ERROR.Println("[get data] stmt rows", err)
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	revel.INFO.Println("[get data] stmt", stmt)

	columnNames, err := rows.Columns()
	ret.Columns = columnNames
	//ret.Columns = listColumns(p.Table)
	var retRows [][]string
	for rows.Next() {
		var retRow []string
		columnPointers := make([]interface{}, len(columnNames))
		for i := 0; i < len(columnNames); i++ {
			columnPointers[i] = new(sql.RawBytes)
		}
		if err := rows.Scan(columnPointers...); err != nil {
			ret.Error = err.Error()
			return c.RenderJson(ret)
		}
		for i := 0; i < len(columnNames); i++ {
			if rb, ok := columnPointers[i].(*sql.RawBytes); ok {
				retRow = append(retRow, string(*rb))
			} else {
				ret.Error = "erorr get row"
				return c.RenderJson(ret)
			}
		}
		retRows = append(retRows, retRow)
	}
	ret.Rows = retRows

	return c.RenderJson(ret)
}