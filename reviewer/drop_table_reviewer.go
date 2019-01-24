package reviewer

import (
	"fmt"

	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/config"
)

type DropTableReviewer struct {
	ReviewMSG *ReviewMSG

	StmtNode     *ast.DropTableStmt
	ReviewConfig *config.ReviewConfig
	DBConfig     *config.DBConfig
}

func (this *DropTableReviewer) Init() {
	this.ReviewMSG = NewReivewMSG(config.StmtTypeDropTable, "", "")
}

func (this *DropTableReviewer) Review() *ReviewMSG {
	this.Init()

	if !this.ReviewConfig.RuleAllowDropTable {
		this.ReviewMSG.AppendMSG(true, config.MSG_ALLOW_DROP_TABLE_ERROR)

		return this.ReviewMSG
	}

	// 链接实例检测表相关信息(所有)
	haveError := this.DetectInstanceTables()
	if haveError {
		return this.ReviewMSG
	}

	return this.ReviewMSG
}

// 链接指定实例检测相关表信息(所有)
func (this *DropTableReviewer) DetectInstanceTables() (haveError bool) {
	for _, tableStmt := range this.StmtNode.Tables {
		var schemaName string
		if tableStmt.Schema.String() != "" {
			schemaName = tableStmt.Schema.String()
		} else {
			schemaName = this.DBConfig.Database
		}
		haveError = this.DetectInstanceTable(schemaName, tableStmt.Name.String())
		if haveError {
			return
		}
	}

	return
}

/* 链接指定实例检测相关表信息
Params:
    _dbName: 数据库名
    _tableName: 原表名
*/
func (this *DropTableReviewer) DetectInstanceTable(_dbName, _tableName string) (haveError bool) {
	var msg string

	tableInfo := NewTableInfo(this.DBConfig, _tableName)
	err := tableInfo.OpenInstance()
	if err != nil {
		msg = fmt.Sprintf("警告: 无法链接到指定实例. 无法删除表 %v. %v",
			_tableName, err)
		this.ReviewMSG.AppendMSG(false, msg)
		return
	}

	// 检测表是否存在
	haveError, msg = DetectTableNotExistsByName(tableInfo, _dbName, _tableName)
	haveMSG := this.ReviewMSG.AppendMSG(haveError, msg)
	if haveError || haveMSG {
		tableInfo.CloseInstance()
		return
	}

	tableInfo.CloseInstance()
	return
}
