package controllers

import (
	"github.com/robfig/revel"
	"encoding/json"
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

//func (c App) GetEntity(id string, returnData bool) revel.Result {
//	cfg := readXML()
//	ret := make(map[string]string)
//	for _, ent := range cfg.Entities {
//		if (ent.Id == id) {
//			if (returnData) {
//				data, err := getEntitySql(ent)
//				if err != nil {
//					ret["error"] = "get entity error";
//					return c.RenderJson(ret)
//				}
//				return c.RenderJson(data)
//			} else {
//				return c.RenderJson(ent)
//			}
//		}
//	}
//	ret["error"] = "entity not found";
//	return c.RenderJson(ret)
//}

func (c App) GetEntities() revel.Result {
	return c.RenderJson(entityList())
}

func (c App) GetEntity(id string) revel.Result {
	ret := make(map[string]string)
	data, err := entityGet(id)
	if err != nil {
		ret["error"] = "get entity error";
		return c.RenderJson(ret)
	}
	return c.RenderJson(data)
}

func (c App) CreateEntity(id string, name string, props string) revel.Result {
	var err error
	var properties []Property
	var entityNew Entity
	ret := make(map[string]string)
	revel.INFO.Println("create entity", id)
	err = json.Unmarshal([]byte(props), &properties)
	if err != nil {
		ret["error"] = "properties error format";
		return c.RenderJson(ret)
	}
	cfg := readXML()
	entityNew.Name = name
	entityNew.Id = id
	entityNew.Properties = properties
	err = createTable(entityNew)
	if err == nil {
		cfg.Entities = append(cfg.Entities, entityNew)
		generateXML(cfg)
	} else {
		ret["error"] = err.Error();
	}
	return c.RenderJson(ret)
}

func (c App) UpdateEntity(id string, name string, props string) revel.Result {
	revel.INFO.Println("update entity", id)
	var err error
	var entities []Entity
	var properties []Property
	ret := make(map[string]string)
	cfg := readXML()
	err = json.Unmarshal([]byte(props), &properties)
	if err != nil {
		ret["error"] = "properties error format";
		return c.RenderJson(ret)
	}
	for _, ent := range cfg.Entities {
		if (ent.Id == id) {
			entity := Entity{}
			entity.Id = id
			entity.Name = name
			entity.Properties = properties
			entities = append(entities, entity)
		} else {
			entities = append(entities, ent)
		}
	}
	cfg.Entities = entities
	generateXML(cfg)
	return c.RenderJson(ret)
}

func (c App) RemoveEntity() revel.Result {
	var id string = c.Params.Get("id")
	ret := make(map[string]string)
	revel.INFO.Println("remove entity", id)
	cfg := readXML()
	var entities []Entity
	for _, ent := range cfg.Entities {
		if (ent.Id != id) {
			entities = append(entities, ent)
		}
	}
	cfg.Entities = entities;
	err := dropTable(id)
	if err == nil {
		cfg.Entities = entities;
		generateXML(cfg)
	} else {
		ret["error"] = err.Error()
	}
	return c.RenderJson(ret)
}

func (c App) GenerateDB() revel.Result {
	revel.INFO.Println("generate db")
	var err error
	ret := make(map[string]string)
	cfg := readXML()
	for _, ent := range cfg.Entities {
		err = createTable(ent)
		if err != nil {
			ret["error"] = "can't create table" + ent.Id;
			return c.RenderJson(ret)
		}
	}
	return c.RenderJson(ret)
}

func (c App) Auth(user string, token string) revel.Result {
	ret := make(map[string]string)
	_, err := DB.Exec("INSERT INTO sys_tokens(user_id, token) VALUES ($1, $2)", user, token)
	if err != nil {
		ret["error"] = err.Error()
	}
	users := listUsers()
	exist := false
	for _, u := range users {
		if u == user {
			exist = true
		}
	}
	if !exist {
		_, err := DB.Exec("CREATE ROLE " + user + " LOGIN")
		if err != nil {
			ret["error"] = err.Error()
		}
	}
	return c.RenderJson(ret)
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

func (c App) Protect(table string) revel.Result {
	ret := make(map[string]string)
	err := makeProtect(table)
	if err != nil {
		ret["error"] = "can't protect table " + table;
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