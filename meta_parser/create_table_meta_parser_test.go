package meta_parser

import (
	"encoding/json"
	"fmt"
	"github.com/daiguadaidai/blingbling/parser"
	"testing"
)

func TestCreateTableMetaParser_MetaParse(t *testing.T) {
	sql := `
CREATE TABLE test.t1 (
  id bigint(18) NOT NULL AUTO_INCREMENT COMMENT '主键',
  dep varchar(3) NOT NULL DEFAULT '' Comment '注释',
  arr varchar(3) NOT NULL DEFAULT '' Comment '注释',
  flightNo varchar(10) NOT NULL DEFAULT '' Comment '注释',
  flightDate date NOT NULL DEFAULT '1000-10-10' Comment '注释',
  flightTime varchar(20) NOT NULL DEFAULT '' Comment '注释',
  isCodeShare tinyint(1) Comment '注释',
  tax int(11) NOT NULL DEFAULT '0' Comment '注释',
  yq int(11) NOT NULL DEFAULT '0' Comment '注释',
  cabin char(2) NOT NULL default '' Comment '注释',
  ibe_price int(11) NOT NULL DEFAULT '0' Comment '注释',
  ctrip_price int(11) NOT NULL DEFAULT '0' Comment '注释',
  official_price int(11) NOT NULL DEFAULT '0' Comment '注释',
  uptime datetime NOT NULL DEFAULT '1000-10-10 10:10:10' Comment '注释',
  PRIMARY KEY (id),
  UNIQUE KEY udx_uid (dep, arr, flightNo, flightDate, cabin),
  Index idx_uptime (uptime),
  KEY idx_flight (dep,arr),
  KEY idx_flightdate (flightDate)
) ENGINE=InnoDb  DEFAULT CHARSET=utF8 COLLATE=Utf8mb4_general_ci comment="你号";
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
		mt := mi.MD.(*MetaCreateTable)
		fmt.Println(mt.Schema, mt.Table, mt.Comment, mt.AutoIncrement, mt.Engine, mt.Charset, mt.Collate)

		for _, col := range mt.Columns {
			fmt.Println(col)
		}

		for _, cons := range mt.Constraints {
			fmt.Println(cons)
		}

		jsonBytes, err := json.Marshal(mi)
		if err != nil {
			t.Fatal(err.Error())
		}
		fmt.Println(string(jsonBytes))
	}
}
