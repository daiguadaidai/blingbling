package reviewer

import (
	"fmt"

	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/config"
)

type UpdateReviewer struct {
	ReviewMSG *ReviewMSG

	StmtNode     *ast.UpdateStmt
	ReviewConfig *config.ReviewConfig
	DBConfig     *config.DBConfig
	visitor      *UpdateVisitor
}

func (this *UpdateReviewer) Init() {
	this.ReviewMSG = NewReivewMSG()

	this.visitor = NewUpdateVisitor()
}

func (this *UpdateReviewer) Review() *ReviewMSG {
	this.Init()
	this.StmtNode.Accept(this.visitor)

	// 通过 visitor检测相关值
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

// 通过 visitor 来检测相关信息
func (this *UpdateReviewer) DetectWithVisitor() (haveError bool) {
	// 检测是否允许使用 join
	haveError = this.DetectUpdateHasJoin()
	if haveError {
		return
	}

	// 检测是否允许使用 sub clause
	haveError = this.DetectUpdateHasSubClause()
	if haveError {
		return
	}

	// 检测是否允许使用 sub clause
	haveError = this.DetectUpdateNoWhere()
	if haveError {
		return
	}

	// 检测是否允许使用 limit
	haveError = this.DetectUpdateLimit()
	if haveError {
		return
	}

	return
}

// 检测删除语句中是否允许Join
func (this *UpdateReviewer) DetectUpdateHasJoin() (haveError bool) {
	// 不允许join, 有多个表 ref 代表有join操作
	if !this.ReviewConfig.RuleAllowUpdateHasJoin && len(this.visitor.RefTables) > 1 {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_UPDATE_HAS_JOIN_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 检测删除语句中是否允许有子句
func (this *UpdateReviewer) DetectUpdateHasSubClause() (haveError bool) {
	// 不允许有子句
	if !this.ReviewConfig.RuleAllowUpdateHasSubClause && this.visitor.HasSubClause {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_UPDATE_HAS_SUB_CLAUSE_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 检测删除语句中是否允许有where子句
func (this *UpdateReviewer) DetectUpdateNoWhere() (haveError bool) {
	// 不允许没有where条件
	if !this.ReviewConfig.RuleAllowUpdateNoWhere && !this.visitor.HasWhereClause {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_UPDATE_NO_WHERE_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 检测删除语句中是否允许有where子句
func (this *UpdateReviewer) DetectUpdateLimit() (haveError bool) {
	// 不允许没有where条件
	if !this.ReviewConfig.RuleAllowUpdateLimit && this.visitor.HasLimitClause {
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_UPDATE_LIMIT_ERROR)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

// 和原表结合进行检测
func (this *UpdateReviewer) DetectFromInstance() (haveError bool) {
	tableInfo := NewTableInfo(this.DBConfig, "")
	err := tableInfo.OpenInstance()
	if err != nil {
		msg := fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在. %v", err)
		this.ReviewMSG.AppendMSG(haveError, msg)
		tableInfo.CloseInstance()
		return
	}

	// 检测update影响的行数
	haveError = this.DetectUpdateRowCount(tableInfo)
	if haveError {
		tableInfo.CloseInstance()
		return
	}

	tableInfo.CloseInstance()
	return
}

func (this *UpdateReviewer) DetectUpdateRowCount(tableInfo *TableInfo) (haveError bool) {
	explainSelectSql, err := GetExplainSelectSqlByUpdateSql(this.StmtNode.Text(),
		this.visitor.SetSubClauseWhereCount, this.visitor.HasWhereClause)
	if err != nil {
		msg := fmt.Sprintf("检测失败. %v", err)
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	updateRowCount, err := tableInfo.GetExplainMaxRows(explainSelectSql)
	if err != nil {
		msg := fmt.Sprintf("检测失败. 执行explain sql获取sql影响行数失败: %v", err)
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}
	if updateRowCount > this.ReviewConfig.RuleUpdateLessThan {
		msg := fmt.Sprintf("检测失败. %v",
			fmt.Sprintf(config.MSG_UPDATE_LESS_THAN_ERROR, this.ReviewConfig.RuleUpdateLessThan))
		haveError = true
		this.ReviewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}
