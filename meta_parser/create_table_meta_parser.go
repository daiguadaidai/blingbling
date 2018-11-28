package meta_parser

import (
	"fmt"
	"github.com/daiguadaidai/blingbling/ast"
)

type CreateTableMetaParser struct {
	StmtNode    *ast.CreateTableStmt
	MI          *MetaInfo
	CreateTable *MetaCreateTable
}

func (this *CreateTableMetaParser) init() {
	this.CreateTable = NewMetaCreateTable()

	this.MI = new(MetaInfo)
	this.MI.Type = META_INFO_TYPE_CREATE_TABLE
	this.MI.MD = this.CreateTable
}

func (this *CreateTableMetaParser) MetaParse() (*MetaInfo, error) {
	this.init()

	if err := this.parse(); err != nil {
		return nil, err
	}

	return this.MI, nil
}

// 开始解析
func (this *CreateTableMetaParser) parse() error {
	if err := this.unsupport(); err != nil {
		return err
	}

	this.addBaseMeta()     // 表的基础信息
	this.addTableOptions() // 表的选项基础信息
	this.addColumns()      // 添加列
	this.addConstraints()  // 添加约束

	return nil
}

// 检测不支持的情况
func (this *CreateTableMetaParser) unsupport() error {
	if this.StmtNode.Select != nil {
		return fmt.Errorf("不支持create select")
	}
	if this.StmtNode.Partition != nil {
		return fmt.Errorf("不支持分区表")
	}
	return nil
}

// 获取基础的信息
func (this *CreateTableMetaParser) addBaseMeta() {
	this.CreateTable.Schema = this.StmtNode.Table.Schema.String()
	this.CreateTable.Table = this.StmtNode.Table.Name.String()
	this.CreateTable.IfNotExists = this.StmtNode.IfNotExists
}

// 获取标级别选项
func (this *CreateTableMetaParser) addTableOptions() {
	for _, option := range this.StmtNode.Options {
		switch option.Tp {
		case ast.TableOptionEngine:
			this.CreateTable.Engine = option.StrValue
		case ast.TableOptionCharset:
			this.CreateTable.Charset = option.StrValue
		case ast.TableOptionCollate:
			this.CreateTable.Collate = option.StrValue
		case ast.TableOptionComment:
			this.CreateTable.Comment = option.StrValue
		case ast.TableOptionAutoIncrement:
			this.CreateTable.AutoIncrement = option.StrValue
		}
	}
}

// 获取列元数据信息
func (this *CreateTableMetaParser) addColumns() {
	for _, col := range this.StmtNode.Cols {
		metaColumn := new(MetaColumn)
		metaColumn.Name = col.Name.Name.String()
		metaColumn.Type = col.Tp.String()
		for _, option := range col.Options {
			switch option.Tp {
			case ast.ColumnOptionNotNull:
				metaColumn.NotNull = true
			case ast.ColumnOptionDefaultValue:
				metaColumn.Default = option.Expr.GetValue()
			case ast.ColumnOptionComment:
				metaColumn.Comment = option.Expr.GetValue()
			case ast.ColumnOptionAutoIncrement:
				metaColumn.AutoIncrement = true
			}
		}
		this.CreateTable.Columns = append(this.CreateTable.Columns, metaColumn)
	}
}

// 添加约束
func (this *CreateTableMetaParser) addConstraints() {
	for _, cons := range this.StmtNode.Constraints {
		metaConstraint := NewMetaConstraint()
		switch cons.Tp {
		case ast.ConstraintPrimaryKey:
			metaConstraint.Type = CONSTRAINT_TYPE_PK
		case ast.ConstraintKey, ast.ConstraintIndex:
			metaConstraint.Type = CONSTRAINT_TYPE_IDX
		case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
			metaConstraint.Type = CONSTRAINT_TYPE_UK
		case ast.ConstraintForeignKey:
			metaConstraint.Type = CONSTRAINT_TYPE_FK
		case ast.ConstraintFulltext:
			metaConstraint.Type = CONSTRAINT_TYPE_FT
		}
		metaConstraint.Name = cons.Name
		for _, key := range cons.Keys {
			metaConstraint.ColumnNames = append(metaConstraint.ColumnNames, key.Column.Name.String())
		}

		this.CreateTable.Constraints = append(this.CreateTable.Constraints, metaConstraint)
	}
}
