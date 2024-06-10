package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	OldDSN string
	NewDSN string
}

func main() {
	config := loadConfig()
	oldDB := openDatabase(config.OldDSN)
	newDB := openDatabase(config.NewDSN)

	tables := []string{"abilities", "channels", "logs", "options", "redemptions", "tokens", "users"}
	fmt.Println("🚩数据处理开始🚩")
	fmt.Println("======================")
	for _, table := range tables {
		fmt.Printf("🚀 正在处理表: %s\n", table)
		migrateTable(oldDB, newDB, table)
		fmt.Printf("✅ 完成处理表: %s\n", table)
	}
	fmt.Println("======================")
	fmt.Println("🚩数据处理完成🚩")
}

func loadConfig() Config {
	return Config{
		OldDSN: os.Getenv("ONEAPI_OLD_SQL_DSN"),
		NewDSN: os.Getenv("ONEAPI_NEW_SQL_DSN"),
	}
}

func openDatabase(dsn string) *sql.DB {
	db, err := sql.Open(detectDriver(dsn), dsn)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	return db
}

func detectDriver(dsn string) string {
	if strings.Contains(dsn, "postgres") {
		return "postgres"
	} else if strings.Contains(dsn, "sqlite") {
		return "sqlite3"
	}
	return "mysql"
}

func migrateTable(oldDB, newDB *sql.DB, table string) {
	oldColumns := getColumns(oldDB, table)
	newColumns := getColumns(newDB, table)

	if len(newColumns) == 0 {
		fmt.Printf("⚠️ 新库中没有找到表: %s\n", table)
		return
	}

	missingColumns := findMissingColumns(oldColumns, newColumns)
	if len(missingColumns) > 0 {
		fmt.Printf("⚠️ 旧库中的表 %s 存在新库中没有的字段: %v\n", table, missingColumns)
	}

	rows, err := oldDB.Query(fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		log.Fatalf("查询旧库表 %s 失败: %v", table, err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	insertSQL := buildInsertSQL(table, newColumns, oldColumns)

	tx, err := newDB.Begin()
	if err != nil {
		log.Fatalf("开启事务失败: %v", err)
	}

	count := 0
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatalf("扫描行数据失败: %v", err)
		}

		insertValues := buildInsertValues(values, oldColumns, newColumns)
		_, err = tx.Exec(insertSQL, insertValues...)
		if err != nil {
			log.Fatalf("插入新库表 %s 失败: %v", table, err)
		}
		count++
		if count%100 == 0 {
			fmt.Printf("⏳ 已处理 %d 行数据\n", count)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("提交事务失败: %v", err)
	}

	fmt.Printf("✅ 表 %s 迁移完成，共处理 %d 行数据\n", table, count)
}

func getColumns(db *sql.DB, table string) []string {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table))
	if err != nil {
		return nil
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("获取表 %s 列信息失败: %v", table, err)
	}

	return columns
}

func findMissingColumns(oldColumns, newColumns []string) []string {
	missingColumns := []string{}
	for _, col := range oldColumns {
		if !contains(newColumns, col) {
			missingColumns = append(missingColumns, col)
		}
	}
	return missingColumns
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func buildInsertSQL(table string, newColumns, oldColumns []string) string {
	columns := []string{}
	for _, col := range newColumns {
		if contains(oldColumns, col) {
			columns = append(columns, fmt.Sprintf("`%s`", col))
		} else {
			columns = append(columns, fmt.Sprintf("`%s`", col))
		}
	}
	placeholders := strings.Repeat("?,", len(columns))
	placeholders = placeholders[:len(placeholders)-1]
	return fmt.Sprintf("INSERT IGNORE INTO `%s` (%s) VALUES (%s)", table, strings.Join(columns, ","), placeholders)
}

func buildInsertValues(values []interface{}, oldColumns, newColumns []string) []interface{} {
	insertValues := []interface{}{}
	for _, col := range newColumns {
		if idx := indexOf(oldColumns, col); idx != -1 {
			insertValues = append(insertValues, values[idx])
		} else {
			insertValues = append(insertValues, getDefaultForType(reflect.TypeOf(values[0])))
		}
	}
	return insertValues
}

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

func getDefaultForType(t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.String:
		return ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return 0
	case reflect.Float32, reflect.Float64:
		return 0.0
	case reflect.Bool:
		return false
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface:
		return nil
	default:
		return ""
	}
}
