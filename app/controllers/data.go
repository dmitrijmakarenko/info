package controllers

import (
	"github.com/robfig/revel"
	"database/sql"
	"encoding/json"
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

	//sql query
	ret.Id = p.Table
	rows, err := DB.Query("SELECT * FROM " + p.Table)
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