package meta_parser

const (
	META_INFO_TYPE_CREATE_TABLE   = "ct"
	META_INFO_TYPE_ALTER_TABLE    = "at"
	META_INFO_TYPE_DROP_TABLE     = "dt"
	META_INFO_TYPE_TRUNCATE_TABLE = "tt"
)

type MetaData interface {
	MetaDataInterfaceFunc()
}

type MetaInfo struct {
	Type string   `json:"type"`
	MD   MetaData `json:"md"`
}

type MetaColumn struct {
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	Default       interface{} `json:"default"`
	Comment       interface{} `json:"comment"`
	NotNull       bool        `json:"not_null"`
	AutoIncrement bool        `json:"auto_increment"`
}

const (
	CONSTRAINT_TYPE_PK  = "pk"
	CONSTRAINT_TYPE_UK  = "uk"
	CONSTRAINT_TYPE_IDX = "idx"
	CONSTRAINT_TYPE_FK  = "fk"
	CONSTRAINT_TYPE_FT  = "ft"
)

type MetaConstraint struct {
	Name        string   `json:"name"`
	ColumnNames []string `json:"column_names"`
	Type        string   `json:"type"`
}

func NewMetaConstraint() *MetaConstraint {
	return &MetaConstraint{
		ColumnNames: make([]string, 0, 1),
	}
}

type MetaCreateTable struct {
	Schema        string            `json:"schema"`
	Table         string            `json:"table"`
	Columns       []*MetaColumn     `json:"columns"`
	Constraints   []*MetaConstraint `json:"constraints"`
	IfNotExists   bool              `json:"if_not_exists"`
	Engine        string            `json:"engine"`
	Charset       string            `json:"charset"`
	Collate       string            `json:"collate"`
	Comment       string            `json:"comment"`
	AutoIncrement string            `json:"auto_increment"`
}

func NewMetaCreateTable() *MetaCreateTable {
	return &MetaCreateTable{
		Columns:     make([]*MetaColumn, 0, 1),
		Constraints: make([]*MetaConstraint, 0, 1),
	}
}

func (this *MetaCreateTable) MetaDataInterfaceFunc() {}

const (
	ALTER_TABLE_TYPE_ADD_COLUMN      = "add_column"
	ALTER_TABLE_TYPE_DROP_COLUMN     = "drop_column"
	ALTER_TABLE_TYPE_CHANGE_COLUMN   = "change_column"
	ALTER_TABLE_TYPE_MODIFY_COLUMN   = "modify_column"
	ALTER_TABLE_TYPE_ADD_CONSTRAINT  = "add_constraint"
	ALTER_TABLE_TYPE_DROP_CONSTRAINT = "drop_constraint"
	ALTER_TABLE_TYPE_RENAME_TABLE    = "rename_table"
	ALTER_TABLE_TYPE_OPTIONS         = "options"
	ALTER_TABLE_TYPE_RENAME_INDEX    = "rename_index"
)

type MetaAddColumn struct {
	Type    string        `json:"type"`
	After   string        `json:"after"`
	Columns []*MetaColumn `json:"columns"`
}

func (this *MetaAddColumn) MetaDataInterfaceFunc() {}

func NewMetaAddColumn() *MetaAddColumn {
	return &MetaAddColumn{
		Type:    ALTER_TABLE_TYPE_ADD_COLUMN,
		Columns: make([]*MetaColumn, 0, 1),
	}
}

type MetaDropColumn struct {
	Type       string `json:type`
	ColumnName string `json:"column_name"`
}

func (this *MetaDropColumn) MetaDataInterfaceFunc() {}

func NewMetaDropColumn() *MetaDropColumn {
	return &MetaDropColumn{
		Type: ALTER_TABLE_TYPE_DROP_COLUMN,
	}
}

type MetaAddConstraint struct {
	Type       string          `json:"type"`
	Constraint *MetaConstraint `json:"constraint"`
}

func (this *MetaAddConstraint) MetaDataInterfaceFunc() {}

func NewMetaAddConstraint() *MetaAddConstraint {
	return &MetaAddConstraint{
		Type:       ALTER_TABLE_TYPE_ADD_CONSTRAINT,
		Constraint: NewMetaConstraint(),
	}
}

type MetaDropConstraint struct {
	Type           string `json:"type"`
	ConstraintType string `json:"constraint_type"`
	Name           string `json:"name"`
}

func (this *MetaDropConstraint) MetaDataInterfaceFunc() {}

type MetaModifyColumn struct {
	Type    string        `json:"type"`
	After   string        `json:"after"`
	Columns []*MetaColumn `json:"columns"`
}

func (this *MetaModifyColumn) MetaDataInterfaceFunc() {}

func NewMetaModifyColumn() *MetaModifyColumn {
	return &MetaModifyColumn{
		Type:    ALTER_TABLE_TYPE_MODIFY_COLUMN,
		Columns: make([]*MetaColumn, 0, 1),
	}
}

type MetaChangeColumn struct {
	Type      string      `json:"type"`
	OldName   string      `json:"old_name"`
	After     string      `json:"after"`
	NewColumn *MetaColumn `json:"new_column"`
}

func (this *MetaChangeColumn) MetaDataInterfaceFunc() {}

func NewMetaChangeColumn() *MetaChangeColumn {
	return &MetaChangeColumn{
		Type:      ALTER_TABLE_TYPE_CHANGE_COLUMN,
		NewColumn: new(MetaColumn),
	}
}

type MetaRenameTable struct {
	Type   string `json:"type"`
	Schema string `json:"schema"`
	Table  string `json:"table"`
}

func (this *MetaRenameTable) MetaDataInterfaceFunc() {}

func NewMetaRenameTable() *MetaRenameTable {
	return &MetaRenameTable{
		Type: ALTER_TABLE_TYPE_RENAME_TABLE,
	}
}

const (
	TABLE_OPTION_ENGINE         = "engine"
	TABLE_OPTION_CHARSET        = "charset"
	TABLE_OPTION_COLLATE        = "collate"
	TABLE_OPTION_COMMENT        = "comment"
	TABLE_OPTION_AUTO_INCREMENT = "auto_increment"
)

type MetaTableOption struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type MetaTableOptions struct {
	Type    string             `json:"type"`
	Options []*MetaTableOption `json:"options"`
}

func (this *MetaTableOptions) MetaDataInterfaceFunc() {}

func NewMetaTableOptions() *MetaTableOptions {
	return &MetaTableOptions{
		Type:    ALTER_TABLE_TYPE_OPTIONS,
		Options: make([]*MetaTableOption, 0, 1),
	}
}

type MetaRenameIndex struct {
	Type    string `json:"type"`
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

func (this *MetaRenameIndex) MetaDataInterfaceFunc() {}

func NewMetaRenameIndex() *MetaRenameIndex {
	return &MetaRenameIndex{
		Type: ALTER_TABLE_TYPE_RENAME_INDEX,
	}
}

// alter table 元信息
type MetaAlterTable struct {
	Schema string     `json:"schema"`
	Table  string     `json:"table"`
	OPS    []MetaData `json:"ops"`
}

func (this *MetaAlterTable) MetaDataInterfaceFunc() {}

func NewMetaAlterTable() *MetaAlterTable {
	return &MetaAlterTable{
		OPS: make([]MetaData, 0, 1),
	}
}

// drop table 元数据
type MetaDropTable struct {
	Schema   string `json:"schema"`
	Table    string `json:"table"`
	IfExists bool   `json:"if_exists"`
}

func (this *MetaDropTable) MetaDataInterfaceFunc() {}

type MetaDropTables []*MetaDropTable

func (this MetaDropTables) MetaDataInterfaceFunc() {}

// truncate table 元数据
type MetaTruncateTable struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
}

func (this *MetaTruncateTable) MetaDataInterfaceFunc() {}
