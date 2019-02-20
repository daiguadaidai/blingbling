package reviewer

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"

	"github.com/daiguadaidai/blingbling/ast"
	"github.com/daiguadaidai/blingbling/common"
	"github.com/daiguadaidai/blingbling/config"
	"github.com/daiguadaidai/blingbling/dependency/mysql"
	"github.com/dlclark/regexp2"
	"github.com/juju/errors"
)

/* 检测名称长度是否合法
Params:
	_name: 需要检测的名字
*/
func DetectNameLength(_name string, _length int) (haveError bool, msg string) {
	if len(_name) > _length {
		haveError = true
		msg = fmt.Sprintf(
			"检测失败: %v. 名称: %v",
			fmt.Sprintf(config.MSG_NAME_LENGTH_ERROR, _length),
			_name,
		)
	}

	return
}

/* 检测名字是否合法
Params:
	_name: 需要检测的名字
*/
func DetectNameReg(_name string, _reg string) (haveError bool, msg string) {
	// 使用正则表达式匹配名称
	re := regexp2.MustCompile(_reg, 0)
	if isMatch, _ := re.MatchString(_name); !isMatch {
		haveError = true
		msg = fmt.Sprintf("检测失败. %v. 名称: %v, ",
			fmt.Sprintf(config.MSG_NAME_REG_ERROR, _reg),
			_name)
	}

	return
}

/* 检测数据库的字符集
Params:
    _charset: 需要审核的字符集
    _allowCharsetStr: 允许的字符集 字符串 "utf8,gbk,utf8mb4"
*/
func DetectCharset(_charset string, _allowCharsetStr string) (haveError bool, msg string) {
	if len(_charset) == 0 || len(_allowCharsetStr) == 0 {
		return
	}

	allowCharsets := strings.Split(_allowCharsetStr, ",") // 获取允许的字符集数组
	isMatch := false
	// 将需要检测的字符集 和 允许的字符集进行循环比较
	for _, allowCharset := range allowCharsets {
		if strings.ToLower(_charset) == allowCharset {
			isMatch = true
			break
		}
	}

	if !isMatch {
		haveError = true
		msg = fmt.Sprintf(
			"字符类型检测失败: %v",
			fmt.Sprintf(config.MSG_CHARSET_ERROR, _allowCharsetStr),
		)
	}

	return
}

/* 检测数据库的Collate
Params:
    _collate: 需要审核的字符集
    _allowCollateStr: 允许的 collate 字符串 "utf8_general_ci,utf8mb4_general_ci"
*/
func DetectCollate(_collate string, _allowCollateStr string) (haveError bool, msg string) {
	if len(_collate) == 0 || len(_allowCollateStr) == 0 {
		return
	}

	allowCollates := strings.Split(_allowCollateStr, ",") // 获取允许的Collate数组
	isMatch := false
	// 将需要检测的collate 和 允许的字符集进行循环比较
	for _, allowCollate := range allowCollates {
		if strings.ToLower(_collate) == allowCollate {
			isMatch = true
			break
		}
	}

	if !isMatch {
		haveError = true
		msg = fmt.Sprintf(
			"Collate 类型检测失败: %v",
			fmt.Sprintf(config.MSG_COLLATE_ERROR, _allowCollateStr),
		)
	}

	return
}

/* 检测数据库允许的存储引擎
Params:
    _engine: 需要审核的存储引擎
    _allowEngineStr: 允许的存储引擎
*/
func DetectEngine(_engine string, _allowEngineStr string) (haveError bool, msg string) {
	allowEngines := strings.Split(_allowEngineStr, ",") // 获取允许的存储引擎
	isMatch := false
	// 将需要检测的collate 和 允许的字符集进行循环比较
	for _, allowEngine := range allowEngines {
		if strings.ToLower(_engine) == allowEngine {
			isMatch = true
			break
		}
	}

	if !isMatch {
		haveError = true
		msg = fmt.Sprintf(
			"存储引擎 类型检测失败: %v",
			fmt.Sprintf(config.MSG_TABLE_ENGINE_ERROR, _allowEngineStr),
		)
	}

	return
}

/* 通过所有所以和所有的唯一索引获取所有的普通索引
Params:
    _indexes: 所有的索引
	_uniqueIndex: 所有的普通索引
*/
func GetNoUniqueIndexes(_indexes map[string][]string, _uniqueIndex map[string][]string) map[string][]string {
	normalIndexes := make(map[string][]string)

	for indexName, index := range _indexes {
		if _, ok := _uniqueIndex[indexName]; ok { // 过滤掉索引中的 唯一索引
			continue
		}

		normalIndex := make([]string, 0, 1)
		for _, columnName := range index {
			normalIndex = append(normalIndex, columnName)
		}

		normalIndexes[indexName] = normalIndex
	}

	return normalIndexes
}

/* A/B两个索引进行合并
Params:
    _indexesA: 索引A
    _indexesB: 索引B
*/
func CombineIndexes(_indexesA map[string][]string, _indexesB map[string][]string) map[string][]string {
	indexes := make(map[string][]string)

	for key, columns := range _indexesA {
		if len(columns) == 0 {
			continue
		}
		indexes[key] = columns
	}
	for key, columns := range _indexesB {
		if len(columns) == 0 {
			continue
		}
		indexes[key] = columns
	}

	return indexes
}

/* 将索引的字段转化成 hash过后的值
Params:
    _indexes: 需要转化的索引
*/
func GetIndexesHashColumn(_indexes map[string][]string) map[string]string {
	hashIndexes := make(map[string]string)

	for indexName, index := range _indexes {
		hashIndex := make([]string, 0, 1)
		for _, columnName := range index {
			data := []byte(columnName)
			has := md5.Sum(data)
			hashColumn := fmt.Sprintf("%x", has)
			hashIndex = append(hashIndex, hashColumn)
		}

		hashIndexes[indexName] = strings.Join(hashIndex, ",")
	}

	return hashIndexes
}

/* 将一个数组的东西做hash
Params:
    _names: 需要转化的索引
*/
func GetHashNames(_names []string) []string {
	hashIndex := make([]string, 0, 1)
	for _, columnName := range _names {
		data := []byte(columnName)
		has := md5.Sum(data)
		hashColumn := fmt.Sprintf("%x", has)
		hashIndex = append(hashIndex, hashColumn)
	}

	return hashIndex
}

/* 数据库存在返回错误
Params:
    _tableInfo: 库相关信息
    _dbName: 数据库名
*/
func DetectDatabaseExistsByName(_tableInfo *TableInfo, _dbName string) (haveError bool, msg string) {
	// 检测实例中数据库是否存在
	exists, err := _tableInfo.DatabaseExistsByName(_dbName)
	if err != nil {
		msg = fmt.Sprintf("警告: 检测目标实例的数据库是否存在出错. %v", err)
		return
	}
	if exists {
		haveError = true
		msg = fmt.Sprintf("检测失败: 目标数据库 %v 已经存在.", _dbName)
		return
	}

	return
}

/* 数据库不存在返回错误
Params:
    _tableInfo: 库相关信息
    _dbName: 数据库名
*/
func DetectDatabaseNotExistsByName(_tableInfo *TableInfo, _dbName string) (haveError bool, msg string) {
	// 检测实例中数据库是否存在
	exists, err := _tableInfo.DatabaseExistsByName(_dbName)
	if err != nil {
		msg := fmt.Sprintf("警告: 检测目标实例的数据库是否存在出错. %v", err)
		return haveError, msg
	}
	if !exists {
		haveError = true
		msg := fmt.Sprintf("检测失败: 目标数据库 %v 不存在.", _dbName)
		return haveError, msg
	}

	return
}

/* 表否存在返回错误
Params:
    _tableInfo: 表相关信息
    _dbName: 数据库名
    _tableName: 需要判断的表名
*/
func DetectTableExistsByName(_tableInfo *TableInfo, _dbName, _tableName string) (haveError bool, msg string) {
	exists, err := _tableInfo.TableExistsByName(_dbName, _tableName)
	if err != nil {
		msg = fmt.Sprintf("警告: 检测目标实例的表是否存在出错. %v", err)
		return
	}
	if exists {
		haveError = true
		msg = fmt.Sprintf("检测失败: 在数据库中表 %v 已经存在.", _tableName)
		return
	}

	return
}

/* 表不否存在返回错误
Params:
    _tableInfo: 表相关信息
    _dbName: 数据库名
    _tableName: 需要判断的表名
*/
func DetectTableNotExistsByName(_tableInfo *TableInfo, _dbName, _tableName string) (haveError bool, msg string) {
	exists, err := _tableInfo.TableExistsByName(_dbName, _tableName)
	if err != nil {
		msg = fmt.Sprintf("警告: 检测目标实例的表是否存在出错. %v", err)
		return
	}
	if !exists {
		haveError = true
		msg = fmt.Sprintf("检测失败: 在数据库中表 %v 不存在.", _tableName)
		return
	}

	return
}

/* 将 delete sql 转化称 explain select sql
Params:
    _deleteSql: 删除sql
*/
func GetExplainSelectSqlByDeleteSql(_deleteSql string) string {
	var explainSelectSql string

	lowerSql := strings.ToLower(_deleteSql)
	sqlItems := strings.Split(lowerSql, " from ")
	explainSelectSql = fmt.Sprintf("%v %v",
		"explain select * from ", strings.Join(sqlItems[1:], " from "))

	return explainSelectSql
}

/* 将 update sql 转化称 explain select sql
Params:
    _updateSql: 更新sql
    _setWhereCount: set字句中where关键字的个数
    _hasWhereClause: 是否有Where 子句
*/
func GetExplainSelectSqlByUpdateSql(
	_updateSql string,
	_setWhereCount int,
	_hasWhereClause bool,
) (string, error) {
	var explainSelectSql string
	var explainSelectSuffix string
	var explainSelectWhere string

	// 通过 set 分开
	setReg := regexp.MustCompile(`(?i)\sSET\s`)
	setSplitItems := setReg.Split(_updateSql, 2)
	if len(setSplitItems) != 2 {
		errMSG := fmt.Sprintf("多个set关键字, 无法将update语句变成explain select语句")
		return "", errors.New(errMSG)
	}

	// 生成 explain select 前缀
	updateSuffixReg := regexp.MustCompile(`(?i)^\s*UPDATE\s`)
	explainSelectSuffix = updateSuffixReg.ReplaceAllString(setSplitItems[0], "explain select * from ")

	if _hasWhereClause {
		// 生成 explain select where 子句
		whereReg := regexp.MustCompile(`(?i)\sWHERE\s`)
		whereItems := whereReg.Split(setSplitItems[1], _setWhereCount+2)
		explainSelectWhere = whereItems[len(whereItems)-1]
	}

	explainSelectSql = fmt.Sprintf("%v where %v", explainSelectSuffix, explainSelectWhere)

	return explainSelectSql, nil
}

/* 匹配是否是 create table like 语句
Params:
    _sql: 建表语句
*/
func IsCreateTableLikeStmt(_sql string) bool {
	// reg := fmt.Sprintf(`(?i)\s*CREATE\s*TABLE\s*[%v\w\d_]+\s*LIKE\s*[%v\w\d_;]+`, "`", "`")
	reg := fmt.Sprintf(`(?i)^\s*CREATE\s*TABLE\s*[0-9a-z_%s\.]+\s*LIKE\s*[0-9a-z_%s\.]+\s*;?\s*$`, "`")

	return common.StrIsMatch(_sql, reg)
}

const (
	SQL_TYPE_TINYINT_MAX_LENGTH    = 1
	SQL_TYPE_SMALLINT_MAX_LENGTH   = 2
	SQL_TYPE_MEDIUMINT_MAX_LENGTH  = 3
	SQL_TYPE_INT_MAX_LENGTH        = 4
	SQL_TYPE_BIGINT_MAX_LENGTH     = 8
	SQL_TYPE_FLOAT_MAX_LENGTH      = 4
	SQL_TYPE_DUOBLE_MAX_LENGTH     = 8
	SQL_TYPE_CHAR_MAX_LENGTH       = 255
	SQL_TYPE_VARCHAR_MAX_LENGTH    = 65535
	SQL_TYPE_TINYTEXT_MAX_LENGTH   = 255
	SQL_TYPE_TEXT_MAX_LENGTH       = 65535
	SQL_TYPE_MEDIUMTEXT_MAX_LENGTH = 16777215
	SQL_TYPE_LONGTEXT_MAX_LENGTH   = 4294967295
	SQL_TYPE_DATE_MAX_LENGTH       = 3
	SQL_TYPE_TIME_MAX_LENGTH       = 3
	SQL_TYPE_YEAR_MAX_LENGTH       = 1
	SQL_TYPE_TIMESTAMP_MAX_LENGTH  = 4
	SQL_TYPE_DATETIME_MAX_LENGTH   = 8
)

// 获取字段的定义长度
func GetColumnDefineCharLen(column *ast.ColumnDef) (int, error) {
	switch column.Tp.Tp {
	case mysql.TypeTiny:
		return SQL_TYPE_TINYINT_MAX_LENGTH, nil
	case mysql.TypeShort:
		return SQL_TYPE_SMALLINT_MAX_LENGTH, nil
	case mysql.TypeInt24:
		return SQL_TYPE_MEDIUMINT_MAX_LENGTH, nil
	case mysql.TypeLong:
		return SQL_TYPE_INT_MAX_LENGTH, nil
	case mysql.TypeLonglong:
		return SQL_TYPE_BIGINT_MAX_LENGTH, nil
	case mysql.TypeFloat:
		return SQL_TYPE_FLOAT_MAX_LENGTH, nil
	case mysql.TypeDouble:
		return SQL_TYPE_DUOBLE_MAX_LENGTH, nil
	case mysql.TypeTimestamp:
		return SQL_TYPE_TIMESTAMP_MAX_LENGTH, nil
	case mysql.TypeDatetime:
		return SQL_TYPE_DATETIME_MAX_LENGTH, nil
	case mysql.TypeYear:
		return SQL_TYPE_YEAR_MAX_LENGTH, nil
	case mysql.TypeDate, mysql.TypeNewDate:
		return SQL_TYPE_DATE_MAX_LENGTH, nil
	case mysql.TypeDuration:
		return SQL_TYPE_TIME_MAX_LENGTH, nil
	case mysql.TypeDecimal, mysql.TypeNewDecimal:
		len := column.Tp.Flen
		if column.Tp.Decimal > 0 {
			len += column.Tp.Decimal
		}
		return len, nil
	case mysql.TypeTinyBlob:
		return SQL_TYPE_TINYTEXT_MAX_LENGTH, nil
	case mysql.TypeBlob:
		return SQL_TYPE_TEXT_MAX_LENGTH, nil
	case mysql.TypeMediumBlob:
		return SQL_TYPE_MEDIUMTEXT_MAX_LENGTH, nil
	case mysql.TypeLongBlob:
		return SQL_TYPE_LONGTEXT_MAX_LENGTH, nil
	case mysql.TypeString, mysql.TypeVarchar:
		if column.Tp.Flen >= SQL_TYPE_VARCHAR_MAX_LENGTH {
			return -1, fmt.Errorf("字段产度")
		}
		return column.Tp.Flen, nil
	case mysql.TypeGeometry, mysql.TypeSet, mysql.TypeEnum, mysql.TypeJSON,
		mysql.TypeBit, mysql.TypeNull:
		return -1, fmt.Errorf("未知类型: %s. 字段名:", column.Tp, column.Name)
	}
	return -1, fmt.Errorf("未知类型: %s. 字段名:", column.Tp, column.Name)
}

// 获取所有字段的长度
func GetColumnsCharLen(nameLenMap map[string]int, names ...string) int {
	total := 0
	for _, name := range names {
		if len, ok := nameLenMap[name]; ok {
			total += len
		}
	}

	return total
}

// 是否是Blob类型, text/blob/longblob ...
func IsBlob(sqlType byte) bool {
	switch sqlType {
	case mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob, mysql.TypeBlob,
		mysql.TypeJSON, mysql.TypeGeometry:
		return true
	}

	return false
}
