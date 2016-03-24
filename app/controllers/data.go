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
	Token string `json:"token"`
	Table string `json:"table"`
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
	revel.INFO.Println("[get data] table", p.Table)

	//check user
	var user string
	err = DB.QueryRow("SELECT COALESCE(acs_get_user($1), '')", p.Token).Scan(&user)
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}

	if user == "" {
		ret.Error = "authorization error"
		return c.RenderJson(ret)
	}

	//get rights
	rows, err := DB.Query("SELECT DISTINCT security_rule FROM "+TABLE_RULES_DATA+" WHERE rule_user=$1 AND rule_action='select'", user)
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

	//sql query
	var stmt *sql.Stmt
	if (hasRights) {
		stmt, err = DB.Prepare("SELECT "+p.Table+".*  FROM "+p.Table+" LEFT OUTER JOIN acs.rule_record AS rules ON ("+p.Table+".uuid_record = rules.uuid_record) WHERE security_rule IS NULL OR security_rule IN " + rulesList)
		//stmt, err = DB.Prepare("SELECT * FROM " + p.Table + " WHERE rule IS NULL OR rule IN " + rulesList)
	} else {
		stmt, err = DB.Prepare("SELECT "+p.Table+".*  FROM "+p.Table+" LEFT OUTER JOIN acs.rule_record AS rules ON ("+p.Table+".uuid_record = rules.uuid_record) WHERE security_rule IS NULL")
		//stmt, err = DB.Prepare("SELECT * FROM " + p.Table + " WHERE rule IS NULL")
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