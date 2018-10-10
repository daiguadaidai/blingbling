package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
	"fmt"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dao"
)

type DeleteReviewer struct {
	StmtNode *ast.DeleteStmt
	ReviewConfig *config.ReviewConfig
	DBConfig *config.DBConfig
	Visitor *DeleteVisitor
}

func (this *DeleteReviewer) Init() {
	this.Visitor = NewDeleteVisitor()
}

func (this *DeleteReviewer) Review() *ReviewMSG {
	var reviewMSG *ReviewMSG

	this.Init()
	this.StmtNode.Accept(this.Visitor)

	// 通过 Visitor检测相关值
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

// 通过 Visitor 来检测相关信息
func (this *DeleteReviewer) DetectWithVisitor() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 检测是否有删除多个表
	reviewMSG = this.DetectDeleteManyTable()
	if reviewMSG != nil {
		return reviewMSG
	}

	// 检测是否允许使用 join
	reviewMSG = this.DetectDeleteHasJoin()
	if reviewMSG != nil {
		return reviewMSG
	}

	// 检测是否允许使用 sub clause
	reviewMSG = this.DetectDeleteHasSubClause()
	if reviewMSG != nil {
		return reviewMSG
	}

	// 检测是否允许使用 sub clause
	reviewMSG = this.DetectDeleteNoWhere()
	if reviewMSG != nil {
		return reviewMSG
	}

	// 检测是否允许使用 limit
	reviewMSG = this.DetectDeleteLimit()
	if reviewMSG != nil {
		return reviewMSG
	}

	return reviewMSG
}

func (this *DeleteReviewer) DetectDeleteManyTable() *ReviewMSG {
	var reviewMSG *ReviewMSG

	switch len(this.Visitor.DeleteTables) {
	case 0: // 没有指定删除表的情况 (delete 后面没有写表)
		if !this.ReviewConfig.RuleAllowDeleteManyTable && len(this.Visitor.RefTables) > 1 {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("检测失败. case 0 %v", config.MSG_ALLOW_DELETE_MANY_TABLE_ERROR)
			return reviewMSG
		}
	case 1:
	default:
		if !this.ReviewConfig.RuleAllowDeleteManyTable {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("检测失败. case(default) %v", config.MSG_ALLOW_DELETE_MANY_TABLE_ERROR)
			return reviewMSG
		}
	}

	return reviewMSG
}

// 检测删除语句中是否允许Join
func (this *DeleteReviewer) DetectDeleteHasJoin() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 不允许join, 有多个表 ref 代表有join操作
	if !this.ReviewConfig.RuleAllowDeleteHasJoin && len(this.Visitor.RefTables) > 1 {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_DELETE_HAS_JOIN_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 检测删除语句中是否允许有子句
func (this *DeleteReviewer) DetectDeleteHasSubClause() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 不允许有子句
	if !this.ReviewConfig.RuleAllowDeleteHasSubClause && this.Visitor.HasSubClause {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_DELETE_HAS_SUB_CLAUSE_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 检测删除语句中是否允许有where子句
func (this *DeleteReviewer) DetectDeleteNoWhere() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 不允许没有where条件
	if !this.ReviewConfig.RuleAllowDeleteNoWhere && !this.Visitor.HasWhereClause {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_DELETE_NO_WHERE_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 是否允许使用limit
func (this *DeleteReviewer) DetectDeleteLimit() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 不允许没有where条件
	if !this.ReviewConfig.RuleAllowDeleteLimit && this.Visitor.HasLimitClause {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_DELETE_LIMIT_ERROR)
		return reviewMSG
	}

	return reviewMSG
}

// 和原表结合进行检测
func (this *DeleteReviewer) DetectFromInstance() *ReviewMSG {
	var reviewMSG *ReviewMSG

	tableInfo := dao.NewTableInfo(this.DBConfig, "")
	err := tableInfo.OpenInstance()
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在. %v", err)
		return reviewMSG
	}

	// 检测执行delete影响的行数
	reviewMSG = this.DetectDeleteRowCount(tableInfo)
	if reviewMSG != nil {
		return CloseTableInstance(reviewMSG, tableInfo)
	}

	return CloseTableInstance(reviewMSG, tableInfo)
}

func (this *DeleteReviewer) DetectDeleteRowCount(tableInfo *dao.TableInfo) *ReviewMSG {
	var reviewMSG *ReviewMSG

	explainSelectSql := GetExplainSelectSqlByDeleteSql(this.StmtNode.Text())

	deleteRowCount, err := tableInfo.GetExplainMaxRows(explainSelectSql)
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. 执行explain sql获取sql影响行数失败: %v", err)
		reviewMSG.Code = REVIEW_CODE_WARNING
		return reviewMSG
	}

	if deleteRowCount > this.ReviewConfig.RuleDeleteLessThan {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v",
			fmt.Sprintf(config.MSG_DELETE_LESS_THAN_ERROR, this.ReviewConfig.RuleDeleteLessThan))
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	return  reviewMSG
}
