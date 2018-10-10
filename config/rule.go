package config

const (
	// 名称长度
	RULE_NAME_LENGTH = 100
	// 通用名字命名规则 正则规则: 以(字母/$/_)开头, 之后任意多个(字母/数字/_/$)
	RULE_NAME_REG = `^[a-z\$_][a-z\$\d_]*$`
	// 通用字符集检测
	RULE_CHARSET = "utf8,utf8mb4"
	// 通用 COLLATE
	RULE_COLLATE = "utf8_general_ci,utf8mb4_general_ci"
	// 是否允许创建
	RULE_ALLOW_CREATE_DATABASE = false
	// 是否允许删除数据库
	RULE_ALLOW_DROP_DATABASE = false
	// 是否允许删除表
	RULE_ALLOW_DROP_TABLE = false
	// 是否允许 rename table
	RULE_ALLOW_RENAME_TABLE = true
	// 是否允许 truncate table
	RULE_ALLOW_TRUNCATE_TABLE = false
	// 建表允许的存储引擎, 多个以逗号隔开
	RULE_TABLE_ENGINE = "innodb"
	// 是否允许大字段: text, blob
	RULE_NOT_ALLOW_COLUMN_TYPE = "tinytext,mediumtext,logtext,tinyblob,mediumblob,longblob"
	// 表是否需要注释
	RULE_NEED_TABLE_COMMENT = true
	// 字段是否需要注释
	RULE_NEED_COLUMN_COMMENT = true
	// 主键需要有子增
	RULE_PK_AUTO_INCREMENT = true
	// 必须有主键
	RULE_NEED_PK = true
	// 索引字段个数
	RULE_INDEX_COLUMN_COUNT = 5
	// 表名  命名规范
	RULE_TABLE_NAME_GRE = `(?i)^(?!taishan)[a-z\$_][a-z\$\d_]*$`
	// 索引命名规范
	RULE_INDEX_NAME_REG = `^idx_[a-z\$\d_]*$`
	// 唯一索引命名规范
	RULE_UNIQUE_INDEX_NAME_REG = `^udx_[a-z\$\d_]*$`
	// 所有字段都 必须有 not null
	RULE_ALL_COLUMN_NOT_NULL = false
	// 是否允许外键
	RULE_ALLOW_FOREIGN_KEY = false
	// 是否允许有全文索引
	RULE_ALLOW_FULL_TEXT = false
	// 必须为NOT NULL 的类型
	RULE_NOT_NULL_COLUMN_TYPE = "varchar"
	// 必须为 NOT NULL 的字段名
	RULE_NOT_NULL_COLUMN_NAME = "created_at,updated_at,create_time,update_time,create_at,update_at,created_time,updated_time"
	// Text 字段类型允许使用个数
	RULE_TEXT_TYPE_COLUMN_COUNT = 2
	// 指定字段名必须有索引
	RULE_NEED_INDEX_COLUMN_NAME = "created_at,updated_at,create_time,update_time,create_at,update_at,created_time,updated_time"
	// 必须有的字段名
	RULE_HAVE_COLUMN_NAME = ""
	// 是否要有默认值
	RULE_NEED_DEFAULT_VALUE = false
	// 必须要有默认值的字段名
	RULE_NEED_DEFAULT_VALUE_NAME = "created_at,updated_at,create_time,update_time,create_at,update_at,created_time,updated_time"
	// 是否允许删除字段
	RULE_ALLOW_DROP_COLUMN = true
	// 是否允许 after 字句
	RULE_ALLOW_AFTER_CLAUSE = true
	// 是否允许 alter change
	RULE_ALLOW_CHANGE_COLUMN = true
	// 是否允许删除索引
	RULE_ALLOW_DROP_INDEX = true
	// 是否允许删除主键
	RULE_ALLOW_DROP_PRIMARY_KEY = true
	// 是否允许重命名索引
	RULE_ALLOW_RENAME_INDEX = true
	// 是否允许删除分区
	RULE_ALLOW_DROP_PARTITION = true
	// 一个表的索引个数
	RULE_INDEX_COUNT = 15
	// 是否允许DELETE多个表
	RULE_ALLOW_DELETE_MANY_TABLE = false
	// 是否允许DELETE 表关联语句
	RULE_ALLOW_DELETE_HAS_JOIN = false
	// 是否允许DELETE 使用子句
	RULE_ALLOW_DELETE_HAS_SUB_CLAUSE = false
	// 是否允许DELETE 没有WHERE
	RULE_ALLOW_DELETE_NO_WHERE = false
	// 是否允许DELETE limit
	RULE_ALLOW_DELETE_LIMIT = false
	// DELETE 行数限制
	RULE_DELETE_LESS_THAN = 10000
	// 是否允许 UPDATE 表关联语句
	RULE_ALLOW_UPDATE_HAS_JOIN = false
	// 是否允许 UPDATE 使用子句
	RULE_ALLOW_UPDATE_HAS_SUB_CLAUSE = false
	// 是否允许 UPDATE 没有WHERE
	RULE_ALLOW_UPDATE_NO_WHERE = false
	// 是否允许 UPDATE limit
	RULE_ALLOW_UPDATE_LIMIT = false
	// UPDATE 行数限制
	RULE_UPDATE_LESS_THAN = 10000
	// 是否允许insert select
	RULE_ALLOW_INSERT_SELECT = true
	// insert每批数量
	RULE_INSERT_ROWS = 1000
	// 是否允许不指定字段
	RULE_ALLOW_INSERT_NO_COLUMN = true
	// 是否允许 insert ignore
	RULE_ALLOW_INSERT_IGNORE = true
	// 是否允许 replace into
	RULE_ALLOW_INSERT_REPLIACE = true
)
