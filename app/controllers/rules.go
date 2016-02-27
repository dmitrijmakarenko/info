package controllers

import (
	"github.com/robfig/revel"
	"encoding/json"
	"github.com/nu7hatch/gouuid"
)

type RuleCntl struct {
	*revel.Controller
}

type RuleItem struct {
	Id string `json:"id"`
	Desc string `json:"desc"`
}

type RulesList struct {
	Rules []RuleItem `json:"rules"`
	Error string `json:"error"`
}

type RuleData struct {
	Users []UserItem `json:"users"`
	Groups []GroupItem `json:"groups"`
	Error string `json:"error"`
}

type RuleSettings struct {
	Error  string `json:"error"`
	Id string `json:"id"`
	Desc string `json:"desc"`
}

func (c RuleCntl) List() revel.Result {
	var ret RulesList
	rows, err := DB.Query("SELECT rule, rule_desc FROM rules")
	if err != nil {
		revel.ERROR.Println("[get rules]", err)
		ret.Error = err.Error()
	} else {
		for rows.Next() {
			var id string
			var desc string
			err := rows.Scan(&id, &desc)
			if err != nil {
				revel.ERROR.Println(err)
			} else {
				rule := RuleItem{}
				rule.Id = id
				rule.Desc = desc
				ret.Rules = append(ret.Rules, rule)
			}
		}
	}
	return c.RenderJson(ret)
}

func (c RuleCntl) Get(id string) revel.Result {
	var ret RuleSettings
	rows, err := DB.Query("SELECT rule_desc FROM rules WHERE rule=$1", id)
	if err != nil {
		revel.ERROR.Println("[get rule]", err)
		ret.Error = err.Error()
	} else {
		for rows.Next() {
			var desc string
			err := rows.Scan(&desc)
			if err != nil {
				revel.ERROR.Println(err)
			} else {
				ret.Desc = desc
			}
		}
	}
	return c.RenderJson(ret)
}

func (c RuleCntl) Data() revel.Result {
	var ret RuleData
	users, err := usersList()
	if err != nil {
		ret.Error = err.Error()
	} else {
		ret.Users = users
	}
	groups, err := groupsList()
	if err != nil {
		ret.Error = err.Error()
	} else {
		ret.Groups = groups
	}
	return c.RenderJson(ret)
}

func (c RuleCntl) Update(id string, data string) revel.Result {
	ret := make(map[string]string)
	var settings RuleSettings
	err := json.Unmarshal([]byte(data), &settings)
	if id == "!new" {
		uuid, _ := uuid.NewV4()
		settings.Id = uuid.String()
	} else {
		settings.Id = id
	}
	if err != nil {
		ret["error"] = "settings error format";
	} else {

		if id == "!new" {
			_, err = DB.Exec("INSERT INTO rules(rule, rule_desc) VALUES ($1, $2)", settings.Id, settings.Desc)
		} else {
			_, err = DB.Exec("UPDATE rules SET rule_desc=$2 WHERE rule=$1", settings.Id, settings.Desc)
		}
		if err != nil {
			ret["error"] = err.Error();
		}
	}
	return c.RenderJson(ret)
}

func (c RuleCntl) Delete(id string) revel.Result {
	ret := make(map[string]string)
	_, err := DB.Exec("DELETE FROM rules WHERE rule=$1", id)
	if err != nil {
		revel.ERROR.Println("[delete rule]", err)
		ret["error"] = err.Error()
	}
	return c.RenderJson(ret)
}