package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
	"fmt"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dao"
)

type DeleteReviewer struct {
	ReviewMSG * ReviewMSG

	StmtNode *ast.DeleteStmt
	ReviewConfig *config.ReviewConfig
	DBConfig *config.DBConfig
	Visitor *DeleteVisitor
}

func (this *DeleteReviewer) Init() {
	this.ReviewMSG = NewReivewMSG()

	this.Visitor = NewDeleteVisitor()
}

func (this *DeleteReviewer) Review() *ReviewMSG {
	this.Init()
	this.StmtNode.Accept(this.Visitor)

	// 通过 Visitor检测相关值
	haveError := this.DetectWithVisitor()
	if haveError {
		return this.ReviewMSG
	}

	// 检测需要查询数据库相关的值
	haveError = this.DetectFromInstance()
	if haveError {
		return this.ReviewMSG
	}

	return this.ReviewMSG
}

// 通过 Visitor 来检测相关信息
func (this *DeleteReviewer) DetectWithVisitor() (haveError bool) {
	// 检测是否有删除多个表
	haveError = this.DetectDeleteManyTable()
	if haveError {
		return
	}

	// 检测是否允许使用 join
	haveError = this.DetectDeleteHasJoin()
	if haveError {
		return
	}

	// 检测是否允许使用 sub clause
	haveError = this.DetectDeleteHasSubClause()
	if haveError {
		return
	}

	// 检测是否允许使用 sub clause
	haveError = this.DetectDeleteNoWhere()
	if haveError {
		return
	}

	// 检测是否允许使用 limit
	haveError = this.DetectDeleteLimit()
	if haveError {
		return
	}

	return
}

func (this *DeleteReviewer) DetectDeleteManyTable() (haveError bool) {
	switch len(this.Visitor.DeleteTables) {
	case 0: // 没有指定删除表的情况 (delete 后面没有写表)
		if !this.ReviewConfig.RuleAllowDeleteManyTable && len(this.Visitor.RefTables) > 1 {
			msg := fmt.Sprintf("检测失败. case 0 %v", config.MSG_ALLOW_DELETE_MANY_TABLE_ERROR)
			haveError = true
			this.ReviewMSG.AppendMSG(haveError, msg)
			return
		}
	case 1:
	default:
		if !this.ReviewConfig.RuleAllowDeleteManyTable {
			msg := fmt.Sprintf("检测失败. case(default) %v", config.MSG_ALLOW_DELETE_MANY_TABLE_ERROR)
			haveError = true
			this.ReviewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	return
}

// 检测删除语句中是否允许Join
func (this *DeleteReviewer) DetectDeleteHasJoin() (haveError bool) {
	// 不允许join, 有多个表 ref 代表有join操作
	if !this.ReviewConfig.RuleAllowDeleteHasJoin && len(this.Visitor.RefTables) > 1 {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_DELETE_HAS_JOIN_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 检测删除语句中是否允许有子句
func (this *DeleteReviewer) DetectDeleteHasSubClause() (haveError bool) {
	// 不允许有子句
	if !this.ReviewConfig.RuleAllowDeleteHasSubClause && this.Visitor.HasSubClause {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_DELETE_HAS_SUB_CLAUSE_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 检测删除语句中是否允许有where子句
func (this *DeleteReviewer) DetectDeleteNoWhere() (haveError bool) {
	// 不允许没有where条件
	if !this.ReviewConfig.RuleAllowDeleteNoWhere && !this.Visitor.HasWhereClause {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_DELETE_NO_WHERE_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 是否允许使用limit
func (this *DeleteReviewer) DetectDeleteLimit() (haveError bool) {
	// 不允许没有where条件
	if !this.ReviewConfig.RuleAllowDeleteLimit && this.Visitor.HasLimitClause {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_DELETE_LIMIT_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 和原表结合进行检测
func (this *DeleteReviewer) DetectFromInstance() (haveError bool) {
	tableInfo := dao.NewTableInfo(this.DBConfig, "")
	err := tableInfo.OpenInstance()
	if err != nil {
		msg := fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在. %v", err)
		this.ReviewMSG.AppendMSG(haveError, msg)
		tableInfo.CloseInstance()
		return
	}

	// 检测执行delete影响的行数
	haveError = this.DetectDeleteRowCount(tableInfo)
	if haveError {
		tableInfo.CloseInstance()
		return
	}

	tableInfo.CloseInstance()
	return
}

func (this *DeleteReviewer) DetectDeleteRowCount(tableInfo *dao.TableInfo) (haveError bool) {
	explainSelectSql := GetExplainSelectSqlByDeleteSql(this.StmtNode.Text())

	deleteRowCount, err := tableInfo.GetExplainMaxRows(explainSelectSql)
	if err != nil {
		msg := fmt.Sprintf("检测失败. 执行explain sql获取sql影响行数失败: %v", err)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	if deleteRowCount > this.ReviewConfig.RuleDeleteLessThan {
		msg := fmt.Sprintf("检测失败. %v",
			fmt.Sprintf(config.MSG_DELETE_LESS_THAN_ERROR, this.ReviewConfig.RuleDeleteLessThan))
		this.ReviewMSG.AppendMSG(haveError, msg)
	}

	return
}
