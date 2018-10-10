package reviewer

import (
"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dao"
	"fmt"
)

type DropTableReviewer struct {
	StmtNode *ast.DropTableStmt
	ReviewConfig *config.ReviewConfig
	DBConfig *config.DBConfig
}

func (this *DropTableReviewer) Review() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if !this.ReviewConfig.RuleAllowDropTable {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_ERROR
		reviewMSG.MSG = config.MSG_ALLOW_DROP_TABLE_ERROR

		return reviewMSG
	}

	// 链接实例检测表相关信息(所有)
	reviewMSG = this.DetectInstanceTables()
	if reviewMSG != nil {
		return reviewMSG
	}

	reviewMSG = new(ReviewMSG)
	reviewMSG.Code = REVIEW_CODE_SUCCESS
	reviewMSG.MSG = "审核成功"

	return reviewMSG
}

// 链接指定实例检测相关表信息(所有)
func (this *DropTableReviewer) DetectInstanceTables() *ReviewMSG {
	var reviewMSG *ReviewMSG

	for _, tableStmt := range this.StmtNode.Tables {
		var schemaName string
		if tableStmt.Schema.String() != "" {
			schemaName = tableStmt.Schema.String()
		} else {
			schemaName = this.DBConfig.Database
		}
		reviewMSG = this.DetectInstanceTable(schemaName, tableStmt.Name.String())
		if reviewMSG != nil {
			return reviewMSG
		}
	}

	return reviewMSG
}

/* 链接指定实例检测相关表信息
Params:
    _dbName: 数据库名
    _tableName: 原表名
 */
func (this *DropTableReviewer) DetectInstanceTable(_dbName, _tableName string) *ReviewMSG {
	var reviewMSG *ReviewMSG

	tableInfo := dao.NewTableInfo(this.DBConfig, _tableName)
	err := tableInfo.OpenInstance()
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 无法链接到指定实例. 无法删除表[%v]. %v",
			_tableName, err)
		return reviewMSG
	}

	// 检测表是否存在
	reviewMSG = DetectTableNotExistsByName(tableInfo, _dbName, _tableName)
	if reviewMSG != nil {
		return CloseTableInstance(reviewMSG, tableInfo)
	}

	return CloseTableInstance(reviewMSG, tableInfo)
}
