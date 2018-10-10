package config

const (
	// ------------- 通用规则 --------------
	// 通用名称字符长度
	MSG_NAME_LENGTH_ERROR = "名称长度不能超过 %v"
	// 通用名称的规范
	MSG_NAME_REG_ERROR = "命名规则: 以 (字母/_/$) 开头, 之后可用字符 (字母/数字/_/$)"
	// 通用字符集
	MSG_CHARSET_ERROR = "使用的字符集只允许 %v"
	// 通用COLLATE
	MSG_COLLATE_ERROR = "使用的collate只允许 %v"
	// 禁用DROP DATABASE 操作
	MSG_ALLOW_DROP_DATABASE_ERROR = "禁止删除数据库"
	// 禁止DROP TABLE操作
	MSG_ALLOW_DROP_TABLE_ERROR = "禁止删除表"
	// 禁止 truncate 表
	MSG_ALLOW_TRUNCATE_TABLE_ERROR = "禁止 truncate 表操作"
	// 禁止 rename 表
	MSG_ALLOW_RENAME_TABLE_ERROR = "禁止 rename 表操作"
	// 允许的存储引擎
	MSG_TABLE_ENGINE_ERROR = "允许的存储引擎 %v"
	// 重复定义字段名
	MSG_TABLE_COLUMN_DUP_ERROR = "一个字段名字不能被定义多次"
	// 不允许的字段
	MSG_NOT_ALLOW_COLUMN_TYPE_ERROR = "这些字段不能使用: %v"
	// 表需要注释
	MSG_NEED_TABLE_COMMENT_ERROR = "新建表必须要有注释"
	// 字段需要注释
	MSG_NEED_COLUMN_COMMENT_ERROR = "字段必须要有注释"
	// 主键需要自增
	MSG_PK_AUTO_INCREMENT_ERROR = "主键字段中必须有一个字段有自增属性"
	// 必须有主键
	MSG_NEED_PK = "必须要有主键"
	// 索引字段个数
	MSG_INDEX_COLUMN_COUNT_ERROR = "索引字段不能超过 %v 个"
	// 表名命名规则
	MSG_TABLE_NAME_GRE_ERROR = "表名, 命名规则不符合规范. 规范为(正则): %v"
	// 所有名称命名规范
	MSG_INDEX_NAME_REG_ERROR = "索引名称命名不符合规范, 索引命名规范为(正则): %v"
	// 唯一索引命名规范
	MSG_UNIQUE_INDEX_NAME_REG_ERROR = `唯一索引名称命名不符合规范, 索引命名规范为(正则): %v`
	// 所有字段必须为not null
	MSG_ALL_COLUMN_NOT_NULL_ERROR = "所有字段都必须为NOT NULL"
	// 是否允许外键
	MSG_ALLOW_FOREIGN_KEY_ERROR = "不允许使用外键"
	// 是否允许有全文索引
	MSG_ALLOW_FULL_TEXT_ERROR = "不允许创建全文索引"
	// 必须为NOT NULL的类型
	MSG_NOT_NULL_COLUMN_TYPE_ERROR = "%v 这些字段类型必须为NOT NULL"
	// 必须为 NOT NULL 的字段名
	MSG_NOT_NULL_COLUMN_NAME_ERROR = "以下字段名必须为NOT NULL: %v"
	// text 字段允许个数
	MSG_TEXT_TYPE_COLUMN_COUNT_ERROR = "text/blob字段使用超过 %v 个"
	// 指定字段名必须有索引
	MSG_NEED_INDEX_COLUMN_NAME_ERROR = "必须指定索引的字段名有: %v"
	// 必须有的字段名
	MSG_HAVE_COLUMN_NAME_ERROR = "必须包含的字段名有: %v"
	// 必须要有默认值
	MSG_NEED_DEFAULT_VALUE_ERROR = "字段定义必须要有默认值"
	// 必须要有默认值的字段名
	MSG_NEED_DEFAULT_VALUE_NAME_ERROR = "必须包含默认值的字段名有: %v"
	// 不否允许删除字段
	MSG_ALLOW_DROP_COLUMN_ERROR = "不允许删除字段"
	// 是否允许使用 after 字句
	MSG_ALLOW_AFTER_CLAUSE_ERROR = "不允许使用after字句"
	// 是否允许 alter change
	MSG_ALLOW_CHANGE_COLUMN_ERROR = "不允许 change column 字段"
	// 是否允许删除索引
	MSG_ALLOW_DROP_INDEX_ERROR = "不允许删除索引"
	// 是否允许删除主键
	MSG_ALLOW_DROP_PRIMARY_KEY_ERROR = "不允许删除主键"
	// 是否允许重命名索引重命名索引
	MSG_ALLOW_RENAME_INDEX_ERROR = "不允许重命名索引"
	// 是否允许删除分区
	MSG_ALLOW_DROP_PARTITION_ERROR = "不允许删除分区"
	// 表的索引个数
	MSG_INDEX_COUNT_ERROR = "索引个数不能超过%v个"
	// 是否允许DELETE多个表
	MSG_ALLOW_DELETE_MANY_TABLE_ERROR = "DELETE语句不允许同时删除多个表的数据"
	// 是否允许DELETE 表关联语句
	MSG_ALLOW_DELETE_HAS_JOIN_ERROR = "DELETE语句不允许有关联操作"
	// 是否允许DELETE 使用子句
	MSG_ALLOW_DELETE_HAS_SUB_CLAUSE_ERROR = "DELETE语句不允许使用子句"
	// 是否允许DELETE 没有WHERE
	MSG_ALLOW_DELETE_NO_WHERE_ERROR = "DELETE语句必须有WHERE条件"
	// 是否允许DELETE limit
	MSG_ALLOW_DELETE_LIMIT_ERROR = "DELETE语句不允许使用limit"
	// DELETE 行数限制
	MSG_DELETE_LESS_THAN_ERROR = "DELETE数据的行数不能超过%v行"
	// 是否允许 UPDATE 表关联语句
	MSG_ALLOW_UPDATE_HAS_JOIN_ERROR = "UPDATE语句不允许有关联操作"
	// 是否允许 UPDATE 使用子句
	MSG_ALLOW_UPDATE_HAS_SUB_CLAUSE_ERROR = "UPDATE语句不允许使用子句"
	// 是否允许 UPDATE 没有WHERE
	MSG_ALLOW_UPDATE_NO_WHERE_ERROR = "UPDATE语句必须有WHERE条件"
	// 是否允许 UPDATE limit
	MSG_ALLOW_UPDATE_LIMIT_ERROR = "UPDATE语句不允许使用limit"
	// UPDATE 行数限制
	MSG_UPDATE_LESS_THAN_ERROR = "UPDATE数据的行数不能超过%v行"
	// 是否允许insert select
	MSG_ALLOW_INSERT_SELECT_ERROR = "不允许使用 insert select 语句"
	// insert每批数量
	MSG_INSERT_ROWS_ERROR = "每批 insert 行数超过指定行数: %v"
	// 是否允许不指定字段
	MSG_ALLOW_INSERT_NO_COLUMN_ERROR = "insert 必须明确指定字段名"
	// 是否允许 insert ignore
	MSG_ALLOW_INSERT_IGNORE_ERROR = "不允许 insert ignore 操作"
	// 是否允许 replace into
	MSG_ALLOW_INSERT_REPLIACE_ERROR = "不允许replace into 操作"
)
