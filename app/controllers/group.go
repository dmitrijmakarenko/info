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

type GroupParentItem struct {
	Id string `json:"id"`
	Level int  `json:"level"`
}

type GroupData struct {
	Users []UserItem `json:"users"`
	Groups []GroupItem `json:"groups"`
	Error string `json:"error"`
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
	Parents []GroupParentItem `json:"parents"`
	Error string `json:"error"`
}

type GroupSettingsUpdate struct {
	Id string `json:"id"`
	Name string  `json:"name"`
	Members []UserItem `json:"members"`
	Users []UserItem `json:"users"`
	ParentsAdd []GroupParentItem `json:"parentsAdd"`
	ParentsRemove []GroupParentItem `json:"parentsRemove"`
}

func CreateGroupTable() (err error) {
	_, err = DB.Exec("CREATE TABLE "+TABLE_GROUPS+" ( id uuid NOT NULL, name name ) WITH (OIDS=FALSE); ALTER TABLE sys_groups OWNER TO postgres;")
	return err
}

func (c GroupCntl) List() revel.Result {
	var ret GroupsList
	rows, err := DB.Query("SELECT id, name FROM "+TABLE_GROUPS+" ORDER BY name")
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
	var settings GroupSettingsUpdate
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
		_, err = DB.Exec("INSERT INTO "+TABLE_GROUPS+"(id, name) VALUES ($1, $2)", settings.Id, settings.Name)
	} else {
		_, err = DB.Exec("UPDATE "+TABLE_GROUPS+" SET name=$2 WHERE id=$1", settings.Id, settings.Name)
	}
	_, err = DB.Exec("DELETE FROM "+TABLE_GROUP_USER+" WHERE group_id=$1", settings.Id)
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	var groupRules []RuleGroupItem
	rows, err := DB.Query("SELECT DISTINCT rule, action FROM "+TABLE_RULES_P+" WHERE rule_group=$1", settings.Id)
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

	_, err = DB.Exec("DELETE FROM "+TABLE_RULES_P+" WHERE rule_group=$1", settings.Id)
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	for _, member := range settings.Members {
		_, err = DB.Exec("INSERT INTO "+TABLE_GROUP_USER+"(group_id, user_id) VALUES ($1, $2)", settings.Id, member.Id)
		for _, rule := range groupRules {
			_, err = DB.Exec("INSERT INTO "+TABLE_RULES_P+"(rule, rule_role, action, rule_group) VALUES ($1, $2, $3, $4)", rule.Id, member.Id, rule.Operation, settings.Id)
		}
	}
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}

	for _, parent := range settings.ParentsRemove {
		_, err = DB.Exec("DELETE FROM "+TABLE_GROUPS_STRUCT+" WHERE group_id=$1 AND parent_id=$2", settings.Id, parent.Id)
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
	}
	for _, parent := range settings.ParentsAdd {
		_, err = DB.Exec("INSERT INTO "+TABLE_GROUPS_STRUCT+"(group_id, parent_id, level) VALUES ($1, $2, $3)", settings.Id, parent.Id, 1)
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
		//rules
		rows, err = DB.Query("SELECT DISTINCT rule FROM "+TABLE_RULES_P+" WHERE rule_group=$1", settings.Id)
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
		for rows.Next() {
			var rule string
			err := rows.Scan(&rule)
			if err != nil {
				ret["error"] = err.Error()
				return c.RenderJson(ret)
			} else {
				revel.INFO.Println("[update group]", rule)
				//_, err = DB.Exec("INSERT INTO "+TABLE_GROUPS_STRUCT+"(group_id, parent_id, level) VALUES ($1, $2, $3)", groupId, parent.Id, level+1)
			}
		}

		rows, err = DB.Query("SELECT group_id,level FROM "+TABLE_GROUPS_STRUCT+" WHERE parent_id=$1", settings.Id)
		if err != nil {
			ret["error"] = err.Error()
			return c.RenderJson(ret)
		}
		for rows.Next() {
			var groupId string
			var level int
			err := rows.Scan(&groupId, &level)
			if err != nil {
				ret["error"] = err.Error()
				return c.RenderJson(ret)
			} else {
				_, err = DB.Exec("INSERT INTO "+TABLE_GROUPS_STRUCT+"(group_id, parent_id, level) VALUES ($1, $2, $3)", groupId, parent.Id, level+1)
			}
		}
		if err != nil {
			ret["error"] = err.Error()
		}
	}

	return c.RenderJson(ret)
}

func (c GroupCntl) Get(id string) revel.Result {
	var ret GroupSettings
	var usersGroup []string
	rows, err := DB.Query("SELECT name FROM "+TABLE_GROUPS+" WHERE id=$1", id)
	rowsUsers, err := DB.Query("SELECT user_id FROM "+TABLE_GROUP_USER+" WHERE group_id=$1", id)
	allUsers, err := usersList()
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	for rowsUsers.Next() {
		var id string
		err := rowsUsers.Scan(&id)
		if err != nil {
			ret.Error = err.Error()
		} else {
			usersGroup = append(usersGroup, id)
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
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			ret.Error = err.Error()
			return c.RenderJson(ret)
		} else {
			ret.Id = id
			ret.Name = name
		}
	}

	//get parents
	rows, err = DB.Query("SELECT parent_id,level FROM "+TABLE_GROUPS_STRUCT+" WHERE group_id=$1", id)
	if err != nil {
		ret.Error = err.Error()
		return c.RenderJson(ret)
	}
	for rows.Next() {
		var id string
		var level int
		err := rows.Scan(&id, &level)
		if err != nil {
			ret.Error = err.Error()
			return c.RenderJson(ret)
		} else {
			parent := GroupParentItem{}
			parent.Id = id
			parent.Level = level
			ret.Parents = append(ret.Parents, parent)
		}
	}

	return c.RenderJson(ret)
}

func (c GroupCntl) Data() revel.Result {
	var ret GroupData
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

func (c GroupCntl) Delete(id string) revel.Result {
	ret := make(map[string]string)
	_, err := DB.Exec("DELETE FROM "+TABLE_GROUPS+" WHERE id=$1", id)
	if err != nil {
		revel.ERROR.Println("[delete group]", err)
		ret["error"] = err.Error()
	}
	return c.RenderJson(ret)
}