package meta_parser

import (
	"fmt"
	"github.com/daiguadaidai/blingbling/parser"
	"testing"
)

func TestDropTableMetaParser_MetaParse(t *testing.T) {
	sql := `
    DROP TABLE a, b, test.c;
    DROP TABLE IF EXISTS a, b, test.c;
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
		for _, metaDropTable := range mi.MD.(MetaDropTables) {
			fmt.Println(metaDropTable)
		}
		fmt.Println("-----------")
	}
}
