package reviewer

import (
	"fmt"

	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/config"
)

type RenameTableReviewer struct {
	ReviewMSG *ReviewMSG

	StmtNode     *ast.RenameTableStmt
	ReviewConfig *config.ReviewConfig
	DBConfig     *config.DBConfig
}

func (this *RenameTableReviewer) Init() {
	this.ReviewMSG = NewReivewMSG(config.StmtTypeRenameTable, "", "")
}

func (this *RenameTableReviewer) Review() *ReviewMSG {
	this.Init()

	// 禁止使用 rename
	if !this.ReviewConfig.RuleAllowRenameTable {
		this.ReviewMSG.AppendMSG(true, config.MSG_ALLOW_RENAME_TABLE_ERROR)
		return this.ReviewMSG
	}

	// 允许使用rename
	// 循环一个语句中需要rename的所有表, 如: rename t1 to tt2, t2 to tt2;
	for _, tableToTable := range this.StmtNode.TableToTables {

		// 检测数据库名称长度
		haveError := this.DetectDBNameLength(tableToTable.NewTable.Schema.String())
		if haveError {
			return this.ReviewMSG
		}

		// 检测数据库命名规则
		if tableToTable.NewTable.Schema.String() != "" {
			haveError = this.DetectDBNameReg(tableToTable.NewTable.Schema.String())
			if haveError {
				return this.ReviewMSG
			}
		}

		// 检测表名称长度
		haveError = this.DetectToTableNameLength(tableToTable.NewTable.Name.String())
		if haveError {
			return this.ReviewMSG
		}

		// 检测表命名规则
		haveError = this.DetectToTableNameReg(tableToTable.NewTable.Name.String())
		if haveError {
			return this.ReviewMSG
		}
	}

	// 链接实例检测表相关信息(所有)
	haveError := this.DetectInstanceTables()
	if haveError {
		return this.ReviewMSG
	}

	return this.ReviewMSG
}

/* 检测数据库名长度
Params:
_name: 需要检测的名称
*/
func (this *RenameTableReviewer) DetectDBNameLength(_name string) (haveError bool) {
	var msg string
	haveError, msg = DetectNameLength(_name, this.ReviewConfig.RuleNameLength)
	if haveError {
		msg = fmt.Sprintf("%v %v", "数据库名", msg)
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}
	return
}

/* 检测数据库命名规范
Params:
_name: 需要检测的名称
*/
func (this *RenameTableReviewer) DetectDBNameReg(_name string) (haveError bool) {
	var msg string
	haveError, msg = DetectNameReg(_name, this.ReviewConfig.RuleNameReg)
	if haveError {
		msg = fmt.Sprintf("%v %v", "数据库名", msg)
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}
	return
}

/* 检测数据库名长度
Params:
    _name: 需要检测的名称
*/
func (this *RenameTableReviewer) DetectToTableNameLength(_name string) (haveError bool) {
	var msg string
	haveError, msg = DetectNameLength(_name, this.ReviewConfig.RuleNameLength)
	if haveError {
		msg = fmt.Sprintf("%v %v", "数据库名", msg)
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}
	return
}

/* 检测数据库命名规范
Params:
    _name: 需要检测的名称
*/
func (this *RenameTableReviewer) DetectToTableNameReg(_name string) (haveError bool) {
	var msg string
	haveError, msg = DetectNameReg(_name, this.ReviewConfig.RuleTableNameReg)
	if haveError {
		msg = fmt.Sprintf("检测失败. %v 表名: %v",
			fmt.Sprintf(config.MSG_TABLE_NAME_GRE_ERROR, this.ReviewConfig.RuleTableNameReg),
			_name)
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 链接指定实例检测相关表信息(所有)
func (this *RenameTableReviewer) DetectInstanceTables() (haveError bool) {
	tableInfo := NewTableInfo(this.DBConfig, "")
	err := tableInfo.OpenInstance()
	if err != nil {
		msg := fmt.Sprintf("警告: 无法链接到指定实例. 删除表sql. %v", err)
		this.ReviewMSG.AppendMSG(false, msg)
		tableInfo.CloseInstance()
		return
	}

	for _, tableStmt := range this.StmtNode.TableToTables {
		var oldSchemaName string
		var newSchemaName string
		var oldTableName string
		var newTableName string

		if tableStmt.OldTable.Schema.String() != "" {
			oldSchemaName = tableStmt.OldTable.Schema.String()
		} else {
			oldSchemaName = this.DBConfig.Database
		}

		if tableStmt.NewTable.Schema.String() != "" {
			newSchemaName = tableStmt.NewTable.Schema.String()
		} else {
			newSchemaName = this.DBConfig.Database
		}

		haveError = this.DetectInstanceTable(tableInfo, oldSchemaName, oldTableName,
			newSchemaName, newTableName)
		if haveError {
			tableInfo.CloseInstance()
			return
		}
	}

	tableInfo.CloseInstance()
	return
}

/* 链接指定实例检测相关表信息
Params:
    _tableInfo: 表信息
	_OldDBName: 库名
	_OldTableName 原表名
	_NewDBName 新库名
	_NewTableName 新表名
*/
func (this *RenameTableReviewer) DetectInstanceTable(
	_tableInfo *TableInfo,
	_OldDBName string,
	_OldTableName string,
	_NewDBName string,
	_NewTableName string,
) (haveError bool) {
	// 老表不存在报错
	var msg string
	haveError, msg = DetectTableNotExistsByName(_tableInfo, _OldDBName, _OldTableName)
	haveMSG := this.ReviewMSG.AppendMSG(haveError, msg)
	if haveError || haveMSG {
		return
	}

	// 新表已经存在报错
	haveError, msg = DetectTableExistsByName(_tableInfo, _NewDBName, _NewTableName)
	haveMSG = this.ReviewMSG.AppendMSG(haveError, msg)
	if haveError || haveMSG {
		return
	}

	return
}
