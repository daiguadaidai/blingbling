package meta_parser

import (
	"github.com/daiguadaidai/blingbling/ast"
)

type DropTableMetaParser struct {
	StmtNode   *ast.DropTableStmt
	MI         *MetaInfo
	DropTables MetaDropTables
}

func (this *DropTableMetaParser) MetaParse() (*MetaInfo, error) {

	dropTables := make([]*MetaDropTable, len(this.StmtNode.Tables))
	for i, table := range this.StmtNode.Tables {
		dropTable := new(MetaDropTable)
		dropTable.IfExists = this.StmtNode.IfExists
		dropTable.Schema = table.Schema.String()
		dropTable.Table = table.Name.String()
		dropTables[i] = dropTable
	}
	this.DropTables = MetaDropTables(dropTables)

	this.MI = new(MetaInfo)
	this.MI.Type = META_INFO_TYPE_DROP_TABLE
	this.MI.MD = this.DropTables

	return this.MI, nil
}
