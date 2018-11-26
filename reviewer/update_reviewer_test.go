package reviewer

import (
	"fmt"
	"testing"

	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/parser"
)

func TestUpdateReviewer_Review(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
UPDATE xcs_user_credit_score a1
    ,xcs_user_credit_score a2
SET a1.user_currday_score = a1.user_currday_increment_score + a2.user_currday_score,
    a1.name = 'HH',
    a2.name = (select name from t1 where id = 1 and name = 'HH' limit 0, 1),
    a2.age = (select (select age from t1 where id = 1 limit 1) from t1 where id = (select id from t2 where id = 1 and name = 'HH'))
WHERE a1.pt_day = '2017-09-20'
    AND a2.pt_day = '2017-09-19'
    AND a1.uid = a2.uid
    AND a1.aid = (select id from t2 where id = 3);
    `
	fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username, password, database)
	reviewConfig := config.NewReviewConfig()
	reviewConfig.RuleAllowUpdateHasJoin = true
	reviewConfig.RuleAllowUpdateHasSubClause = true
	reviewConfig.RuleAllowUpdateNoWhere = true
	reviewConfig.RuleAllowUpdateLimit = true
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Println(reviewMSG.String())

		updateReview := review.(*UpdateReviewer)
		fmt.Printf("SetSubClauseWhereCount: %v, ", updateReview.visitor.SetSubClauseWhereCount)
	}
}

func TestUpdateReviewer_Review_02(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
UPDATE xcs_user_credit_score a1
    ,xcs_user_credit_score a2
SET a1.user_currday_score = a1.user_currday_increment_score + a2.user_currday_score,
    a1.name = 'HH'
WHERE a1.pt_day = '2017-09-20'
    AND a2.pt_day = '2017-09-19'
    AND a1.uid = a2.uid
    AND a1.aid = (select id from t2 where id = 3);
    `
	fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username, password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Println(reviewMSG.String())

		updateReview := review.(*UpdateReviewer)
		fmt.Printf("SetSubClauseWhereCount: %v, ", updateReview.visitor.SetSubClauseWhereCount)
	}
}

func TestUpdateReviewer_Review_NoAllowJoin(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
UPDATE xcs_user_credit_score a1
    ,xcs_user_credit_score a2
SET a1.user_currday_score = a1.user_currday_increment_score + a2.user_currday_score,
    a1.name = 'HH'
WHERE a1.pt_day = '2017-09-20'
    AND a2.pt_day = '2017-09-19'
    AND a1.uid = a2.uid
    `
	fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username, password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Println(reviewMSG.String())

		updateReview := review.(*UpdateReviewer)
		fmt.Printf("SetSubClauseWhereCount: %v, ", updateReview.visitor.SetSubClauseWhereCount)
	}
}

func TestUpdateReviewer_Review_NoAllowSubClause(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
UPDATE xcs_user_credit_score a1
SET a1.user_currday_score = a1.user_currday_increment_score + a2.user_currday_score,
    a1.name = (select 1 from t1 where name = 'HH')
WHERE a1.pt_day = '2017-09-20'
    AND a2.pt_day = '2017-09-19'
    AND a1.uid = a2.uid
    `
	fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username, password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Println(reviewMSG.String())

		updateReview := review.(*UpdateReviewer)
		fmt.Printf("SetSubClauseWhereCount: %v, ", updateReview.visitor.SetSubClauseWhereCount)
	}
}

func TestUpdateReviewer_Review_NoAllowNoWhere(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
UPDATE xcs_user_credit_score a1
SET a1.user_currday_score = a1.user_currday_increment_score + a2.user_currday_score
    `
	fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username, password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Println(reviewMSG.String())

		updateReview := review.(*UpdateReviewer)
		fmt.Printf("SetSubClauseWhereCount: %v, ", updateReview.visitor.SetSubClauseWhereCount)
	}
}

func TestUpdateReviewer_Review_NoAllowLimit(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
UPDATE xcs_user_credit_score a1
SET a1.user_currday_score = a1.user_currday_increment_score + a2.user_currday_score
WHERE id = 1
LIMIT 1
    `
	fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username, password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Println(reviewMSG.String())

		updateReview := review.(*UpdateReviewer)
		fmt.Printf("SetSubClauseWhereCount: %v, ", updateReview.visitor.SetSubClauseWhereCount)
	}
}

func TestUpdateReviewer_Review_AffectRows(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
UpdaTe employees sEt birth_date = '2018-01-01' where emp_no = 10001
    `
	fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username, password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Println(reviewMSG.String())

		updateReview := review.(*UpdateReviewer)
		fmt.Printf("SetSubClauseWhereCount: %v, ", updateReview.visitor.SetSubClauseWhereCount)
	}
}

func TestUpdateReviewer_Review_WhereIn(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
UPDATE t1
SET expired_at = '1531374086000'
WHERE id IN (
    29145, 33493, 33484, 29674, 29516, 33942, 33953, 33947, 38229,
    37757, 36089, 33492, 33483, 33750, 51106, 51105, 51107, 82857,
    59968, 55822, 42854, 42853
);
    `
	fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username, password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Println(reviewMSG.String())

		updateReview := review.(*UpdateReviewer)
		fmt.Printf("SetSubClauseWhereCount: %v, ", updateReview.visitor.SetSubClauseWhereCount)
	}
}

func TestUpdateReviewer_Review_WhereBetween(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"

	sql := `
UPDATE t1
SET expired_at = '1531374086000'
WHERE id between 1 and 10;
    `
	fmt.Sprintf(sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username, password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Println(reviewMSG.String())

		updateReview := review.(*UpdateReviewer)
		fmt.Printf("SetSubClauseWhereCount: %v, ", updateReview.visitor.SetSubClauseWhereCount)
	}
}
