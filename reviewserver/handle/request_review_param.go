package handle

import "github.com/daiguadaidai/blingbling/config"

type RequestReviewParam struct {
	*config.ReviewConfig // 指定的数据库规则

	*config.DBConfig // 链接数据库的参数

	Sqls string // 用于审核的sql

	/////////////////////////////////////////////////////
	// 以下参数主要是为了识别是否使用自定义的规则, 都是bool类型
	/////////////////////////////////////////////////////
	// 是否使用自定义, 通用名字长度
	CustomRuleNameLength bool
	// 是否使用自定义, 通用名字命名规则 正则规则: 以(字母/$/_)开头, 之后任意多个(字母/数字/_/$)
	CustomRuleNameReg bool
	// 是否使用自定义, 通用字符集检测
	CustomRuleCharSet bool
	// 是否使用自定义, 通用 COLLATE
	CustomRuleCollate bool
	// 是否使用自定义, 是否允许创建数据库
	CustomRuleAllowCreateDatabase bool
	// 是否使用自定义, 是否允许删除数据库
	CustomRuleAllowDropDatabase bool
	// 是否使用自定义, 是否允许删除表
	CustomRuleAllowDropTable bool
	// 是否使用自定义, 是否允许 rename table
	CustomRuleAllowRenameTable bool
	// 是否使用自定义, 是否允许 truncate table
	CustomRuleAllowTruncateTable bool
	// 是否使用自定义, 允许的存储引擎
	CustomRuleTableEngine bool
	// 是否使用自定义, 不允许使用的字段
	CustomRuleNotAllowColumnType bool
	// 是否使用自定义, 表是否需要注释
	CustomRuleNeedTableComment bool
	// 是否使用自定义, 字段需要有注释
	CustomRuleNeedColumnComment bool
	// 是否使用自定义, 主键自增
	CustomRulePKAutoIncrement bool
	// 是否使用自定义, 是否使用自定义, 必须要要有主键
	CustomRuleNeedPK bool
	// 是否使用自定义, 索引字段个数
	CustomRuleIndexColumnCount bool
	// 是否使用自定义, 表名 命名规范
	CustomRuleTableNameReg bool
	// 是否使用自定义, 索引命名规范
	CustomRuleIndexNameReg bool
	// 是否使用自定义, 唯一所有命名规范
	CustomRuleUniqueIndexNameReg bool
	// 是否使用自定义, 所有字段都必须为 NOT NULL
	CustomRuleAllColumnNotNull bool
	// 是否使用自定义, 是否允许使用外键
	CustomRuleAllowForeignKey bool
	// 是否使用自定义, 是否允许有全文索引
	CustomRuleAllowFullText bool
	// 是否使用自定义, 必须为NOT NULL的字段
	CustomRuleNotNullColumnType bool
	// 是否使用自定义, 必须为NOT NULL 的字段名
	CustomRuleNotNullColumnName bool
	// 是否使用自定义, text字段允许使用个数
	CustomRuleTextTypeColumnCount bool
	// 是否使用自定义, 必须有索引的字段名
	CustomRuleNeedIndexColumnName bool
	// 是否使用自定义, 必须包含的字段名
	CustomRuleHaveColumnName bool
	// 是否使用自定义, 字段定义必须要有默认值
	CustomRuleNeedDefaultValue bool
	// 是否使用自定义, 必须有默认值的字段名字
	CustomRuleNeedDefaultValueName bool
	// 是否使用自定义, 是否允许删除字段
	CustomRuleAllowDropColumn bool
	// 是否使用自定义, 是否允许 after 子句
	CustomRuleAllowAfterClause bool
	// 是否使用自定义, 是否允许 alter change 语句
	CustomRuleAllowChangeColumn bool
	// 是否使用自定义, 是否允许删除索引
	CustomRuleAllowDropIndex bool
	// 是否使用自定义, 是否允许删除主键
	CustomRuleAllowDropPrimaryKey bool
	// 是否使用自定义, 是否重命名索引
	CustomRuleAllowRenameIndex bool
	// 是否使用自定义, 是否允许删除分区
	CustomRuleAllowDropPartition bool
	// 是否使用自定义, 表的索引个数
	CustomRuleIndexCount bool
	// 是否使用自定义, 是否允许DELETE多个表
	CustomRuleAllowDeleteManyTable bool
	// 是否使用自定义, 是否允许DELETE 表关联语句
	CustomRuleAllowDeleteHasJoin bool
	// 是否使用自定义, 是否允许DELETE 使用子句
	CustomRuleAllowDeleteHasSubClause bool
	// 是否使用自定义, 是否允许DELETE 没有WHERE
	CustomRuleAllowDeleteNoWhere bool
	// 是否使用自定义, 是否允许 delete 使用 limit
	CustomRuleAllowDeleteLimit bool
	// 是否使用自定义, DELETE 行数限制
	CustomRuleDeleteLessThan bool
	// 是否使用自定义, 是否允许 UPDATE 表关联语句
	CustomRuleAllowUpdateHasJoin bool
	// 是否使用自定义, 是否允许 UPDATE 使用子句
	CustomRuleAllowUpdateHasSubClause bool
	// 是否使用自定义, 是否允许 UPDATE 没有WHERE
	CustomRuleAllowUpdateNoWhere bool
	// 是否使用自定义, 是否允许 UPDATE 使用 limit
	CustomRuleAllowUpdateLimit bool
	// 是否使用自定义, UPDATE 行数限制
	CustomRuleUpdateLessThan bool
	// 是否使用自定义, 是否允许insert select
	CustomRuleAllowInsertSelect bool
	// 是否使用自定义, insert每批数量
	CustomRuleInsertRows bool
	// 是否使用自定义, 是否允许不指定字段
	CustomRuleAllowInsertNoColumn bool
	// 是否使用自定义, 是否允许 insert ignore
	CustomRuleAllowInsertIgnore bool
	// 是否使用自定义, 是否允许 replace boolo
	CustomRuleAllowInsertReplace bool
}

func (this *RequestReviewParam) GetReviewConfig() *config.ReviewConfig {
	reviewConfig := config.GetReviewConfig()

	// 是否使用自定义, 通用名字长度
	if this.CustomRuleNameLength {
		reviewConfig.RuleNameLength = this.ReviewConfig.RuleNameLength
	}
	// 是否使用自定义, 通用名字命名规则 正则规则: 以(字母/$/_)开头, 之后任意多个(字母/数字/_/$)
	if this.CustomRuleNameReg {
		reviewConfig.RuleNameReg = this.ReviewConfig.RuleNameReg
	}
	// 是否使用自定义, 通用字符集检测
	if this.CustomRuleCharSet {
		reviewConfig.RuleCharSet = this.ReviewConfig.RuleCharSet
	}
	// 是否使用自定义, 通用 COLLATE
	if this.CustomRuleCollate {
		reviewConfig.RuleCollate = this.ReviewConfig.RuleCollate
	}
	// 是否使用自定义, 是否允许创建数据库
	if this.CustomRuleAllowCreateDatabase {
		reviewConfig.RuleAllowCreateDatabase = this.ReviewConfig.RuleAllowCreateDatabase
	}
	// 是否使用自定义, 是否允许删除数据库
	if this.CustomRuleAllowDropDatabase {
		reviewConfig.RuleAllowDropDatabase = this.ReviewConfig.RuleAllowDropDatabase
	}
	// 是否使用自定义, 是否允许删除表
	if this.CustomRuleAllowDropTable {
		reviewConfig.RuleAllowDropTable = this.ReviewConfig.RuleAllowDropTable
	}
	// 是否使用自定义, 是否允许 rename table
	if this.CustomRuleAllowRenameTable {
		reviewConfig.RuleAllowRenameTable = this.ReviewConfig.RuleAllowRenameTable
	}
	// 是否使用自定义, 是否允许 truncate table
	if this.CustomRuleAllowTruncateTable {
		reviewConfig.RuleAllowTruncateTable = this.ReviewConfig.RuleAllowTruncateTable
	}
	// 是否使用自定义, 允许的存储引擎
	if this.CustomRuleTableEngine {
		reviewConfig.RuleTableEngine = this.ReviewConfig.RuleTableEngine
	}
	// 是否使用自定义, 不允许使用的字段
	if this.CustomRuleNotAllowColumnType {
		reviewConfig.RuleNotAllowColumnType = this.ReviewConfig.RuleNotAllowColumnType
	}
	// 是否使用自定义, 表是否需要注释
	if this.CustomRuleNeedTableComment {
		reviewConfig.RuleNeedTableComment = this.ReviewConfig.RuleNeedTableComment
	}
	// 是否使用自定义, 字段需要有注释
	if this.CustomRuleNeedColumnComment {
		reviewConfig.RuleNeedColumnComment = this.ReviewConfig.RuleNeedColumnComment
	}
	// 是否使用自定义, 主键自增
	if this.CustomRulePKAutoIncrement {
		reviewConfig.RulePKAutoIncrement = this.ReviewConfig.RulePKAutoIncrement
	}
	// 是否使用自定义, 是否使用自定义, 必须要要有主键
	if this.CustomRuleNeedPK {
		reviewConfig.RuleNeedPK = this.ReviewConfig.RuleNeedPK
	}
	// 是否使用自定义, 索引字段个数
	if this.CustomRuleIndexColumnCount {
		reviewConfig.RuleIndexColumnCount = this.ReviewConfig.RuleIndexColumnCount
	}
	// 是否使用自定义, 表名 命名规范
	if this.CustomRuleTableNameReg {
		reviewConfig.RuleTableNameReg = this.ReviewConfig.RuleTableNameReg
	}
	// 是否使用自定义, 索引命名规范
	if this.CustomRuleIndexNameReg {
		reviewConfig.RuleIndexNameReg = this.ReviewConfig.RuleIndexNameReg
	}
	// 是否使用自定义, 唯一所有命名规范
	if this.CustomRuleUniqueIndexNameReg {
		reviewConfig.RuleUniqueIndexNameReg = this.ReviewConfig.RuleUniqueIndexNameReg
	}
	// 是否使用自定义, 所有字段都必须为 NOT NULL
	if this.CustomRuleAllColumnNotNull {
		reviewConfig.RuleAllColumnNotNull = this.ReviewConfig.RuleAllColumnNotNull
	}
	// 是否使用自定义, 是否允许使用外键
	if this.CustomRuleAllowForeignKey {
		reviewConfig.RuleAllowForeignKey = this.ReviewConfig.RuleAllowForeignKey
	}
	// 是否使用自定义, 是否允许有全文索引
	if this.CustomRuleAllowFullText {
		reviewConfig.RuleAllowFullText = this.ReviewConfig.RuleAllowFullText
	}
	// 是否使用自定义, 必须为NOT NULL的字段
	if this.CustomRuleNotNullColumnType {
		reviewConfig.RuleNotNullColumnType = this.ReviewConfig.RuleNotNullColumnType
	}
	// 是否使用自定义, 必须为NOT NULL 的字段名
	if this.CustomRuleNotNullColumnName {
		reviewConfig.RuleNotNullColumnName = this.ReviewConfig.RuleNotNullColumnName
	}
	// 是否使用自定义, text字段允许使用个数
	if this.CustomRuleTextTypeColumnCount {
		reviewConfig.RuleTextTypeColumnCount = this.ReviewConfig.RuleTextTypeColumnCount
	}
	// 是否使用自定义, 必须有索引的字段名
	if this.CustomRuleNeedIndexColumnName {
		reviewConfig.RuleNeedIndexColumnName = this.ReviewConfig.RuleNeedIndexColumnName
	}
	// 是否使用自定义, 必须包含的字段名
	if this.CustomRuleHaveColumnName {
		reviewConfig.RuleHaveColumnName = this.ReviewConfig.RuleHaveColumnName
	}
	// 是否使用自定义, 字段定义必须要有默认值
	if this.CustomRuleNeedDefaultValue {
		reviewConfig.RuleNeedDefaultValue = this.ReviewConfig.RuleNeedDefaultValue
	}
	// 是否使用自定义, 必须有默认值的字段名字
	if this.CustomRuleNeedDefaultValueName {
		reviewConfig.RuleNeedIndexColumnName = this.ReviewConfig.RuleNeedIndexColumnName
	}
	// 是否使用自定义, 是否允许删除字段
	if this.CustomRuleAllowDropColumn {
		reviewConfig.RuleAllowDropColumn = this.ReviewConfig.RuleAllowDropColumn
	}
	// 是否使用自定义, 是否允许 after 子句
	if this.CustomRuleAllowAfterClause {
		reviewConfig.RuleAllowAfterClause = this.ReviewConfig.RuleAllowAfterClause
	}
	// 是否使用自定义, 是否允许 alter change 语句
	if this.CustomRuleAllowChangeColumn {
		reviewConfig.RuleAllowChangeColumn = this.ReviewConfig.RuleAllowChangeColumn
	}
	// 是否使用自定义, 是否允许删除索引
	if this.CustomRuleAllowDropIndex {
		reviewConfig.RuleAllowDropIndex = this.ReviewConfig.RuleAllowDropIndex
	}
	// 是否使用自定义, 是否允许删除主键
	if this.CustomRuleAllowDropPrimaryKey {
		reviewConfig.RuleAllowDropPrimaryKey = this.ReviewConfig.RuleAllowDropPrimaryKey
	}
	// 是否使用自定义, 是否重命名索引
	if this.CustomRuleAllowRenameIndex {
		reviewConfig.RuleAllowRenameIndex = this.ReviewConfig.RuleAllowRenameIndex
	}
	// 是否使用自定义, 是否允许删除分区
	if this.CustomRuleAllowDropPartition {
		reviewConfig.RuleAllowDropPartition = this.ReviewConfig.RuleAllowDropPartition
	}
	// 是否使用自定义, 表的索引个数
	if this.CustomRuleIndexCount {
		reviewConfig.RuleIndexCount = this.ReviewConfig.RuleIndexCount
	}
	// 是否使用自定义, 是否允许DELETE多个表
	if this.CustomRuleAllowDeleteManyTable {
		reviewConfig.RuleAllowDeleteManyTable = this.ReviewConfig.RuleAllowDeleteManyTable
	}
	// 是否使用自定义, 是否允许DELETE 表关联语句
	if this.CustomRuleAllowDeleteHasJoin {
		reviewConfig.RuleAllowDeleteHasJoin = this.ReviewConfig.RuleAllowDeleteHasJoin
	}
	// 是否使用自定义, 是否允许DELETE 使用子句
	if this.CustomRuleAllowDeleteHasSubClause {
		reviewConfig.RuleAllowDeleteHasSubClause = this.ReviewConfig.RuleAllowDeleteHasSubClause
	}
	// 是否使用自定义, 是否允许DELETE 没有WHERE
	if this.CustomRuleAllowDeleteNoWhere {
		reviewConfig.RuleAllowDeleteNoWhere = this.ReviewConfig.RuleAllowDeleteNoWhere
	}
	// 是否使用自定义, 是否允许 delete 使用 limit
	if this.CustomRuleAllowDeleteLimit {
		reviewConfig.RuleAllowDeleteLimit = this.ReviewConfig.RuleAllowDeleteLimit
	}
	// 是否使用自定义, DELETE 行数限制
	if this.CustomRuleDeleteLessThan {
		reviewConfig.RuleDeleteLessThan = this.ReviewConfig.RuleDeleteLessThan
	}
	// 是否使用自定义, 是否允许 UPDATE 表关联语句
	if this.CustomRuleAllowUpdateHasJoin {
		reviewConfig.RuleAllowUpdateHasJoin = this.ReviewConfig.RuleAllowUpdateHasJoin
	}
	// 是否使用自定义, 是否允许 UPDATE 使用子句
	if this.CustomRuleAllowUpdateHasSubClause {
		reviewConfig.RuleAllowUpdateHasSubClause = this.ReviewConfig.RuleAllowUpdateHasSubClause
	}
	// 是否使用自定义, 是否允许 UPDATE 没有WHERE
	if this.CustomRuleAllowUpdateNoWhere {
		reviewConfig.RuleAllowUpdateNoWhere = this.ReviewConfig.RuleAllowUpdateNoWhere
	}
	// 是否使用自定义, 是否允许 UPDATE 使用 limit
	if this.CustomRuleAllowUpdateLimit {
		reviewConfig.RuleAllowUpdateLimit = this.ReviewConfig.RuleAllowUpdateLimit
	}
	// 是否使用自定义, UPDATE 行数限制
	if this.CustomRuleUpdateLessThan {
		reviewConfig.RuleUpdateLessThan = this.ReviewConfig.RuleUpdateLessThan
	}
	// 是否使用自定义, 是否允许insert select
	if this.CustomRuleAllowInsertSelect {
		reviewConfig.RuleAllowInsertSelect = this.ReviewConfig.RuleAllowInsertSelect
	}
	// 是否使用自定义, insert每批数量
	if this.CustomRuleInsertRows {
		reviewConfig.RuleInsertRows = this.ReviewConfig.RuleInsertRows
	}
	// 是否使用自定义, 是否允许不指定字段
	if this.CustomRuleAllowInsertNoColumn {
		reviewConfig.RuleAllowInsertNoColumn = this.ReviewConfig.RuleAllowInsertNoColumn
	}
	// 是否使用自定义, 是否允许 insert ignore
	if this.CustomRuleAllowInsertIgnore {
		reviewConfig.RuleAllowInsertIgnore = this.ReviewConfig.RuleAllowInsertIgnore
	}
	// 是否使用自定义, 是否允许 replace boolo
	if this.CustomRuleAllowInsertReplace {
		reviewConfig.RuleAllowInsertReplace = this.ReviewConfig.RuleAllowInsertReplace
	}

	return reviewConfig
}

// 获取数据库链接信息
func (this *RequestReviewParam) GetDBConfig() *config.DBConfig {
	return config.NewDBConfig(
		this.Host,
		this.Port,
		this.Username,
		this.Password,
		this.Database,
	)
}
