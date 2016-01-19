package controllers

import (
	"github.com/robfig/revel"
	_ "strconv"
	_ "encoding/json"
	"encoding/json"
	"strconv"
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

type AccountItem struct {
	Id string `json:"id"`
}

type AccountsList struct {
	Error string `json:"error"`
	Accounts []AccountItem `json:"accounts"`
}

type RuleItem struct {
	Id string `json:"id"`
	Desc string `json:"desc"`
}

type RulesList struct {
	Error  string `json:"error"`
	Rules []RuleItem `json:"rules"`
}

type TableDB struct {
	Name string `json:"name"`
	Protected bool `json:"protected"`
}

func (c App) Index() revel.Result {
	return c.Render()
}

//func (c App) GetEntities() revel.Result {
//	revel.INFO.Println("get configs")
//	var out []map[string]string
//	cfg := readXML()
//	for _, ent := range cfg.Entities {
//		var entity = make(map[string]string)
//		entity["id"] = ent.Id
//		entity["name"] = ent.Name
//		entity["cntProps"] = strconv.Itoa(len(ent.Properties))
//		out = append(out, entity)
//	}
//	return c.RenderJson(out)
//}

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
//	var out []map[string]string
//	entities := entityList();
//	for _, ent := range entities {
//		var entity = make(map[string]string)
//		entity["id"] = ent.Id
//		entity["name"] = ent.Name
//		out = append(out, entity)
//	}
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

func (c App) ValidateDB() revel.Result {
	revel.INFO.Println("validate db")
	cntCfgTables := 0
	cntConfirmTables := 0
	cntExcessTables := 0
	ret := make(map[string]string)
	tablesDB := listTables()
	revel.INFO.Println("tablesDB", tablesDB)
	cfg := readXML()
	for _, table := range tablesDB {
		usage := false
		for _, ent := range cfg.Entities {
			//listColumns(ent)
			if (table == ent.Id) {
				cntConfirmTables++
				usage = true
				break
			}
		}
		if (!usage) {
			cntExcessTables++
		}
	}
	cntCfgTables = len(cfg.Entities)
	ret["totalCfgTables"] = strconv.Itoa(cntCfgTables)
	ret["cntConfirmTables"] = strconv.Itoa(cntConfirmTables)
	ret["cntExcessTables"] = strconv.Itoa(cntExcessTables)
	return c.RenderJson(ret)
}

func (c App) ClearDB() revel.Result {
	revel.INFO.Println("clear db")
	ret := make(map[string]string)
	//cfg := readXML()
	return c.RenderJson(ret)
}

func (c App) GetTables() revel.Result {
	var ret []TableDB
	tables := listTables()
	//cfg := readXML()
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
	revel.INFO.Println("protect table", table)
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
	//cfg := readXML()
	for _, table := range tables {
		var tableItem TableDB
		protect := isProtect(table)
		tableItem.Protected = protect
		tableItem.Name = table
		ret = append(ret, tableItem)
	}
	return c.RenderJson(ret)
}

func (c App) GetAccounts() revel.Result {
	var ret AccountsList
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
				account := AccountItem{}
				account.Id = id
				ret.Accounts = append(ret.Accounts, account)
			}
		}
	}
	return c.RenderJson(ret)
}

func (c App) GetRules() revel.Result {
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

func (c App) TestSelect() revel.Result {
	user := "user1"
	action := "select"
	rows, err := DB.Query("SELECT rule FROM rules_p WHERE rule_role = $1 AND action = $2", user, action)
	if err != nil {
		revel.ERROR.Println("[test] select ", err)
	} else {
		for rows.Next() {
			var rule string
			err := rows.Scan(&rule)
			if err != nil {
				revel.ERROR.Println(err)
			} else {
				revel.INFO.Println(rule)
			}
		}
	}
	ret := make(map[string]string)
	return c.RenderJson(ret)
}