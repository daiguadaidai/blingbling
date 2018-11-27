package meta_parser

import (
	"encoding/json"
	"fmt"
	"github.com/daiguadaidai/blingbling/parser"
	"testing"
)

func TestAlterTableMetaParser_MetaParse(t *testing.T) {
	sql := `
ALTER TABLE test.t1
  ADD COLUMN (id bigint(18) NOT NULL AUTO_INCREMENT COMMENT '主键',
              dep varchar(3) NOT NULL DEFAULT '' Comment '注释'),
  Add COLUMN arr varchar(3) NOT NULL DEFAULT '' Comment '注释',
  DROP COLUMN name,
  ADD PRIMARY KEY (id, name),
  ADD INDEX idx_id_name(id, name),
  ADD UNIQUE INDEX idx_id_name(id, name),
  DROP PRIMARY KEY,
  DROP INDEX idx_id_name,
  DROP INDEX idx_id_name,
  MODIFY id bigint(18) NOT NULL AUTO_INCREMENT COMMENT '主键',
  MODIFY dep varchar(3) NOT NULL DEFAULT '' Comment '注释',
  MODIFY arr varchar(3) NOT NULL DEFAULT '' Comment '注释',
  CHANGE id id1 bigint(18) NOT NULL AUTO_INCREMENT COMMENT '主键',
  CHARSET='utf8mb4' ENGINE=innodb,
  COMMENT="表注释",
  RENAME TO test2.t2,
  RENAME INDEX idx1 to idx2
;
    `

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	for _, stmtNode := range stmtNodes {
		metaParser := NewMetaParser(stmtNode)
		mi, err := metaParser.MetaParse()
		if err != nil {
			t.Fatal(err.Error())
		}
		mt := mi.MD.(*MetaAlterTable)
		fmt.Println(mt.Schema, mt.Table)

		for _, op := range mt.OPS {
			switch spec := op.(type) {
			case *MetaTableOptions:
				fmt.Println(spec.Type)
				for _, option := range spec.Options {
					fmt.Println(option.Type, option.Value)
				}
			case *MetaAddColumn:
				fmt.Println(spec.Type)
				for _, col := range spec.Columns {
					fmt.Println(col.Name, col.Type, col.Default, col.AutoIncrement, col.NotNull, col.Comment)
				}
				fmt.Println(spec.After)
			case *MetaDropColumn:
				fmt.Println(spec.Type)
				fmt.Println(spec.ColumnName)
			case *MetaAddConstraint:
				fmt.Println(spec.Type)
				fmt.Println(spec.Constraint.Type)
				fmt.Println(spec.Constraint.Name)
				fmt.Println(spec.Constraint.ColumnNames)
			case *MetaDropConstraint:
				fmt.Println(spec.Type)
				fmt.Println(spec.ConstraintType)
				fmt.Println(spec.Name)
			case *MetaModifyColumn:
				fmt.Println(spec.Type)
				for _, col := range spec.Columns {
					fmt.Println(col.Name, col.Type, col.Default, col.AutoIncrement, col.NotNull, col.Comment)
				}
				fmt.Println(spec.After)
			case *MetaChangeColumn:
				fmt.Println(spec.Type)
				fmt.Println(spec.NewColumn)
				fmt.Println(spec.After)
			case *MetaRenameTable:
				fmt.Println(spec.Type)
				fmt.Println(spec.Schema, spec.Table)
			case *MetaRenameIndex:
				fmt.Println(spec.Type)
				fmt.Println(spec.OldName, spec.NewName)
			}
			fmt.Println("--------------------------------------")
		}

		jsonBytes, err := json.Marshal(mi)
		if err != nil {
			t.Fatal(err.Error())
		}
		fmt.Println(string(jsonBytes))
	}
}
