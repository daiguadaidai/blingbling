package meta_parser

import (
	"fmt"
	"github.com/daiguadaidai/blingbling/parser"
	"testing"
)

func TestTruncateTableMetaParser_MetaParse(t *testing.T) {
	sql := `
    TRUNCATE TABLE test.c;
    TRUNCATE TABLE b;
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
		runcateTable, _ := mi.MD.(*MetaTruncateTable)
		fmt.Println(runcateTable)

	}
}
