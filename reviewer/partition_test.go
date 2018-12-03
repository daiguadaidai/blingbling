package reviewer

import (
	"fmt"
	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/parser"
	"testing"
)

func Test_Parition(t *testing.T) {
	/*
			sql := `
		CREATE TABLE tblist (
		    id INT NOT NULL,
		    store_id INT
		)
		PARTITION BY LIST(store_id) (
		    PARTITION a VALUES IN (1,5,6),
		    PARTITION b VALUES IN (2,7,8),
		    PARTITION c VALUES IN (3,9,10),
		    PARTITION d VALUES IN (4,11,12)
		);
		    `
	*/

	sql := `
CREATE TABLE tblist (
    id INT NOT NULL,
    store_id INT
)
PARTITION BY LIST COLUMNS(id) (
    PARTITION a VALUES IN (1,5,6),
    PARTITION b VALUES IN (2,7,8)
);
`

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	for _, stmtNode := range stmtNodes {
		switch stmt := stmtNode.(type) {
		case *ast.CreateTableStmt:
			// stmt.Partition.ColumnNames
			fmt.Println("table name:", stmt.Table.Name)
			fmt.Println("partition type:", stmt.Partition.Tp)
			fmt.Println("partition list columns:", stmt.Partition.ColumnNames)
			fmt.Println("partition defs:", stmt.Partition.Definitions)
			for _, def := range stmt.Partition.Definitions {
				fmt.Println("partion def name: ", def.Name)
				for _, value := range def.LessThan {
					fmt.Printf("%v, ", value.GetValue())
				}
				fmt.Println()
			}
		}
	}
}

func Test_Parition_Range(t *testing.T) {
	sql := `
CREATE TABLE test.mf_fd_cache (
  id bigint(18) NOT NULL AUTO_INCREMENT COMMENT '注释',
  dep varchar(3) NOT NULL DEFAULT '' Comment '注释',
  arr varchar(3) NOT NULL DEFAULT '' Comment '注释',
  flightNo varchar(10) NOT NULL DEFAULT '' Comment '注释',
  flightDate date NOT NULL DEFAULT '1000-10-10' Comment '注释',
  flightTime varchar(20) NOT NULL DEFAULT '' Comment '注释',
  isCodeShare tinyint(1) Comment '注释',
  tax int(11) NOT NULL DEFAULT '0' Comment '注释',
  yq int(11) NOT NULL DEFAULT '0' Comment '注释',
  cabin char(2) NOT NULL DEFAULT '' Comment '注释',
  ibe_price int(11) NOT NULL DEFAULT '0' Comment '注释',
  ctrip_price int(11) NOT NULL DEFAULT '0' Comment '注释',
  official_price int(11) NOT NULL DEFAULT '0' Comment '注释',
  uptime datetime NOT NULL DEFAULT '1000-10-10 10:10:10' Comment '注释',
  PRIMARY KEY (id, uptime),
  UNIQUE KEY udx_uid (dep, arr, flightNo, uptime),
  -- UNIQUE KEY udx_uid_1 (cabin),
  Index idx_uptime (uptime),
  KEY idx_flight (dep,arr),
  KEY idx_flightdate (flightDate)
) ENGINE=InnoDb  DEFAULT CHARSET=Utf8mb4 COLLATE=Utf8mb4_general_ci comment="注释"
/*!50100 PARTITION BY RANGE(TO_DAYS (uptime1))
(
    PARTITION p0 VALUES LESS THAN (TO_DAYS('2010-04-15')),
    PARTITION p1 VALUES LESS THAN (TO_DAYS('2010-05-01')),
    PARTITION p2 VALUES LESS THAN (TO_DAYS('2010-05-15')),
    PARTITION p3 VALUES LESS THAN (TO_DAYS('2010-05-31')),
    PARTITION p4 VALUES LESS THAN (TO_DAYS('2010-06-15')),
    PARTITION p19 VALUES LESS ThAN  MAXVALUE
)*/;
`

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	for _, stmtNode := range stmtNodes {
		switch stmt := stmtNode.(type) {
		case *ast.CreateTableStmt:
			// stmt.Partition.ColumnNames
			fmt.Println(stmt.Table.Name)
		}
	}
}

func Test_Parition_ListValues(t *testing.T) {
	/*
			sql := `
		CREATE TABLE tblist (
		    id INT NOT NULL,
		    store_id INT
		)
		PARTITION BY LIST(store_id) (
		    PARTITION a VALUES IN (1,5,6),
		    PARTITION b VALUES IN (2,7,8),
		    PARTITION c VALUES IN (3,9,10),
		    PARTITION d VALUES IN (4,11,12)
		);
		    `
	*/

	sql := `
CREATE TABLE tblist (
    id INT NOT NULL,
    store_id INT
)
PARTITION BY LIST COLUMNS(id) (
    PARTITION a VALUES IN (1,5,6),
    PARTITION b VALUES IN (2,7,8)
);
`

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	for _, stmtNode := range stmtNodes {
		switch stmt := stmtNode.(type) {
		case *ast.CreateTableStmt:
			// stmt.Partition.ColumnNames
			fmt.Println("table name:", stmt.Table.Name)
			fmt.Println("partition type:", stmt.Partition.Tp)
			fmt.Println("partition list columns:", stmt.Partition.ColumnNames)
			fmt.Println("partition defs:", stmt.Partition.Definitions)
			for _, def := range stmt.Partition.Definitions {
				fmt.Println(def.Name)
				for _, values := range def.ValuesList {
					for _, value := range values {
						fmt.Printf("%v, ", value.GetValue())
					}
					fmt.Println()
				}
			}
		}
	}
}
