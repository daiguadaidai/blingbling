package reviewer

import (
	"testing"
	"fmt"
	"github.com/daiguadaidai/blingbling/parser"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/ast"
)

func TestAlterTableReviewer_Review(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "test"
	sql := `
ALTER TABLE test.t3
    ADD COLUMN add_col_1 varchar(50) NOT NULL COMMENT '毕业院校' AFTER a,
    ADD COLUMN add_col_2 varchar(50) NOT NULL COMMENT '毕业院校',
    ADD COLUMN add_col_3 int COMMENT '毕业院校',
    ADD COLUMN update_time int not null default 1 COMMENT '毕业院校',
    DROP COLUMN drop_col_1,
    MODIFY COLUMN modify_col_1 varchar(50) NOT NULL DEFAULT '0' COMMENT '毕业院校',
    CHANGE change_col_1 change_col_2 varchar(50) NOT NULL DEFAULT '0' COMMENT '毕业院校',
    RENAME TO test.tt1,
    ADD INDEX idx_name(name, created_at),
    ADD INDEX idx_update_time(update_time),
    -- ADD INDEX idx_id_create_update_time(id, create_time, update_time),
    -- ADD INDEX idx_id(id),
    DROP INDEX idx_xx,
    -- ADD UNIQUE INDEX udx_name(name, created_at),
    -- ADD PRIMARY KEY(id),
    DROP PRIMARY KEY,
    -- ADD CONSTRAINT fk_id FOREIGN KEY(user_id) REFERENCES tb_user(id),
    -- DROP FOREIGN KEY fk_id,
    RENAME INDEX old_index_name TO new_index_name,
    drop partition p1,
    ADD PARTITION (
        PARTITION p3 VALUES LESS THAN (2000),
        PARTITION p4 VALUES LESS THAN (3000)
    ),
    ADD COLUMN text01 text comment '1',
    ADD COLUMN text07 text comment '1',
    engine=innodb;
    `
    fmt.Sprintf("%v", sql)

	sql1 := `
ALTER TABLE test.t3
add (
    scene_type int(10) unsigned NOT NULL DEFAULT '0' COMMENT '订单来源产品用',
    is_risk tinyint(2) unsigned NOT NULL DEFAULT '0' COMMENT '是否刷单，0：未检测；1：正常；2：刷单'
);
    `
	fmt.Sprintf("%v", sql1)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql1, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username ,password, database)
	reviewConfig := config.NewReviewConfig()
	reviewMSGs := make([]*ReviewMSG, 0, 1)
	for _, stmtNode := range stmtNodes {
		alterTableStmt := stmtNode.(*ast.AlterTableStmt)
		fmt.Printf("Schema: %v, Table: %v\n", alterTableStmt.Table.Schema.String(),
			alterTableStmt.Table.Name.String())

		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		reviewMSGs = append(reviewMSGs, reviewMSG)

		for _, reviewMSG := range reviewMSGs {
			if reviewMSG != nil {
				fmt.Printf("Code: %v, MSG: %v\n", reviewMSG.Code, reviewMSG.MSG)
			}
		}

		for i, spec := range alterTableStmt.Specs {
			switch spec.Tp {
			case ast.AlterTableOption:
				fmt.Printf("--- %v: AlterTableOption ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
			case ast.AlterTableAddColumns:
				fmt.Printf("--- %v: AlterTableAddColumns ---: ", i)
				for _, column := range spec.NewColumns {
					fmt.Printf("Name: %v, Type: %v\n", column.Name, column.Tp.String())
					for _, option := range column.Options {
						switch option.Tp {
						case ast.ColumnOptionNoOption:
							fmt.Println("    ---- ColumnOptionNoOption ----")
						case ast.ColumnOptionPrimaryKey:
							fmt.Println("    ---- ColumnOptionPrimaryKey ----")
						case ast.ColumnOptionNotNull:
							fmt.Println("    ---- ColumnOptionNotNull ----")
						case ast.ColumnOptionAutoIncrement:
							fmt.Println("    ---- ColumnOptionAutoIncrement ----")
						case ast.ColumnOptionDefaultValue:
							fmt.Println("    ---- ColumnOptionDefaultValue ----")
						case ast.ColumnOptionUniqKey:
							fmt.Println("    ---- ColumnOptionUniqKey ----")
						case ast.ColumnOptionNull:
							fmt.Println("    ---- ColumnOptionNull ----")
						case ast.ColumnOptionOnUpdate:
							fmt.Println("    ---- ColumnOptionOnUpdate ----")
						case ast.ColumnOptionFulltext:
							fmt.Println("    ---- ColumnOptionFulltext ----")
						case ast.ColumnOptionComment:
							fmt.Println("    ---- ColumnOptionComment ----")
						case ast.ColumnOptionGenerated:
							fmt.Println("    ---- ColumnOptionGenerated ----")
						case ast.ColumnOptionReference:
							fmt.Println("    ---- ColumnOptionReference ----")
						default:
							fmt.Println("    ---- default ----")
						}

					}

					if spec.Position != nil {
						switch spec.Position.Tp {
						case ast.ColumnPositionNone:
							fmt.Println("    --- ColumnPositionNone ---")
						case ast.ColumnPositionFirst:
							fmt.Println("    --- ColumnPositionFirst ---")
						case ast.ColumnPositionAfter:
							fmt.Println("    --- ColumnPositionAfter ---")
						}
					}
				}
			case ast.AlterTableAddConstraint:
				fmt.Printf("--- %v: AlterTableAddConstraint ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
				fmt.Println(spec.Constraint.Name)
				for _, keyName := range spec.Constraint.Keys {
					fmt.Printf("%v, ", keyName.Column.String())
				}
				fmt.Println("")
				switch spec.Constraint.Tp {
				case ast.ConstraintNoConstraint:
					fmt.Println("    --- ConstraintNoConstraint ---")
				case ast.ConstraintPrimaryKey:
					fmt.Println("    --- ConstraintPrimaryKey ---")
				case ast.ConstraintKey:
					fmt.Println("    --- ConstraintKey ---")
				case ast.ConstraintIndex:
					fmt.Println("    --- ConstraintIndex ---")
				case ast.ConstraintUniq:
					fmt.Println("    --- ConstraintUniq ---")
				case ast.ConstraintUniqKey:
					fmt.Println("    --- ConstraintUniqKey ---")
				case ast.ConstraintUniqIndex:
					fmt.Println("    --- ConstraintUniqIndex ---")
				case ast.ConstraintForeignKey:
					fmt.Println("    --- ConstraintForeignKey ---")
				case ast.ConstraintFulltext:
					fmt.Println("    --- ConstraintFulltext ---")
				}
			case ast.AlterTableDropColumn:
				fmt.Printf("--- %v: AlterTableDropColumn ---: name: %v, comment: %v, column: %v\n",
					i, spec.Name, spec.Comment, spec.OldColumnName.String())
			case ast.AlterTableDropPrimaryKey:
				fmt.Printf("--- %v: AlterTableDropPrimaryKey ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
				fmt.Println("    ----:", )
			case ast.AlterTableDropIndex:
				fmt.Printf("--- %v: AlterTableDropIndex ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
				fmt.Println("    ----:", )
			case ast.AlterTableDropForeignKey:
				fmt.Printf("--- %v: AlterTableDropForeignKey ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
			case ast.AlterTableModifyColumn:
				fmt.Printf("--- %v: AlterTableModifyColumn ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
				for _, column := range spec.NewColumns {
					fmt.Println("    ----:", column.Name.String())
				}
			case ast.AlterTableChangeColumn:
				fmt.Printf("--- %v: AlterTableChangeColumn ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
				fmt.Println(spec.NewColumns, spec.OldColumnName)
				for _, column := range spec.NewColumns{
					fmt.Println("    ----:", column.Name.String())
					// for
				}
			case ast.AlterTableRenameTable:
				fmt.Printf("--- %v: AlterTableRenameTable ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
				fmt.Println(" -> ", spec.NewTable.Name.String())
			case ast.AlterTableAlterColumn:
				fmt.Printf("--- %v: AlterTableAlterColumn ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
			case ast.AlterTableLock:
				fmt.Printf("--- %v: AlterTableLock ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
			case ast.AlterTableAlgorithm:
				fmt.Printf("--- %v: AlterTableAlgorithm ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
			case ast.AlterTableRenameIndex:
				fmt.Printf("--- %v: AlterTableRenameIndex ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
				fmt.Println("    ---:", spec.FromKey.String(), spec.ToKey.String())
			case ast.AlterTableForce:
				fmt.Printf("--- %v: AlterTableForce ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
			case ast.AlterTableAddPartitions:
				fmt.Printf("--- %v: AlterTableAddPartitions ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
				for i, partition := range spec.PartDefinitions {
					fmt.Println("    ----:", i, partition.Name)
				}
			case ast.AlterTableDropPartition:
				fmt.Printf("--- %v: AlterTableDropPartition ---: name: %v, comment: %v\n",
					i, spec.Name, spec.Comment)
			}

		}
	}

}

func TestAlterTableReviewer_Review_Normal(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"
	sql := `
alter table employees add column age1 int not null default 0 comment 'aaaa'
    `
	fmt.Sprintf("%v", sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username ,password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Printf("Code: %v, MSG: %v\n", reviewMSG.Code, reviewMSG.MSG)
	}

}

func TestAlterTableReviewer_Review_DropAddPK(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"
	sql := `
ALTER TABLE emp
    DROP PRIMARY KEY,
    ADD COLUMN id bigint NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'id主键',
    ADD COLUMN id1 bigint NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'id主键';
    `
	fmt.Sprintf("%v", sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username ,password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Printf("Code: %v, MSG: %v\n", reviewMSG.Code, reviewMSG.MSG)
	}

}

func TestAlterTableReviewer_Review_1(t *testing.T) {
	var host string = "10.10.10.21"
	var port int = 3307
	var username string = "HH"
	var password string = "oracle12"
	var database string = "employees"
	sql := "alter table `aaa` add `ttt` int not null default 0 comment '检测'"
	fmt.Sprintf("%v", sql)

	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(sql, "", "")
	if err != nil {
		fmt.Printf("Syntax Error: %v", err)
	}

	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	dbConfig := config.NewDBConfig(host, port, username ,password, database)
	reviewConfig := config.NewReviewConfig()
	for _, stmtNode := range stmtNodes {
		review := NewReviewer(stmtNode, reviewConfig, dbConfig)
		reviewMSG := review.Review()
		fmt.Printf("Code: %v, MSG: %v\n", reviewMSG.Code, reviewMSG.MSG)
	}

}
