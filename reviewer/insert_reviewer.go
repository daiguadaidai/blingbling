package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
	"fmt"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dao"
)

type InsertReviewer struct {
	ReviewMSG *ReviewMSG

	StmtNode *ast.InsertStmt
	ReviewConfig *config.ReviewConfig
	DBConfig *config.DBConfig
	SchemaName string
	TableName string
}

func (this *InsertReviewer) Init() {
	this.ReviewMSG = NewReivewMSG()

	tableName := this.StmtNode.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName)
	this.SchemaName = tableName.Schema.String()
	this.TableName = tableName.Name.String()
}

func (this *InsertReviewer) Review() *ReviewMSG {
	this.Init()

	// 是否允许不指定字段
	haveError := this.DetectAllowNoColumns()
	if haveError {
		return this.ReviewMSG
	}

	// 是否允许 insert ignore
	haveError = this.DetectAllowInsertIgnore()
	if haveError {
		return this.ReviewMSG
	}

	// 是否允许 replace into
	haveError = this.DetectAllowInsertReplace()
	if haveError {
		return this.ReviewMSG
	}

	// 检测是否允许 insert select
	haveError = this.DetectAllowInsertSelect()
	if haveError {
		return this.ReviewMSG
	}

	// 检测每批insert是否超过指定行数
	haveError = this.DetectInsertRows()
	if haveError {
		return this.ReviewMSG
	}

	// 检测insert的值的个数是否和字段个数相等
	haveError = this.DetectInsertValueCount()
	if haveError {
		return this.ReviewMSG
	}

	// 需要查询数据库相关的信息
	haveError = this.DetectFromInstance()
	if haveError {
		return this.ReviewMSG
	}

	return this.ReviewMSG
}

// 检测是否允许不指定字段
func (this *InsertReviewer) DetectAllowNoColumns() (haveError bool) {
	if !this.ReviewConfig.RuleAllowInsertNoColumn && len(this.StmtNode.Columns) == 0 {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_INSERT_NO_COLUMN_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 检测是否允许使用 insert ignore
func (this * InsertReviewer) DetectAllowInsertIgnore() (haveError bool) {
	if !this.ReviewConfig.RuleAllowInsertIgnore && this.StmtNode.IgnoreErr {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_INSERT_IGNORE_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 检测是否允许使用 replace into
func (this * InsertReviewer) DetectAllowInsertReplace() (haveError bool) {
	if !this.ReviewConfig.RuleAllowInsertReplace && this.StmtNode.IsReplace {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_INSERT_REPLIACE_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 检测是否允许 insert select
func (this *InsertReviewer) DetectAllowInsertSelect() (haveError bool) {
	if !this.ReviewConfig.RuleAllowInsertSelect && this.StmtNode.Select != nil {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_INSERT_SELECT_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 检测Insert
func (this *InsertReviewer) DetectInsertRows() (haveError bool) {
	if len(this.StmtNode.Lists) > this.ReviewConfig.RuleInsertRows {
		msg := fmt.Sprintf("检测失败. %v",
			fmt.Sprintf(config.MSG_INSERT_ROWS_ERROR, this.ReviewConfig.RuleInsertRows))
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 检测每个值是否和字段数一样
func (this *InsertReviewer) DetectInsertValueCount() (haveError bool) {
	if columnLen := len(this.StmtNode.Columns); columnLen > 0 { // 有指定字段则直接和字段长度做比较
		for i, list := range this.StmtNode.Lists {
			if len(list) != columnLen {
				msg := fmt.Sprintf("检测失败. 第%v行值的个数和字段个数不一样", i+1)
				haveError = true
				this.ReviewMSG.AppendMSG(haveError, msg)
				return
			}
		}
	} else { // 没有指定字段名就逐个进行计算了
		isFrist := true // 定义是否是第一行
		beforeColumnLen := 0
		for i, list := range this.StmtNode.Lists {
			if isFrist {
				isFrist = false
				beforeColumnLen = len(list)
			} else {
				if beforeColumnLen != len(list) {
					msg := fmt.Sprintf("检测失败. 第%v行值的个数和其他行不一样", i+1)
					haveError = true
					this.ReviewMSG.AppendMSG(haveError, msg)
					return
				}

				beforeColumnLen = len(list)
			}
		}
	}

	return
}

// 和原表结合进行检测
func (this *InsertReviewer) DetectFromInstance() (haveError bool) {
	var msg string

	tableInfo := dao.NewTableInfo(this.DBConfig, "")
	err := tableInfo.OpenInstance()
	if err != nil {
		msg = fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在. %v", err)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		tableInfo.CloseInstance()
		return
	}

	// 检测数据库不存在报错
	if this.SchemaName != "" {
		haveError, msg = DetectDatabaseNotExistsByName(tableInfo, this.SchemaName)
		haveMSG := this.ReviewMSG.AppendMSG(haveError, msg)
		if haveError || haveMSG {
			tableInfo.CloseInstance()
			return
		}
	}

	// 检测表不存在则报错
	haveError, msg = DetectTableNotExistsByName(tableInfo, this.SchemaName, this.TableName)
	haveMSG := this.ReviewMSG.AppendMSG(haveError, msg)
	if haveError || haveMSG {
		tableInfo.CloseInstance()
		return
	}

	// 获取原表的建表语句
	err = tableInfo.InitCreateTableSql(this.SchemaName, this.TableName)
	if err != nil {
		msg := fmt.Sprintf("警告: 该Insert语法正确. " +
			"但是无法获取到源表建表sql. %v", err)
		this.ReviewMSG.AppendMSG(false, msg)
		tableInfo.CloseInstance()
		return
	}
	tableInfo.CloseInstance() // 该查询数据库的地方已经完成, 关闭相关链接

	err = tableInfo.ParseCreateTableInfo()
	if err != nil {
		msg := fmt.Sprintf("警告: 该Insert语法正确. 解析原表建表语句错误. 无法对比" +
			"%v", err)
		this.ReviewMSG.AppendMSG(false, msg)
		tableInfo.CloseInstance()
		return
	}

	// 检测字段是否存在
	haveError = this.DetectColumnExists(tableInfo)
	if haveError {
		return
	}

	return
}

/* 检测字段是否存在
Params:
    _tableInfo: 解析的数据库表信息
 */
func (this *InsertReviewer) DetectColumnExists(_tableInfo *dao.TableInfo) (haveError bool) {
	for _, column := range this.StmtNode.Columns {
		if _, ok := _tableInfo.ColumnNameMap[column.Name.String()]; !ok {
			msg := fmt.Sprintf("警告: Insert 指定的字段 %v 不存在", column.Name.String())
			haveError = true
			this.ReviewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	return
}
