package reviewer

import (
	"fmt"
	"strings"

	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/common"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dependency/mysql"
)

type AlterTableReviewer struct {
	ReViewMSG *ReviewMSG

	StmtNode          *ast.AlterTableStmt
	ReviewConfig      *config.ReviewConfig
	DBConfig          *config.DBConfig
	AddColumns        map[string]bool     // 新添加的字段
	DropColumns       map[string]bool     // 要删除的字段
	AddIndexes        map[string][]string // 新添加的索引
	DropIndexes       map[string]bool     //  要删除的索引
	AddUniqueIndexes  map[string][]string // 新添加的唯一索引
	IsDropPrimaryKey  bool                // 要删除的主键
	PKName            string              // 主键名称
	PKColumnNames     map[string]bool     // 主键字段
	AfterColumnNames  map[string]bool     // 有出现在 after 字句中的字段
	AutoIncrementName string              // 自增字段名字
	AddPartitions     map[string]bool     // 需要添加的分区
	DropPartitions    map[string]bool     // 需要删除的分区
	ColumnsCharLenMap map[string]int      // 每个字段

	NotAllowColumnTypeMap   map[string]bool // 不允许的字段类型
	NotNullColumnTypeMap    map[string]bool // 必须为not null的字段类型
	NotNullColumnNameMap    map[string]bool // 必须为 not null的字段名称
	ColumnTypeCount         map[byte]int    // 保存字段类型出现的个数
	NeedDefaultValueNameMap map[string]bool // 必须要有默认值的字段名

	OldSchemaName string // 原数据库名
	NewSchemaName string // 新数据库名
	OldTableName  string // 原表名
	NewTableName  string // 新表名
}

func (this *AlterTableReviewer) Init() {
	this.ReViewMSG = NewReivewMSG(config.StmtTypeAlterTable,
		this.StmtNode.Table.Schema.String(), this.StmtNode.Table.Name.String())

	this.AddColumns = make(map[string]bool)
	this.DropColumns = make(map[string]bool)
	this.AddIndexes = make(map[string][]string)
	this.DropIndexes = make(map[string]bool)
	this.AddUniqueIndexes = make(map[string][]string)
	this.PKColumnNames = make(map[string]bool)
	this.AfterColumnNames = make(map[string]bool)
	this.AddPartitions = make(map[string]bool)
	this.DropPartitions = make(map[string]bool)
	this.ColumnsCharLenMap = make(map[string]int)

	this.NotAllowColumnTypeMap = this.ReviewConfig.GetNotAllowColumnTypeMap()
	this.NotNullColumnTypeMap = this.ReviewConfig.GetNotNullColumnTypeMap()
	this.NotNullColumnNameMap = this.ReviewConfig.GetNotNullColumnNameMap()
	this.ColumnTypeCount = make(map[byte]int)
	this.NeedDefaultValueNameMap = this.ReviewConfig.GetNeedDefaultValueNameMap()

	this.OldSchemaName = this.StmtNode.Table.Schema.String()
	this.OldTableName = this.StmtNode.Table.Name.String()
}

func (this *AlterTableReviewer) Review() *ReviewMSG {
	this.Init()

	var haveError bool

	// 循环每个
	for i, spec := range this.StmtNode.Specs {
		switch spec.Tp {
		case ast.AlterTableOption:
			haveError = this.DetectTableOptions(spec)
		case ast.AlterTableAddColumns:
			haveError = this.DetectAddColumn(spec)
		case ast.AlterTableAddConstraint:
			haveError = this.DetectAddConstraint(spec)
		case ast.AlterTableDropColumn:
			haveError = this.DetectDropColumn(spec)
		case ast.AlterTableDropPrimaryKey:
			haveError = this.DetectDropPrimaryKey(spec)
		case ast.AlterTableDropIndex:
			haveError = this.DetectDropIndex(spec)
		case ast.AlterTableDropForeignKey:
			msg := fmt.Sprintf("检测失败, 第%v个字句AlterTableDropForeignKey, 碰到不能识别的语句. 请联系DBA", i)
			this.ReViewMSG.AppendMSG(true, msg)
			return this.ReViewMSG
		case ast.AlterTableModifyColumn:
			haveError = this.DetectModifyColumn(spec)
		case ast.AlterTableChangeColumn:
			haveError = this.DetectChangeColumn(spec)
		case ast.AlterTableRenameTable:
			haveError = this.DetectRenameTable(spec)
		case ast.AlterTableAlterColumn:
			msg := fmt.Sprintf("检测失败, 第%v个字句AlterTableAlterColumn, 碰到不能识别的语句. 请联系DBA", i)
			this.ReViewMSG.AppendMSG(true, msg)
			return this.ReViewMSG
		case ast.AlterTableLock:
			msg := fmt.Sprintf("检测失败, 第%v个字句AlterTableLock, 碰到不能识别的语句. 请联系DBA", i)
			this.ReViewMSG.AppendMSG(true, msg)
			return this.ReViewMSG
		case ast.AlterTableAlgorithm:
			msg := fmt.Sprintf("检测失败, 第%v个字句AlterTableAlgorithm, 碰到不能识别的语句. 请联系DBA", i)
			this.ReViewMSG.AppendMSG(true, msg)
			return this.ReViewMSG
		case ast.AlterTableRenameIndex:
			haveError = this.DetectRenameIndex(spec)
		case ast.AlterTableForce:
			msg := fmt.Sprintf("检测失败, 第%v个字句AlterTableForce, 碰到不能识别的语句. 请联系DBA", i)
			this.ReViewMSG.AppendMSG(true, msg)
			return this.ReViewMSG
		case ast.AlterTableAddPartitions:
			haveError = this.DetectAddPartition(spec)
		case ast.AlterTableDropPartition:
			haveError = this.DetectDropPartition(spec)
		default:
			msg := fmt.Sprintf("检测失败, 第%v个字句, 碰到不能识别的语句. 请联系DBA", i)
			this.ReViewMSG.AppendMSG(true, msg)
			return this.ReViewMSG
		}

		if haveError {
			return this.ReViewMSG
		}
	}

	// 在没有链接数据库的时候也需要检测一下索引长度是否超过
	haveError = this.DetectIndexCharLengthNoInstance()
	if haveError {
		return this.ReViewMSG
	}

	// 和实例中的表进行检测
	haveError = this.DetectInstanceTable()
	if haveError {
		return this.ReViewMSG
	}

	return this.ReViewMSG
}

/* 检测添加字段字句
Params:
    _spec: 添加字段字句
*/
func (this *AlterTableReviewer) DetectAddColumn(_spec *ast.AlterTableSpec) (haveError bool) {
	for _, column := range _spec.NewColumns {
		// 检测新增字段
		haveError = this.DetectNewColumn(column, _spec, "alter add")
		if haveError {
			return
		}

		// 对每个字段添加类型个数
		this.IncrColumnTypeCount(column)
	}

	return
}

// 增加字段个数
func (this *AlterTableReviewer) IncrColumnTypeCount(_column *ast.ColumnDef) {
	switch _column.Tp.Tp {
	case mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob, mysql.TypeBlob:
		// 4种大字段都设置为是 Blob
		this.ColumnTypeCount[mysql.TypeBlob]++
	default:
		this.ColumnTypeCount[_column.Tp.Tp]++
	}
}

/* 检测删除字段语句
Params:
    _spec: 删除字段字句
*/
func (this *AlterTableReviewer) DetectDropColumn(_spec *ast.AlterTableSpec) (haveError bool) {
	// 不允许删除字段
	if !this.ReviewConfig.RuleAllowDropColumn {
		haveError = true
		msg := fmt.Sprintf("检测失败. [%v], %v",
			_spec.OldColumnName.String(), config.MSG_ALLOW_DROP_COLUMN_ERROR)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	// 添加删除字段
	if _, ok := this.DropColumns[_spec.OldColumnName.String()]; ok {
		haveError = true
		msg := fmt.Sprintf("检测失败. 语句出现重复删除字段[%v]",
			_spec.OldColumnName.String())
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}
	this.DropColumns[_spec.OldColumnName.String()] = true

	return
}

/* 检测modify 子句
Params:
    _spec: modify 子句
*/
func (this *AlterTableReviewer) DetectModifyColumn(_spec *ast.AlterTableSpec) (haveError bool) {
	for _, column := range _spec.NewColumns {
		// 检测新增字段
		haveError = this.DetectNewColumn(column, _spec, "alter modify")
		if haveError {
			return
		}

		// 因为modify的字段是已经存在的所以不应该添加到 添加字段的列表中, 而应该添加到删除字段的列表中
		delete(this.AddColumns, column.Name.String())
		this.DropColumns[column.Name.String()] = true

		// 对每个字段添加类型个数
		this.IncrColumnTypeCount(column)
	}

	return
}

func (this *AlterTableReviewer) DetectChangeColumn(_spec *ast.AlterTableSpec) (haveError bool) {
	// 是否允许使用 alter change 子句
	if !this.ReviewConfig.RuleAllowChangeColumn {
		haveError = true
		msg := fmt.Sprintf("检测失败. %v", config.MSG_ALLOW_CHANGE_COLUMN_ERROR)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	// 加入到 drop column(需要删除的列) 列表中
	if _, ok := this.DropColumns[_spec.OldColumnName.String()]; ok {
		haveError = true
		msg := fmt.Sprintf("检测失败. alter change 的字段名和需要删除或其他alter change字段重复[%v]",
			_spec.OldColumnName.String())
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}
	this.DropColumns[_spec.OldColumnName.String()] = true

	// 检测修改为的新字段
	for _, column := range _spec.NewColumns {
		// 检测新增字段
		haveError = this.DetectNewColumn(column, _spec, "alter change")
		if haveError {
			return
		}

		// 对每个字段添加类型个数
		this.IncrColumnTypeCount(column)
	}

	return
}

/* 检测新增字段定义
Params:
    _column: 字段定义
    _spec: 子句
    _state: 是那个场景的字段. alter add/modify/change
*/
func (this *AlterTableReviewer) DetectNewColumn(
	_column *ast.ColumnDef,
	_spec *ast.AlterTableSpec,
	_state string,
) (haveError bool) {
	// 添加字段, 并检测字段是否在本次添加中有重复
	if _, ok := this.AddColumns[_column.Name.String()]; ok {
		haveError = true
		msg := fmt.Sprintf("检测失败: 本次添加字段语句中有重复的字段[%v]",
			_column.Name.String())
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}
	this.AddColumns[_column.Name.String()] = true

	// 检测字段名字长度
	var msg string
	haveError, msg = DetectNameLength(_column.Name.String(), this.ReviewConfig.RuleNameLength)
	if haveError {
		msg = fmt.Sprintf("%v %v", "字段名", msg)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}
	this.ReViewMSG.AppendMSG(haveError, msg)

	// 检测字段名字规则
	haveError, msg = DetectNameReg(_column.Name.String(), this.ReviewConfig.RuleNameReg)
	if haveError {
		msg = fmt.Sprintf("%v %v", "字段名", msg)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}
	this.ReViewMSG.AppendMSG(haveError, msg)

	var isNotNull bool = false        // 该字段是否为 not null
	var hasDefaultValue bool = false  // 是否有默认值
	var hasColumnComment bool = false // 用于检测字段的注释是否指定
	// 获取字段是否 not null, 是否有默认值
	for _, option := range _column.Options {
		switch option.Tp {
		case ast.ColumnOptionPrimaryKey:
			// 将主键添加到唯一键中
			if this.PKName != "" {
				haveError = true
				msg = fmt.Sprintf("检测失败: %v 语句中有重复定义主键[%v]",
					_state, _column.Name.String())
				this.ReViewMSG.AppendMSG(haveError, msg)
				return
			}
			this.PKName = "PRIMARY KEY"

			// 检测本次sql中是否有添加重复主键
			if _, ok := this.PKColumnNames[_column.Name.String()]; ok {
				haveError = true
				msg = fmt.Sprintf("检测失败: %v 语句中有重复定义主键字段[%v]",
					_state, _column.Name.String())
				this.ReViewMSG.AppendMSG(haveError, msg)
				return
			}
			this.PKColumnNames[_column.Name.String()] = true

			if uniqueIndexColumnNames, ok := this.AddUniqueIndexes[this.PKName]; ok {
				if uniqueIndexColumnNames == nil {
					this.AddUniqueIndexes[this.PKName] = make([]string, 0, 1)
				}
				this.AddUniqueIndexes[this.PKName] = append(this.AddUniqueIndexes[this.PKName], _column.Name.String())
			}
		case ast.ColumnOptionNotNull:
			isNotNull = true
		case ast.ColumnOptionDefaultValue:
			hasDefaultValue = true
		case ast.ColumnOptionComment:
			if strings.Trim(option.Expr.GetValue().(string), " ") != "" {
				hasColumnComment = true
			}
		case ast.ColumnOptionAutoIncrement:
			this.AutoIncrementName = _column.Name.String()
		}
	}

	// 1.检测不允许的字段类型
	if _, ok := this.NotAllowColumnTypeMap[config.GetTextSqlTypeByByte(_column.Tp.Tp)]; ok {
		haveError = true
		msg = fmt.Sprintf("%v 字段: %v, 类型: %v 不被允许. %v",
			_state, _column.Name.String(), _column.Tp.String(),
			fmt.Sprintf(config.MSG_NOT_ALLOW_COLUMN_TYPE_ERROR, this.ReviewConfig.RuleNotAllowColumnType))
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	// 2. 检测字段是否有注释
	if this.ReviewConfig.RuleNeedColumnComment { // 字段需要都有注释
		if !hasColumnComment {
			haveError = true
			msg = fmt.Sprintf("检测失败. %v. %v 字段: %v ", config.MSG_NEED_COLUMN_COMMENT_ERROR,
				_state, _column.Name.String())
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	// 3. 检测是否设置都为 NOT NULL
	if this.ReviewConfig.RuleAllColumnNotNull { // 需要所有字段为 NOT NULL
		if !isNotNull {
			haveError = true
			msg = fmt.Sprintf("检测失败. %v. %v 字段: %v ",
				config.MSG_ALL_COLUMN_NOT_NULL_ERROR,
				_state, _column.Name.String())
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	// 4. 主键必须为not null
	if _, ok := this.PKColumnNames[_column.Name.String()]; ok { // 该字段是主键
		if !isNotNull {
			haveError = true
			msg = fmt.Sprintf("检测失败. 主键必须定义为NOT NULL. %v 主键: %v",
				_state, _column.Name.String())
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	// 5. 必须为NOT NULL的字段类型
	if _, ok := this.NotNullColumnTypeMap[config.GetTextSqlTypeByByte(_column.Tp.Tp)]; ok {
		if !isNotNull {
			haveError = true
			msg = fmt.Sprintf("检测失败. %v. %v 字段: %v", _state,
				fmt.Sprintf(config.MSG_NOT_NULL_COLUMN_TYPE_ERROR, this.ReviewConfig.RuleNotNullColumnType),
				_column.Name.String())
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	// 6. 必须为NOT NULL的字段名称
	if _, ok := this.NotNullColumnNameMap[strings.ToLower(_column.Name.String())]; ok {
		if !isNotNull {
			haveError = true
			msg = fmt.Sprintf("检测失败. %v 字段: %v 必须为NOT NULL. %v. ",
				_column.Name.String(), _state,
				fmt.Sprintf(config.MSG_NOT_NULL_COLUMN_NAME_ERROR, this.ReviewConfig.RuleNotNullColumnName))
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	// 7. 必须有默认值
	if this.ReviewConfig.RuleNeedDefaultValue && !hasDefaultValue {
		haveError = true
		msg = fmt.Sprintf("检测失败. %v 字段: %v %v. ",
			_state, _column.Name.String(), config.MSG_NEED_DEFAULT_VALUE_ERROR)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	// 8. 必须要有默认值的字段
	if _, ok := this.NeedDefaultValueNameMap[strings.ToLower(_column.Name.String())]; ok {
		if !hasDefaultValue {
			haveError = true
			msg = fmt.Sprintf("检测失败. %v 字段: %v 必须要有默认值. %v. ",
				_column.Name.String(), _state,
				fmt.Sprintf(config.MSG_NEED_INDEX_COLUMN_NAME_ERROR, this.ReviewConfig.RuleNeedDefaultValueName))
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	// 检测是否允许after 字句
	if _spec.Position != nil {
		switch _spec.Position.Tp {
		case ast.ColumnPositionFirst:
		case ast.ColumnPositionAfter:
			if !this.ReviewConfig.RuleAllowAfterClause {
				haveError = true
				msg = fmt.Sprintf("检测失败. %v 字段[%v], %v",
					_column.Name.String(), _state, config.MSG_ALLOW_AFTER_CLAUSE_ERROR)
				this.ReViewMSG.AppendMSG(haveError, msg)
				return
			}
			this.AfterColumnNames[_spec.Position.RelativeColumn.Name.String()] = true
		}
	}

	// 检测字段定义的 Charset
	if this.ReviewConfig.RuleAllowColumnCharset { // 允许collate
		haveError, msg = DetectCharset(_column.Tp.Charset, this.ReviewConfig.RuleCollate)
		if haveError {
			this.ReViewMSG.AppendMSG(haveError, fmt.Sprintf("字段:%s . %s", _column.Name.String(), msg))
			return
		}
	} else { // 不允许charset
		if len(_column.Tp.Charset) > 0 {
			haveError = true
			this.ReViewMSG.AppendMSG(haveError, fmt.Sprintf(config.MSG_ALLOW_COLUMN_CHARSET,
				_column.Name.String()))
			return
		}
	}

	// 检测字段定义 Collect
	if this.ReviewConfig.RuleAllowColumnCollate { // 允许字段使用charset
		haveError, msg = DetectCollate(_column.Tp.Charset, this.ReviewConfig.RuleCollate)
		if haveError {
			this.ReViewMSG.AppendMSG(haveError, fmt.Sprintf("字段:%s . %s", _column.Name.String(), msg))
			return
		}
	} else { // 不允许字段使用collate
		if len(_column.Tp.Collate) > 0 {
			haveError = true
			this.ReViewMSG.AppendMSG(haveError, fmt.Sprintf(config.MSG_ALLOW_COLUMN_COLLATE,
				_column.Name.String()))
			return
		}
	}

	// 检测 blob 类型不能有默认值
	if IsBlob(_column.Tp.Tp) && hasDefaultValue {
		haveError = true
		this.ReViewMSG.AppendMSG(haveError, fmt.Sprintf("字段:%s. Text/Blob/JSON/GEO类型 不能有默认值",
			_column.Name.String()))
		return
	}

	// 添加字段定义长度
	charLen, err := GetColumnDefineCharLen(_column)
	if err != nil {
		haveError = true
		this.ReViewMSG.AppendMSG(haveError, err.Error())
		return
	}
	this.ColumnsCharLenMap[_column.Name.String()] = charLen

	return
}

/* 检测修改表名称
Params:
    _spec: 子句
*/
func (this *AlterTableReviewer) DetectRenameTable(_spec *ast.AlterTableSpec) (haveError bool) {
	if !this.ReviewConfig.RuleAllowRenameTable {
		haveError = true
		msg := fmt.Sprintf("检测失败. alter rename. %v",
			config.MSG_ALLOW_RENAME_TABLE_ERROR)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	// 是否有重复对表进行rename
	if this.NewTableName != "" {
		haveError = true
		msg := fmt.Sprintf("检测失败. alter rename. 重复 %v", this.NewTableName)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	this.NewSchemaName = _spec.NewTable.Schema.String()
	this.NewTableName = _spec.NewTable.Name.String()

	// 检测字段名字长度
	var msg string
	haveError, msg = DetectNameLength(this.NewTableName, this.ReviewConfig.RuleNameLength)
	if haveError {
		msg = fmt.Sprintf("检测失败. alter rename 新表名: %v, 长度 %v",
			this.NewTableName, msg)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	// 检测字段名字规则
	haveError, msg = DetectNameReg(this.NewTableName, this.ReviewConfig.RuleTableNameReg)
	if haveError {
		msg = fmt.Sprintf("%v. alter rename 新表名: %v", msg, this.NewTableName)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

/* 检测添加约束

 */
func (this *AlterTableReviewer) DetectAddConstraint(_spec *ast.AlterTableSpec) (haveError bool) {
	// 检测索引/约束名是否重复
	if _, ok := this.AddIndexes[_spec.Constraint.Name]; ok {
		haveError = true
		msg := fmt.Sprintf("检测失败. alter 语句中有索引/约束名称重复: %v",
			_spec.Constraint.Name)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}
	indexColumnNameMap := make(map[string]bool)
	this.AddIndexes[_spec.Constraint.Name] = make([]string, 0, 1)

	// 检测一个 索引/约束中不能有重复字段, 并保存 索引/约束 中
	for _, indexName := range _spec.Constraint.Keys {
		// 检测 索引/约束 中有重复字段
		if _, ok := indexColumnNameMap[indexName.Column.String()]; ok { // 警告
			msg := fmt.Sprintf("检测失败. alter 语句中 同一个 索引/约束 中有同一个重复字段. "+
				"索引/约束: %v, 重复的字段名: %v",
				_spec.Constraint.Name, indexName.Column.String())
			this.ReViewMSG.AppendMSG(haveError, msg)
		}

		if _spec.Constraint.Name == "" { // 如果主键没有名字, 不存为索引
			continue
		}
		this.AddIndexes[_spec.Constraint.Name] = append(this.AddIndexes[_spec.Constraint.Name], indexName.Column.String())
		indexColumnNameMap[indexName.Column.String()] = true // 保存 索引/约束中的字段名
	}

	// 约束名称长度
	var msg string
	haveError, msg = DetectNameLength(_spec.Constraint.Name, this.ReviewConfig.RuleNameLength)
	if haveError {
		msg = fmt.Sprintf("检测失败. alter 语句 %v. 索引/约束: %v",
			fmt.Sprintf(config.MSG_NAME_LENGTH_ERROR, this.ReviewConfig.RuleNameLength),
			_spec.Constraint.Name)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	switch _spec.Constraint.Tp {
	case ast.ConstraintPrimaryKey:
		haveError = this.DectectConstraintPrimaryKey(_spec.Constraint)
		if haveError {
			return
		}
	case ast.ConstraintKey, ast.ConstraintIndex:
		haveError = this.DectectConstraintIndex(_spec.Constraint)
		if haveError {
			return haveError
		}
	case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
		haveError = this.DectectConstraintUniqIndex(_spec.Constraint)
		if haveError {
			return
		}
	case ast.ConstraintForeignKey:
		// 检测是否允许有外键
		if !this.ReviewConfig.RuleAllowForeignKey { // 不允许有外键, 报错
			haveError = true
			msg := fmt.Sprintf("检测失败. alter 语句, %v. 表名: %v",
				config.MSG_ALLOW_FOREIGN_KEY_ERROR, this.StmtNode.Table.Name.String())
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	case ast.ConstraintFulltext:
		// 检测是否允许使用全文索引
		if !this.ReviewConfig.RuleAllowFullText { // 不允许使用全文索引, 报错
			haveError = true
			msg := fmt.Sprintf("检测失败. alter 语句 %v. 表名: %v",
				config.MSG_ALLOW_FULL_TEXT_ERROR, this.StmtNode.Table.Name.String())
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	return
}

/* 检测主键约束相关东西
Params:
	_constraint: 约束信息
*/
func (this *AlterTableReviewer) DectectConstraintPrimaryKey(_constraint *ast.Constraint) (haveError bool) {
	// 检测在字段定义字句中和约束定义字句中是否有重复定义 主键
	if len(this.PKColumnNames) > 0 {
		haveError = true
		msg := fmt.Sprintf("检测失败. alter 语句主键有重复定义. 表: %v",
			this.StmtNode.Table.Name.String())
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}
	this.PKName = "PRIMARY KEY"

	// 添加唯一索引, 赋值主键列名
	uniqueIndex := make([]string, 0, 1)
	for _, pkName := range _constraint.Keys {
		uniqueIndex = append(uniqueIndex, pkName.Column.String())
		this.PKColumnNames[pkName.Column.String()] = true
	}
	this.AddUniqueIndexes[this.PKName] = uniqueIndex
	this.AddIndexes[this.PKName] = uniqueIndex

	return
}

/* 检测索引相关信息
_constraint: 约束信息
*/
func (this *AlterTableReviewer) DectectConstraintIndex(_constraint *ast.Constraint) (haveError bool) {
	// 检测索引命名规范
	haveError, _ = DetectNameReg(_constraint.Name, this.ReviewConfig.RuleIndexNameReg)
	if haveError {
		msg := fmt.Sprintf("检测失败. alter 语句, %v 索引/约束: %v",
			fmt.Sprintf(config.MSG_INDEX_NAME_REG_ERROR, this.ReviewConfig.RuleIndexNameReg),
			_constraint.Name)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

/* 检测索引相关信息
_constraint: 约束信息
*/
func (this *AlterTableReviewer) DectectConstraintUniqIndex(_constraint *ast.Constraint) (haveError bool) {
	// 间隔唯一索引命名规范
	haveError, _ = DetectNameReg(_constraint.Name, this.ReviewConfig.RuleUniqueIndexNameReg)
	if haveError {
		msg := fmt.Sprintf("检测失败. alter 语句, %v 唯一索引: %v",
			fmt.Sprintf(config.MSG_UNIQUE_INDEX_NAME_REG_ERROR, this.ReviewConfig.RuleUniqueIndexNameReg),
			_constraint.Name)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	// 添加唯一索引
	uniqueIndex := make([]string, 0, 1)
	for _, column := range _constraint.Keys {
		uniqueIndex = append(uniqueIndex, column.Column.String())
	}
	this.AddUniqueIndexes[_constraint.Name] = uniqueIndex

	return
}

/* 检测删除索引
Params:
    _spec: 子句
*/
func (this *AlterTableReviewer) DetectDropIndex(_spec *ast.AlterTableSpec) (haveError bool) {
	if !this.ReviewConfig.RuleAllowDropIndex {
		haveError = true
		msg := fmt.Sprintf("检测失败. alter语句, %v. %v",
			config.MSG_ALLOW_DROP_INDEX_ERROR, _spec.Name)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	if _, ok := this.DropIndexes[_spec.Name]; ok {
		haveError = true
		msg := fmt.Sprintf("检测失败. alter语句, 重复删除索引: %v",
			_spec.Name)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}
	this.DropIndexes[_spec.Name] = true

	return
}

/* 检测删除主键
Params:
    _spec: 子句
*/
func (this *AlterTableReviewer) DetectDropPrimaryKey(_spec *ast.AlterTableSpec) (haveError bool) {
	if !this.ReviewConfig.RuleAllowDropPrimaryKey {
		haveError = true
		msg := fmt.Sprintf("检测失败. alter语句, %v", config.MSG_ALLOW_DROP_PRIMARY_KEY_ERROR)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	this.IsDropPrimaryKey = true

	return
}

/* 检测从命名索引
   _spec: 子句
*/
func (this *AlterTableReviewer) DetectRenameIndex(_spec *ast.AlterTableSpec) (haveError bool) {
	// 是否允许重命名索引
	if !this.ReviewConfig.RuleAllowRenameIndex {
		msg := fmt.Sprintf("检测失败. alter语句, %v", config.MSG_ALLOW_RENAME_INDEX_ERROR)
		this.ReViewMSG.AppendMSG(haveError, msg)
	}

	// 索引是否有重复
	if _, ok := this.AddIndexes[_spec.ToKey.String()]; ok {
		haveError = true
		msg := fmt.Sprintf("检测失败. alter rename index语句, 新索引名称有重复: %v",
			_spec.ToKey.String())
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	// 将新索引名称添加到需要添加的索引里面
	this.AddIndexes[_spec.ToKey.String()] = make([]string, 0, 1)
	// 将老索引名称添加到需要删除的索引
	this.DropIndexes[_spec.FromKey.String()] = true

	return
}

/* 检测需要添加的分区
Params:
    _spec: 子句
*/
func (this *AlterTableReviewer) DetectAddPartition(_spec *ast.AlterTableSpec) (haveError bool) {
	for _, partition := range _spec.PartDefinitions {
		var msg string

		if _, ok := this.AddPartitions[partition.Name.String()]; ok {
			haveError = true
			msg := fmt.Sprintf("检测失败. alter 语句添加分区名重复: %v",
				partition.Name.String())
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
		this.AddPartitions[partition.Name.String()] = true

		// 检测字段名字长度
		haveError, msg = DetectNameLength(partition.Name.String(), this.ReviewConfig.RuleNameLength)
		if haveError {
			msg = fmt.Sprintf("检测失败. alter 添加分区名称长度 %v", msg)
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}

		// 检测字段名字规则
		haveError, msg = DetectNameReg(partition.Name.String(), this.ReviewConfig.RuleTableNameReg)
		if haveError {
			msg = fmt.Sprintf("检测失败. alter 添加分区名称规则 %v", msg)
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}

	}

	return
}

/* 检测需要删除的分区
Params:
    _spec: 子句
*/
func (this *AlterTableReviewer) DetectDropPartition(_spec *ast.AlterTableSpec) (haveError bool) {
	// 是否允许删除分区
	if !this.ReviewConfig.RuleAllowDropPartition {
		haveError = true
		msg := fmt.Sprintf("检测失败. alter 语句%v: %v",
			config.MSG_ALLOW_DROP_PARTITION_ERROR, _spec.Name)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	if _, ok := this.DropPartitions[_spec.Name]; ok {
		haveError = true
		msg := fmt.Sprintf("检测失败. alter 语句删除分区名重复: %v", _spec.Name)
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	this.DropPartitions[_spec.Name] = true

	return
}

/* 检测表操作
Params:
	_spec: 子句
*/
func (this *AlterTableReviewer) DetectTableOptions(_spec *ast.AlterTableSpec) (haveError bool) {
	for _, option := range _spec.Options {
		var msg string

		switch option.Tp {
		case ast.TableOptionEngine:
			haveError, msg = DetectEngine(option.StrValue, this.ReviewConfig.RuleTableEngine)
		case ast.TableOptionCharset:
			haveError, msg = DetectCharset(option.StrValue, this.ReviewConfig.RuleCharSet)
		case ast.TableOptionCollate:
			haveError, msg = DetectCollate(option.StrValue, this.ReviewConfig.RuleCollate)
		}
		// 一检测到有问题键停止下面检测, 返回检测错误值
		if haveError {
			this.ReViewMSG.ErrorMSGs = append(this.ReViewMSG.ErrorMSGs, msg)
			return
		} else {
			if msg != "" {
				this.ReViewMSG.WarningMSGs = append(this.ReViewMSG.WarningMSGs, msg)
			}
		}
	}

	return
}

func (this *AlterTableReviewer) DetectInstanceTable() (haveError bool) {
	if this.OldSchemaName == "" {
		this.OldSchemaName = this.DBConfig.Database
	}
	if this.OldTableName == "" {
		this.OldTableName = this.StmtNode.Table.Name.String()
	}
	if this.NewSchemaName == "" {
		this.NewSchemaName = this.OldSchemaName
	}

	tableInfo := NewTableInfo(this.DBConfig, this.StmtNode.Table.Name.String())
	err := tableInfo.OpenInstance()
	if err != nil { // 无法链接原实例就没必要在继续往下检测了, 不过这只是警告
		msg := fmt.Sprintf("警告: sql语法正确, 但无法链接到指定实例. 无法检测数据库是否存在.")
		this.ReViewMSG.AppendMSG(false, msg)
		tableInfo.CloseInstance()
		return
	}

	// 源数据库不存在报错
	var msg string
	tableInfo.DBName = this.OldSchemaName
	haveError, msg = DetectDatabaseNotExistsByName(tableInfo, this.OldSchemaName)
	haveMSG := this.ReViewMSG.AppendMSG(haveError, msg)
	if haveError || haveMSG { // 就算是 haveError = false, 只要有信息返回, 代表就有警告, 该警告可能是使用数据库的时候哦也出错了.就不要继续检查下去了
		tableInfo.CloseInstance()
		return
	}

	// 新数据库不错在报错
	if this.NewSchemaName != "" {
		tableInfo.DBName = this.NewSchemaName
		haveError, msg = DetectDatabaseNotExistsByName(tableInfo, this.NewSchemaName)
		haveMSG := this.ReViewMSG.AppendMSG(haveError, msg)
		if haveError || haveMSG {
			tableInfo.CloseInstance()
			return
		}
	}

	// 老表不存在报错
	haveError, msg = DetectTableNotExistsByName(tableInfo, this.OldSchemaName, this.OldTableName)
	haveMSG = this.ReViewMSG.AppendMSG(haveError, msg)
	if haveError || haveMSG {
		tableInfo.CloseInstance()
		return
	}

	// 新表不存在报错
	// 如果有设置新的表需要检测新表不存在
	if this.NewTableName != "" {
		haveError, msg = DetectTableExistsByName(tableInfo, this.NewSchemaName, this.NewTableName)
		haveMSG = this.ReViewMSG.AppendMSG(haveError, msg)
		if haveError || haveMSG {
			tableInfo.CloseInstance()
			return
		}
	}

	// 获取原表的建表语句
	err = tableInfo.InitCreateTableSql(this.OldSchemaName, this.OldTableName)
	if err != nil {
		msg = fmt.Sprintf("警告: 该alter sql语法正确. "+
			"但是无法获取到源表建表sql. %v", err)
		haveMSG = this.ReViewMSG.AppendMSG(false, msg)
		tableInfo.CloseInstance()
		return
	}
	tableInfo.CloseInstance() // 该查询数据库的地方已经完成, 关闭相关链接

	// 对源表键表语句进行解析, 得到字段, 约束, 分区 等信息
	err = tableInfo.ParseCreateTableInfo()
	if err != nil {
		msg = fmt.Sprintf("警告: 该alter sql语法正确. "+
			"但解析源表信息失败, 以至于无法检测相关信息. %v", err)
		haveMSG = this.ReViewMSG.AppendMSG(false, msg)
		tableInfo.CloseInstance()
		return
	}
	tableInfo.CloseInstance()

	// 检测索引个数是否超过指定
	haveError = this.DetectIndexCount(tableInfo)
	if haveError {
		return
	}

	// 检测添加的索引是否已经存在
	haveError = this.DetectIndexExists(tableInfo)
	if haveError {
		return
	}

	// 检测after 子句的字段是否存在
	haveError = this.DetectAfterColumnExists(tableInfo)
	if haveError {
		return
	}

	// 检测主键信息(和原表一起)
	haveError = this.DetectInstancePKInfo(tableInfo)
	if haveError {
		return
	}

	// 检测列信息(和原表一起)
	haveError = this.DetectInstanceColumnInfo(tableInfo)
	if haveError {
		return
	}

	// 检测Text字段类型使用个数是否超过限制
	haveError = this.DetectTextColumnTypeCount(tableInfo)
	if haveError {
		return
	}

	// 指定字段必须有索引
	haveError = this.DetectNeedIndexColumnName(tableInfo)
	if haveError {
		return
	}

	// 检测分区相关信息(结合原表)
	haveError = this.DetectInstancePartition(tableInfo)
	if haveError {
		return
	}

	// 检测必须要有的字段名
	haveError = this.DetectHaveColumnName(tableInfo)
	if haveError {
		return
	}

	// 检测索引中是否有唯一索引
	haveError = this.DetectAllIndexHaveUniqueIndex(tableInfo)
	if haveError {
		return
	}

	// 检测重复索引
	haveError = this.DetectDuplecateIndex(tableInfo)
	if haveError {
		return
	}

	// 检测索引中的字段是否存在
	haveError = this.DetectIndexColumnExists(tableInfo)
	if haveError {
		return
	}

	// 检测索引长度
	haveError = this.DetectIndexCharLength(this.ColumnsCharLenMap, tableInfo.ColumnsCharLenMap)
	if haveError {
		return
	}

	return
}

/* 检测实例主键信息
Params:
    _tableInfo: 原表信息
*/
func (this *AlterTableReviewer) DetectInstancePKInfo(_tableInfo *TableInfo) (haveError bool) {
	// 检测主键相关
	if len(this.PKColumnNames) > 0 { // 本次 alter 语句有 add primary key
		// 检测主键是否有重复定义, 并且该alter语句没有删除主键语句
		if len(_tableInfo.PKColumnNameList) > 0 && !this.IsDropPrimaryKey {
			msg := fmt.Sprintf("检测失败. 表: %v. 主键已经存在.",
				this.StmtNode.Table.Name.String())
			haveError = true
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}

		// 主键是否有autoincrement
		if this.ReviewConfig.RulePKAutoIncrement {
			var autoIncrementName string
			if len(this.AutoIncrementName) != 0 {
				autoIncrementName = this.AutoIncrementName
			} else if len(_tableInfo.AutoIncrementName) != 0 {
				autoIncrementName = _tableInfo.AutoIncrementName
			}
			if _, ok := this.PKColumnNames[autoIncrementName]; !ok {
				msg := fmt.Sprintf("检测失败. 表: %v. 主键必须有自增列.",
					this.StmtNode.Table.Name.String())
				haveError = true
				this.ReViewMSG.AppendMSG(haveError, msg)
				return
			}
		}
	}

	return
}

func (this *AlterTableReviewer) DetectInstanceColumnInfo(_tableInfo *TableInfo) (haveError bool) {
	// 需要删除/重命名的字段是否存在
	for columnName, _ := range this.DropColumns {
		if _, ok := _tableInfo.ColumnNameMap[columnName]; !ok {
			msg := fmt.Sprintf("检测失败. 字段 %v 不存在.", columnName)
			haveError = true
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	// 新增字段是否已经存在
	for columnName, _ := range this.AddColumns {
		if _, ok := _tableInfo.ColumnNameMap[columnName]; ok {
			msg := fmt.Sprintf("检测失败. 字段 %v 已经存在.", columnName)
			haveError = true
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	return
}

/* 检测text字段个数是否超过指定
Params:
    _tableInfo: 原表信息
*/
func (this *AlterTableReviewer) DetectTextColumnTypeCount(_tableInfo *TableInfo) (haveError bool) {
	// 获取新添加的text字段个数
	addTextCount, ok := this.ColumnTypeCount[mysql.TypeBlob]
	if !ok {
		return
	}

	oriTextCount, ok := _tableInfo.ColumnTypeCount[mysql.TypeBlob]
	if !ok {
		oriTextCount = 0
	}

	if addTextCount+oriTextCount > this.ReviewConfig.RuleTextTypeColumnCount {
		msg := fmt.Sprintf("检测失败. 表: %v. %v.",
			this.StmtNode.Table.Name.String(),
			fmt.Sprintf(config.MSG_TEXT_TYPE_COLUMN_COUNT_ERROR, this.ReviewConfig.RuleTextTypeColumnCount))
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

/* 检测partition 相关信息
Params:
    _tableInfo: 原表信息
*/
func (this *AlterTableReviewer) DetectInstancePartition(_tableInfo *TableInfo) (haveError bool) {
	// 新增的 partition是否已经存在
	for partition, _ := range this.AddPartitions {
		if _, ok := _tableInfo.PartitionNames[partition]; ok {
			msg := fmt.Sprintf("检测失败. 新添加分区 %v 已经存在.", partition)
			haveError = true
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	// 需要删除的 partition 是否已经存在
	for partition, _ := range this.DropPartitions {
		if _, ok := _tableInfo.PartitionNames[partition]; !ok {
			msg := fmt.Sprintf("检测失败. 要删除的分区 %v 不存在.", partition)
			haveError = true
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}

	return
}

/* 检测必须包含的字段名
Params:
    _tableInfo: 原表信息
*/
func (this *AlterTableReviewer) DetectHaveColumnName(_tableInfo *TableInfo) (haveError bool) {
	// 检测必须要有的字段
	for haveColumnName, _ := range this.ReviewConfig.GetHaveColumnNameMap() {
		// 先判断是否有该字段
		if _, ok := _tableInfo.ColumnNameMap[haveColumnName]; !ok {
			if _, ok := this.AddColumns[haveColumnName]; !ok {
				msg := fmt.Sprintf("检测失败. 表: %v. 没有指定字段: %v. %v.",
					this.StmtNode.Table.Name.String(),
					haveColumnName,
					fmt.Sprintf(config.MSG_HAVE_COLUMN_NAME_ERROR, this.ReviewConfig.RuleHaveColumnName))
				this.ReViewMSG.AppendMSG(haveError, msg)
				continue
			}
		}
	}

	return
}

/* 检测必须指定索引的字段名, 并且必须是索引的第一个字段
Params:
    _tableInfo: 原表信息
*/
func (this *AlterTableReviewer) DetectNeedIndexColumnName(_tableInfo *TableInfo) (haveError bool) {
	// 循环必须要有索引的字段.
	for needIndexColumnName, _ := range this.ReviewConfig.GetNeedIndexColumnNameMap() {
		// 先判断是否有该字段
		if _, ok := this.AddColumns[needIndexColumnName]; !ok {
			continue
		}

		exists := false // 初始化该字段不存在索引

		for _, index := range this.AddIndexes { // 只需要检查索引中的第一个字段名就好
			if len(index) == 0 { // 没有字段的index应该是rename index操作添加到AddIndexes列表的
				continue
			}
			if needIndexColumnName == strings.ToLower(index[0]) {
				exists = true
				break
			}
		}

		if !exists {
			msg := fmt.Sprintf("检测失败. 字段 %v 必须要有索引. %v.",
				needIndexColumnName,
				fmt.Sprintf(config.MSG_NEED_INDEX_COLUMN_NAME_ERROR, this.ReviewConfig.RuleNeedIndexColumnName))
			this.ReViewMSG.AppendMSG(haveError, msg)
			continue
		}
	}

	return
}

/* 索引中是否包含唯一索引
Params:
    _tableInfo: 原表信息
*/
func (this *AlterTableReviewer) DetectAllIndexHaveUniqueIndex(_tableInfo *TableInfo) (haveError bool) {
	normalIndexes := CombineIndexes(_tableInfo.Indexes, this.AddIndexes)
	uniqueIndexes := CombineIndexes(_tableInfo.UniqueIndexes, this.AddUniqueIndexes)
	hashNormalIndex := GetIndexesHashColumn(normalIndexes)
	hashUniqueIndex := GetIndexesHashColumn(uniqueIndexes)

	// 循环 唯一键和 索引进行匹配, 看看唯一索引的字段是否都包含在普通索引中
	for uniqueIndexName, hashUniqueIndexStr := range hashUniqueIndex {

		for normalIndexName, hashNormalIndexStr := range hashNormalIndex {
			if normalIndexName == uniqueIndexName { // 同一个索引不进行比较
				continue
			}
			if isMatch := common.StrIsMatch(hashNormalIndexStr, hashUniqueIndexStr); isMatch {
				msg := fmt.Sprintf("检测失败. 普通索引: %v, 包含了唯一索引: %v 的字段.",
					normalIndexName, uniqueIndexName)
				this.ReViewMSG.AppendMSG(haveError, msg)
				continue
			}
		}
	}

	return
}

/* 检测是否有重复索引
Params:
    _tableInfo: 原表信息
*/
func (this *AlterTableReviewer) DetectDuplecateIndex(_tableInfo *TableInfo) (haveError bool) {
	normalIndexes := CombineIndexes(_tableInfo.Indexes, this.AddIndexes)
	hashNormalIndex := GetIndexesHashColumn(normalIndexes)

	// 循环 唯一键和 索引进行匹配, 看看唯一索引的字段是否都包含在普通索引中
	for normalIndexName1, hashNormalIndexStr1 := range hashNormalIndex {

		for normalIndexName2, hashNormalIndexStr2 := range hashNormalIndex {
			if normalIndexName1 == normalIndexName2 { // 同一个索引不进行比较
				continue
			}
			if isMatch := strings.HasPrefix(hashNormalIndexStr1, hashNormalIndexStr2); isMatch {
				msg := fmt.Sprintf("检测失败. 检测到重复索引: %v <=> %v.",
					normalIndexName1, normalIndexName2)
				this.ReViewMSG.AppendMSG(haveError, msg)
				continue
			}
		}
	}

	return
}

/* 检测索引中的字段是否都存在
Params:
    _tableInfo: 原表信息
*/
func (this *AlterTableReviewer) DetectIndexColumnExists(_tableInfo *TableInfo) (haveError bool) {
	for indexName, columns := range this.AddIndexes {
		if len(columns) == 0 {
			continue
		}

		for _, columnName := range columns {
			if _, ok := _tableInfo.ColumnNameMap[columnName]; !ok {
				if _, ok := this.AddColumns[columnName]; !ok {
					haveError = true
					msg := fmt.Sprintf("检测失败. 索引中字段不存在. 索引: %v, 字段: %v",
						indexName, columnName)
					this.ReViewMSG.AppendMSG(haveError, msg)
					return
				}
			}
		}
	}

	return
}

/* 检测索引个数是否超过指定个数
Params:
    _tableInfo: 原表信息
*/
func (this *AlterTableReviewer) DetectIndexCount(_tableInfo *TableInfo) (haveError bool) {
	addIndexCount := 0
	for _, columns := range this.AddIndexes { // 剔除rename索引操作加入的索引
		if len(columns) == 0 {
			continue
		}
		addIndexCount++
	}

	if addIndexCount+len(_tableInfo.Indexes) > this.ReviewConfig.RuleIndexCount {
		haveError = true
		msg := fmt.Sprintf("检测失败. %v",
			fmt.Sprintf(config.MSG_INDEX_COUNT_ERROR, this.ReviewConfig.RuleIndexCount))
		this.ReViewMSG.AppendMSG(haveError, msg)
		return
	}

	return
}

/* 检测After子句的字段是否存在
Params:
    _tableInfo: 原表信息
*/
func (this *AlterTableReviewer) DetectAfterColumnExists(_tableInfo *TableInfo) (haveError bool) {
	for afterColumnName, _ := range this.AfterColumnNames {
		if _, ok := _tableInfo.ColumnNameMap[afterColumnName]; !ok {
			if _, ok := this.AddColumns[afterColumnName]; !ok {
				msg := fmt.Sprintf("检测失败. after 字句指定的字段 %v 不存在",
					afterColumnName)
				this.ReViewMSG.AppendMSG(haveError, msg)
				return
			}
		}
	}

	return
}

// 检测索引的字符长度, 没有获取数据库中的源数据
func (this *AlterTableReviewer) DetectIndexCharLengthNoInstance() (haveError bool) {
	for name, columnNames := range this.AddIndexes {
		len := GetColumnsCharLen(this.ColumnsCharLenMap, columnNames...)
		if len > this.ReviewConfig.RuleIndexCharLength {
			msg := fmt.Sprintf("检测失败(没有链接数据库). 索引:%v. %v", name,
				fmt.Sprintf(config.MSG_INDEX_CHAR_LENGTH_ERROR, this.ReviewConfig.RuleIndexCharLength))
			haveError = true
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}
	return
}

// 检测索引的字符长度, 没有获取数据库中的源数据
func (this *AlterTableReviewer) DetectIndexCharLength(maps ...map[string]int) (haveError bool) {
	lenMap := common.CombindMapStrInt(maps...)
	for name, columnNames := range this.AddIndexes {
		len := GetColumnsCharLen(lenMap, columnNames...)
		if len > this.ReviewConfig.RuleIndexCharLength {
			msg := fmt.Sprintf("检测失败. 索引:%v. %v", name,
				fmt.Sprintf(config.MSG_INDEX_CHAR_LENGTH_ERROR, this.ReviewConfig.RuleIndexCharLength))
			haveError = true
			this.ReViewMSG.AppendMSG(haveError, msg)
			return
		}
	}
	return
}

/* 检测索引是否已经存在
Params:
    _tableInfo: 原表信息
*/
func (this *AlterTableReviewer) DetectIndexExists(_tableInfo *TableInfo) (haveError bool) {
	for name, _ := range this.AddIndexes { // 剔除rename索引操作加入的索引
		if _, ok := _tableInfo.Indexes[name]; ok { // 源表中已经有该索引名称
			if _, dropOK := this.DropIndexes[name]; !dropOK { // 判断该索引名称是否是否有被删除
				haveError = true
				msg := fmt.Sprintf("索引 %s 已经存在. 不允许添加", name)
				this.ReViewMSG.AppendMSG(haveError, msg)
				return
			}
		}
	}

	return
}
