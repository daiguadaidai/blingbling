# blingbling

MySQL SQL 解析审核工具

没错!!! 我的名字就是闪闪发亮的那个`blingbling`

- [安装](#安装)
- [启动和参数介绍](#启动和参数介绍)
    - [--help](#--help)
    - [启动](#启动)
    - [检测是否启动](#检测是否启动)
- [客户端使用](#客户端使用)
    - [可以指定的参数](#可以指定的参数)
    - [CURL访问](#CURL访问)
    - [Python客户端使用](#Python客户端使用)
    - [Golang客户端使用](#Golang客户端使用)
    - [Jquery-Ajax客户端](#Jquery-Ajax客户端)
    - [VUE-RESOURCE客户端](#VUE-RESOURCE客户端)
    - [axios客户端](#axios客户端)
- [高级玩法](#高级玩法)
- [获取元数据](#获取元数据)
    - [Python获取元数据](#Python获取元数据)
- [supervisor管理进程](#supervisor管理进程)
    - [基本使用方法](#基本使用方法)
    - [主要的配置](主要的配置)
    - [启动supervisor](#启动supervisor)
    - [小兴趣](#小兴趣)

## 安装

1. 下载 blingbling

```
git clone git@github.com:daiguadaidai/blingbling.git
```

2. 编译生成`go`语法解析文件

我使用的`Golang`版本是:`1.9.2`

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
      --listen-port int                       启动服务使用的端口 (default 19527)
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
    --listen-port=19527 \
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

检测启动方法可以通过 `ps -ef | grep blingbling` 查看, 也可以使用`netstat -natpl | grep 19527`查看

```
# ps 方法查看
ps -ef | grep blingbling
hh       11943  7909  0 20:15 pts/3    00:00:00 ./blingbling

# netstat 方法查看
netstat -natpl | grep 19527
tcp6       0      0 :::19527                :::*                    LISTEN      11943/blingbling
``` 

## 客户端使用

各种程序是如何使用的, 在下面都有实例. 相关的程序在项目中的`client_sample`目录里面都能够找的到.

使用客户端访问时. 服务端会以JSON的格式返回相关审核结果. 如下返回结果

```
{
    "Code": 0,
    "ReviewMSGs": [
        {
            "Sql": "alter table employees add column age1 int not null;",
            "HaveError": true,
            "HaveWarning": false,
            "ErrorMSGs":[
                "检测失败. 字段必须要有注释. alter add 字段: age1 "
            ],
            "WarningMSGs":[]
        }, {
            "Sql": " delete from employees WHERE id = 1;",
            "HaveError": true,
            "HaveWarning": false,
            "ErrorMSGs": [
                "检测失败. 执行explain sql获取sql影响行数失败: 执行explain失败: 10.10.10.21:3307. explain select * from  employees where id = 1; Error 1054: Unknown column 'id' in 'where clause'"
            ],
            "WarningMSGs": []
        }
    ]
}
```

**参数解释**

1. **Code:** 有3中值

    - **0:** 审核程序成功执行

    - **1:** 审核程序警告

    - **2:** 审核程序失败

2. **ReviewMSGs:** 所有审核程序, 当输入多个sql语句是, 则有多个审核消息.

    - **Sql:** 审核的sql

    - **HaveError:** 是否有错误

    - **HaveWarning:** 是否有警告

    - **ErrorMSGs:** 错误的消息(是一个数组)

    - **WarningMSGs:** 警告的消息(是一个数组)
    
### 可以指定的参数

通过访问 `http://127.0.0.1:19527/ClientParams` 可以获取客户端可以指定的参数

```
curl http://127.0.0.1:19527/ClientParams

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

### CURL访问

1. POST 请求

```
curl -X POST http://10.10.10.55:19527/sqlReview -d '{"Host":"10.10.10.21", "Port":3307, "Username":"root", "Password":"root", "Database":"employees", "Sqls":"alter table employees add column age1 int not null; delete from employees WHERE id = 1;"}'
{
    "Code": 0,
    "ReviewMSGs": [
        {
            "Sql": "alter table employees add column age1 int not null;",
            "HaveError": true,
            "HaveWarning": false,
            "ErrorMSGs":[
                "检测失败. 字段必须要有注释. alter add 字段: age1 "
            ],
            "WarningMSGs":[]
        }, {
            "Sql": " delete from employees WHERE id = 1;",
            "HaveError": true,
            "HaveWarning": false,
            "ErrorMSGs": [
                "检测失败. 执行explain sql获取sql影响行数失败: 执行explain失败: 10.10.10.21:3307. explain select * from  employees where id = 1; Error 1054: Unknown column 'id' in 'where clause'"
            ],
            "WarningMSGs": []
        }
    ]
}
```

2. GET 请求

由于在 GET 方法中使用分号(`;`), 一次性输入多个`sql`会导致`url`参数解析错误. 所以暂时就一个个来吧

```
curl "http://10.10.10.55:19527/sqlReview?Host=10.10.10.21&Port=3307&Username=HH&Password=oracle12&Database=employees&Sqls=alter%20table%20employees%20add%20column%20age1%20int%20not%20null"
{
    "Code": 0,
    "ReviewMSGs": [
        {
            "Sql": "alter table employees add column age1 int not null",
            "HaveError": true,
            "HaveWarning": false,
            "ErrorMSGs": [
                "检测失败. 字段必须要有注释. alter add 字段: age1 "
            ],
            "WarningMSGs": []
        }
    ]
}
```

### Python客户端使用

只演示使用POST的方法

```
import requests

data = {
    'Host': '10.10.10.21',
    'Port': 3307,
    'Username': 'root',
    'Password': 'root',
    'Database': 'employees',
    'Sqls': 'alter table employees add column age1 int not null; delete from employees WHERE id = 1;',
}

url = 'http://10.10.10.55:19527/sqlReview'

r = requests.post(url, json = data)

print(r.text)
```

### Golang客户端使用

只演示使用POST的方法

```
package main

import (
    "encoding/json"
    "fmt"
    "bytes"
    "net/http"
    "github.com/daiguadaidai/blingbling/reviewer"
    "github.com/liudng/godump"
)

func main() {
    // 
    params := make(map[string]interface{})
    params["Host"] = "10.10.10.21"
    params["Port"] = 3307
    params["Username"] = "root"
    params["Password"] = "root"
    params["Database"] = "employees"
    params["Sqls"] = "alter table employees add column age1 int not null; delete from employees WHERE id = 1;"

    // 将参数转化为Json
    jsonParams, err := json.Marshal(params)
    if err != nil {
        fmt.Println(err.Error() )
        return
    }

    // 设置POST请求参数
    reader := bytes.NewReader(jsonParams)
    url := "http://10.10.10.55:19527/sqlReview"
    request, err := http.NewRequest("POST", url, reader)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    request.Header.Set("Content-Type", "application/json;charset=UTF-8")

    // 执行
    client := http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // 获取指定返回值
    respData := new(reviewer.ResponseReviewData)
    err = json.NewDecoder(resp.Body).Decode(respData)
    if err != nil {
        fmt.Printf("json decode err: %v", err)
        return
    }

    godump.Dump(respData)
}
```

### Jquery-Ajax客户端

只演示`POST`请求

```
<html>
<head>
    <script type="text/javascript" src="../js/jquery-3.3.1.min.js"></script>
    <script type="text/javascript">
        $(document).ready(function(){
            var url = 'http://10.10.10.55:19527/sqlReview';
            var data = {
                Host: '10.10.10.21',
                Port: 3307,
                Username: 'root',
                Password: 'root',
                Database: 'employees',
                Sqls: 'alter table employees add column age1 int not null; delete from employees WHERE id = 1;'
            };
            $.ajax({
                url: url,
                contentType: "application/json; charset=utf-8",
                type: 'POST',
                data: JSON.stringify(data),
                dataType: 'json',
                success: function(data, textStatus, jqXHR){
                    console.log(data);
                    console.log(textStatus);
                    console.log(jqXHR);
                },
                error: function(xhr,textStatus){
                    console.log('错误');
                    console.log(xhr);
                    console.log(textStatus);
                },
            })
        });
    </script>
</head>

<body>
</body>

</html>
```

### VUE-RESOURCE客户端

只演示`POST`请求

```
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <script type="text/javascript" src="../js/jquery-3.3.1.min.js"></script>
    <script type="text/javascript" src="../js/vue-2.1.8.min.js"></script>
    <script type="text/javascript" src="../js/vue-resource-1.5.1.min.js"></script>
</head>
<body>

<div id="app">
    <input type="submit" value="sql" @click="sqlReview()">
    <hr />
    {{ reviewData }}
</div>

<script>
    new Vue({
        el: '#app',
        data: {
            reviewData: ""
        },
        methods: {
            sqlReview: function() {
                var url = 'http://10.10.10.55:19527/sqlReview';
                var data = {
                    Host: '10.10.10.21',
                    Port: 3307,
                    Username: 'root',
                    Password: 'root',
                    Database: 'employees',
                    Sqls: 'alter table employees add column age1 int not null; delete from employees WHERE id = 1;'
                };
                this.$http.post(
                    url,
                    JSON.stringify(data),
                    {emulateJSON: true}
                ).then(
                    (response)=>{
                        console.log(response);
                        this.reviewData = response.data;
                    },
                    (error)=>{
                        console.log(error);
                    }
                );
            }
        }
    })
</script>
</body>
</html>
```

### axios客户端

只演示`POST`请求

```
<html>
<head>
    <script type="text/javascript" src="../js/jquery-3.3.1.min.js"></script>
    <script type="text/javascript" src="../js/axios.min.js"></script>
</head>

<body>

<div>
    <button id="sql-review">审核sql</button>
    <hr />
    <div id="review-data"></div>
</div>

</body>

<script type="text/javascript">
    $(document).ready(function(){
        $("#sql-review").click(function() {
            sqlReview();
        });
    });

    function sqlReview() {
        var url = 'http://10.10.10.55:19527/sqlReview';
        var reviewData = {
            Host: '10.10.10.21',
            Port: 3307,
            Username: 'root',
            Password: 'root',
            Database: 'employees',
            Sqls: 'alter table employees add column age1 int not null; delete from employees WHERE id = 1;'
        };
        axios({
            headers: {
                'Content-Type': 'application/json'
            },
            transformRequest: [function(data) {
                data = JSON.stringify(data)
                return data
            }],
            url: url,
            method: 'post',
            params: {},
            data: reviewData
        }).then(function (response) {
            $("#review-data").html(JSON.stringify(response.data));
            console.log(response.data);
        }).catch(function (error) {
            console.log(error);
        });
    }
</script>
</html>
```

## 高级玩法

其实也没有那么高级, 无非就是搞搞自定义参数

下面就以 `python` 的使用方法来做演示. 其他客户端使用的方法大同小异.

演示添加字段, 必须以下划线(`_`)开头, 其他的按正则的规范.

```
import requests

data = {
    'Host': '10.10.10.21',
    'Port': 3307,
    'Username': 'root',
    'Password': 'root',
    'Database': 'employees',
    'Sqls': 'alter table employees add column age1 int not null comment "年龄"',
    'CustomRuleNameReg': True,
    'RuleNameReg': '^_[a-z\$_][a-z\$\d_]*$',
}

url = 'http://10.10.10.55:19527/sqlReview'

r = requests.post(url, json = data)

print(r.text)
# 输出
{"Code":0,"ReviewMSGs":[{"Sql":"alter table employees add column age1 int not null comment \"年龄\"","HaveError":true,"HaveWarning":false,"ErrorMSGs":["字段名 检测失败. 命名规则: ^_[a-z\\$_][a-z\\$\\d_]*$. 名称: age1, "],"WarningMSGs":[]}]}


# 将sql语句该为以下划线(_)命名的将审核通过
data['Sqls'] = 'alter table employees add column _age1 int not null comment "年龄"'
r = requests.post(url, json = data)
print(r.text)
# 输出
{"Code":0,"ReviewMSGs":[{"Sql":"alter table employees add column _age1 int not null comment \"年龄\"","HaveError":false,"HaveWarning":true,"ErrorMSGs":[],"WarningMSGs":["警告: 检测目标实例的数据库是否存在出错. Error 1045: Access denied for user 'root'@'10.10.10.55' (using password: YES)"]}]}
```

> **注意:**上面`CustomRuleNameReg`, `RuleNameReg`这两个参数--必须--是同时出现的.
> 千万别只出现了`CustomRuleNameReg=True`而`RuleNameReg`不设置值. 这样的后果是会使用`Golang`的字符串的默认值. 审核的时候将会遇到奇葩结果.

## 获取元数据

**接口:** `http://10.10.10.55:19527/meta`

可以通过该接口获取到sql的相关原数据.

我们这边只演示`Python`的方法来访问该接口

### Python获取元数据

```
#!/usr/bin/evn python3
#-*- coding:utf-8 -*-

import requests

sqls = '''
CREATE TABLE test.t1 (
  id bigint(18) NOT NULL AUTO_INCREMENT COMMENT '主键',
  PRIMARY KEY (id),
  UNIQUE KEY udx_uid (dep, arr, flightNo, flightDate, cabin),
  Index idx_uptime (uptime),
) ENGINE=InnoDb  DEFAULT CHARSET=utF8 COLLATE=Utf8mb4_general_ci comment="你号";

ALTER TABLE test.t1
  Add COLUMN arr varchar(3) NOT NULL DEFAULT '' Comment '注释',
  DROP COLUMN name,
  ADD PRIMARY KEY (id, name),
  ADD INDEX idx_id_name(id, name),
  ADD UNIQUE INDEX idx_id_name(id, name),
  DROP PRIMARY KEY,
  DROP INDEX idx_id_name,
  MODIFY arr varchar(3) NOT NULL DEFAULT '' Comment '注释',
  CHANGE id id1 bigint(18) NOT NULL AUTO_INCREMENT COMMENT '主键',
  CHARSET='utf8mb4' ENGINE=innodb COMMENT="表注释",
  RENAME TO test2.t2,
  RENAME INDEX idx1 to idx2;
'''

data = {
    'sqls': sqls,
}

url = 'http://10.10.10.55:19527/meta'

r = requests.post(url, json = data)

print(r.text)
```

输出

```
{
    "Status": true,
    "Message": "",
    "Data": [
        {
            "type": "ct",
            "md": {
                "schema": "test",
                "table": "t1",
                "columns": [
                    {
                        "name": "id",
                        "type": "bigint(18)",
                        "default": "",
                        "comment": "主键",
                        "not_null": true,
                        "auto_increment": true
                    }
                ],
                "constraints": [
                    {
                        "name": "",
                        "column_names": [
                            "id"
                        ],
                        "type": "pk"
                    },
                    {
                        "name": "udx_uid",
                        "column_names": [
                            "dep",
                            "arr",
                            "flightNo",
                            "flightDate",
                            "cabin"
                        ],
                        "type": "uk"
                    },
                    {
                        "name": "idx_uptime",
                        "column_names": [
                            "uptime"
                        ],
                        "type": "idx"
                    }
                ],
                "if_not_exists": false,
                "engine": "InnoDb",
                "charset": "utF8",
                "collate": "Utf8mb4_general_ci",
                "comment": "你号",
                "auto_increment": ""
            }
        },
        {
            "type": "at",
            "md": {
                "schema": "test",
                "table": "t1",
                "ops": [
                    {
                        "type": "add_column",
                        "after": "",
                        "columns": [
                            {
                                "name": "arr",
                                "type": "varchar(3)",
                                "default": "",
                                "comment": "注释",
                                "not_null": true,
                                "auto_increment": false
                            }
                        ]
                    },
                    {
                        "Type": "drop_column",
                        "column_name": "name"
                    },
                    {
                        "type": "add_constraint",
                        "constraint": {
                            "name": "",
                            "column_names": [
                                "id",
                                "name"
                            ],
                            "type": "pk"
                        }
                    },
                    {
                        "type": "add_constraint",
                        "constraint": {
                            "name": "idx_id_name",
                            "column_names": [
                                "id",
                                "name"
                            ],
                            "type": "idx"
                        }
                    },
                    {
                        "type": "add_constraint",
                        "constraint": {
                            "name": "idx_id_name",
                            "column_names": [
                                "id",
                                "name"
                            ],
                            "type": "uk"
                        }
                    },
                    {
                        "type": "drop_constraint",
                        "constraint_type": "pk",
                        "name": ""
                    },
                    {
                        "type": "drop_constraint",
                        "constraint_type": "idx",
                        "name": "idx_id_name"
                    },
                    {
                        "type": "modify_column",
                        "after": "",
                        "columns": [
                            {
                                "name": "arr",
                                "type": "varchar(3)",
                                "default": "",
                                "comment": "注释",
                                "not_null": true,
                                "auto_increment": false
                            }
                        ]
                    },
                    {
                        "type": "change_column",
                        "old_name": "id",
                        "after": "",
                        "new_column": {
                            "name": "id1",
                            "type": "bigint(18)",
                            "default": "",
                            "comment": "主键",
                            "not_null": true,
                            "auto_increment": true
                        }
                    },
                    {
                        "type": "options",
                        "options": [
                            {
                                "type": "charset",
                                "value": "utf8mb4"
                            },
                            {
                                "type": "engine",
                                "value": "innodb"
                            },
                            {
                                "type": "comment",
                                "value": "表注释"
                            }
                        ]
                    },
                    {
                        "type": "rename_table",
                        "schema": "test2",
                        "table": "t2"
                    },
                    {
                        "type": "rename_index",
                        "old_name": "idx1",
                        "new_name": "idx2"
                    }
                ]
            }
        }
    ]
}
```

**返回数据的解释**

1. **Status:** `true`: 访问成功. `false`: 访问失败

2. **Message:** 返回信息

3. **Data:** 返回数据, 是一个数组. 数组中存放的就是元信息

    - **Type:** `ct`: 建表元数据. `at`: Alter table元数据

    - **md:** MetaData 元数据信息

4. **建表的元数据**

    - **schema:** 数据库

    - **table:** 表

    - **if_not_exists:** `true/false` 是否有 `if not exists` 子句

    - **engine:** 存储引擎

    - **charset:** 字符集

    - **collate:** 排序

    - **comment:** 注释

    - **auto_increment:** `AUTO_INTREMENT` 值

    - **columns:** 字段信息, 是一个数组

        - **name:** 字段名

        - **type:** 字段类型

        - **default:** 字段默认值

        - **comment:** 字段注释

        - **not_null:** `true/false`是否为 `NOT NULL`

        - **auto_increment:** `true/false`是否为 `AUTO_INCREMENT`

    - **constraints:** 约束信息. 是一个数组
        
        - **type:** 约束类型 `pk/uk/idx/fk/ft`. `pk`:主键. `uk`:唯一键. `idx`:普通索引. `fk`:外键. `ft`:全文索引

        - **name:** 约束名称
        
        - **column_names:** 约束字段名

5. **ALTER TABLE元数据**
    
    - **schema:** 数据库
    
    - **table:** 表

    - **ops:** 是一个数组保存了各种alter 元数据

    - **添加列:**

        - **type:** "add_column"
        
        - **after:** after值
        
        - **columns:** 字段信息, 是一个数组
        
            - **name:** 字段名
        
            - **type:** 字段类型
        
            - **default:** 字段默认值
        
            - **comment:** 字段注释
        
            - **not_null:** `true/false`是否为 `NOT NULL`
        
            - **auto_increment:** `true/false`是否为 `AUTO_INCREMENT`

    - **删除列:**

        - **Type:** "drop_column",

        - **column_name:** 删除的 column的名称

    - **添加约束:**

        - **type**: "add_constraint",
        
        - **constraint**: 约束
        
            - **name:** 约束名称
        
            - **column_names:** 字段名称, 是一个数组
        
            - **type:** 约束类型 `pk/idx/fk`

    - **删除约束:**

        - **type:** "drop_constraint",

        - **constraint_type:** 约束类型 `pk/idx/fk`

        - **name**: 删除的约束名称

    - **modify列:**

        - **type:** "modify_column"
        
        - **after:** after值
        
        - **columns:** 字段信息, 是一个数组
        
            - **name:** 字段名
        
            - **type:** 字段类型
        
            - **default:** 字段默认值
        
            - **comment:** 字段注释
        
            - **not_null:** `true/false`是否为 `NOT NULL`
        
            - **auto_increment:** `true/false`是否为 `AUTO_INCREMENT`

    - **change列:**

        - **type:** "change_column"
        
        - **after:** after值

        - **old_name:** 旧字段名

        - **new_column:** 字段信息
        
            - **name:** 新字段名
        
            - **type:** 字段类型
        
            - **default:** 字段默认值
        
            - **comment:** 字段注释
        
            - **not_null:** `true/false`是否为 `NOT NULL`
        
            - **auto_increment:** `true/false`是否为 `AUTO_INCREMENT`

    - **修改表名:**

        - **type:** "rename_table",

        - **schema:** 新数据库名

        - **table:** 新表名

    - **修改索引名:**

        - **type:** "rename_index",
        
        - **old_name:** 老索引名
        
        - **new_name:** 新索引名

    - **修改表级别的选项:**

        - **type:** "options",

        - **options:** 表的选项, 是一个数组

        - **type**: `engine/comment/charset/collate/auto_increment`,

        - **value:** 上面`type`对应的值
        
## supervisor管理进程

`supervisor`可以很好的管理你的进程. 自己就不必再写一个相关守护进程的东西了. 为他点给赞.

如何安装的我就不说了. 请自行去`百度`, `google`, `bing`

### 基本使用方法

先搞一波基本使用的语法, 在`supervisor`文件夹下面有一个配置文件

```
# 指定配置文件启动
supervisord -c /u01/supervisor/supervisor.conf
# 停止 supervisord
supervisorctl shutdown
# 重新加载配置文件
supervisorctl reload
# 启动所有进程
supervisorctl start all
# 停止所有进程
supervisorctl stop all
# 启动某个进程
supervisorctl start program-name
# 停止某个进程
supervisorctl stop program-name
# 重启所有进程或所有进程
supervisorctl restart all
supervisorctl reatart program-name
# 查看supervisord当前管理的所有进程的状态
supervisorctl status
```

### 主要的配置
    
下面列出了主要的一些配置, 配置的注释使用分号开头(`;`).

大家需要注意看的主要是`[program:blingbling]`这个模块的东西, 该模块主要是指定了管理程序的名称是什么, 这边我们写的程序名称是`blingbling`.

其他的参数主要大家主要应该要修改一些目录相关的东西

```
[unix_http_server]
file=/u01/supervisor/supervisor.sock   ; the path to the socket file

[inet_http_server]         ; inet (TCP) server disabled by default
port=127.0.0.1:9001        ; ip_address:port specifier, *:port for all iface

[supervisord]
logfile=/u01/supervisor/supervisord.log ; main log file; default $CWD/supervisord.log
logfile_maxbytes=50MB        ; max main logfile bytes b4 rotation; default 50MB
logfile_backups=10           ; # of main logfile backups; 0 means none, default 10
loglevel=info                ; log level; default info; others: debug,warn,trace
pidfile=/u01/supervisor/supervisord.pid ; supervisord pidfile; default supervisord.pid
nodaemon=false               ; start in foreground if true; default false
minfds=1024                  ; min. avail startup file descriptors; default 1024
minprocs=200                 ; min. avail process descriptors;default 200

[supervisorctl]
serverurl=unix:///u01/supervisor/supervisor.sock ; use a unix:// URL  for a unix socket

[program:blingbling]
command=/u01/supervisor/blingbling/bin/blingbling    ; the program (relative uses PATH, can take args)
process_name=%(program_name)s                        ; process_name expr (default %(program_name)s)
numprocs=1                                           ; number of processes copies to start (def 1)
directory=/u01/supervisor/blingbling/log             ; directory to cwd to before exec (def no cwd)
user=root                                            ; setuid to this UNIX account to run the program
autostart=true
autorestart=true
startsecs=5
startretries=3
redirect_stderr=true
stdout_logfile=/u01/supervisor/blingbling/log/gen_blingbling.log
stdout_logfile_maxbytes=100MB
stdout_logfile_backups=10
stderr_logfile=/u01/supervisor/blingbling/log/err_blingbling.log
stderr_logfile_maxbytes=100MB
stderr_logfile_backups=10
stopasgroup=true
```

### 启动supervisor

启动查看`supervisor`, 并且查看`blingbling`是否自动启动

```
ps -ef | grep blingbling

supervisord -c /u01/supervisor/supervisor.conf

ps -ef | grep supervisor
root      6406     1  0 22:42 ?        00:00:00 /usr/bin/python /usr/local/bin/supervisord -c /u01/supervisor/supervisor.conf
root      6407  6406  0 22:42 ?        00:00:00 /u01/supervisor/blingbling/bin/blingbling
```

上面可以看到`blingbling`被supervisor拉起来了. 可以观察进程`pid`

### 小兴趣

有兴趣的朋友可以将启动的`blingbling kill`掉. 再次查看`blingbling`的进程. 发现`blingbling重启了`.

要想停止`blingbling`可以使用`supervisorctl stop blingbling`命令
