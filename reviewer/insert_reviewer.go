package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
	"fmt"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dao"
)

type InsertReviewer struct {
	StmtNode *ast.InsertStmt
	ReviewConfig *config.ReviewConfig
	DBConfig *config.DBConfig
	SchemaName string
	TableName string
}

func (this *InsertReviewer) Init() {
	tableName := this.StmtNode.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName)
	this.SchemaName = tableName.Schema.String()
	this.TableName = tableName.Name.String()
}

func (this *InsertReviewer) Review() *ReviewMSG {
	var reviewMSG *ReviewMSG

	this.Init()

	// 是否允许不指定字段
	reviewMSG = this.DetectAllowNoColumns()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 是否允许 insert ignore
	reviewMSG = this.DetectAllowInsertIgnore()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}
	// 是否允许 replace into
	reviewMSG = this.DetectAllowInsertReplace()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测是否允许 insert select
	reviewMSG = this.DetectAllowInsertSelect()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测每批insert是否超过指定行数
	reviewMSG = this.DetectInsertRows()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测insert的值的个数是否和字段个数相等
	reviewMSG = this.DetectInsertValueCount()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 需要查询数据库相关的信息
	reviewMSG = this.DetectFromInstance()
	if reviewMSG != nil {
		return reviewMSG
	}

	// 能走到这里说明写的语句审核成功
	reviewMSG = new(ReviewMSG)
	reviewMSG.Code = REVIEW_CODE_SUCCESS
	reviewMSG.MSG = "审核成功"

	return reviewMSG
}

// 检测是否允许不指定字段
func (this *InsertReviewer) DetectAllowNoColumns() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if !this.ReviewConfig.RuleAllowInsertNoColumn && len(this.StmtNode.Columns) == 0 {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_INSERT_NO_COLUMN_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 检测是否允许使用 insert ignore
func (this * InsertReviewer) DetectAllowInsertIgnore() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if !this.ReviewConfig.RuleAllowInsertIgnore && this.StmtNode.IgnoreErr {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_INSERT_IGNORE_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 检测是否允许使用 replace into
func (this * InsertReviewer) DetectAllowInsertReplace() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if !this.ReviewConfig.RuleAllowInsertReplace && this.StmtNode.IsReplace {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_INSERT_REPLIACE_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 检测是否允许 insert select
func (this *InsertReviewer) DetectAllowInsertSelect() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if !this.ReviewConfig.RuleAllowInsertSelect && this.StmtNode.Select != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_INSERT_SELECT_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 检测Insert
func (this *InsertReviewer) DetectInsertRows() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if len(this.StmtNode.Lists) > this.ReviewConfig.RuleInsertRows {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v",
			fmt.Sprintf(config.MSG_INSERT_ROWS_ERROR, this.ReviewConfig.RuleInsertRows))
		return reviewMSG
	}

	return reviewMSG
}

// 检测每个值是否和字段数一样
func (this *InsertReviewer) DetectInsertValueCount() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if columnLen := len(this.StmtNode.Columns); columnLen > 0 { // 有指定字段则直接和字段长度做比较
		for i, list := range this.StmtNode.Lists {
			if len(list) != columnLen {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败. 第%v行值的个数和字段个数不一样", i+1)
				return reviewMSG
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
					reviewMSG = new(ReviewMSG)
					reviewMSG.MSG = fmt.Sprintf("检测失败. 第%v行值的个数和其他行不一样", i+1)
					return reviewMSG
				}

				beforeColumnLen = len(list)
			}
		}
	}

	return reviewMSG
}

// 和原表结合进行检测
func (this *InsertReviewer) DetectFromInstance() *ReviewMSG {
	var reviewMSG *ReviewMSG

	tableInfo := dao.NewTableInfo(this.DBConfig, "")
	err := tableInfo.OpenInstance()
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在. %v", err)
		return reviewMSG
	}

	// 检测数据库不存在报错
	if this.SchemaName != "" {
		reviewMSG = DetectDatabaseNotExistsByName(tableInfo, this.SchemaName)
		if reviewMSG != nil {
			return CloseTableInstance(reviewMSG, tableInfo)
		}
	}

	// 检测表不存在则报错
	reviewMSG = DetectTableNotExistsByName(tableInfo, this.SchemaName, this.TableName)
	if reviewMSG != nil {
		return CloseTableInstance(reviewMSG, tableInfo)
	}

	// 获取原表的建表语句
	err = tableInfo.InitCreateTableSql(this.SchemaName, this.TableName)
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 该Insert语法正确. " +
			"但是无法获取到源表建表sql. %v", err)
		return CloseTableInstance(reviewMSG, tableInfo)
	}
	CloseTableInstance(nil, tableInfo) // 该查询数据库的地方已经完成, 关闭相关链接

	// 检测字段是否存在
	reviewMSG = this.DetectColumnExists(tableInfo)
	if reviewMSG != nil {
		return CloseTableInstance(reviewMSG, tableInfo)
	}

	return CloseTableInstance(reviewMSG, tableInfo)
}

/* 检测字段是否存在
Params:
    _tableInfo: 解析的数据库表信息
 */
func (this *InsertReviewer) DetectColumnExists(_tableInfo *dao.TableInfo) *ReviewMSG {
	var reviewMSG *ReviewMSG

	for _, column := range this.StmtNode.Columns {
		if _, ok := _tableInfo.ColumnNameMap[column.Name.String()]; !ok {
			reviewMSG = new(ReviewMSG)
			reviewMSG.Code = REVIEW_CODE_ERROR
			reviewMSG.MSG = fmt.Sprintf("警告: Insert 指定的字段[%v]不存在", column.Name.String())
			return reviewMSG
		}
	}

	return reviewMSG
}
