package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dao"
	"fmt"
)

type TruncateTableReviewer struct {
	StmtNode *ast.TruncateTableStmt
	ReviewConfig *config.ReviewConfig
	DBConfig *config.DBConfig

	SchemaName string
}

func (this *TruncateTableReviewer) Init() {
	if this.StmtNode.Table.Schema.String() != "" {
		this.SchemaName = this.StmtNode.Table.Schema.String()
	} else {
		this.SchemaName = this.DBConfig.Database
	}

}

func (this *TruncateTableReviewer) Review() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if !this.ReviewConfig.RuleAllowTruncateTable {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_ERROR
		reviewMSG.MSG = config.MSG_ALLOW_TRUNCATE_TABLE_ERROR

		return reviewMSG
	}

	// 链接数据库检测实例相关信息
	reviewMSG = this.DetectInstanceTable()
	if reviewMSG != nil {
		return reviewMSG
	}

	reviewMSG = new(ReviewMSG)
	reviewMSG.Code = REVIEW_CODE_SUCCESS
	reviewMSG.MSG = "审核成功"

	return reviewMSG
}

// 链接到实例检测相关信息
func (this *TruncateTableReviewer) DetectInstanceTable() *ReviewMSG {
	var reviewMSG *ReviewMSG

	tableInfo := dao.NewTableInfo(this.DBConfig, this.StmtNode.Table.Name.String())
	err := tableInfo.OpenInstance()
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在.")
		return reviewMSG
	}

	// 表不存在报错
	reviewMSG = DetectTableNotExistsByName(tableInfo, this.SchemaName, this.StmtNode.Table.Name.String())
	if reviewMSG != nil {
		return CloseTableInstance(reviewMSG, tableInfo)
	}

	return CloseTableInstance(reviewMSG, tableInfo)
}
