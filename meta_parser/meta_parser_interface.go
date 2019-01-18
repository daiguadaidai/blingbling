package meta_parser

import "github.com/daiguadaidai/blingbling/ast"

type MetaParser interface {
	MetaParse() (*MetaInfo, error)
}

func NewMetaParser(stmtNode ast.Node) MetaParser {
	switch stmt := stmtNode.(type) {
	case *ast.CreateTableStmt:
		return &CreateTableMetaParser{StmtNode: stmt}
	case *ast.AlterTableStmt:
		return &AlterTableMetaParser{StmtNode: stmt}
	case *ast.DropTableStmt:
		return &DropTableMetaParser{StmtNode: stmt}
	case *ast.TruncateTableStmt:
		return &TruncateTableMetaParser{StmtNode: stmt}
	}

	return nil
}
