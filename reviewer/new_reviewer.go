package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/config"
)

/* 新建一个需要进行review的实例
Params:
    _stmtNode: 解析完成的语句节点
    _reviewConfig: 审核的配置信息
    _dbConfig: 对哪个数据库进行审核
 */
func NewReviewer(
	_stmtNode ast.Node,
	_reviewConfig *config.ReviewConfig,
	_dbConfig *config.DBConfig,
) Reviewer {
	switch stmt := _stmtNode.(type) {
	case *ast.CreateDatabaseStmt:
		return &CreateDatabaseReviewer{StmtNode: stmt, ReviewConfig: _reviewConfig, DBConfig: _dbConfig}
	case *ast.DropDatabaseStmt:
		return &DropDatabaseReviewer{StmtNode: stmt, ReviewConfig: _reviewConfig, DBConfig: _dbConfig}
	case *ast.CreateTableStmt:
		return &CreateTableReviewer{StmtNode: stmt, ReviewConfig: _reviewConfig, DBConfig: _dbConfig}
	case *ast.DropTableStmt:
		return &DropTableReviewer{StmtNode: stmt, ReviewConfig: _reviewConfig, DBConfig: _dbConfig}
	case *ast.RenameTableStmt:
		return &RenameTableReviewer{StmtNode: stmt, ReviewConfig: _reviewConfig, DBConfig: _dbConfig}
	case *ast.CreateViewStmt:
	case *ast.CreateIndexStmt:
	case *ast.DropIndexStmt:
	case *ast.AlterTableStmt:
		return &AlterTableReviewer{StmtNode: stmt, ReviewConfig: _reviewConfig, DBConfig: _dbConfig}
	case *ast.TruncateTableStmt:
		return &TruncateTableReviewer{StmtNode: stmt, ReviewConfig: _reviewConfig, DBConfig: _dbConfig}
	case *ast.SelectStmt:
	case *ast.UnionStmt:
	case *ast.LoadDataStmt:
	case *ast.InsertStmt:
		return &InsertReviewer{StmtNode: stmt, ReviewConfig: _reviewConfig, DBConfig: _dbConfig}
	case *ast.DeleteStmt:
		return &DeleteReviewer{StmtNode: stmt, ReviewConfig: _reviewConfig, DBConfig: _dbConfig}
	case *ast.UpdateStmt:
		return &UpdateReviewer{StmtNode: stmt, ReviewConfig: _reviewConfig, DBConfig: _dbConfig}
	case *ast.ShowStmt:
	case *ast.TraceStmt:
	case *ast.ExplainStmt:
	case *ast.PrepareStmt:
	case *ast.DeallocateStmt:
	case *ast.ExecuteStmt:
	case *ast.BeginStmt:
	case *ast.BinlogStmt:
	case *ast.CommitStmt:
	case *ast.RollbackStmt:
	case *ast.UseStmt:
	case *ast.FlushStmt:
	case *ast.KillStmt:
	case *ast.SetStmt:
	case *ast.SetPwdStmt:
	case *ast.CreateUserStmt:
	case *ast.AlterUserStmt:
	case *ast.DropUserStmt:
	case *ast.DoStmt:
	case *ast.AdminStmt:
	case *ast.RevokeStmt:
	case *ast.GrantStmt:
	case *ast.AnalyzeTableStmt:
	case *ast.DropStatsStmt:
	case *ast.LoadStatsStmt:
	}

	return nil
}
