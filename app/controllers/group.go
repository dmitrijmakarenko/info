package controllers

import (
	"github.com/robfig/revel"
	"strings"
)

type GroupCntl struct {
	*revel.Controller
}

type GroupItem struct {
	Id string `json:"id"`
	Name string  `json:"name"`
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

func (c GroupCntl) Update(id string, name string) revel.Result {
	ret := make(map[string]string)
	var err error
	if id == "!new" {
		_, err = DB.Exec("INSERT INTO sys_groups(id, name) VALUES (uuid_generate_v4(), $1)", name)
	} else {
		_, err = DB.Exec("UPDATE sys_groups SET name=$2 WHERE id=$1", id, name)
	}
	if err != nil {
		revel.ERROR.Println("[group update]", err)
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
	revel.INFO.Println("[users list]", usersGroup)
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

func (c GroupCntl) Delete(id string) revel.Result {
	ret := make(map[string]string)
	_, err := DB.Exec("DELETE FROM sys_groups WHERE id=$1", id)
	if err != nil {
		revel.ERROR.Println("[delete group]", err)
		ret["error"] = err.Error()
	}
	return c.RenderJson(ret)
}