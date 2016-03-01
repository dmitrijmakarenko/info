package controllers

import (
	"github.com/robfig/revel"
	"strings"
	"encoding/json"
	"github.com/nu7hatch/gouuid"
)

type GroupCntl struct {
	*revel.Controller
}

type GroupItem struct {
	Id string `json:"id"`
	Name string  `json:"name"`
}

type GroupData struct {
	Users []UserItem `json:"users"`
}

type GroupsList struct {
	Error string `json:"error"`
	Groups []GroupItem `json:"groups"`
}

type GroupSettings struct {
	Id string `json:"id"`
	Name string  `json:"name"`
	Members []UserItem `json:"members"`
	Users []UserItem `json:"users"`
}

func CreateGroupTable() (err error) {
	_, err = DB.Exec("CREATE TABLE sys_groups ( id uuid NOT NULL, name name ) WITH (OIDS=FALSE); ALTER TABLE sys_groups OWNER TO postgres;")
	return err
}

func (c GroupCntl) List() revel.Result {
	var ret GroupsList
	rows, err := DB.Query("SELECT id, name FROM sys_groups")
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			err = CreateGroupTable();
			if err != nil {
				ret.Error = err.Error()
			}
		} else {
			ret.Error = err.Error()
		}
	} else {
		for rows.Next() {
			var id string
			var name string
			err := rows.Scan(&id, &name)
			if err != nil {
				revel.ERROR.Println(err)
			} else {
				group := GroupItem{}
				group.Id = id
				group.Name = name
				ret.Groups = append(ret.Groups, group)
			}
		}
	}
	return c.RenderJson(ret)
}

func (c GroupCntl) Update(id string, data string) revel.Result {
	ret := make(map[string]string)
	var err error
	var settings GroupSettings
	err = json.Unmarshal([]byte(data), &settings)
	if id == "!new" {
		uuid, _ := uuid.NewV4()
		settings.Id = uuid.String()
	} else {
		settings.Id = id
	}
	if err != nil {
		ret["error"] = "settings error format";
		return c.RenderJson(ret)
	}

	if id == "!new" {
		_, err = DB.Exec("INSERT INTO sys_groups(id, name) VALUES ($1, $2)", settings.Id, settings.Name)
	} else {
		_, err = DB.Exec("UPDATE sys_groups SET name=$2 WHERE id=$1", settings.Id, settings.Name)
	}
	_, err = DB.Exec("DELETE FROM sys_group_user WHERE group_id=$1", settings.Id)
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	var groupRules []RuleGroupItem
	rows, err := DB.Query("SELECT DISTINCT rule, action FROM rules_p WHERE rule_group=$1", settings.Id)
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	for rows.Next() {
		var ruleId string
		var action string
		err := rows.Scan(&ruleId, &action)
		if err != nil {
			ret["error"] = err.Error()
		} else {
			rule := RuleGroupItem{}
			rule.Id = ruleId
			rule.Operation = action
			groupRules = append(groupRules, rule)
		}
	}

	_, err = DB.Exec("DELETE FROM rules_p WHERE rule_group=$1", settings.Id)
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	for _, member := range settings.Members {
		_, err = DB.Exec("INSERT INTO sys_group_user(group_id, user_id) VALUES ($1, $2)", settings.Id, member.Id)
		for _, rule := range groupRules {
			revel.INFO.Println("[add member]", member.Id)
			_, err = DB.Exec("INSERT INTO rules_p(rule, rule_role, action, rule_group) VALUES ($1, $2, $3, $4)", rule.Id, member.Id, rule.Operation, settings.Id)
		}
	}
	if err != nil {
		ret["error"] = err.Error()
	}

	return c.RenderJson(ret)
}

func (c GroupCntl) Get(id string) revel.Result {
	var ret GroupSettings
	var usersGroup []string
	rows, err := DB.Query("SELECT name FROM sys_groups WHERE id=$1", id)
	rowsUsers, err := DB.Query("SELECT user_id FROM sys_group_user WHERE group_id=$1", id)
	allUsers, err := usersList()
	if err != nil {
		revel.ERROR.Println(err)
	} else {
		for rowsUsers.Next() {
			var id string
			err := rowsUsers.Scan(&id)
			if err != nil {
				revel.ERROR.Println(err)
			} else {
				usersGroup = append(usersGroup, id)
			}
		}
	}
	var isMember bool
	for _, user := range allUsers {
		isMember = false
		for _, userGroup := range usersGroup {
			if userGroup == user.Id {
				isMember = true
			}
		}
		if isMember {
			ret.Members = append(ret.Members, user)
		} else {
			ret.Users = append(ret.Users, user)
		}
	}
	if err != nil {
		revel.ERROR.Println("[get group]", err)
		ret := make(map[string]string)
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	} else {
		for rows.Next() {
			var name string
			err := rows.Scan(&name)
			if err != nil {
				revel.ERROR.Println(err)
			} else {
				ret.Id = id
				ret.Name = name
			}
		}
	}
	return c.RenderJson(ret)
}

func (c GroupCntl) Data() revel.Result {
	var ret GroupData
	allUsers, err := usersList()
	if err != nil {
		ret := make(map[string]string)
		ret["error"] = err.Error()
	} else {
		ret.Users = allUsers
	}
	return c.RenderJson(ret)
}

func (c GroupCntl) Delete(id string) revel.Result {
	ret := make(map[string]string)
	_, err := DB.Exec("DELETE FROM sys_groups WHERE id=$1", id)
	if err != nil {
		revel.ERROR.Println("[delete group]", err)
		ret["error"] = err.Error()
	}
	return c.RenderJson(ret)
}