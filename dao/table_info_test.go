package dao

import (
	"testing"
	"fmt"
	"github.com/daiguadaidai/blingbling/config"
)


func TestTableInfo_DatabaseExists(t *testing.T) {
	host := "10.10.10.12"
	port := 3306
	username := "HH"
	password := "oracle"
	database := "test"
	table := "t1"

	dbConfig := config.NewDBConfig(host, port, username, password, database)
	tableInfo := NewTableInfo(dbConfig, table)

	if err := tableInfo.OpenInstance(); err != nil {
		t.Fatalf("打开数据库出错: %v", err)
	}

	exist, err := tableInfo.DatabaseExists()
	if err != nil {
		t.Fatalf("执行数据库是否存在报错: %v", err)
	}
	fmt.Println("表是否存在:", exist)

	if err := tableInfo.CloseInstance(); err != nil {
		t.Fatalf("关闭数据库出错: %v", err)
	}
}

func TestTableInfo_TableExists(t *testing.T) {
	host := "10.10.10.12"
	port := 3306
	username := "HH"
	password := "oracle"
	database := "test"
	table := "t1"

	dbConfig := config.NewDBConfig(host, port, username, password, database)
	tableInfo := NewTableInfo(dbConfig, table)

	if err := tableInfo.OpenInstance(); err != nil {
		t.Fatalf("打开数据库出错: %v", err)
	}

	exist, err := tableInfo.TableExists()
	if err != nil {
		t.Fatalf("执行表是否存在报错: %v", err)
	}
	fmt.Println("表是否存在:", exist)

	if err := tableInfo.CloseInstance(); err != nil {
		t.Fatalf("关闭数据库出错: %v", err)
	}
}

func TestTableInfo_FindColumnNameMap(t *testing.T) {
	host := "10.10.10.12"
	port := 3306
	username := "HH"
	password := "oracle"
	database := "test"
	table := "t1"

	dbConfig := config.NewDBConfig(host, port, username, password, database)
	tableInfo := NewTableInfo(dbConfig, table)

	if err := tableInfo.OpenInstance(); err != nil {
		t.Fatalf("打开数据库出错: %v", err)
	}

	columnName, err := tableInfo.FindColumnNameMap()
	if err != nil {
		t.Fatalf("执行表是否存在报错: %v", err)
	}
	fmt.Println("表所有列名:", columnName)

	if err := tableInfo.CloseInstance(); err != nil {
		t.Fatalf("关闭数据库出错: %v", err)
	}
}

func TestTableInfo_FindPrimaryKey(t *testing.T) {
	host := "10.10.10.12"
	port := 3306
	username := "HH"
	password := "oracle"
	database := "test"
	table := "t1"

	dbConfig := config.NewDBConfig(host, port, username, password, database)
	tableInfo := NewTableInfo(dbConfig, table)

	if err := tableInfo.OpenInstance(); err != nil {
		t.Fatalf("打开数据库出错: %v", err)
	}

	pkList, pkNameMap, err := tableInfo.FindPrimaryKey()
	if err != nil {
		t.Fatalf("执行表是否存在报错: %v", err)
	}
	fmt.Println("主键List:", pkList)
	fmt.Println("主键Map:", pkNameMap)

	if err := tableInfo.CloseInstance(); err != nil {
		t.Fatalf("关闭数据库出错: %v", err)
	}
}

func TestTableInfo_FindUniqueIndexes(t *testing.T) {
	host := "10.10.10.12"
	port := 3306
	username := "HH"
	password := "oracle"
	database := "test"
	table := "t1"

	dbConfig := config.NewDBConfig(host, port, username, password, database)
	tableInfo := NewTableInfo(dbConfig, table)

	if err := tableInfo.OpenInstance(); err != nil {
		t.Fatalf("打开数据库出错: %v", err)
	}

	uniqueIndexes, err := tableInfo.FindUniqueIndexes()
	if err != nil {
		t.Fatalf("执行表是否存在报错: %v", err)
	}
	fmt.Println("唯一键:", uniqueIndexes)

	if err := tableInfo.CloseInstance(); err != nil {
		t.Fatalf("关闭数据库出错: %v", err)
	}
}

func TestTableInfo_FindNormalIndexes(t *testing.T) {
	host := "10.10.10.12"
	port := 3306
	username := "HH"
	password := "oracle"
	database := "test"
	table := "t1"

	dbConfig := config.NewDBConfig(host, port, username, password, database)
	tableInfo := NewTableInfo(dbConfig, table)

	if err := tableInfo.OpenInstance(); err != nil {
		t.Fatalf("打开数据库出错: %v", err)
	}

	indexes, err := tableInfo.FindAllIndexes()
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println("普通索引:", indexes)
	fmt.Println("主键:", tableInfo.PKColumnNameList)
	fmt.Println("唯一索引:", tableInfo.UniqueIndexes)


	if err := tableInfo.CloseInstance(); err != nil {
		t.Fatalf("关闭数据库出错: %v", err)
	}
}
