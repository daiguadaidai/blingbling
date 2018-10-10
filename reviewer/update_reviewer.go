package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
	"fmt"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dao"
)

type UpdateReviewer struct {
	StmtNode *ast.UpdateStmt
	ReviewConfig *config.ReviewConfig
	DBConfig *config.DBConfig
	visitor *UpdateVisitor
}

func (this *UpdateReviewer) Init() {
	this.visitor = NewUpdateVisitor()
}

func (this *UpdateReviewer) Review() *ReviewMSG {
	var reviewMSG *ReviewMSG

	this.Init()
	this.StmtNode.Accept(this.visitor)

	// 通过 visitor检测相关值
	reviewMSG = this.DetectWithVisitor()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测需要查询数据库相关的值
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

// 通过 visitor 来检测相关信息
func (this *UpdateReviewer) DetectWithVisitor() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 检测是否允许使用 join
	reviewMSG = this.DetectUpdateHasJoin()
	if reviewMSG != nil {
		return reviewMSG
	}

	// 检测是否允许使用 sub clause
	reviewMSG = this.DetectUpdateHasSubClause()
	if reviewMSG != nil {
		return reviewMSG
	}

	// 检测是否允许使用 sub clause
	reviewMSG = this.DetectUpdateNoWhere()
	if reviewMSG != nil {
		return reviewMSG
	}

	// 检测是否允许使用 limit
	reviewMSG = this.DetectUpdateLimit()
	if reviewMSG != nil {
		return reviewMSG
	}

	return reviewMSG
}

// 检测删除语句中是否允许Join
func (this *UpdateReviewer) DetectUpdateHasJoin() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 不允许join, 有多个表 ref 代表有join操作
	if !this.ReviewConfig.RuleAllowUpdateHasJoin && len(this.visitor.RefTables) > 1 {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_UPDATE_HAS_JOIN_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 检测删除语句中是否允许有子句
func (this *UpdateReviewer) DetectUpdateHasSubClause() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 不允许有子句
	if !this.ReviewConfig.RuleAllowUpdateHasSubClause && this.visitor.HasSubClause {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_UPDATE_HAS_SUB_CLAUSE_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 检测删除语句中是否允许有where子句
func (this *UpdateReviewer) DetectUpdateNoWhere() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 不允许没有where条件
	if !this.ReviewConfig.RuleAllowUpdateNoWhere && !this.visitor.HasWhereClause {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_UPDATE_NO_WHERE_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 检测删除语句中是否允许有where子句
func (this *UpdateReviewer) DetectUpdateLimit() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 不允许没有where条件
	if !this.ReviewConfig.RuleAllowUpdateLimit && this.visitor.HasLimitClause {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_UPDATE_LIMIT_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 和原表结合进行检测
func (this *UpdateReviewer) DetectFromInstance() *ReviewMSG {
	var reviewMSG *ReviewMSG

	tableInfo := dao.NewTableInfo(this.DBConfig, "")
	err := tableInfo.OpenInstance()
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在. %v", err)
		return reviewMSG
	}

	// 检测update影响的行数
	reviewMSG = this.DetectUpdateRowCount(tableInfo)
	if reviewMSG != nil {
		return CloseTableInstance(reviewMSG, tableInfo)
	}

	return CloseTableInstance(reviewMSG, tableInfo)
}

func (this *UpdateReviewer) DetectUpdateRowCount(tableInfo *dao.TableInfo) *ReviewMSG {
	var reviewMSG *ReviewMSG

	explainSelectSql, err := GetExplainSelectSqlByUpdateSql(this.StmtNode.Text(),
		this.visitor.SetSubClauseWhereCount, this.visitor.HasWhereClause)
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", err)
		reviewMSG.Code = REVIEW_CODE_WARNING
		return reviewMSG
	}

	updateRowCount, err := tableInfo.GetExplainMaxRows(explainSelectSql)
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. 执行explain sql获取sql影响行数失败: %v", err)
		reviewMSG.Code = REVIEW_CODE_WARNING
		return reviewMSG
	}
	if updateRowCount > this.ReviewConfig.RuleUpdateLessThan {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v",
			fmt.Sprintf(config.MSG_UPDATE_LESS_THAN_ERROR, this.ReviewConfig.RuleUpdateLessThan))
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	return  reviewMSG
}
