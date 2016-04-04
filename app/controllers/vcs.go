package controllers

import "github.com/robfig/revel"

type VcsCntl struct {
	*revel.Controller
}

type VcsTables struct {
	TablesAll []string `json:"tablesAll"`
	TablesVcs []string `json:"tablesVcs"`
	Schema string `json:"schema"`
	Error  string `json:"error"`
}

func (c VcsCntl) Tables() revel.Result {
	var ret VcsTables

	rows, err := DB.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			ret.Error = err.Error()
			return c.RenderJson(ret)
		} else {
			ret.TablesAll = append(ret.TablesAll, tableName)
		}
	}

	rows, err = DB.Query("SELECT table_name FROM acs.tables WHERE schema_name = 'public'")
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			ret.Error = err.Error()
			return c.RenderJson(ret)
		} else {
			ret.TablesVcs = append(ret.TablesVcs, tableName)
		}
	}
	ret.Schema = "public"

	return c.RenderJson(ret)
}

func (c VcsCntl) Add(table string) revel.Result {
	ret := make(map[string]string)
	_, err := DB.Query("SELECT acs_vcs_table_add($1)", table)
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	return c.RenderJson(ret)
}

func (c VcsCntl) Delete(table string) revel.Result {
	ret := make(map[string]string)
	_, err := DB.Query("SELECT acs_vcs_table_rm($1)", table)
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	return c.RenderJson(ret)
}