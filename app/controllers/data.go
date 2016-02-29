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
}

func (c DataCntl) Get(params string) revel.Result {
	var ret DataRecords
	var p DataParams

	err := json.Unmarshal([]byte(params), &p)
	if err != nil {
		ret.Error = "invalid params"
		return c.RenderJson(ret)
	}
	revel.INFO.Println("[get data] user", p.User)
	revel.INFO.Println("[get data] table", p.Table)

	//check user

	//get rights
	rows, err := DB.Query("SELECT DISTINCT rule FROM rules_p WHERE rule_role=$1", p.User)
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	var rules []string
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
	stmt, err := DB.Prepare("SELECT * FROM " + p.Table + " WHERE rule IN ($1, $2)")
	if err != nil {
		revel.ERROR.Println("[get data] stmt", err)
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	rows, err = stmt.Query(rules)
	if err != nil {
		revel.ERROR.Println("[get data] stmt rows", err)
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	revel.INFO.Println("[get data] stmt", stmt)

	/*query := fmt.Sprintf("SELECT * FROM t WHERE id Iargs := []int{1, 2, 3} qN (%s)", strings.Join(strings.Split(strings.Repeat("?", len(rules)), ""), ","))
	stmt, _ := DB.Prepare(query)
	revel.INFO.Println("[get data] stmt", stmt)
	rows, _ = stmt.Query(rules...)*/
	//ret.Id = p.Table
	//stmt, err := DB.Prepare("SELECT * FROM " + p.Table + "WHERE rule IN " + rulesList)
	//args := []int{1, 2, 3}
	//revel.INFO.Println("[get data] query", stmt)
	//rows, _ = stmt.Query([]string{"1", "2", "3"})
	rows, err = DB.Query("SELECT * FROM " + p.Table)

	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	columnNames, err := rows.Columns()
	ret.Columns = listColumns(p.Table)
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