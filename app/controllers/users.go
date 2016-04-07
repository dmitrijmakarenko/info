package controllers

import (
	"github.com/robfig/revel"
	"encoding/json"
)

type UsersCntl struct {
	*revel.Controller
}

type UserItem struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type UsersList struct {
	Error string `json:"error"`
	Users []UserItem `json:"accounts"`
}

type UsersRuleItem struct {
	Table string `json:"table"`
	Rule string `json:"rule"`
}

type UsersTempItem struct {
	Table string `json:"table"`
	Time int `json:"time"`
}

type UserData struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Password string `json:"password"`
	Position string `json:"position"`
	TableRule string `json:"tableRule"`
	TableRules []UsersRuleItem `json:"tableRules"`
	TempRule int `json:"tempRule"`
	TempRules []UsersTempItem `json:"tempRules"`
	Error string `json:"error"`
}

type UserSettings struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Password string `json:"password"`
	Position string `json:"position"`
	TableRule string `json:"tableRule"`
	TableRules []UsersRuleItem `json:"tableRules"`
	TempRule int `json:"tempRule"`
	TempRules []UsersTempItem `json:"tempRules"`
}

func (c UsersCntl) List() revel.Result {
	var ret UsersList
	rows, err := DB.Query("SELECT id, COALESCE(realname, '') as name FROM "+TABLE_USERS+" ORDER BY id")
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	} else {
		for rows.Next() {
			var id string
			var name string
			err := rows.Scan(&id, &name)
			if err != nil {
				revel.ERROR.Println(err)
			} else {
				user := UserItem{}
				user.Id = id
				user.Name = name
				ret.Users = append(ret.Users, user)
			}
		}
	}
	return c.RenderJson(ret)
}

func (c UsersCntl) Get(id string) revel.Result {
	var ret UserData
	var ruleTable string
	var tempTime int

	rows, err := DB.Query("SELECT realname, position_user FROM "+TABLE_USERS+" WHERE id=$1", id)
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	for rows.Next() {
		var name string
		var position string
		err := rows.Scan(&name, &position)
		if err != nil {
			ret.Error = err.Error()
			return c.RenderJson(ret)
		} else {
			ret.Id = id
			ret.Name = name
			ret.Position = position
		}
	}
	//get rules settings
	_ = DB.QueryRow("SELECT security_rule FROM acs.user_rules WHERE user_id=$1 AND temp_use=false AND table_all=true", id).Scan(&ruleTable)
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	if ruleTable != "" {
		ret.TableRule = ruleTable
	} else {
		rows, err := DB.Query("SELECT table_name, security_rule FROM acs.user_rules WHERE user_id=$1 AND temp_use=false AND table_all=false", id)
		if err != nil {
			ret.Error = err.Error()
			return c.RenderJson(ret)
		}
		for rows.Next() {
			var table string
			var rule string
			err := rows.Scan(&table, &rule)
			if err != nil {
				ret.Error = err.Error()
				return c.RenderJson(ret)
			} else {
				s := UsersRuleItem{}
				s.Rule = rule
				s.Table = table
				ret.TableRules = append(ret.TableRules, s)
			}
		}
	}
	//get temp rules
	_ = DB.QueryRow("SELECT EXTRACT(epoch FROM temp_time)/60 AS time FROM acs.user_rules WHERE user_id=$1 AND temp_use=true AND table_all=true", id).Scan(&tempTime)
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	if tempTime != 0 {
		ret.TempRule = tempTime
	} else {
		rows, err := DB.Query("SELECT table_name, EXTRACT(epoch FROM temp_time)/60 AS time FROM acs.user_rules WHERE user_id=$1 AND temp_use=true AND table_all=false", id)
		if err != nil {
			ret.Error = err.Error()
			return c.RenderJson(ret)
		}
		for rows.Next() {
			var table string
			var time int
			err := rows.Scan(&table, &time)
			if err != nil {
				ret.Error = err.Error()
				return c.RenderJson(ret)
			} else {
				s := UsersTempItem{}
				s.Time = time
				s.Table = table
				ret.TempRules = append(ret.TempRules, s)
			}
		}
	}

	return c.RenderJson(ret)
}

func (c UsersCntl) Update(id string, data string) revel.Result {
	ret := make(map[string]string)
	var settings UserSettings
	changePass := false

	err := json.Unmarshal([]byte(data), &settings)
	if err != nil {
		ret["error"] = "settings error format"
		return c.RenderJson(ret)
	}
	//update user settings
	if id == "!new" {
		_, err = DB.Exec("INSERT INTO "+TABLE_USERS+"(id, pass, realname, position_user) VALUES ($1, crypt($2, gen_salt('bf')), $3, $4)", settings.Id, settings.Password, settings.Name, settings.Position)
	} else {
		if settings.Password != "" {
			changePass = true
		}
		if changePass {
			_, err = DB.Exec("UPDATE "+TABLE_USERS+" SET id=$2, pass=crypt($3, gen_salt('bf')), realname=$4, position_user=$5 WHERE id=$1", id, settings.Id, settings.Password, settings.Name, settings.Position)
		} else {
			_, err = DB.Exec("UPDATE "+TABLE_USERS+" SET id=$2, realname=$3, position_user=$4 WHERE id=$1", id, settings.Id, settings.Name, settings.Position)
		}
	}
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	//table rules
	_, err = DB.Exec("DELETE FROM acs.user_rules WHERE user_id=$1", settings.Id)
	if settings.TableRule != "" {
		_, err = DB.Exec("INSERT INTO acs.user_rules(user_id, table_all, security_rule) VALUES($1, true, $2)", settings.Id, settings.TableRule)
	}
	if settings.TableRules != nil {
		for _, s := range settings.TableRules {
			_, err = DB.Exec("INSERT INTO acs.user_rules(user_id, table_name, security_rule) VALUES($1, $2, $3)", settings.Id, s.Table, s.Rule)
		}
	}
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	//temp rules
	if settings.TempRule != 0 {
		_, err = DB.Exec("INSERT INTO acs.user_rules(user_id, table_all, temp_use, temp_time) VALUES($1, true, true, $2)", settings.Id, settings.TempRule * 60)
	}
	if settings.TempRules != nil {
		for _, s := range settings.TempRules {
			_, err = DB.Exec("INSERT INTO acs.user_rules(user_id, table_name, temp_use, temp_time) VALUES($1, $2, true, $3)", settings.Id, s.Table, s.Time * 60)
		}
	}
	return c.RenderJson(ret)
}

func (c UsersCntl) Delete(id string) revel.Result {
	ret := make(map[string]string)
	_, err := DB.Exec("DELETE FROM "+TABLE_USERS+" WHERE id=$1", id)
	if err != nil {
		ret["error"] = err.Error()
	}
	return c.RenderJson(ret)
}