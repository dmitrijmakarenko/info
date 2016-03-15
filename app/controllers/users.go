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
	Error string `json:"error"`
}

func (c UsersCntl) List() revel.Result {
	var ret UsersList
	rows, err := DB.Query("SELECT id, COALESCE(realname, '') as name FROM " + TABLE_USERS+" ORDER BY id")
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
	var ret UserSettings
	rows, err := DB.Query("SELECT realname, position_user FROM "+TABLE_USERS+" WHERE id=$1", id)
	if err != nil {
		ret.Error = err.Error()
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
	ret := make(map[string]string)
	var settings UserSettings
	err := json.Unmarshal([]byte(data), &settings)
	if err != nil {
		ret["error"] = "settings error format";
	} else {
		if id == "!new" {
			_, err = DB.Exec("INSERT INTO "+TABLE_USERS+"(record_uuid, id, realname, position_user) VALUES (uuid_generate_v4(), $1, $2, $3)", settings.Id, settings.Name, settings.Position)
		} else {
			_, err = DB.Exec("UPDATE "+TABLE_USERS+" SET id=$2, realname=$3, position_user=$4 WHERE id=$1", id, settings.Id, settings.Name, settings.Position)
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
		ret["error"] = err.Error()
	}
	return c.RenderJson(ret)
}