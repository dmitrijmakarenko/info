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

type UserSettings struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Position string `json:"position"`
}

func (c UsersCntl) List() revel.Result {
	var ret UsersList
	rows, err := DB.Query("SELECT id, COALESCE(name, '') as name FROM " + TABLE_USERS+" ORDER BY id")
	if err != nil {
		revel.ERROR.Println("[get accounts]", err)
		ret.Error = err.Error()
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
	var ret UserSettings
	rows, err := DB.Query("SELECT name,position FROM "+TABLE_USERS+" WHERE id=$1", id)
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
				ret.Id = id
				ret.Name = name
				ret.Position = position
			}
		}
	}
	return c.RenderJson(ret)
}

func (c UsersCntl) Update(id string, data string) revel.Result {
	revel.INFO.Println("[update user data]", data)
	ret := make(map[string]string)
	var settings UserSettings
	err := json.Unmarshal([]byte(data), &settings)
	if err != nil {
		ret["error"] = "settings error format";
	} else {
		if id == "!new" {
			_, err = DB.Exec("INSERT INTO "+TABLE_USERS+"(id, name, position) VALUES ($1, $2, $3)", settings.Id, settings.Name, settings.Position)
		} else {
			_, err = DB.Exec("UPDATE "+TABLE_USERS+" SET id=$2, name=$3, position=$4 WHERE id=$1", id, settings.Id, settings.Name, settings.Position)
		}
		if err != nil {
			ret["error"] = err.Error();
		}
	}
	return c.RenderJson(ret)
}

func (c UsersCntl) Delete(id string) revel.Result {
	ret := make(map[string]string)
	_, err := DB.Exec("DELETE FROM "+TABLE_USERS+" WHERE id=$1", id)
	if err != nil {
		revel.ERROR.Println("[delete user]", err)
		ret["error"] = err.Error()
	}
	return c.RenderJson(ret)
}