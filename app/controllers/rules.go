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

type RuleAction struct {
	Object string `json:"object"`
	Operation string `json:"operation"`
	IsUser bool `json:"isUser"`
}

type RuleData struct {
	Users []UserItem `json:"users"`
	Groups []GroupItem `json:"groups"`
	Error string `json:"error"`
}

type RuleSettings struct {
	Id string `json:"id"`
	Desc string `json:"desc"`
	Actions []RuleAction `json:"actions"`
	Error  string `json:"error"`
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
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	for rows.Next() {
		var desc string
		err := rows.Scan(&desc)
		if err != nil {
			ret.Error = err.Error()
		} else {
			ret.Desc = desc
		}
	}

	//get groups
	rows, err = DB.Query("SELECT DISTINCT COALESCE(rule_group, '') as rule_group,action FROM rules_p WHERE rule=$1 AND rule_group IS NOT NULL", id)
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	for rows.Next() {
		var group string
		var operation string
		err := rows.Scan(&group, &operation)
		if err != nil {
			ret.Error = err.Error()
		} else {
			revel.INFO.Println("group", group);
			act := RuleAction{}
			act.Object = group
			act.Operation = operation
			act.IsUser = false
			ret.Actions = append(ret.Actions, act)
		}
	}

	//get users
	rows, err = DB.Query("SELECT rule_role,action FROM rules_p WHERE rule=$1 AND rule_group IS NULL", id)
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	for rows.Next() {
		var user string
		var operation string
		err := rows.Scan(&user, &operation)
		if err != nil {
			ret.Error = err.Error()
		} else {
			revel.INFO.Println("user", user);
			act := RuleAction{}
			act.Object = user
			act.Operation = operation
			act.IsUser = true
			ret.Actions = append(ret.Actions, act)
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
	if err != nil {
		ret["error"] = "settings error format";
		return c.RenderJson(ret)
	}

	if id == "!new" {
		uuid, _ := uuid.NewV4()
		settings.Id = uuid.String()
		_, err = DB.Exec("INSERT INTO rules(rule, rule_desc) VALUES ($1, $2)", settings.Id, settings.Desc)
	} else {
		settings.Id = id
		_, err = DB.Exec("UPDATE rules SET rule_desc=$2 WHERE rule=$1", settings.Id, settings.Desc)
		_, err = DB.Exec("DELETE FROM rules_p WHERE rule=$1", settings.Id)
	}
	if err != nil {
		ret["error"] = err.Error();
		return c.RenderJson(ret)
	}
	if err != nil {
		ret["error"] = err.Error();
	}

	for _, action := range settings.Actions {
		if action.IsUser {
			_, err = DB.Exec("INSERT INTO rules_p(rule, rule_role, action) VALUES ($1, $2, $3)", settings.Id, action.Object, action.Operation)
		} else {
			users, err := getUsersByGroup(action.Object)
			if err != nil {
				ret["error"] = err.Error();
				return c.RenderJson(ret)
			}
			for _, user := range users {
				_, err = DB.Exec("INSERT INTO rules_p(rule, rule_role, action, rule_group) VALUES ($1, $2, $3, $4)", settings.Id, user, action.Operation, action.Object)
			}
		}
	}

	if err != nil {
		ret["error"] = err.Error();
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