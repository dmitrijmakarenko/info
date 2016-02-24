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
}

type UsersList struct {
	Error string `json:"error"`
	Users []UserItem `json:"accounts"`
}

type UserSettings struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Position string `json:"position"`
}

func (c UsersCntl) List() revel.Result {
	var ret UsersList
	rows, err := DB.Query("SELECT id FROM sys_users")
	if err != nil {
		revel.ERROR.Println("[get accounts]", err)
		ret.Error = err.Error()
	} else {
		for rows.Next() {
			var id string
			err := rows.Scan(&id)
			if err != nil {
				revel.ERROR.Println(err)
			} else {
				user := UserItem{}
				user.Id = id
				ret.Users = append(ret.Users, user)
			}
		}
	}
	return c.RenderJson(ret)
}

func (c UsersCntl) Get(id string) revel.Result {
	revel.INFO.Println("[get user]", id)
	var ret UserSettings
	rows, err := DB.Query("SELECT name,position FROM sys_users WHERE id=$1", id)
	if err != nil {
		revel.ERROR.Println("[get user]", err)
		ret := make(map[string]string)
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	} else {
		for rows.Next() {
			var name string
			var position string
			err := rows.Scan(&name, &position)
			if err != nil {
				revel.ERROR.Println(err)
			} else {
				revel.INFO.Println("[get user name]", name)
				revel.INFO.Println("[get user position]", position)
				ret.Id = id
				ret.Name = name
				ret.Position = position
			}
		}
	}
	return c.RenderJson(ret)
}

func (c UsersCntl) Update(create bool,data string) revel.Result {
	revel.INFO.Println("[update user data]", data)
	revel.INFO.Println("[update user create]", create)
	ret := make(map[string]string)
	var settings UserSettings
	err := json.Unmarshal([]byte(data), &settings)
	if err != nil {
		ret["error"] = "settings error format";
	} else {
		if create {
			_, err = DB.Exec("INSERT INTO sys_users(id, name, position) VALUES ($1, $2, $3)", settings.Id, settings.Name, settings.Position)
		} else {
			_, err = DB.Exec("UPDATE sys_users SET id=$1, name=$2, position=$3 WHERE id=$1", settings.Id, settings.Name, settings.Position)
		}
		if err != nil {
			ret["error"] = err.Error();
		}
	}
	return c.RenderJson(ret)
}