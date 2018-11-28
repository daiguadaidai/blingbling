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

func TestCreateTableMetaParser_MetaParse1(t *testing.T) {
	sql := "" +
		"create table `ad_merchant_market_list` ( `id` bigint(20) unsigned not null auto_increment comment 'id', `duo_id` bigint(20) unsigned not null default '0' comment 'duoid', `create_at` datetime default null comment '创建时间', `update_at` datetime default null comment '更新时间', `is_deleted` tinyint(2) unsigned not null default '0' comment '0:未删除;1：已删除', primary key (`id`) ) engine=innodb default charset=utf8mb4 comment='商家寻推列表'"

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
