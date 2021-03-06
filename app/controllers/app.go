package controllers

import (
	"github.com/robfig/revel"
	"database/sql"
	"io/ioutil"
)

type App struct {
	*revel.Controller
}

type Property struct {
	Name string `json:"id"`
	Type string `json:"type"`
	Desc string `json:"desc"`
}

type Entity struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Properties []Property `json:"props"`
}

type Config struct {
	Name string
	Entities []Entity
}

type TableDB struct {
	Name string `json:"name"`
	Protected bool `json:"protected"`
}

func (c App) Index() revel.Result {
	return c.Render()
}

var TABLE_USERS string
var TABLE_GROUPS string
var TABLE_GROUP_USER string
var TABLE_GROUPS_STRUCT string
var TABLE_RULES string
var TABLE_RULES_DATA string

func InstallFunc(db *sql.DB) error {
	sqlPath := "gocode/src/info/sql/"
	files, _ := ioutil.ReadDir(sqlPath)
	for _, file := range files {
		b, err := ioutil.ReadFile(sqlPath + file.Name())
		if err != nil {
			return err
		}
		query := b[3:len(b)]
		_, err = db.Exec(string(query))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c App) Auth(user string, pass string) revel.Result {
	ret := make(map[string]string)
	var token string
	err := DB.QueryRow("SELECT acs_auth($1, $2)", user, pass).Scan(&token)
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	if token != "" {
		ret["token"] = token
	} else {
		ret["error"] = "wrong login or password"
	}
	return c.RenderJson(ret)
}

func (c App) Protect(table string) revel.Result {
	ret := make(map[string]string)
	_, err := DB.Query("SELECT acs_protect_table($1)", table)
	if err != nil {
		ret["error"] = err.Error()
	}
	return c.RenderJson(ret)
}

func (c App) GetEntities() revel.Result {
	return c.RenderJson(entityList())
}

func (c App) GetTable(id string) revel.Result {
	ret := make(map[string]string)
	data, err := entityGet(id)
	if err != nil {
		ret["error"] = err.Error()
		return c.RenderJson(ret)
	}
	return c.RenderJson(data)
}

func (c App) GetTables() revel.Result {
	var ret []TableDB
	tables := listTables()
	for _, table := range tables {
		var tableItem TableDB
		protect := isProtect(table)
		tableItem.Protected = protect
		tableItem.Name = table
		ret = append(ret, tableItem)
	}
	return c.RenderJson(ret)
}

func (c App) GetViews() revel.Result {
	var ret []TableDB
	tables := listViews()
	for _, table := range tables {
		var tableItem TableDB
		protect := isProtect(table)
		tableItem.Protected = protect
		tableItem.Name = table
		ret = append(ret, tableItem)
	}
	return c.RenderJson(ret)
}

func usersList() (users []UserItem, err error) {
	rows, err := DB.Query("SELECT id, COALESCE(realname, '') as name FROM " + TABLE_USERS)
	if err != nil {
		return users, err
	}
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
			users = append(users, user)
		}
	}
	return users, err
}

func getUsersByGroup(group string) (users []string, err error) {
	rows, err := DB.Query("SELECT user_id FROM "+TABLE_GROUP_USER+" WHERE group_id=$1", group)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return users, err
		} else {
			users = append(users, id)
		}
	}
	return users, err
}

func groupsList() (groups []GroupItem, err error) {
	rows, err := DB.Query("SELECT group_id, realname FROM "+TABLE_GROUPS)
	if err != nil {
		return groups, err
	}
	for rows.Next() {
		var id string
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return groups, err
		} else {
			group := GroupItem{}
			group.Id = id
			group.Name = name
			groups = append(groups, group)
		}
	}
	return groups, err
}