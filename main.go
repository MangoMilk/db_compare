package main

import (
	"database/sql"
	"fmt"

	// "log"

	_ "github.com/go-sql-driver/mysql"
)

const dateFormat = "2006-01-02"

func main() {
	InitConfig()
	InitLog(conf.APP_NAME)

	DBCompare()
}

// ================================================ DB SCHEMA COMPARE TOOL  ================================================
type DB_SCHEMA struct {
	TABLES map[string]DB_TABLE
}

type DB_TABLE struct {
	COLUMNS map[string]T_COLUMN
}

/*
 * COLUMNS: Columns of table info are in this table.
 */
type T_COLUMN struct {
	TABLE_CATALOG            string `db:"TABLE_CATALOG"`
	TABLE_SCHEMA             string `db:"TABLE_SCHEMA"`
	TABLE_NAME               string `db:"TABLE_NAME"`
	COLUMN_NAME              string `db:"COLUMN_NAME"`
	ORDINAL_POSITION         uint64 `db:"ORDINAL_POSITION"`
	COLUMN_DEFAULT           string `db:"COLUMN_DEFAULT"`
	IS_NULLABLE              string `db:"IS_NULLABLE"`
	DATA_TYPE                string `db:"DATA_TYPE"`
	CHARACTER_MAXIMUM_LENGTH uint64 `db:"CHARACTER_MAXIMUM_LENGTH"`
	CHARACTER_OCTET_LENGTH   uint64 `db:"CHARACTER_OCTET_LENGTH"`
	NUMERIC_PRECISION        uint64 `db:"NUMERIC_PRECISION"`
	NUMERIC_SCALE            uint64 `db:"NUMERIC_SCALE"`
	DATETIME_PRECISION       uint64 `db:"DATETIME_PRECISION"`
	CHARACTER_SET_NAME       string `db:"CHARACTER_SET_NAME"`
	COLLATION_NAME           string `db:"COLLATION_NAME"`
	COLUMN_TYPE              string `db:"COLUMN_TYPE"`
	COLUMN_KEY               string `db:"COLUMN_KEY"`
	EXTRA                    string `db:"EXTRA"`
	PRIVILEGES               string `db:"PRIVILEGES"`
	COLUMN_COMMENT           string `db:"COLUMN_COMMENT"`
	GENERATION_EXPRESSION    string `db:"GENERATION_EXPRESSION"`
}

/*
 * GetDBSchema
 */
func GetDBSchema(dbConf DBConf) map[string]DB_SCHEMA {
	var Database string = "information_schema"
	// set config
	var dsn string = dbConf.User + ":" + dbConf.Password + "@tcp(" + dbConf.Host + ":" + dbConf.Port + ")/" + Database + "?charset=utf8mb4"

	// open connection
	MysqlDB, MysqlDBErr := sql.Open("mysql", dsn)
	if MysqlDBErr != nil {
		// log.Println("dsn: " + dsn)
		panic("连接配置错误: " + MysqlDBErr.Error())
	}
	// defer MysqlDB.Close()

	var sqlStr = "select * from COLUMNS WHERE TABLE_SCHEMA not in ('information_schema','performance_schema','mysql','sys')"
	rows, queryErr := MysqlDB.Query(sqlStr)

	if queryErr != nil {
		panic(queryErr)
	}

	var dbSchemaMap = make(map[string]DB_SCHEMA)
	for rows.Next() {

		row := T_COLUMN{}

		rows.Scan(
			&row.TABLE_CATALOG,
			&row.TABLE_SCHEMA,
			&row.TABLE_NAME,
			&row.COLUMN_NAME,
			&row.ORDINAL_POSITION,
			&row.COLUMN_DEFAULT,
			&row.IS_NULLABLE,
			&row.DATA_TYPE,
			&row.CHARACTER_MAXIMUM_LENGTH,
			&row.CHARACTER_OCTET_LENGTH,
			&row.NUMERIC_PRECISION,
			&row.NUMERIC_SCALE,
			&row.DATETIME_PRECISION,
			&row.CHARACTER_SET_NAME,
			&row.COLLATION_NAME,
			&row.COLUMN_TYPE,
			&row.COLUMN_KEY,
			&row.EXTRA,
			&row.PRIVILEGES,
			&row.COLUMN_COMMENT,
			&row.GENERATION_EXPRESSION)
		// log.Println(row)

		// db
		if dbSchemaMap[row.TABLE_SCHEMA].TABLES == nil {
			dbSchemaMap[row.TABLE_SCHEMA] = DB_SCHEMA{
				TABLES: make(map[string]DB_TABLE)}
		}

		// table
		if dbSchemaMap[row.TABLE_SCHEMA].TABLES[row.TABLE_NAME].COLUMNS == nil {
			dbSchemaMap[row.TABLE_SCHEMA].TABLES[row.TABLE_NAME] = DB_TABLE{
				COLUMNS: make(map[string]T_COLUMN)}
		}

		// column
		dbSchemaMap[row.TABLE_SCHEMA].TABLES[row.TABLE_NAME].COLUMNS[row.TABLE_SCHEMA+"."+row.TABLE_NAME+"."+row.COLUMN_NAME] = row

	}

	rows.Close()

	return dbSchemaMap

	/*
	 * TABLE_CONSTRAINTS: Indexs of tables are in this table.
	 */
	// type T_CONSTRAINT struct {
	// 	CONSTRAINT_CATALOG string `db:"CONSTRAINT_CATALOG"`
	// 	CONSTRAINT_SCHEMA  string `db:"CONSTRAINT_SCHEMA"`
	// 	CONSTRAINT_NAME    string `db:"CONSTRAINT_NAME"`
	// 	TABLE_SCHEMA       string `db:"TABLE_SCHEMA"`
	// 	TABLE_NAME         string `db:"TABLE_NAME"`
	// 	CONSTRAINT_TYPE    string `db:"CONSTRAINT_TYPE"`
	// }

	/*
	 * TABLES: Table's info is in this table.
	 */
	// type TABLE struct {
	// }
}

/*
 * DBCompare
 * Compare of A and B, Base on A.
 */
func DBCompare() {
	// get db schema
	A := GetDBSchema(conf.DBA)
	B := GetDBSchema(conf.DBB)

	// compare
	for dbName, tables := range A {
		if B[dbName].TABLES == nil {
			// log.Println("There is no DB:" + dbName)
			fmt.Println("Missing DB:【" + dbName + "】")
			continue
		}

		for tableName, columns := range tables.TABLES {
			if B[dbName].TABLES[tableName].COLUMNS == nil {
				// log.Println("There is no Table:" + tableName)
				fmt.Println("Missing Table:【" + dbName + "." + tableName + "】")
				continue
			}

			for columnName, v := range columns.COLUMNS {

				var tColumn T_COLUMN = B[dbName].TABLES[tableName].COLUMNS[columnName]

				if tColumn.COLUMN_NAME == "" {
					// log.Println("There is no Column:" + columnName)
					fmt.Println("Missing Column:【" + columnName + "】")
					continue
				}

				if tColumn.DATA_TYPE != v.DATA_TYPE {
					fmt.Println("Mismatched Column Data Type:【" + columnName + "】B is [" + tColumn.DATA_TYPE + "]" + " A is [" + v.DATA_TYPE + "]")
					continue
				}
			}
		}
	}
}
