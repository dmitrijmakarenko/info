package controllers

import (
	"github.com/robfig/revel"
	"database/sql"
	_ "github.com/lib/pq"
	_ "strconv"
	"errors"
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
	revel.INFO.Println("make protect", table)
	_, err = DB.Exec("ALTER TABLE " + table + " ADD COLUMN rule uuid")
	if err != nil {
		revel.ERROR.Println("[make protect] add column ", err)
		return err
	}
	_, err = DB.Exec("ALTER TABLE " + table + " RENAME TO " + table + "_protected")
	if err != nil {
		revel.ERROR.Println("[make protect] rename ", err)
		return err
	}
	_, err = DB.Exec("CREATE OR REPLACE VIEW " + table + " AS SELECT * FROM " + table + "_protected")
	if err != nil {
		revel.ERROR.Println("[make protect] create view ", err)
		return err
	}
	users := listUsers()
	for _, user := range users {
		_, err = DB.Exec("GRANT ALL PRIVILEGES ON " + table + " TO " + user)
		if err != nil {
			revel.ERROR.Println("[make protect] grant for " + user, err)
		}
		_, err = DB.Exec("REVOKE ALL PRIVILEGES ON " + table + "_protected FROM " + user)
		if err != nil {
			revel.ERROR.Println("[make protect] revoke for " + user, err)
		}
	}
	return nil
}

func isProtect(table string) (protect bool) {

	return false
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
