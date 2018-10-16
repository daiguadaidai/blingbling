package reviewer

import (
	"github.com/daiguadaidai/blingbling/ast"
	"fmt"
	"github.com/daiguadaidai/blingbling/config"
	"strings"
	"github.com/daiguadaidai/blingbling/dependency/mysql"
	"github.com/daiguadaidai/blingbling/common"
	"github.com/daiguadaidai/blingbling/dao"
)

type CreateTableReviewer struct {
	StmtNode *ast.CreateTableStmt
	ReviewConfig *config.ReviewConfig
	DBConfig *config.DBConfig
	ColumnNames map[string]bool
	PKColumnNames map[string]bool // 所有主键列名
	PKname string // 主键名
	AutoIncrementName string  // 子增字段名
	/* 定义所有索引
	map {
		idx_xxx: ["id", "name"]
	}
	*/
	Indexes map[string][]string
	UniqueIndexes map[string][]string // 所有的唯一索引
	HasTableComment bool // 有表注释
	NotAllowColumnTypeMap map[string]bool // 不允许的字段类型
	NotNullColumnTypeMap map[string]bool // 必须为not null的字段类型
	NotNullColumnNameMap map[string]bool // 必须为 not null的字段名称
	ColumnTypeCount map[byte]int // 保存字段类型出现的个数
	PartitionColumns []string
	NeedDefaultValueNameMap map[string]bool // 必须要有默认值的字段名

	SchemaName string
}

// 初始化一些变量
func (this *CreateTableReviewer) Init() {
	this.ColumnNames = make(map[string]bool)
	this.PKColumnNames = make(map[string]bool)
	this.Indexes = make(map[string][]string)
	this.UniqueIndexes = make(map[string][]string)
	this.NotAllowColumnTypeMap = this.ReviewConfig.GetNotAllowColumnTypeMap()
	this.NotNullColumnTypeMap = this.ReviewConfig.GetNotNullColumnTypeMap()
	this.NotNullColumnNameMap = this.ReviewConfig.GetNotNullColumnNameMap()
	this.ColumnTypeCount = make(map[byte]int)
	this.PartitionColumns = make([]string, 0, 1)
	this.NeedDefaultValueNameMap = this.ReviewConfig.GetNeedDefaultValueNameMap()

	if this.StmtNode.Table.Schema.String() != "" {
		this.SchemaName = this.StmtNode.Table.Schema.String()
	} else {
		this.SchemaName = this.DBConfig.Database
	}
}

func (this *CreateTableReviewer) Review() *ReviewMSG {
	this.Init()
	var reviewMSG *ReviewMSG

	// 检测数据库名称长度
	if this.StmtNode.Table.Schema.String() != "" {
		reviewMSG = this.DetectDBNameLength(this.SchemaName)
		if reviewMSG != nil {
			reviewMSG.MSG = fmt.Sprintf("%v %v", "数据库名", reviewMSG.MSG)
			reviewMSG.Code = REVIEW_CODE_ERROR
			return reviewMSG
		}

		// 检测数据库命名规则
		reviewMSG = this.DetectDBNameReg(this.SchemaName)
		if reviewMSG != nil {
			reviewMSG.MSG = fmt.Sprintf("%v %v", "数据库名", reviewMSG.MSG)
			reviewMSG.Code = REVIEW_CODE_ERROR
			return reviewMSG
		}
	}

	// 检测表名称长度
	reviewMSG = this.DetectTableNameLength(this.StmtNode.Table.Name.String())
	if reviewMSG != nil {
		reviewMSG.MSG = fmt.Sprintf("%v %v", "表名", reviewMSG.MSG)
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测表命名规则
	reviewMSG = this.DetectTableNameReg(this.StmtNode.Table.Name.String())
	if reviewMSG != nil {
		reviewMSG.MSG = fmt.Sprintf("%v %v", "表名", reviewMSG.MSG)
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 判断是否是 create table like语句
	if IsCreateTableLikeStmt(this.StmtNode.Text()) {
		// 检测 create table like 语句
		reviewMSG = this.DetectCreateTableLike()
		if reviewMSG != nil {
			reviewMSG.Code = REVIEW_CODE_ERROR
			return reviewMSG
		}
	} else {
		// 非 create table like语句
		reviewMSG = this.DetectNoCreateTableLike()
		if reviewMSG != nil {
			reviewMSG.Code = REVIEW_CODE_ERROR
			return reviewMSG
		}
	}


	// 能走到这里说明写的语句审核成功
	reviewMSG = new(ReviewMSG)
	reviewMSG.Code = REVIEW_CODE_SUCCESS
	reviewMSG.MSG = "审核成功"

	return reviewMSG
}

// 检测数据库名长度
func (this *CreateTableReviewer) DetectDBNameLength(_name string) *ReviewMSG {
	return DetectNameLength(_name, this.ReviewConfig.RuleNameLength)
}

// 检测数据库命名规范
func (this *CreateTableReviewer) DetectDBNameReg(_name string) *ReviewMSG {
	return DetectNameReg(_name, this.ReviewConfig.RuleNameReg)
}

// 检测表名长度
func (this *CreateTableReviewer) DetectTableNameLength(_name string) *ReviewMSG {
	return DetectNameLength(_name, this.ReviewConfig.RuleNameLength)
}

// 检测表命名规范
func (this *CreateTableReviewer) DetectTableNameReg(_name string) *ReviewMSG {
	var reviewMSG *ReviewMSG

	reviewMSG = DetectNameReg(_name, this.ReviewConfig.RuleTableNameReg)
	if reviewMSG != nil {
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v 表名: %v",
			fmt.Sprintf(config.MSG_TABLE_NAME_GRE_ERROR, this.ReviewConfig.RuleTableNameReg),
			_name)
	}

	return reviewMSG
}

// 检测create table like相关操作
func (this *CreateTableReviewer) DetectCreateTableLike() *ReviewMSG {
	var reviewMSG *ReviewMSG

	tableInfo := dao.NewTableInfo(this.DBConfig, this.StmtNode.Table.Name.String())
	err := tableInfo.OpenInstance()
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在.")
		return reviewMSG
	}

	// 检测 create table like 的原表是否不存在
	reviewMSG = DetectTableNotExistsByName(tableInfo, this.StmtNode.ReferTable.Schema.String(),
		this.StmtNode.ReferTable.Name.String())
	if reviewMSG != nil {
		return CloseTableInstance(reviewMSG, tableInfo)
	}

	// 检测需要新建的表是否存在
	reviewMSG = DetectTableExistsByName(tableInfo, this.SchemaName, this.StmtNode.Table.Name.String())
	if reviewMSG != nil {
		return CloseTableInstance(reviewMSG, tableInfo)
	}

	return CloseTableInstance(reviewMSG, tableInfo)
}

// 检测不是create table like语句的相关信息
func (this *CreateTableReviewer) DetectNoCreateTableLike() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 检测建表选项
	reviewMSG = this.DetectTableOptions()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测表字段信息
	reviewMSG = this.DetectColumns()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测字段中定义了多个主键, 这中定义是在字段定义后面添加 primary key, 而不是在添加索引中定义的
	reviewMSG = this.DetectColumnPKReDefine()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测表的约束
	reviewMSG = this.DetectConstraints()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测是否有主键
	reviewMSG = this.DetectHasPK()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测主键是否有使用 auto increment
	reviewMSG = this.DetectPKAutoIncrement()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测 字段 相关项
	reviewMSG = this.DetectColumnOptions()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测Text字段类型使用个数是否超过限制
	reviewMSG = this.DetectTextColumnTypeCount()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测必须要有索引的字段
	reviewMSG = this.DetectNeedIndexColumnName()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测必须要有的字段名
	reviewMSG = this.DetectHaveColumnName()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测普通索引中不能有唯一索引联合在一起的字段
	reviewMSG = this.DetectNormalIndexHaveUniqueIndex()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测重复索引
	reviewMSG = this.DetectDuplecateIndex()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测分区表相关信息
	reviewMSG = this.DetectPartition()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 检测索引个数是否超过指定数
	reviewMSG = this.DetectIndexCount()
	if reviewMSG != nil {
		reviewMSG.Code = REVIEW_CODE_ERROR
		return reviewMSG
	}

	// 链接实例检测表相关信息
	reviewMSG = this.DetectInstanceTable()
	if reviewMSG != nil {
		return reviewMSG
	}

	return reviewMSG
}

// 检测创建数据库其他选项值
func (this *CreateTableReviewer) DetectTableOptions() *ReviewMSG {
	var reviewMSG *ReviewMSG

	for _, option := range this.StmtNode.Options {
		switch option.Tp {
		case ast.TableOptionEngine:
			reviewMSG = DetectEngine(option.StrValue, this.ReviewConfig.RuleTableEngine)
		case ast.TableOptionCharset:
			reviewMSG = DetectCharset(option.StrValue, this.ReviewConfig.RuleCharSet)
		case ast.TableOptionCollate:
			reviewMSG = DetectCollate(option.StrValue, this.ReviewConfig.RuleCollate)
		case ast.TableOptionComment:
			// 有设置表注释, 并且不是空字符串, 才代表有设置注释
			if strings.Trim(option.StrValue, " ") != "" {
				this.HasTableComment = true
			}
		}
		// 一检测到有问题键停止下面检测, 返回检测错误值
		if reviewMSG != nil {
			return reviewMSG
		}
	}

	// 检测表是否有注释
	if this.ReviewConfig.RuleNeedTableComment {
		if !this.HasTableComment {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("表: %v 检测失败. %v", this.StmtNode.Table.Name.String(),
				config.MSG_NEED_TABLE_COMMENT_ERROR)
		}
	}

	return reviewMSG
}

// 循环检测表的字段
func (this *CreateTableReviewer) DetectColumns() *ReviewMSG {
	var reviewMSG *ReviewMSG

	for _, column := range this.StmtNode.Cols {
		// 1. 检测字段名是否有重复
		if _, ok := this.ColumnNames[column.Name.String()]; ok {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("字段: %v. %v",
				column.Name.String(), config.MSG_TABLE_COLUMN_DUP_ERROR)
			return reviewMSG
		}
		this.ColumnNames[column.Name.String()] = true // 缓存字段名

		// 2. 检测字段名长度
		reviewMSG = DetectNameLength(column.Name.String(), this.ReviewConfig.RuleNameLength)
		if reviewMSG != nil {
			reviewMSG.MSG = fmt.Sprintf("字段名 %v", reviewMSG.MSG)
			return reviewMSG
		}

		// 3. 检测字段名规则
		reviewMSG = DetectNameReg(column.Name.String(), this.ReviewConfig.RuleNameReg)
		if reviewMSG != nil {
			reviewMSG.MSG = fmt.Sprintf("字段名 %v", reviewMSG.MSG)
			return reviewMSG
		}

		// 4. 增加字段类型使用个数
		this.IncrColumnTypeCount(column)

		// 5. 字段定义选项
		this.SetReviewPkInfo(column)
	}

	return reviewMSG
}

// 增加字段个数
func (this *CreateTableReviewer) IncrColumnTypeCount(_column *ast.ColumnDef) {
	switch _column.Tp.Tp {
	case mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob, mysql.TypeBlob:
		// 4种大字段都设置为是 Blob
		this.ColumnTypeCount[mysql.TypeBlob] ++
	default:
		this.ColumnTypeCount[_column.Tp.Tp] ++
	}
}

// 设置 createTableReviewer 主键的相关信息, 主键字段有哪些, 是否有自增
func (this *CreateTableReviewer) SetReviewPkInfo(_column *ast.ColumnDef) {
	for _, option := range _column.Options {
		switch option.Tp {
		case ast.ColumnOptionPrimaryKey:
			this.PKColumnNames[_column.Name.String()] = true
		case ast.ColumnOptionAutoIncrement:
			this.AutoIncrementName = _column.Name.String()
		}
	}
}

// 检测在定义字段中有多个 primary key出现
func (this *CreateTableReviewer) DetectColumnPKReDefine() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if len(this.PKColumnNames) > 1 {
		columnNames := make([]string, 0, 1)
		for name, _ := range this.PKColumnNames {
			columnNames = append(columnNames, name)
		}
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. 有两个字段都定义了主键(%v). " +
			"请考虑使用定于约束字句定义组合主键", strings.Join(columnNames, ", "))
	}

	return reviewMSG
}

// 检测是否有主键
func (this *CreateTableReviewer) DetectHasPK() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if this.ReviewConfig.RuleNeedPK {
		if len(this.PKColumnNames) < 1 {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("检测失败. 表: %v. 没有主键. %v",
				this.StmtNode.Table.Name.String(), config.MSG_NEED_PK)
		}
	}

	return reviewMSG
}

// 检测主键需要自增
func (this *CreateTableReviewer) DetectPKAutoIncrement() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if this.ReviewConfig.RulePKAutoIncrement {
		// 有主键才检查主键中需要有 auto_increment 选项
		if len(this.PKColumnNames) > 0 { // 有主键字段
			var pkHasAutoIncrement bool = false // 主键中是否有 auto_increment
			if strings.Trim(this.AutoIncrementName, " ") != "" {
				if _, ok := this.PKColumnNames[this.AutoIncrementName]; ok {
					pkHasAutoIncrement = true
				}
			}
			if !pkHasAutoIncrement {
				columnNames := make([]string, 0, 1)
				for name, _ := range this.PKColumnNames {
					columnNames = append(columnNames, name)
				}
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败. 主键字段: %v. %v",
					strings.Join(columnNames, ", "), config.MSG_PK_AUTO_INCREMENT_ERROR)
			}
		}
	}

	return reviewMSG
}

// 循环检测数据库的相关索引信息
func (this *CreateTableReviewer) DetectConstraints() *ReviewMSG {
	var reviewMSG *ReviewMSG

	for _, constraint := range this.StmtNode.Constraints {
		// 检测索引/约束名是否重复
		if _, ok := this.Indexes[constraint.Name]; ok {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("检测失败. 有索引/约束名称重复: %v", constraint.Name)
			return reviewMSG
		}
		indexColumnNameMap := make(map[string]bool)
		this.Indexes[constraint.Name] = make([]string, 0, 1)

		// 检测一个 索引/约束中不能有重复字段, 并保存 索引/约束 中
		for _, indexName := range constraint.Keys {
			// 检测 索引/约束 中有重复字段
			if _, ok := indexColumnNameMap[indexName.Column.String()]; ok {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败. 同一个 索引/约束 中有同一个重复字段. " +
					"索引/约束: %v, 重复的字段名: %v",
					constraint.Name, indexName.Column.String())
				return reviewMSG
			}
			this.Indexes[constraint.Name] = append(this.Indexes[constraint.Name], indexName.Column.String())
			indexColumnNameMap[indexName.Column.String()] = true // 保存 索引/约束中的字段名

			// 检测索引字段需要在表的字段中
			if _, ok := this.ColumnNames[indexName.Column.String()]; !ok {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败. 索引字段没定义. 索引/约束: %v, " +
					"字段: %v, 不存在表: %v 中 ",
					constraint.Name, indexName.Column.String(), this.StmtNode.Table.Name.String())
				return reviewMSG
			}
		}

		// 检测索引中字段个数是否符合 指定
		if len(indexColumnNameMap) > this.ReviewConfig.RuleIndexColumnCount {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("检测失败. 索引/约束: %v. %v", constraint.Name,
				fmt.Sprintf(config.MSG_INDEX_COLUMN_COUNT_ERROR, this.ReviewConfig.RuleIndexColumnCount))
			return reviewMSG
		}

		// 约束名称长度
		reviewMSG = DetectNameLength(constraint.Name, this.ReviewConfig.RuleNameLength)
		if reviewMSG != nil {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("检测失败. %v. 索引/约束: %v",
				fmt.Sprintf(config.MSG_NAME_LENGTH_ERROR, this.ReviewConfig.RuleNameLength),
				constraint.Name)
			return reviewMSG
		}

		switch constraint.Tp {
		case ast.ConstraintPrimaryKey:
			this.PKname = constraint.Name
			reviewMSG = this.DectectConstraintPrimaryKey(constraint)
			if reviewMSG != nil {
				return reviewMSG
			}

			// 添加唯一索引, 赋值主键列名
			uniqueIndex := make([]string, 0, 1)
			for _, pkName := range constraint.Keys {
				uniqueIndex = append(uniqueIndex, pkName.Column.String())
				this.PKColumnNames[pkName.Column.String()] = true
			}
			this.UniqueIndexes[constraint.Name] = uniqueIndex

		case ast.ConstraintKey, ast.ConstraintIndex:
			reviewMSG = this.DectectConstraintIndex(constraint)
			if reviewMSG != nil {
				return reviewMSG
			}
		case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
			reviewMSG = this.DectectConstraintUniqIndex(constraint)
			if reviewMSG != nil {
				return reviewMSG
			}

			// 添加唯一索引
			uniqueIndex := make([]string, 0, 1)
			for _, column := range constraint.Keys {
				uniqueIndex = append(uniqueIndex, column.Column.String())
			}
			this.UniqueIndexes[constraint.Name] = uniqueIndex

		case ast.ConstraintForeignKey:
			// 检测是否允许有外键
			if !this.ReviewConfig.RuleAllowForeignKey { // 不允许有外键, 报错
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败. %v. 表名: %v",
					config.MSG_ALLOW_FOREIGN_KEY_ERROR, this.StmtNode.Table.Name.String())
			}
		case ast.ConstraintFulltext:
			// 检测是否允许使用全文索引
			if !this.ReviewConfig.RuleAllowFullText { // 不允许使用全文索引, 报错
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败. %v. 表名: %v",
					config.MSG_ALLOW_FULL_TEXT_ERROR, this.StmtNode.Table.Name.String())
			}
		}
	}

	return reviewMSG
}

/* 检测主键约束相关东西
Params:
	_constraint: 约束信息
 */
func (this *CreateTableReviewer) DectectConstraintPrimaryKey(_constraint *ast.Constraint) *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 检测在字段定义字句中和约束定义字句中是否有重复定义 主键
	if len(this.PKColumnNames) > 0 {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. 表: %v 主键有重复定义. ",
			this.StmtNode.Table.Name.String())
		return reviewMSG
	}

	return reviewMSG
}

/* 检测索引相关信息
	_constraint: 约束信息
 */
func (this *CreateTableReviewer) DectectConstraintIndex(_constraint *ast.Constraint) *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 检测索引命名规范
	reviewMSG = DetectNameReg(_constraint.Name, this.ReviewConfig.RuleIndexNameReg)
	if reviewMSG != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v 索引/约束: %v",
			fmt.Sprintf(config.MSG_INDEX_NAME_REG_ERROR, this.ReviewConfig.RuleIndexNameReg),
			_constraint.Name)
		return reviewMSG
	}

	return reviewMSG
}

/* 检测索引相关信息
	_constraint: 约束信息
 */
func (this *CreateTableReviewer) DectectConstraintUniqIndex(_constraint *ast.Constraint) *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 间隔唯一索引命名规范
	reviewMSG = DetectNameReg(_constraint.Name, this.ReviewConfig.RuleUniqueIndexNameReg)
	if reviewMSG != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v 唯一索引: %v",
			fmt.Sprintf(config.MSG_UNIQUE_INDEX_NAME_REG_ERROR, this.ReviewConfig.RuleUniqueIndexNameReg),
			_constraint.Name)
		return reviewMSG
	}

	return reviewMSG
}

// 检测字段 相关信息
func (this *CreateTableReviewer) DetectColumnOptions() *ReviewMSG {
	var reviewMSG *ReviewMSG


	for _, column := range this.StmtNode.Cols {
		var isNotNull bool = false // 该字段是否为 not null
		var hasDefaultValue bool = false // 是否有默认值
		var hasColumnComment bool = false // 用于检测字段的注释是否指定

		// 获取字段是否 not null, 是否有默认值
		for _, option := range column.Options {
			switch option.Tp {
			case ast.ColumnOptionPrimaryKey:
			case ast.ColumnOptionNotNull:
				isNotNull = true
			case ast.ColumnOptionDefaultValue:
				hasDefaultValue = true
			case ast.ColumnOptionComment:
				if strings.Trim(option.Expr.GetValue().(string), " ") != "" {
					hasColumnComment = true
				}
			}
		}

		// 1.检测不允许的字段类型
		if _, ok := this.NotAllowColumnTypeMap[config.GetTextSqlTypeByByte(column.Tp.Tp)]; ok {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("字段: %v, 类型: %v 不被允许. %v",
				column.Name.String(), column.Tp.String(),
				fmt.Sprintf(config.MSG_NOT_ALLOW_COLUMN_TYPE_ERROR, this.ReviewConfig.RuleNotAllowColumnType))
			return reviewMSG
		}

		// 2. 检测字段是否有注释
		if this.ReviewConfig.RuleNeedColumnComment { // 字段需要都有注释
			if !hasColumnComment {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("字段: %v 检测失败. %v", column.Name.String(),
					config.MSG_NEED_COLUMN_COMMENT_ERROR)
				return reviewMSG
			}
		}

		// 3. 检测是否设置都为 NOT NULL
		if this.ReviewConfig.RuleAllColumnNotNull { // 需要所有字段为 NOT NULL
			if !isNotNull {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("字段: %v 检测失败. %v", column.Name.String(),
					config.MSG_ALL_COLUMN_NOT_NULL_ERROR)
				return reviewMSG
			}
		}

		// 4. 主键必须为not null
		if _, ok := this.PKColumnNames[column.Name.String()]; ok { // 该字段是主键
			if !isNotNull {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败. 主键必须定义为NOT NULL. 主键: %v", column.Name.String())
				return reviewMSG
			}
		}

		// 5. 必须为NOT NULL的字段类型
		if _, ok := this.NotNullColumnTypeMap[config.GetTextSqlTypeByByte(column.Tp.Tp)]; ok {
			if !isNotNull {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败. %v. 字段: %v",
					fmt.Sprintf(config.MSG_NOT_NULL_COLUMN_TYPE_ERROR, this.ReviewConfig.RuleNotNullColumnType),
					column.Name.String())
				return reviewMSG
			}
		}

		// 6. 必须为NOT NULL的字段名称
		if _, ok := this.NotNullColumnNameMap[strings.ToLower(column.Name.String())]; ok {
			if !isNotNull {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败.字段: %v 必须为NOT NULL. %v. ",
					column.Name.String(),
					fmt.Sprintf(config.MSG_NOT_NULL_COLUMN_NAME_ERROR, this.ReviewConfig.RuleNotNullColumnName))
				return reviewMSG
			}
		}

		// 7. 必须有默认值
		if this.ReviewConfig.RuleNeedDefaultValue && !hasDefaultValue {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("检测失败.字段: %v %v. ",
				column.Name.String(), config.MSG_NEED_DEFAULT_VALUE_ERROR)
			return reviewMSG
		}

		// 8. 必须要有默认值的字段
		if _, ok := this.NeedDefaultValueNameMap[strings.ToLower(column.Name.String())]; ok {
			if !hasDefaultValue {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败.字段: %v 必须要有默认值. %v. ",
					column.Name.String(),
					fmt.Sprintf(config.MSG_NEED_DEFAULT_VALUE_NAME_ERROR, this.ReviewConfig.RuleNeedDefaultValueName))
				return reviewMSG
			}
		}
	}

	return reviewMSG
}


// 检测Text字段类型使用个数
func (this *CreateTableReviewer) DetectTextColumnTypeCount() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if count, ok := this.ColumnTypeCount[mysql.TypeBlob]; ok {
		if count > this.ReviewConfig.RuleTextTypeColumnCount {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("检测失败. 表: %v. %v.",
				this.StmtNode.Table.Name.String(),
				fmt.Sprintf(config.MSG_TEXT_TYPE_COLUMN_COUNT_ERROR, this.ReviewConfig.RuleTextTypeColumnCount))
			return reviewMSG
		}
	}

	return reviewMSG
}

// 检测必须指定索引的字段名, 并且必须是索引的第一个字段
func (this *CreateTableReviewer) DetectNeedIndexColumnName() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 循环必须要有索引的字段.
	for needIndexColumnName, _ := range this.ReviewConfig.GetNeedIndexColumnNameMap() {
		// 先判断是否有该字段
		if _, ok := this.ColumnNames[needIndexColumnName]; !ok {
			continue
		}

		exists := false // 初始化该字段不存在索引

		for _, index := range this.Indexes { // 只需要检查索引中的第一个字段名就好
			if needIndexColumnName == strings.ToLower(index[0]) {
				exists = true
				break
			}
		}

		if !exists {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("检测失败. 字段 %v 必须要有索引. %v.",
				needIndexColumnName,
				fmt.Sprintf(config.MSG_NEED_INDEX_COLUMN_NAME_ERROR, this.ReviewConfig.RuleNeedIndexColumnName))
			return reviewMSG
		}
	}

	return reviewMSG
}

// 检测必须包含的字段名
func (this *CreateTableReviewer) DetectHaveColumnName() *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 循环必须要有的字段.
	for haveColumnName, _ := range this.ReviewConfig.GetHaveColumnNameMap() {
		// 先判断是否有该字段
		if _, ok := this.ColumnNames[haveColumnName]; !ok {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("检测失败. 表: %v. 没有指定字段: %v. %v.",
				this.StmtNode.Table.Name.String(),
				haveColumnName,
				fmt.Sprintf(config.MSG_HAVE_COLUMN_NAME_ERROR, this.ReviewConfig.RuleHaveColumnName))
			return reviewMSG
		}
	}

	return reviewMSG
}

// 检测普通索引中不能有唯一索引的字段(唯一索引要连在一起)
func (this *CreateTableReviewer) DetectNormalIndexHaveUniqueIndex() *ReviewMSG {
	var reviewMSG *ReviewMSG

	normalIndexes := GetNoUniqueIndexes(this.Indexes, this.UniqueIndexes)
	hashNormalIndex := GetIndexesHashColumn(normalIndexes)
	hashUniqueIndex := GetIndexesHashColumn(this.UniqueIndexes)

	// 循环 唯一键和 索引进行匹配, 看看唯一索引的字段是否都包含在普通索引中
	for uniqueIndexName, hashUniqueIndexStr := range hashUniqueIndex {
		if uniqueIndexName == "" {
			uniqueIndexName = "PRIMARY KEY"
		}

		for normalIndexName, hashNormalIndexStr := range hashNormalIndex {
			if isMatch := common.StrIsMatch(hashNormalIndexStr, hashUniqueIndexStr) ; isMatch {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败. 普通索引: %v, 包含了唯一索引: %v 的字段.",
					normalIndexName, uniqueIndexName)
				return reviewMSG
			}
		}
	}

	return reviewMSG
}

// 检测分区表
func (this *CreateTableReviewer) DetectPartition() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if this.StmtNode.Partition != nil {
		// 获取分区表字段
		if len(this.StmtNode.Partition.ColumnNames) > 0 {
			for _, columnName := range this.StmtNode.Partition.ColumnNames {
				this.PartitionColumns = append(this.PartitionColumns, columnName.Name.String())
			}
		} else {
			switch expr1 := this.StmtNode.Partition.Expr.(type) {
			case *ast.ColumnNameExpr:
				this.PartitionColumns = append(this.PartitionColumns, expr1.Name.String())
			case *ast.FuncCallExpr:
				for _, arg := range expr1.Args {
					switch expr2 := arg.(type) {
					case *ast.ColumnNameExpr:
						this.PartitionColumns = append(this.PartitionColumns, expr2.Name.String())
					default:
						reviewMSG = new(ReviewMSG)
						reviewMSG.MSG = fmt.Sprintf("接测分区表错误. 不能识别指定的分区字段类型, " +
							"请联系DBA. 第二层: %T", this.StmtNode.Partition.Expr)
						return reviewMSG
					}
				}
			default:
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("接测分区表错误. 不能识别指定的分区字段类型, " +
					"请联系DBA. 第一层: %T", this.StmtNode.Partition.Expr)
				return reviewMSG
			}
		}

		// 没有发现 partition相关字段
		if len(this.PartitionColumns) == 0 {
			reviewMSG = new(ReviewMSG)
			reviewMSG.MSG = fmt.Sprintf("检测失败. 表: %v. 是分区表, 却没有发现分区字段.",
				this.StmtNode.Table.Name.String())
			return reviewMSG
		}

		// 检测分区表中的字段必须包含在字段中
		for _, partitionColumnName := range this.PartitionColumns {
			if _, ok := this.ColumnNames[partitionColumnName]; !ok {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败. 表: %v. %v 分区字段没有定义",
					this.StmtNode.Table.Name.String(), partitionColumnName)
				return reviewMSG
			}
		}

		// 所有的唯一键中必须包含所有的分区键
		if len(this.UniqueIndexes) > 0 {
			// 获取索引字段 hash 组合
			hashUniqueIndexes := GetIndexesHashColumn(this.UniqueIndexes)
			partitionHashNames := GetHashNames(this.PartitionColumns)
			for uniqueIndexName, hashUniqueIndex := range hashUniqueIndexes {
				if uniqueIndexName == this.PKname { // 过滤主键
					continue
				}
				for _, partitionHashName := range partitionHashNames {
					if !common.StrIsMatch(hashUniqueIndex, partitionHashName) {
						reviewMSG = new(ReviewMSG)
						reviewMSG.MSG = fmt.Sprintf("检测失败. 表: %v. 唯一索引: %v 没有包含分区字段: %v",
							this.StmtNode.Table.Name.String(), uniqueIndexName, this.PartitionColumns)
						return reviewMSG
					}
				}
			}
		} else { // 主键中需要包含分区的所有字段
			for _, partitionColumnName := range this.PartitionColumns {
				if _, ok := this.PKColumnNames[partitionColumnName]; !ok {
					reviewMSG = new(ReviewMSG)
					reviewMSG.MSG = fmt.Sprintf("检测失败. 表: %v. 主键没有包含分区字段: %v",
						this.StmtNode.Table.Name.String(), partitionColumnName, this.PartitionColumns)
					return reviewMSG
				}
			}
		}
	}

	return reviewMSG
}

// 链接指定实例检测相关表信息
func (this *CreateTableReviewer) DetectInstanceTable() *ReviewMSG {
	var reviewMSG *ReviewMSG

	tableInfo := dao.NewTableInfo(this.DBConfig, this.StmtNode.Table.Name.String())
	err := tableInfo.OpenInstance()
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 无法链接到指定实例. 无法检测数据库是否存在.")
		return reviewMSG
	}

	// 检测表是否存在
	reviewMSG = DetectTableExistsByName(tableInfo, this.SchemaName, this.StmtNode.Table.Name.String())
	if reviewMSG != nil {
		return CloseTableInstance(reviewMSG, tableInfo)
	}

	return CloseTableInstance(reviewMSG, tableInfo)
}

/* 检测是否有重复索引
Params:
    _tableInfo: 原表信息
 */
func (this *CreateTableReviewer) DetectDuplecateIndex() *ReviewMSG {
	var reviewMSG *ReviewMSG

	hashNormalIndex := GetIndexesHashColumn(this.Indexes)

	// 循环 唯一键和 索引进行匹配, 看看唯一索引的字段是否都包含在普通索引中
	for normalIndexName1, hashNormalIndexStr1 := range hashNormalIndex {

		for normalIndexName2, hashNormalIndexStr2 := range hashNormalIndex {
			if normalIndexName1 == normalIndexName2 { // 同一个索引不进行比较
				continue
			}
			if isMatch := common.StrIsMatch(hashNormalIndexStr1, hashNormalIndexStr2) ; isMatch {
				reviewMSG = new(ReviewMSG)
				reviewMSG.MSG = fmt.Sprintf("检测失败. 检测到重复索引: %v <=> %v.",
					normalIndexName1, normalIndexName2)
				return reviewMSG
			}
		}
	}

	return reviewMSG
}

/* 检测索引个数是否超过指定个数
Params:
    _tableInfo: 原表信息
 */
func (this *CreateTableReviewer) DetectIndexCount() *ReviewMSG {
	var reviewMSG *ReviewMSG

	if len(this.Indexes) > this.ReviewConfig.RuleIndexCount {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf("检测失败. %v",
			fmt.Sprintf(config.MSG_INDEX_COUNT_ERROR, this.ReviewConfig.RuleIndexCount))
		return reviewMSG
	}

	return reviewMSG
}
