package api

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *SQL

type SQL struct {
	connection *sql.DB
}

type Column struct {
	Name    string
	Type    string
	Null    bool
	Key     string
	Default string
	Unique  bool
	Extra   string
}

type Table struct {
	Name    string
	Columns []Column
}

func NewSQL(dataSourceName string) *SQL {
	//db, err := sql.Open("mysql", dataSourceName)
	db, err := sql.Open("sqlite3", "./databases/accounts.db")

	if err != nil {
		log.Println(err)
	}
	if err = db.Ping(); err != nil {
		log.Println(err)
	}
	return &SQL{connection: db}
}

func (s *SQL) CreateTable(table Table) error {
	query := "CREATE TABLE IF NOT EXISTS `" + table.Name + "` ("
	for i, column := range table.Columns {
		query += "`" + column.Name + "` " + column.Type
		if !column.Null {
			query += " NOT NULL"
		}
		if column.Key != "" {
			query += " " + column.Key
		}
		if column.Default != "" {
			query += " DEFAULT '" + column.Default + "'"
		}
		if column.Unique {
			query += " UNIQUE"
		}
		if column.Extra != "" {
			query += " " + column.Extra
		}
		if i < len(table.Columns)-1 {
			query += ", "
		}
	}
	query += ")"
	_, err := s.connection.Exec(query)
	return err
}

func (s *SQL) InsertInto(tableName string, columns []string, values []interface{}) error {
	query := "INSERT INTO `" + tableName + "` ("
	for i, column := range columns {
		query += "`" + column + "`"
		if i < len(columns)-1 {
			query += ", "
		}
	}
	query += ") VALUES ("
	for i := 0; i < len(values); i++ {
		query += "?"
		if i < len(values)-1 {
			query += ", "
		}
	}
	query += ")"
	stmt, err := s.connection.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(values...)
	return err
}

func (s *SQL) Update(tableName string, set map[string]interface{}, where map[string]interface{}) error {
	query := "UPDATE `" + tableName + "` SET "
	i := 0
	for column := range set {
		query += "`" + column + "`=?"
		if i < len(set)-1 {
			query += ", "
		}
		i++
	}
	query += " WHERE "
	i = 0
	for column := range where {
		query += "`" + column + "`=?"
		if i < len(where)-1 {
			query += " AND "
		}
		i++
	}
	stmt, err := s.connection.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var values []interface{}
	for _, value := range set {
		values = append(values, value)
	}
	for _, value := range where {
		values = append(values, value)
	}
	_, err = stmt.Exec(values...)
	return err
}

func (s *SQL) Get(tableName string, columns []string, where map[string]interface{}) ([]map[string]interface{}, error) {
	query := "SELECT "
	for i, column := range columns {
		query += "`" + column + "`"
		if i < len(columns)-1 {
			query += ", "
		}
	}
	query += " FROM `" + tableName + "` WHERE "
	i := 0
	for column := range where {
		query += "`" + column + "`=?"
		if i < len(where)-1 {
			query += " AND "
		}
		i++
	}

	stmt, err := s.connection.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var values []interface{}
	for _, value := range where {
		values = append(values, value)
	}

	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	columnsData, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columnsData)
	values = make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for i := 0; i < count; i++ {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}
		row := make(map[string]interface{})
		for i, col := range columnsData {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			row[col] = v
		}
		result = append(result, row)
	}

	return result, nil
}

func (s *SQL) Close() {
	s.connection.Close()
}

func (db *SQL) InitTables() {
	table := Table{
		Name: "users",
		Columns: []Column{
			{Name: "id", Type: "INTEGER PRIMARY KEY AUTOINCREMENT"},
			{Name: "name", Type: "VARCHAR(255)", Null: false},
			{Name: "email", Type: "VARCHAR(255)", Null: false, Unique: true},
			{Name: "password", Type: "VARCHAR(255)", Null: false},
		},
	}

	err := db.CreateTable(table)
	if err != nil {
		panic(err)
	}
}
