package reviewer

import (
"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/config"
	"fmt"
	"github.com/daiguadaidai/blingbling/dao"
)

type RenameTableReviewer struct {
	StmtNode *ast.RenameTableStmt
	ReviewConfig *config.ReviewConfig
	DBConfig *config.DBConfig
}

func (this *RenameTableReviewer) Review() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 禁止使用 rename
	if !this.ReviewConfig.RuleAllowRenameTable {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_ERROR
		reviewMSG.MSG = config.MSG_ALLOW_RENAME_TABLE_ERROR

		return reviewMSG
	}

	// 允许使用rename
	// 循环一个语句中需要rename的所有表, 如: rename t1 to tt2, t2 to tt2;
	for _, tableToTable := range this.StmtNode.TableToTables{

		// 检测数据库名称长度
		reviewMSG = this.DetectDBNameLength(tableToTable.NewTable.Schema.String())
		if reviewMSG != nil {
			reviewMSG.MSG = fmt.Sprintf("%v %v", "数据库名", reviewMSG.MSG)
			reviewMSG.Code = REVIEW_CODE_ERROR
			return reviewMSG
		}

		// 检测数据库命名规则
		if tableToTable.NewTable.Schema.String() != "" {
			reviewMSG = this.DetectDBNameReg(tableToTable.NewTable.Schema.String())
			if reviewMSG != nil {
				reviewMSG.MSG = fmt.Sprintf("%v %v", "数据库名", reviewMSG.MSG)
				reviewMSG.Code = REVIEW_CODE_ERROR
				return reviewMSG
			}
		}

		// 检测表名称长度
		reviewMSG = this.DetectToTableNameLength(tableToTable.NewTable.Name.String())
		if reviewMSG != nil {
			reviewMSG.MSG = fmt.Sprintf("%v %v", "表名", reviewMSG.MSG)
			reviewMSG.Code = REVIEW_CODE_ERROR
			return reviewMSG
		}

		// 检测表命名规则
		reviewMSG = this.DetectToTableNameReg(tableToTable.NewTable.Name.String())
		if reviewMSG != nil {
			reviewMSG.MSG = fmt.Sprintf("%v %v", "表名", reviewMSG.MSG)
			reviewMSG.Code = REVIEW_CODE_ERROR
			return reviewMSG
		}
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

/* 检测数据库名长度
Params:
_name: 需要检测的名称
*/
func (this *RenameTableReviewer) DetectDBNameLength(_name string) *ReviewMSG {
	return DetectNameLength(_name, this.ReviewConfig.RuleNameLength)
}

/* 检测数据库命名规范
Params:
_name: 需要检测的名称
*/
func (this *RenameTableReviewer) DetectDBNameReg(_name string) *ReviewMSG {
	return DetectNameReg(_name, this.ReviewConfig.RuleNameReg)
}

/* 检测数据库名长度
Params:
    _name: 需要检测的名称
 */
func (this *RenameTableReviewer) DetectToTableNameLength(_name string) *ReviewMSG {
	return DetectNameLength(_name, this.ReviewConfig.RuleNameLength)
}

/* 检测数据库命名规范
Params:
    _name: 需要检测的名称
 */
func (this *RenameTableReviewer) DetectToTableNameReg(_name string) *ReviewMSG {
	var reviewMSG *ReviewMSG

	reviewMSG = DetectNameReg(_name, this.ReviewConfig.RuleTableNameReg)
	if reviewMSG != nil {
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v 表名: %v",
			fmt.Sprintf(config.MSG_TABLE_NAME_GRE_ERROR, this.ReviewConfig.RuleTableNameReg),
			_name)
	}

	return reviewMSG
}

// 链接指定实例检测相关表信息(所有)
func (this *RenameTableReviewer) DetectInstanceTables() *ReviewMSG {
	var reviewMSG *ReviewMSG

	tableInfo := dao.NewTableInfo(this.DBConfig, "")
	err := tableInfo.OpenInstance()
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 无法链接到指定实例. 删除表sql. %v", err)
		return reviewMSG
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

		reviewMSG = this.DetectInstanceTable(tableInfo, oldSchemaName, oldTableName,
			newSchemaName, newTableName)

		if reviewMSG != nil {
			return CloseTableInstance(reviewMSG, tableInfo)
		}
	}

	return CloseTableInstance(reviewMSG, tableInfo)
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
	_tableInfo *dao.TableInfo,
	_OldDBName string,
	_OldTableName string,
	_NewDBName string,
	_NewTableName string,
) *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 老表不存在报错
	reviewMSG = DetectTableNotExistsByName(_tableInfo, _OldDBName, _OldTableName)
	if reviewMSG != nil {
		return reviewMSG
	}

	// 新表已经存在报错
	reviewMSG = DetectTableExistsByName(_tableInfo, _NewDBName, _NewTableName)
	if reviewMSG != nil {
		return reviewMSG
	}

	return reviewMSG
}
