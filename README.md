# blingbling

MySQL SQL 解析审核工具

[TOC]

## 安装

1. 下载 blingbling

```
git clone git@github.com:daiguadaidai/blingbling.git
```

2. 编译生成`go`语法解析文件

由于该工具用到了`yacc`, `lex`进行sql语法解析, 因此需要先使用相关工具解析相关语法成为`go`能识别的语法

```
# 进入软件目录
cd blingbling
# 编译
make
```

编译后在`bin`目录下会生成`goyacc`二进制文件. 并且在`parser`目录下会生成一系列的`go`文件.

> **注意:** 这些`go`文件不要去手动修改

3. 生成最终执行工具

在项目目录下执行

```
go build
```

之后会在项目目录下生成 blingbling 二进制文件

```
ll
... 省略 ...
-rwxrwxr-x  1 hh hh 19967515 Oct 10 17:21 blingbling*
... 省略 ...
```

## 启动和参数介绍

`blingbling`启动后会启动一个`HTTP`服务用于客户端审核使用.

当然`blingbling`使用方法和大多数软件一样, 可以使用`--help`来查看大致的使用方法

### --help

```
./blingbling --help
    一款SQL审核工具, 主要用于MySQL SQL 相关审核. 启动工具后会提供一个http接口为用户实时链接并且审核相关SQL.
    启动工具:
    ./blingbling \
        --rule-name-length=100 \
        --rule-name-reg="^[a-zA-Z\$_][a-zA-Z\$\d_]*$" \
        --rule-charset="utf8,utf8mb4" \
        --rule-collate="utf8_general_ci,utf8mb4_general_ci" \
        --rule-allow-drop-database=false \
        --rule-allow-drop-table=false \
        --rule-allow-rename-table=false \
        --rule-allow-truncate-table=false
        --rule-table-engine="innodb"

Usage:
  blingbling [flags]

Flags:
  -h, --help                                  help for blingbling
      --listen-host string                    通用名称匹配规则 (default "0.0.0.0")
      --listen-port int                       启动服务使用的端口 (default 18080)
      --rule-all-column-not-null              是否所有字段. 默认: false
      --rule-allow-after-clause               是否允许after子句. 默认: true (default true)
      --rule-allow-change-column              是否允许Alter Change子句. 默认: true (default true)
      --rule-allow-delete-has-join            是否允许DELETE语句中使用JOIN. 默认: false
      --rule-allow-delete-has-sub-clause      是否允许DELETE语句中使用子查询. 默认: false
      --rule-allow-delete-limit               是否允许DELETE语句使用LIMIT. 默认: false
      --rule-allow-delete-many-table          是否允许同时删除多个表数据. 默认: false
      --rule-allow-delete-no-where            是否允许DELETE没有WHERE条件. 默认: false
      --rule-allow-drop-column                是否允许删除字段. 默认: true (default true)
      --rule-allow-drop-database              是否允许删除数据库, 默认: false
      --rule-allow-drop-index                 是否允许删除索引. 默认: true (default true)
      --rule-allow-drop-partition             是否允许删除分区. 默认: true (default true)
      --rule-allow-drop-primary-key           是否允许删除主键. 默认: true (default true)
      --rule-allow-drop-table                 是否允许删除表, 默认: false
      --rule-allow-foreign-key                是否允许使用外键. 默认: false
      --rule-allow-full-text                  是否允许使用全文索引. 默认: false
      --rule-allow-insert-ignore              是否允许 INSERT IGNORE 语句. 默认: true (default true)
      --rule-allow-insert-no-column           是否允许 INSERT 不明确指定字段名. 默认: true (default true)
      --rule-allow-insert-replace             是否允许 INSERT REPLACE 语句. 默认: true (default true)
      --rule-allow-insert-select              是否允许 INSERT SELECT 语句. 默认: true (default true)
      --rule-allow-rename-index               是否允许重命名索引. 默认: true (default true)
      --rule-allow-rename-table               是否允许重命名表, 默认: true (default true)
      --rule-allow-truncate-table             是否允许truncate表, 默认: false
      --rule-allow-update-has-join            是否允许UPDATE语句中使用JOIN. 默认: false
      --rule-allow-update-has-sub-clause      是否允许UPDATE语句中使用子查询. 默认: false
      --rule-allow-update-limit               是否允许UPDATE语句使用LIMIT. 默认: false
      --rule-allow-update-no-where            是否允许UPDATE没有WHERE条件. 默认: false
      --rule-charset string                   通用允许的字符集, 默认(多个用逗号隔开) (default "utf8,utf8mb4")
      --rule-collate string                   通用允许的collate, 默认(多个用逗号隔开) (default "utf8_general_ci,utf8mb4_general_ci")
      --rule-delete-less-than int             允许一次性删除多少行数据. 使用explain计算出来 (default 10000)
      --rule-have-column-name string          必须要的字段, 默认(多个用逗号隔开)
      --rule-index-column-count int           索引允许字段个数 (default 5)
      --rule-index-count int                  表允许有几个索引 (default 15)
      --rule-index-name-reg string            索引名命名规范(正则) (default "^idx_[a-z\\$\\d_]*$")
      --rule-insert-rows int                  每批允许 insert 的行数 (default 1000)
      --rule-name-length int                  通用名称长度 (default 100)
      --rule-name-reg string                  通用名称匹配规则 (default "^[a-z\\$_][a-z\\$\\d_]*$")
      --rule-need-column-comment              字段是否需要注释 默认: true (default true)
      --rule-need-default-value               是否需要有默认字段. 默认: false
      --rule-need-default-value-name string   必须要有默认值的字段名, 默认(多个用逗号隔开) (default "created_at,updated_at,create_time,update_time,create_at,update_at,created_time,updated_time")
      --rule-need-index-column-name string    必须要有索引的字段名, 默认(多个用逗号隔开) (default "created_at,updated_at,create_time,update_time,create_at,update_at,created_time,updated_time")
      --rule-need-pk                          建表是否需要主键 默认: true (default true)
      --rule-need-table-comment               表是否需要注释 默认: true (default true)
      --rule-not-allow-column-type string     不允许的字段类型, 至此的类型: decimal, tinyint, smallint, int, float, double, timestamp, bigint, mediumint, date, time, datetime, year, newdate, varchar, bit, json, newdecimal, enum, set, tinyblob, mediumblob, longblob, blob, tinytext, mediumtext, longtext, text, geometry (default "tinytext,mediumtext,logtext,tinyblob,mediumblob,longblob")
      --rule-not-null-column-name string      必须为not null 的索引名, 默认(多个用逗号隔开) (default "created_at,updated_at,create_time,update_time,create_at,update_at,created_time,updated_time")
      --rule-not-null-column-type string      必须为not null的字段类型, 默认(多个用逗号隔开). 可填写的类型有: decimal, tinyint, smallint, int, float, double, timestamp, bigint, mediumint, date, time, datetime, year, newdate, varchar, bit, json, newdecimal, enum, set, tinyblob, mediumblob, longblob, blob, tinytext, mediumtext, longtext, text, geometry (default "varchar")
      --rule-pk-auto-increment                主键字段中是否需要有自增字段 默认: true (default true)
      --rule-table-engine string              允许的存储引擎 默认(多个用逗号隔开) (default "innodb")
      --rule-table-name-reg string            表名, 名命名规范(正则) (default "(?i)^(?!taishan)[a-z\\$_][a-z\\$\\d_]*$")
      --rule-text-type-column-count int       允许使用text/blob字段个数. 如果在rule-not-allow-column-type相关text字段.该参数将不其作用 (default 2)
      --rule-unique-index-name-reg string     唯一索引名命名规范(正则) (default "^udx_[a-z\\$\\d_]*$")
      --rule-update-less-than int             允许一次性删除多少行数据. 使用explain计算出来 (default 10000)
```

`blingbling`启动可选参数比较多, 主要是设置一些默认审核的规则. 我自己看眼睛都看花掉了 `- -.`

> 需要说明一下的是关于一些命名规则显示有两个反斜杠`\\`其实在使用正则的时候只有一个反斜杠`\`. 显示成这样主要是因为字符串输出的原因.
> 如: `^udx_[a-z\\$\\d_]*$` 你指定的时候应该是 `^udx_[a-z\$\d_]*$`

### 启动

这里我提供两种启动方法

1. 都使用默认参数

```
./blingbling
```

2. 手动指定参数

```
./blingbling \
    --listen-host=0.0.0.0 \
    --listen-port=18080 \
    --rule-name-length=100 \
    --rule-name-reg="^[a-zA-Z\$_][a-zA-Z\$\d_]*$" \
    --rule-charset="utf8,utf8mb4" \
    --rule-collate="utf8_general_ci,utf8mb4_general_ci" \
    --rule-allow-drop-database=false \
    --rule-allow-drop-table=false \
    --rule-allow-rename-table=false \
    --rule-allow-truncate-table=false
    --rule-table-engine="innodb"
    
... 更多的参数 ...
```

`--listen-host`, `--listen-port` 这两个参数是指定`HTTP`服务监听的`IP`和`端口`的.

需要使用其他参数自行再指定

### 检测是否启动

检测启动方法可以通过 `ps -ef | grep blingbling` 查看, 也可以使用`netstat -natpl | grep 18080`查看

```
# ps 方法查看
ps -ef | grep blingbling
hh       11943  7909  0 20:15 pts/3    00:00:00 ./blingbling

# netstat 方法查看
netstat -natpl | grep 18080
tcp6       0      0 :::18080                :::*                    LISTEN      11943/blingbling
``` 

## 客户端使用

### 可以指定的参数

通过访问 `http://127.0.0.1:18080/ClientParams` 可以获取客户端可以指定的参数

```
curl http://127.0.0.1:18080/ClientParams

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
```

可以将网址输入到流量器中查看

上面的参数: `Username`, `Password`, `Database`, `Host`, `Port`, `Sqls` 这几个参数按道理来说是必须填写的.

> **注意:** `自定义审核规则参数` 和 `是否自定义, 自定义审核规则参数` 是一一对应的, 如果需要指定自定义参数必须要指定对应的`Custom`开头的参数.
> 如: 我需要自定义 `RuleNameLength=1000` 那么必须指定 `CustomRuleNameLength=1000` 不然将视为无效.
> 这样使用自定参数是有点 `傻B` 但是感觉没办法 因为字符串默认值是:`""`, 数字是:`0`. 有些有值的参数. 在`new`一个对象的时候会被清空. 没有指定自定义值也会覆盖服务启动的默认值. 因此我这边使用显示指定是否使用自定义参数来做.

 
 
 
     
     
