package reviewer

import (
	"testing"
	"fmt"
	"github.com/daiguadaidai/blingbling/parser"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/ast"
)

func TestInsertReviewer_Review(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
INSERT INTO test.t1
VALUES(1,2,3,4),(1,2,3,4),(1,2,3,4, 5)
ON DUPLICATE KEY UPDATE field1 = 10, field2 = 20, field3 = 30
    `
    fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username ,password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Printf("Code: %v, MSG: %v \n", reviewMSG.Code, reviewMSG.MSG)

		insertReview := review.(*InsertReviewer)
		fmt.Printf("Table: %T, %[1]v \n", insertReview.StmtNode.Table.TableRefs.Left)
		tableName := insertReview.StmtNode.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName)
		fmt.Println("Schema: ", tableName.Schema.String(), "Table: ", tableName.Name.String())

		fmt.Println("IsIgnore:", insertReview.StmtNode.IgnoreErr)
		fmt.Println("IsReplace:", insertReview.StmtNode.IsReplace)

		fmt.Println("ColumnNames:", insertReview.StmtNode.Columns)
		for i, column := range insertReview.StmtNode.Columns {
			fmt.Println("    ", i, "->", column.Name.String())
		}

		fmt.Println("Value Len:")
		for _, list := range insertReview.StmtNode.Lists {
			fmt.Println("    len:", len(list))
		}

		fmt.Println("Set list:")
		for _, list := range insertReview.StmtNode.Setlist {
			fmt.Println("    ", list.Column.String(), " -> ", list.Expr.GetType(), list.Expr.GetValue())
		}

		fmt.Println("Onduplicate:")
		for _, item := range insertReview.StmtNode.OnDuplicate {
			fmt.Println("    ", item.Column.String(), " -> ", item.Expr.GetType(), item.Expr.GetValue())
		}

	}
}

func TestInsertReviewer_Review_Set(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
INSERT INTO test.t1
SET field1 = 1, field2 = 2, field3 = 3, field4 = 4
ON DUPLICATE KEY UPDATE field1 = 10, field2 = 20, field3 = 30
    `
	fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username ,password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Printf("Code: %v, MSG: %v \n", reviewMSG.Code, reviewMSG.MSG)

		insertReview := review.(*InsertReviewer)
		fmt.Println("Table:", insertReview.StmtNode.Table.TableRefs)

		fmt.Println("IsIgnore:", insertReview.StmtNode.IgnoreErr)
		fmt.Println("IsReplace:", insertReview.StmtNode.IsReplace)

		fmt.Println("ColumnNames:")
		for i, column := range insertReview.StmtNode.Columns {
			fmt.Println("    ", i, "->", column.Name.String())
		}

		fmt.Println("Value Len:")
		for _, list := range insertReview.StmtNode.Lists {
			fmt.Println("    len:", len(list))
		}

		fmt.Println("Set list:")
		for _, list := range insertReview.StmtNode.Setlist {
			fmt.Println("    ", list.Column.String(), " -> ", list.Expr.GetType(), list.Expr.GetValue())
		}

		fmt.Println("Onduplicate:")
		for _, item := range insertReview.StmtNode.OnDuplicate {
			fmt.Println("    ", item.Column.String(), " -> ", item.Expr.GetType(), item.Expr.GetValue())
		}

	}
}

func TestInsertReviewer_Review_Select(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
INSERT INTO test.t1
SELECT field1, field2, field3, field4
FROM test.t2
ON DUPLICATE KEY UPDATE field1 = 10, field2 = 20, field3 = 30
    `
	fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username ,password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Printf("Code: %v, MSG: %v \n", reviewMSG.Code, reviewMSG.MSG)

		insertReview := review.(*InsertReviewer)
		fmt.Printf("Table: %T %[1]v\n", insertReview.StmtNode.Table)

		fmt.Println("IsIgnore:", insertReview.StmtNode.IgnoreErr)
		fmt.Println("IsReplace:", insertReview.StmtNode.IsReplace)

		fmt.Println("ColumnNames:")
		for i, column := range insertReview.StmtNode.Columns {
			fmt.Println("    ", i, "->", column.Name.String())
		}

		fmt.Println("Value Len:")
		for _, list := range insertReview.StmtNode.Lists {
			fmt.Println("    len:", len(list))
		}

		fmt.Println("Set list:")
		for _, list := range insertReview.StmtNode.Setlist {
			fmt.Println("    ", list.Column.String(), " -> ", list.Expr.GetType(), list.Expr.GetValue())
		}

		fmt.Println("Onduplicate:")
		for _, item := range insertReview.StmtNode.OnDuplicate {
			fmt.Println("    ", item.Column.String(), " -> ", item.Expr.GetType(), item.Expr.GetValue())
		}

		fmt.Println("Select:")
		selectStmt := insertReview.StmtNode.Select.(*ast.SelectStmt)
		fmt.Println("    ", selectStmt)
	}
}
