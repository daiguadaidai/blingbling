package handle

import "github.com/daiguadaidai/blingbling/config"

type RequestReviewParam struct {
	config.ReviewConfig // 指定的数据库规则

	config.DBConfig // 链接数据库的参数

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
	// 索引允许长度
	CustomRuleIndexCharLength bool
	// 是否使用自定义, 是否字段允许 COLLATE
	CustomRuleAllowColumnCollate bool
	// 是否使用自定义, 是否字段允许 CHARSET
	CustomRuleAllowColumnCharset bool
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
	// 是否使用自定义, 索引长度
	if this.CustomRuleIndexCharLength {
		reviewConfig.RuleIndexCharLength = this.ReviewConfig.RuleIndexCharLength
	}
	// 是否使用自定义, 是否字段允许 COLLATE
	if this.CustomRuleAllowColumnCollate {
		reviewConfig.RuleAllowColumnCollate = this.ReviewConfig.RuleAllowColumnCollate
	}
	// 是否使用自定义, 是否字段允许 CHARSET
	if this.CustomRuleAllowColumnCharset {
		reviewConfig.RuleAllowColumnCharset = this.ReviewConfig.RuleAllowColumnCharset
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

func (this *RequestReviewParam) ClientParams() string {
	return `
    可选参数                           参数类型         干什么用的
    ------------------------ 需要审核的数据库相关参数 --------------------------
    Username                           string          数据库用户名
    Password                           string          数据库密码
    Database                           string          数据库名称
    Host                               string          数据库IP
    Port                               int             数据库端口

    ------------------------- 主角参数 ----------------------------------------
    Sqls                               String          需要审核的sql, 多个使用逗号(,)隔开

    ------------------------- 自定义审核规则参数 -------------------------------
    RuleNameLength                     int             通用名字长度
    RuleNameReg                        string          通用名字命名规则 正则规则: 以(字母/$/_)开头, 之后任意多个(字母/数字/_/$)
    RuleCharSet                        string          通用字符集检测
    RuleCollate                        string          通用 COLLATE
    RuleAllowCreateDatabase            bool            是否允许创建数据库
    RuleAllowDropDatabase              bool            是否允许删除数据库
    RuleAllowDropTable                 bool            是否允许删除表
    RuleAllowRenameTable               bool            是否允许 rename table
    RuleAllowTruncateTable             bool            是否允许 truncate table
    RuleTableEngine                    string          允许的存储引擎
    RuleNotAllowColumnType             string          不允许使用的字段
    RuleNeedTableComment               bool            表是否需要注释
    RuleNeedColumnComment              bool            字段需要有注释
    RulePKAutoIncrement                bool            主键自增
    RuleNeedPK                         bool            必须要要有主键
    RuleIndexColumnCount               int             索引字段个数
    RuleTableNameReg                   string          表名 命名规范
    RuleIndexNameReg                   string          索引命名规范
    RuleUniqueIndexNameReg             string          唯一所有命名规范
    RuleAllColumnNotNull               bool            所有字段都必须为 NOT NULL
    RuleAllowForeignKey                bool            是否允许使用外键
    RuleAllowFullText                  bool            是否允许有全文索引
    RuleNotNullColumnType              string          必须为NOT NULL的字段
    RuleNotNullColumnName              string          必须为NOT NULL 的字段名
    RuleTextTypeColumnCount            int             text字段允许使用个数
    RuleNeedIndexColumnName            string          必须有索引的字段名
    RuleHaveColumnName                 string          必须包含的字段名
    RuleNeedDefaultValue               bool            字段定义必须要有默认值
    RuleNeedDefaultValueName           string          必须有默认值的字段名字
    RuleAllowDropColumn                bool            是否允许删除字段
    RuleAllowAfterClause               bool            是否允许 after 子句
    RuleAllowChangeColumn              bool            是否允许 alter change 语句
    RuleAllowDropIndex                 bool            是否允许删除索引
    RuleAllowDropPrimaryKey            bool            是否允许删除主键
    RuleAllowRenameIndex               bool            是否重命名索引
    RuleAllowDropPartition             bool            是否允许删除分区
    RuleIndexCount                     int             表的索引个数
    RuleAllowDeleteManyTable           bool            是否允许DELETE多个表
    RuleAllowDeleteHasJoin             bool            是否允许DELETE 表关联语句
    RuleAllowDeleteHasSubClause        bool            是否允许DELETE 使用子句
    RuleAllowDeleteNoWhere             bool            是否允许DELETE 没有WHERE
    RuleAllowDeleteLimit               bool            是否允许 delete 使用 limit
    RuleDeleteLessThan                 int             DELETE 行数限制
    RuleAllowUpdateHasJoin             bool            是否允许 UPDATE 表关联语句
    RuleAllowUpdateHasSubClause        bool            是否允许 UPDATE 使用子句
    RuleAllowUpdateNoWhere             bool            是否允许 UPDATE 没有WHERE
    RuleAllowUpdateLimit               bool            是否允许 UPDATE 使用 limit
    RuleUpdateLessThan                 int             UPDATE 行数限制
    RuleAllowInsertSelect              bool            是否允许insert select
    RuleInsertRows                     int             insert每批数量
    RuleAllowInsertNoColumn            bool            是否允许不指定字段
    RuleAllowInsertIgnore              bool            是否允许 insert ignore
    RuleAllowInsertReplace             bool            是否允许 replace into
    RuleIndexCharLength                int             索引允许的长度. 默认:767
    RuleAllowColumnCollate             bool            是否字段允许 COLLATE
    RuleAllowColumnCharset             bool            是否字段允许 CHARSET

    ------------------------- 是否自定义, 自定义审核规则参数 -------------------------------
    CustomRuleNameLength               bool            是否自定义, 通用名字长度
    CustomRuleNameReg                  bool            是否自定义, 通用名字命名规则 正则规则: 以(字母/$/_)开头, 之后任意多个(字母/数字/_/$)
    CustomRuleCharSet                  bool            是否自定义, 通用字符集检测
    CustomRuleCollate                  bool            是否自定义, 通用 COLLATE
    CustomRuleAllowCreateDatabase      bool            是否自定义, 是否允许创建数据库
    CustomRuleAllowDropDatabase        bool            是否自定义, 是否允许删除数据库
    CustomRuleAllowDropTable           bool            是否自定义, 是否允许删除表
    CustomRuleAllowRenameTable         bool            是否自定义, 是否允许 rename table
    CustomRuleAllowTruncateTable       bool            是否自定义, 是否允许 truncate table
    CustomRuleTableEngine              bool            是否自定义, 允许的存储引擎
    CustomRuleNotAllowColumnType       bool            是否自定义, 不允许使用的字段
    CustomRuleNeedTableComment         bool            是否自定义, 表是否需要注释
    CustomRuleNeedColumnComment        bool            是否自定义, 字段需要有注释
    CustomRulePKAutoIncrement          bool            是否自定义, 主键自增
    CustomRuleNeedPK                   bool            是否自定义, 必须要要有主键
    CustomRuleIndexColumnCount         bool            是否自定义, 索引字段个数
    CustomRuleTableNameReg             bool            是否自定义, 表名 命名规范
    CustomRuleIndexNameReg             bool            是否自定义, 索引命名规范
    CustomRuleUniqueIndexNameReg       bool            是否自定义, 唯一所有命名规范
    CustomRuleAllColumnNotNull         bool            是否自定义, 所有字段都必须为 NOT NULL
    CustomRuleAllowForeignKey          bool            是否自定义, 是否允许使用外键
    CustomRuleAllowFullText            bool            是否自定义, 是否允许有全文索引
    CustomRuleNotNullColumnType        bool            是否自定义, 必须为NOT NULL的字段
    CustomRuleNotNullColumnName        bool            是否自定义, 必须为NOT NULL 的字段名
    CustomRuleTextTypeColumnCount      bool            是否自定义, text字段允许使用个数
    CustomRuleNeedIndexColumnName      bool            是否自定义, 必须有索引的字段名
    CustomRuleHaveColumnName           bool            是否自定义, 必须包含的字段名
    CustomRuleNeedDefaultValue         bool            是否自定义, 字段定义必须要有默认值
    CustomRuleNeedDefaultValueName     bool            是否自定义, 必须有默认值的字段名字
    CustomRuleAllowDropColumn          bool            是否自定义, 是否允许删除字段
    CustomRuleAllowAfterClause         bool            是否自定义, 是否允许 after 子句
    CustomRuleAllowChangeColumn        bool            是否自定义, 是否允许 alter change 语句
    CustomRuleAllowDropIndex           bool            是否自定义, 是否允许删除索引
    CustomRuleAllowDropPrimaryKey      bool            是否自定义, 是否允许删除主键
    CustomRuleAllowRenameIndex         bool            是否自定义, 是否重命名索引
    CustomRuleAllowDropPartition       bool            是否自定义, 是否允许删除分区
    CustomRuleIndexCount               bool            是否自定义, 表的索引个数
    CustomRuleAllowDeleteManyTable     bool            是否自定义, 是否允许DELETE多个表
    CustomRuleAllowDeleteHasJoin       bool            是否自定义, 是否允许DELETE 表关联语句
    CustomRuleAllowDeleteHasSubClause  bool            是否自定义, 是否允许DELETE 使用子句
    CustomRuleAllowDeleteNoWhere       bool            是否自定义, 是否允许DELETE 没有WHERE
    CustomRuleAllowDeleteLimit         bool            是否自定义, 是否允许 delete 使用 limit
    CustomRuleDeleteLessThan           bool            是否自定义, DELETE 行数限制
    CustomRuleAllowUpdateHasJoin       bool            是否自定义, 是否允许 UPDATE 表关联语句
    CustomRuleAllowUpdateHasSubClause  bool            是否自定义, 是否允许 UPDATE 使用子句
    CustomRuleAllowUpdateNoWhere       bool            是否自定义, 是否允许 UPDATE 没有WHERE
    CustomRuleAllowUpdateLimit         bool            是否自定义, 是否允许 UPDATE 使用 limit
    CustomRuleUpdateLessThan           bool            是否自定义, UPDATE 行数限制
    CustomRuleAllowInsertSelect        bool            是否自定义, 是否允许insert select
    CustomRuleInsertRows               bool            是否自定义, insert每批数量
    CustomRuleAllowInsertNoColumn      bool            是否自定义, 是否允许不指定字段
    CustomRuleAllowInsertIgnore        bool            是否自定义, 是否允许 insert ignore
    CustomRuleAllowInsertReplace       bool            是否自定义, 是否允许 replace boolo
    CustomRuleIndexCharLength          bool            是否自定义, 索引长度
    CustomRuleAllowColumnCollate       bool            是否使用自定义, 是否字段允许 COLLATE
    CustomRuleAllowColumnCharset       bool            是否使用自定义, 是否字段允许 CHARSET
`
}
