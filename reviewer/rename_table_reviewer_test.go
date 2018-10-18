package reviewer

import (
	"testing"
	"fmt"
	"github.com/daiguadaidai/blingbling/parser"
	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/config"
)

func TestRenameTableReviewer_Review(t *testing.T) {
	var host string = "10.10.10.12"
	var port int = 3306
	var username string = "HH"
	var password string = "oracle"

	sql := `
		rename table test.table to test1.t1, t2 to tt2;
    `

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}


	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username ,password, "")
	reviewMSGs := make([]*ReviewMSG, 0, 1)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		renameStmt := stmtNode.(*ast.RenameTableStmt)
		for i, subStmt := range renameStmt.TableToTables {
			fmt.Printf(
				"%v: %v -> %v\n",
				i,
				subStmt.OldTable.Name.String(),
				subStmt.NewTable.Name.String(),
			)
		}

		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		reviewMSGs = append(reviewMSGs, reviewMSG)
	}

	for _, reviewMSG := range reviewMSGs {
		if reviewMSG != nil {
			fmt.Println(reviewMSG.String())
		}
	}
}
