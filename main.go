package main

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

type Config struct {
	OldDSN string
	NewDSN string
}

var config Config

func main() {
	if len(os.Args) >= 2 {

		config.OldDSN = os.Args[0]
		config.NewDSN = os.Args[1]
	} else {
		fmt.Println("⚠️命令参数中未查询到数据库连接信息，将从环境变量获取⚠️")
		fmt.Println("⚠️环境变量ONEAPI_OLD_SQL_DSN:songquanpeng/one-api数据库的连接字符串⚠️")
		fmt.Println("⚠️环境变量ONEAPI_NEW_SQL_DSN:MartialBE/one-api数据库的连接字符串⚠️")
		config = loadConfig()
	}

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
	fmt.Scanln()
}

func loadConfig() Config {
	return Config{
		OldDSN: os.Getenv("ONEAPI_OLD_SQL_DSN"),
		NewDSN: os.Getenv("ONEAPI_NEW_SQL_DSN"),
	}
}

func openDatabase(dsn string) *sql.DB {
	driver, dsn := detectDriver(dsn)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	return db
}

func detectDriver(dsn string) (string, string) {
	if strings.Contains(dsn, "postgres://") {
		return "postgres", strings.Split(dsn, "postgres://")[1]
	} else if strings.Contains(dsn, "mysql://") {
		return "mysql", strings.Split(dsn, "mysql://")[1]
	}
	return "sqlite", dsn
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
	driver, _ := detectDriver(config.NewDSN)
	insertSQL := buildInsertSQL(table, newColumns, oldColumns, driver)

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
		insertValues := buildInsertValues(values, oldColumns, newColumns, table)
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

func buildInsertSQL(table string, newColumns, oldColumns []string, driver string) string {
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
	switch driver {
	case "mysql":
		return fmt.Sprintf("INSERT IGNORE INTO `%s` (%s) VALUES (%s)", table, strings.Join(columns, ","), placeholders)
	case "sqlite":
		return fmt.Sprintf("INSERT OR IGNORE INTO `%s` (%s) VALUES (%s)", table, strings.Join(columns, ","), placeholders)
	case "postgres":
		return fmt.Sprintf("INSERT INTO \"%s\" (%s) VALUES (%s) ON CONFLICT DO NOTHING", table, strings.Join(columns, ","), placeholders)
	default:
		log.Fatalf("不支持的数据库驱动: %s", driver)
		return ""
	}
}

func buildInsertValues(values []interface{}, oldColumns, newColumns []string, table string) []interface{} {
	insertValues := []interface{}{}
	for _, col := range newColumns {
		if idx := indexOf(oldColumns, col); idx != -1 {
			value := values[idx]
			if table == "channels" && col == "type" {
				fmt.Println("🔗 处理渠道类别数据")
				value = upgradeChannelType(value)
			}
			insertValues = append(insertValues, value)
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
func BytesToInt(b []uint8) int {
	if len(b) < 4 {
		return 0
	}
	return int(binary.BigEndian.Uint32(b))
}
