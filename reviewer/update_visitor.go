package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
	"fmt"
)

type UpdateVisitor struct{
	IsMultiTable bool
	RefTables map[string]*ReviewTable
	WhereSubClauseTables map[string]*ReviewTable
	StmtBlockType int
	HasWhereClause bool
	SubClauseLevel int
	HasSubClause bool
	HasLimitClause bool
	IsInSetClause bool // 当前便利是否在set子句中
	SetSubClauseCount int // set子句的个数
	SetSubClauseWhereCount int // set子句中的where有几个
	SetSubClauseLevelMeetFirstWhere map[int]bool // set字句中子句层级是否有碰到了第一个where
	TableNameAlias string
}

func NewUpdateVisitor() *UpdateVisitor {
	updateVisitor := new(UpdateVisitor)

	updateVisitor.RefTables = make(map[string]*ReviewTable)
	updateVisitor.WhereSubClauseTables = make(map[string]*ReviewTable)
	updateVisitor.SetSubClauseLevelMeetFirstWhere = make(map[int]bool)

	return  updateVisitor
}

func (this *UpdateVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	fmt.Printf("Enter: %T, %[1]v\n", in)

	// 解析和设置正在解析的语句块是哪块
	switch stmt := in.(type) {
	case *ast.UpdateStmt:
		this.StmtBlockType = UPDATE_STMT
		this.IsMultiTable = stmt.MultipleTable
	case *ast.TableRefsClause:
		if this.SubClauseLevel == 0 && this.StmtBlockType == UPDATE_STMT {
			this.StmtBlockType = TABLE_REFS_CLAUSE
		}
	case *ast.BinaryOperationExpr:
		if this.SubClauseLevel == 0 && this.StmtBlockType == ASSIGNMENT && !this.IsInSetClause { // update where语句
			this.HasWhereClause = true
			this.StmtBlockType = BINARY_OPERATON_EXPR
		} else if this.SubClauseLevel > 0 && this.StmtBlockType == ASSIGNMENT &&
			this.IsInSetClause { // 设置set子句中的第一次遇到where
			if _, ok := this.SetSubClauseLevelMeetFirstWhere[this.SetSubClauseCount]; !ok {
				this.SetSubClauseLevelMeetFirstWhere[this.SetSubClauseCount] = true
				this.SetSubClauseWhereCount ++
			}
		}
	case *ast.PatternInExpr:
		if this.SubClauseLevel == 0 && this.StmtBlockType == ASSIGNMENT && !this.IsInSetClause { // update where语句
			this.HasWhereClause = true
			this.StmtBlockType = BINARY_OPERATON_EXPR
		}
	case *ast.BetweenExpr:
		if this.SubClauseLevel == 0 && this.StmtBlockType == ASSIGNMENT && !this.IsInSetClause { // update where语句
			this.HasWhereClause = true
			this.StmtBlockType = BINARY_OPERATON_EXPR
		}
	case *ast.Assignment:
		if this.StmtBlockType != ASSIGNMENT {
			this.StmtBlockType = ASSIGNMENT
		}
		this.IsInSetClause = true
	case *ast.TableSource:
		this.TableNameAlias = stmt.AsName.String()
	case *ast.TableName:
		if stmt != nil {
			reviewTable := &ReviewTable{
				SchemaName: stmt.Schema.String(),
				TableName:  stmt.Name.String(),
				Alias: this.TableNameAlias,
			}
			if this.StmtBlockType == TABLE_REFS_CLAUSE {
				if _, ok := this.RefTables[reviewTable.ToString()]; ok {
					this.RefTables[reviewTable.ToLongString()] = reviewTable
				} else {
					this.RefTables[reviewTable.ToString()] = reviewTable
				}
			} else if this.SubClauseLevel != 0 && this.StmtBlockType == BINARY_OPERATON_EXPR {
				this.WhereSubClauseTables[reviewTable.ToString()] = reviewTable
			}
		}
	case *ast.ColumnNameExpr:
	case *ast.ColumnName:
	case *ast.ValueExpr:
	case *ast.SelectStmt:
		this.SetSubClauseCount ++
		this.SubClauseLevel ++
		if !this.HasSubClause {
			this.HasSubClause = true
		}
	case *ast.Limit:
		if this.SubClauseLevel == 0 && this.StmtBlockType == BINARY_OPERATON_EXPR {
			this.StmtBlockType = LIMIT
			this.HasLimitClause = true
		}
	}

	return in, false
}

func (this *UpdateVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	// fmt.Printf("Leave: %T\n", in)
	switch in.(type) {
	case *ast.SelectStmt:
		this.SubClauseLevel --
	case *ast.Assignment:
		this.IsInSetClause = false
		this.SetSubClauseLevelMeetFirstWhere = make(map[int]bool) // 重新开始计算每个set中的where
		this.SetSubClauseCount = 0 // 重新计算set中子句个数
	case *ast.TableSource:
		this.TableNameAlias = ""
	}
	return in, true
}
