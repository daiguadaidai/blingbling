package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
)

type DeleteVisitor struct {
	IsMultiTable         bool
	RefTables            map[string]*ReviewTable
	DeleteTables         map[string]*ReviewTable
	WhereSubClauseTables map[string]*ReviewTable
	StmtBlockType        int
	HasWhereClause       bool
	SubClauseLevel       int
	HasSubClause         bool
	HasLimitClause       bool
	TableNameAlias       string
}

func NewDeleteVisitor() *DeleteVisitor {
	deleteVisitor := new(DeleteVisitor)

	deleteVisitor.RefTables = make(map[string]*ReviewTable)
	deleteVisitor.DeleteTables = make(map[string]*ReviewTable)
	deleteVisitor.WhereSubClauseTables = make(map[string]*ReviewTable)

	return deleteVisitor
}

func (this *DeleteVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	// fmt.Printf("Enter: %T, %[1]v\n", in)

	// 解析和设置正在解析的语句块是哪块
	switch stmt := in.(type) {
	case *ast.DeleteStmt:
		this.StmtBlockType = DELETE_STMT
		this.IsMultiTable = stmt.IsMultiTable
	case *ast.TableRefsClause:
		if this.SubClauseLevel == 0 && this.StmtBlockType == DELETE_STMT {
			this.StmtBlockType = TABLE_REFS_CLAUSE
		}
	case *ast.BinaryOperationExpr:
		if this.SubClauseLevel == 0 && this.StmtBlockType == DELETE_TABLE_LIST {
			this.HasWhereClause = true
			this.StmtBlockType = BINARY_OPERATON_EXPR
		}
	case *ast.PatternInExpr:
		if this.SubClauseLevel == 0 && this.StmtBlockType == DELETE_TABLE_LIST {
			this.HasWhereClause = true
			this.StmtBlockType = BINARY_OPERATON_EXPR
		}
	case *ast.BetweenExpr:
		if this.SubClauseLevel == 0 && this.StmtBlockType == DELETE_TABLE_LIST {
			this.HasWhereClause = true
			this.StmtBlockType = BINARY_OPERATON_EXPR
		}
	case *ast.DeleteTableList:
		this.StmtBlockType = DELETE_TABLE_LIST
		if stmt != nil {
			for _, table := range stmt.Tables {
				refTable := &ReviewTable{
					SchemaName: table.Schema.String(),
					TableName:  table.Name.String(),
				}
				this.DeleteTables[refTable.ToString()] = refTable
			}
		}
	case *ast.TableSource:
		this.TableNameAlias = stmt.AsName.String()
	case *ast.TableName:
		if stmt != nil {
			reviewTable := &ReviewTable{
				SchemaName: stmt.Schema.String(),
				TableName:  stmt.Name.String(),
				Alias:      this.TableNameAlias,
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
		this.SubClauseLevel ++
		if !this.HasSubClause {
			this.HasSubClause = true
		}
	case *ast.Limit:
		if this.SubClauseLevel == 0 {
			this.StmtBlockType = LIMIT
			this.HasLimitClause = true
		}
	}

	return in, false
}

func (this *DeleteVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	// fmt.Printf("Leave: %T\n", in)
	switch in.(type) {
	case *ast.SelectStmt:
		this.SubClauseLevel --
	case *ast.TableSource:
		this.TableNameAlias = ""
	}
	return in, true
}
