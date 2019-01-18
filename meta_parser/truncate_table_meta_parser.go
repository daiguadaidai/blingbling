package meta_parser

import (
	"github.com/daiguadaidai/blingbling/ast"
)

type TruncateTableMetaParser struct {
	StmtNode     *ast.TruncateTableStmt
	MI           *MetaInfo
	TuncateTable *MetaTruncateTable
}

func (this *TruncateTableMetaParser) MetaParse() (*MetaInfo, error) {

	truncateTable := &MetaTruncateTable{
		Schema: this.StmtNode.Table.Schema.String(),
		Table:  this.StmtNode.Table.Name.String(),
	}
	this.TuncateTable = truncateTable

	this.MI = new(MetaInfo)
	this.MI.Type = META_INFO_TYPE_TRUNCATE_TABLE
	this.MI.MD = this.TuncateTable

	return this.MI, nil
}
