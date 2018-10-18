package reviewer

import (
	"testing"
	"fmt"
	"github.com/daiguadaidai/blingbling/parser"
	"github.com/daiguadaidai/blingbling/config"
)

func TestDeleteReviewer_Review(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
    DELETE t1, t2
    FROM t3, t2, t1, (
        SELECT
            id,
            name
        FROM tmp
    ) AS t10
    WHERE t1.id = t2.id
        AND t1.id = t3.id
        AND t1.id = t10.id
        AND (
            t1.id = 1
            OR t2.name = 'HH'
        )
        AND t2.id IN (
           SELECT t4.id
           FROM t4, t5
           WHERE t4.id = t5.id
               AND t5.is_delete = 0
        )
    `
    fmt.Sprintf(sql)

    sql1 := `
    DELETE FROM t1 WHERE name = 1 and age = 2 and (ttt = 'a' or ttt = 'b') and kinn > 'HH'
    `
	fmt.Sprintf(sql1)

	sql2 := `
    DELETE FROM employees WHERE emp_no = '10001'
    `
	fmt.Sprintf(sql2)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql2, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username ,password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Println(reviewMSG.String())
	}
}

func TestDeleteReviewer_Review_Limit(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
delete from employees WHERE emp_no = 1
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
		fmt.Println(reviewMSG.String())
	}
}
