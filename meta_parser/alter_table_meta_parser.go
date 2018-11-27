package meta_parser

import (
	"github.com/daiguadaidai/blingbling/ast"
)

type AlterTableMetaParser struct {
	StmtNode   *ast.AlterTableStmt
	MI         *MetaInfo
	AlterTable *MetaAlterTable
}

func (this *AlterTableMetaParser) init() {
	this.AlterTable = NewMetaAlterTable()

	this.MI = new(MetaInfo)
	this.MI.Type = META_INFO_TYPE_ALTER_TABLE
	this.MI.MD = this.AlterTable
}

func (this *AlterTableMetaParser) MetaParse() (*MetaInfo, error) {
	this.init()

	this.parse()

	return this.MI, nil
}

func (this *AlterTableMetaParser) parse() {
	this.addBaseMeta()

	for _, spec := range this.StmtNode.Specs {
		switch spec.Tp {
		case ast.AlterTableOption:
			this.addTableOptions(spec)
		case ast.AlterTableAddColumns:
			this.addColumns(spec)
		case ast.AlterTableDropColumn:
			this.addDropColumns(spec)
		case ast.AlterTableAddConstraint:
			this.addConstraint(spec)
		case ast.AlterTableDropPrimaryKey:
			this.addDropConstraint(spec, CONSTRAINT_TYPE_PK)
		case ast.AlterTableDropIndex:
			this.addDropConstraint(spec, CONSTRAINT_TYPE_IDX)
		case ast.AlterTableDropForeignKey:
			this.addDropConstraint(spec, CONSTRAINT_TYPE_FK)
		case ast.AlterTableModifyColumn:
			this.addModifyColumn(spec)
		case ast.AlterTableChangeColumn:
			this.addChangeColumn(spec)
		case ast.AlterTableRenameTable:
			this.addRenameTable(spec)
		case ast.AlterTableRenameIndex:
			this.addRenameIndex(spec)
		}
	}
}

func (this *AlterTableMetaParser) getAfter(spec *ast.AlterTableSpec) string {
	if spec.Position != nil {
		switch spec.Position.Tp {
		case ast.ColumnPositionAfter:
			return spec.Position.RelativeColumn.String()
		}
	}

	return ""
}

// 获取列
func (this *AlterTableMetaParser) getColumns(spec *ast.AlterTableSpec) []*MetaColumn {
	mcs := make([]*MetaColumn, 0, 1)
	for _, col := range spec.NewColumns {
		mc := new(MetaColumn)
		mc.Name = col.Name.String()
		mc.Type = col.Tp.String()
		for _, option := range col.Options {
			switch option.Tp {
			case ast.ColumnOptionNotNull:
				mc.NotNull = true
			case ast.ColumnOptionDefaultValue:
				mc.Default = option.Expr.GetValue().(string)
			case ast.ColumnOptionComment:
				mc.Comment = option.Expr.GetValue().(string)
			case ast.ColumnOptionAutoIncrement:
				mc.AutoIncrement = true
			}
		}
		mcs = append(mcs, mc)
	}

	return mcs
}

// 获取基础的信息
func (this *AlterTableMetaParser) addBaseMeta() {
	this.AlterTable.Schema = this.StmtNode.Table.Schema.String()
	this.AlterTable.Table = this.StmtNode.Table.Name.String()
}

// 添加修改表级别的选项
func (this *AlterTableMetaParser) addTableOptions(spec *ast.AlterTableSpec) {
	mos := NewMetaTableOptions()
	for _, option := range spec.Options {
		mo := new(MetaTableOption)
		mo.Value = option.StrValue
		switch option.Tp {
		case ast.TableOptionEngine:
			mo.Type = TABLE_OPTION_ENGINE
		case ast.TableOptionCharset:
			mo.Type = TABLE_OPTION_CHARSET
		case ast.TableOptionCollate:
			mo.Type = TABLE_OPTION_COLLATE
		case ast.TableOptionComment:
			mo.Type = TABLE_OPTION_COMMENT
		case ast.TableOptionAutoIncrement:
			mo.Type = TABLE_OPTION_AUTO_INCREMENT
		}
		mos.Options = append(mos.Options, mo)
	}
	this.AlterTable.OPS = append(this.AlterTable.OPS, mos)
}

// 添加列
func (this *AlterTableMetaParser) addColumns(spec *ast.AlterTableSpec) {
	mac := NewMetaAddColumn()
	mac.After = this.getAfter(spec)
	mac.Columns = this.getColumns(spec)
	this.AlterTable.OPS = append(this.AlterTable.OPS, mac)
}

// 添加 删除的列
func (this *AlterTableMetaParser) addDropColumns(spec *ast.AlterTableSpec) {
	mdc := NewMetaDropColumn()
	mdc.ColumnName = spec.OldColumnName.String()
	this.AlterTable.OPS = append(this.AlterTable.OPS, mdc)
}

// 添加约束
func (this *AlterTableMetaParser) addConstraint(spec *ast.AlterTableSpec) {
	mac := NewMetaAddConstraint()
	mac.Constraint.Name = spec.Constraint.Name
	switch spec.Constraint.Tp {
	case ast.ConstraintPrimaryKey:
		mac.Constraint.Type = CONSTRAINT_TYPE_PK
	case ast.ConstraintKey, ast.ConstraintIndex:
		mac.Constraint.Type = CONSTRAINT_TYPE_IDX
	case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
		mac.Constraint.Type = CONSTRAINT_TYPE_UK
	case ast.ConstraintForeignKey:
		mac.Constraint.Type = CONSTRAINT_TYPE_FK
	case ast.ConstraintFulltext:
		mac.Constraint.Type = CONSTRAINT_TYPE_FT
	}

	for _, key := range spec.Constraint.Keys {
		mac.Constraint.ColumnNames = append(mac.Constraint.ColumnNames, key.Column.String())
	}
	this.AlterTable.OPS = append(this.AlterTable.OPS, mac)
}

// 添加删除的约束
func (this *AlterTableMetaParser) addDropConstraint(spec *ast.AlterTableSpec, constraintType string) {
	mdc := new(MetaDropConstraint)
	mdc.Name = spec.Name
	mdc.ConstraintType = constraintType
	mdc.Type = ALTER_TABLE_TYPE_DROP_CONSTRAINT
	this.AlterTable.OPS = append(this.AlterTable.OPS, mdc)
}

// 添加modify column
func (this *AlterTableMetaParser) addModifyColumn(spec *ast.AlterTableSpec) {
	mmc := NewMetaModifyColumn()
	mmc.After = this.getAfter(spec)
	mmc.Columns = this.getColumns(spec)
	this.AlterTable.OPS = append(this.AlterTable.OPS, mmc)
}

// 添加changecolumn
func (this *AlterTableMetaParser) addChangeColumn(spec *ast.AlterTableSpec) {
	mcc := NewMetaChangeColumn()
	mcc.OldName = spec.OldColumnName.Name.String()
	mcc.After = this.getAfter(spec)
	columns := this.getColumns(spec)
	if len(columns) > 0 {
		mcc.NewColumn = columns[0]
	}
	this.AlterTable.OPS = append(this.AlterTable.OPS, mcc)
}

// 添加rename table信息
func (this *AlterTableMetaParser) addRenameTable(spec *ast.AlterTableSpec) {
	mrt := NewMetaRenameTable()
	mrt.Schema = spec.NewTable.Schema.String()
	mrt.Table = spec.NewTable.Name.String()
	this.AlterTable.OPS = append(this.AlterTable.OPS, mrt)
}

// 添加Rename index
func (this *AlterTableMetaParser) addRenameIndex(spec *ast.AlterTableSpec) {
	mri := NewMetaRenameIndex()
	mri.OldName = spec.FromKey.String()
	mri.NewName = spec.ToKey.String()
	this.AlterTable.OPS = append(this.AlterTable.OPS, mri)
}
