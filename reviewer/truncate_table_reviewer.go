package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dao"
	"fmt"
)

type TruncateTableReviewer struct {
	ReviewMSG *ReviewMSG

	StmtNode *ast.TruncateTableStmt
	ReviewConfig *config.ReviewConfig
	DBConfig *config.DBConfig

	SchemaName string
}

func (this *TruncateTableReviewer) Init() {
	this.ReviewMSG = NewReivewMSG()

	if this.StmtNode.Table.Schema.String() != "" {
		this.SchemaName = this.StmtNode.Table.Schema.String()
	} else {
		this.SchemaName = this.DBConfig.Database
	}

}

func (this *TruncateTableReviewer) Review() *ReviewMSG {
	this.Init()

	if !this.ReviewConfig.RuleAllowTruncateTable {
		msg := config.MSG_ALLOW_TRUNCATE_TABLE_ERROR
		this.ReviewMSG.AppendMSG(true, msg)

		return this.ReviewMSG
	}

	// 链接数据库检测实例相关信息
	haveError := this.DetectInstanceTable()
	if haveError {
		return this.ReviewMSG
	}

	return this.ReviewMSG
}

// 链接到实例检测相关信息
func (this *TruncateTableReviewer) DetectInstanceTable() (haveError bool) {
	var msg string

	tableInfo := dao.NewTableInfo(this.DBConfig, this.StmtNode.Table.Name.String())
	err := tableInfo.OpenInstance()
	if err != nil {
		msg = fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在.")
		this.ReviewMSG.AppendMSG(haveError, msg)
		tableInfo.CloseInstance()
		return
	}

	// 表不存在报错
	haveError, msg = DetectTableNotExistsByName(tableInfo, this.SchemaName, this.StmtNode.Table.Name.String())
	haveMSG := this.ReviewMSG.AppendMSG(haveError, msg)
	if haveError || haveMSG {
		tableInfo.CloseInstance()
		return
	}

	tableInfo.CloseInstance()
	return
}
