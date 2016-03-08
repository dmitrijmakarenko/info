package controllers

import (
	"github.com/robfig/revel"
	"database/sql"
	_ "github.com/lib/pq"
	_ "strconv"
	"errors"
	"strings"
)

var DB *sql.DB

type DataEntity struct {
	Id string `json:"id"`
	Columns []string `json:"columns"`
	Rows [][]string `json:"rows"`
}

func InitDB() {
	connstring := "host=" + DB_HOST + " port=" + DB_PORT + " user=" + DB_USER+ " dbname=" + DB_NAME + " password=" + DB_PASSWORD + " sslmode=disable"

	var err error
	DB, err = sql.Open(DB_DRIVER, connstring)
	if err != nil {
		revel.ERROR.Println("DB Error", err)
	}
	TABLE_USERS = "sys_users"
	TABLE_GROUPS = "sys_groups"
	TABLE_GROUP_USER = "sys_group_user"
	TABLE_GROUPS_STRUCT = "sys_groups_struct"
	TABLE_RULES = "rules"
	TABLE_RULES_P = "rules_p"
	revel.INFO.Println("DB Connected")
}

func createTable(entity Entity) (err error) {
	props := ""
	for i, property := range entity.Properties {
		props += property.Name + " " + property.Type
		if (i != len(entity.Properties) - 1) {
			props += " , "
		}
	}
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS " + entity.Id + " ( " + props + " )")
	if err != nil {
		revel.ERROR.Println("create table error", err)
	}
	return err
}

func dropTable(entityId string) (err error) {
	_, err = DB.Exec("DROP TABLE "+ entityId)
	if err != nil {
		revel.ERROR.Println("drop table error", err)
	}
	return err
}

func listTables() ([]string) {
	var tableArray []string
	rows, err := DB.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
	if err != nil {
		revel.ERROR.Println("list table error", err)
	}
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			revel.ERROR.Println(err)
		} else {
			tableArray = append(tableArray, tableName)
		}
	}
	return tableArray
}

func listViews() ([]string) {
	var tableArray []string
	rows, err := DB.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'VIEW'")
	if err != nil {
		revel.ERROR.Println("list table error", err)
	}
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			revel.ERROR.Println(err)
		} else {
			tableArray = append(tableArray, tableName)
		}
	}
	return tableArray
}

func entityList() ([]string) {
	var tableArray []string
	rows, err := DB.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		revel.ERROR.Println("list table error", err)
	}
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			revel.ERROR.Println(err)
		} else {
			tableArray = append(tableArray, tableName)
		}
	}
	return tableArray
}

func entityGet(table string) (DataEntity, error) {
	var ret DataEntity
	ret.Id = table
	rows, err := DB.Query("SELECT * FROM " + table)
	if err != nil {
		revel.ERROR.Println("get table error", err)
		return ret, err
	}
	columnNames, err := rows.Columns()
	ret.Columns = listColumns(table)
	var retRows [][]string
	for rows.Next() {
		var retRow []string
		columnPointers := make([]interface{}, len(columnNames))
		for i := 0; i < len(columnNames); i++ {
			columnPointers[i] = new(sql.RawBytes)
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return ret, err
		}
		for i := 0; i < len(columnNames); i++ {
			if rb, ok := columnPointers[i].(*sql.RawBytes); ok {
				retRow = append(retRow, string(*rb))
			} else {
				return ret, errors.New("erorr get row")
			}
		}
		retRows = append(retRows, retRow)
	}
	ret.Rows = retRows
	return ret, nil
}

func makeProtect(table string) (err error) {
	rows, err := DB.Query("SELECT protect_table($1) AS num", table)
	if err != nil {
		revel.ERROR.Println("[make protect]", err)
		return err
	}
	for rows.Next() {
		var num int
		err := rows.Scan(&num)
		if err != nil {
			revel.ERROR.Println(err.Error())
		} else {
			revel.INFO.Println("[make protect] result", num)
		}
	}
	/*users := listUsers()
	for _, user := range users {
		_, err = DB.Exec("GRANT ALL PRIVILEGES ON " + table + " TO " + user)
		if err != nil {
			revel.ERROR.Println("[make protect] grant for " + user, err)
		}
		_, err = DB.Exec("REVOKE ALL PRIVILEGES ON " + table + "_protected FROM " + user)
		if err != nil {
			revel.ERROR.Println("[make protect] revoke for " + user, err)
		}
	}*/
	return nil
}

func isProtect(table string) (bool) {
	protect := false
	viewName := strings.Replace(table, "_protected", "", -1)
	if strings.Contains(table, "_protected") && viewExist(viewName) && columnExist(table, "rule") {
		protect = true
	}
	return protect
}

func listUsers() ([]string) {
	var usersArray []string
	rows, err := DB.Query("SELECT usename FROM pg_user")
	if err != nil {
		revel.ERROR.Println("[list users] error ", err)
	}
	for rows.Next() {
		var user string
		err := rows.Scan(&user)
		if err != nil {
			revel.ERROR.Println(err)
		} else {
			usersArray = append(usersArray, user)
		}
	}
	return usersArray
}

func viewExist(view string) (bool) {
	viewExist := false
	rows, err := DB.Query("select table_name from INFORMATION_SCHEMA.views WHERE table_schema = ANY (current_schemas(false))")
	if err != nil {
		revel.ERROR.Println(err)
	}
	for rows.Next() {
		var table_name string
		err := rows.Scan(&table_name)
		if err != nil {
			revel.ERROR.Println(err)
		} else {
			if table_name == view {
				viewExist = true
			}
		}
	}
	return viewExist
}

func columnExist(table string, column string) (bool) {
	columnExist := false
	rows, err := DB.Query("SELECT  c.column_name FROM information_schema.tables t " +
	"JOIN information_schema.columns c ON t.table_name = c.table_name " +
	"WHERE t.table_schema = 'public' AND t.table_catalog = current_database() AND t.table_name = $1", table)
	if err != nil {
		revel.ERROR.Println(err)
	}
	for rows.Next() {
		var columnName string
		err := rows.Scan(&columnName)
		if err != nil {
			revel.ERROR.Println(err)
		} else {
			if columnName == column {
				columnExist = true
			}
		}
	}
	return columnExist
}

func listColumns(table string) ([]string) {
	var columnsArray []string
	rows, err := DB.Query("SELECT  c.column_name, c.data_type FROM information_schema.tables t " +
	"JOIN information_schema.columns c ON t.table_name = c.table_name " +
	"WHERE t.table_schema = 'public' AND t.table_catalog = current_database() AND t.table_name = $1", table)
	if err != nil {
		revel.ERROR.Println("list columns error", err)
	}
	for rows.Next() {
		var columnName string
		var columnType string
		err := rows.Scan(&columnName, &columnType)
		if err != nil {
			revel.ERROR.Println(err)
		} else {
			columnsArray = append(columnsArray, columnName)
		}
	}
	return columnsArray
}
