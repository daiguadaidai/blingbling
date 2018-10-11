package reviewer

import (
	"fmt"
	"regexp"
	"github.com/daiguadaidai/blingbling/config"
	"strings"
	"strconv"
	"github.com/dlclark/regexp2"
	"crypto/md5"
	"github.com/daiguadaidai/blingbling/dao"
	"github.com/juju/errors"
)

/* 检测名称长度是否合法
Params:
	_name: 需要检测的名字
 */
func DetectNameLength(_name string, _length int) *ReviewMSG {
	if len(_name) > _length {
		return &ReviewMSG{
			MSG: fmt.Sprintf(
				"检测失败: %v. 名称: %v",
				fmt.Sprintf(config.MSG_NAME_LENGTH_ERROR, _length),
				_name,
			),
		}
	}

	return nil
}

/* 检测名字是否合法
Params:
	_name: 需要检测的名字
 */
func DetectNameReg(_name string, _reg string) *ReviewMSG {
	// 使用正则表达式匹配名称
	re := regexp2.MustCompile(_reg, 0)
	if isMatch, _ := re.MatchString(_name); !isMatch {
		return &ReviewMSG{
			MSG: fmt.Sprintf("检测失败. %v. 名称: %v, ",
				fmt.Sprintf(config.MSG_NAME_REG_ERROR, _reg),
				_name),
		}
	}

	return nil
}

/* 检测数据库的字符集
Params:
    _charset: 需要审核的字符集
    _allowCharsetStr: 允许的字符集 字符串 "utf8,gbk,utf8mb4"
 */
func DetectCharset(_charset string, _allowCharsetStr string) *ReviewMSG {
	var reviewMSG *ReviewMSG

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
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf(
			"字符类型检测失败: %v",
			fmt.Sprintf(config.MSG_CHARSET_ERROR, _allowCharsetStr),
		)
	}

	return reviewMSG
}

/* 检测数据库的Collate
Params:
    _collate: 需要审核的字符集
    _allowCollateStr: 允许的 collate 字符串 "utf8_general_ci,utf8mb4_general_ci"
 */
func DetectCollate(_collate string, _allowCollateStr string) *ReviewMSG {
	var reviewMSG *ReviewMSG

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
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf(
			"Collate 类型检测失败: %v",
			fmt.Sprintf(config.MSG_COLLATE_ERROR, _allowCollateStr),
		)
	}

	return reviewMSG
}

/* 检测数据库允许的存储引擎
Params:
    _engine: 需要审核的存储引擎
    _allowEngineStr: 允许的存储引擎
 */
func DetectEngine(_engine string, _allowEngineStr string) *ReviewMSG {
	var reviewMSG *ReviewMSG

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
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf(
			"存储引擎 类型检测失败: %v",
			fmt.Sprintf(config.MSG_TABLE_ENGINE_ERROR, _allowEngineStr),
		)
	}

	return reviewMSG
}

/* 检测不允许的字段类型
Params:
    _type: 需要审核的字段类型
    _notAllowTypeStr: 不允许的字段类型
 */
func DetectNotAllowColumnType(_type byte, _notAllowTypeSrt string) *ReviewMSG {
	var reviewMSG *ReviewMSG

	notAllowTypes := strings.Split(_notAllowTypeSrt, ",") // 获取不允许的字段类型
	isMatch := false
	fmt.Println(strconv.Itoa(int(_type)), notAllowTypes)
	// 将需要检测的collate 和 允许的字符集进行循环比较
	for _, notAllowType := range notAllowTypes {
		if strconv.Itoa(int(_type)) == notAllowType {
			isMatch = true
			break
		}
	}

	if isMatch {
		reviewMSG = new(ReviewMSG)
		reviewMSG.MSG = fmt.Sprintf(
			"字段检测失败: %v",
			fmt.Sprintf(config.MSG_NOT_ALLOW_COLUMN_TYPE_ERROR, _notAllowTypeSrt),
		)
	}

	return reviewMSG
}

/* 通过所有所以和所有的唯一索引获取所有的普通索引
Params:
    _indexes: 所有的索引
	_uniqueIndex: 所有的普通索引
 */
func GetNoUniqueIndexes(_indexes map[string][]string, _uniqueIndex map[string][]string) map[string][]string{
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
func DetectDatabaseExistsByName(_tableInfo *dao.TableInfo, _dbName string) *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 检测实例中数据库是否存在
	exists, err := _tableInfo.DatabaseExistsByName(_dbName)
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 检测目标实例的数据库是否存在出错. %v", err)
		return reviewMSG
	}
	if exists {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_ERROR
		reviewMSG.MSG = fmt.Sprintf("检测失败: 目标数据库[%v]已经存在.",
			_tableInfo.DBName)
		return reviewMSG
	}

	return reviewMSG
}

/* 数据库不存在返回错误
Params:
    _tableInfo: 库相关信息
    _dbName: 数据库名
 */
func DetectDatabaseNotExistsByName(_tableInfo *dao.TableInfo, _dbName string) *ReviewMSG {
	var reviewMSG *ReviewMSG

	// 检测实例中数据库是否存在
	exists, err := _tableInfo.DatabaseExistsByName(_dbName)
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 检测目标实例的数据库是否存在出错. %v", err)
		return reviewMSG
	}
	if !exists {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_ERROR
		reviewMSG.MSG = fmt.Sprintf("检测失败: 目标数据库[%v]已经存在.",
			_tableInfo.DBName)
		return reviewMSG
	}

	return reviewMSG
}


/* 表否存在返回错误
Params:
    _tableInfo: 表相关信息
    _dbName: 数据库名
    _tableName: 需要判断的表名
 */
func DetectTableExistsByName(_tableInfo *dao.TableInfo, _dbName, _tableName string) *ReviewMSG {
	var reviewMSG *ReviewMSG

	exists, err := _tableInfo.TableExistsByName(_dbName, _tableName)
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 检测目标实例的表是否存在出错. %v", err)
		return reviewMSG
	}
	if exists {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_ERROR
		reviewMSG.MSG = fmt.Sprintf("检测失败: 在数据库中表[%v]已经存在.",
			_tableInfo.TableName)
		return reviewMSG
	}

	return reviewMSG
}

/* 表不否存在返回错误
Params:
    _tableInfo: 表相关信息
    _dbName: 数据库名
    _tableName: 需要判断的表名
 */
func DetectTableNotExistsByName(_tableInfo *dao.TableInfo, _dbName, _tableName string) *ReviewMSG {
	var reviewMSG *ReviewMSG

	exists, err := _tableInfo.TableExistsByName(_dbName, _tableName)
	if err != nil {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_WARNING
		reviewMSG.MSG = fmt.Sprintf("警告: 检测目标实例的表是否存在出错. %v", err)
		return reviewMSG
	}
	if !exists {
		reviewMSG = new(ReviewMSG)
		reviewMSG.Code = REVIEW_CODE_ERROR
		reviewMSG.MSG = fmt.Sprintf("检测失败: 在数据库中表[%v]已经存在.",
			_tableInfo.TableName)
		return reviewMSG
	}

	return reviewMSG
}

/* 关闭链接并返回相关审核信息
Params:
    _reviewMSG: 审核信息
    _tableInfo: 链接数据库的表
 */
func CloseTableInstance(_reviewMSG *ReviewMSG, _tableInfo *dao.TableInfo) *ReviewMSG {
	err := _tableInfo.CloseInstance()
	if err != nil {
		if _reviewMSG == nil {
			_reviewMSG = new(ReviewMSG)
		}
		_reviewMSG.Code = REVIEW_CODE_WARNING
		_reviewMSG.MSG = fmt.Sprintf("警告: 链接实例检测表相关信息. 关闭连接出错. %v",
			_reviewMSG.MSG)
	}

	return _reviewMSG
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
		whereItems := whereReg.Split(setSplitItems[1], _setWhereCount + 2)
		explainSelectWhere = whereItems[len(whereItems) - 1]
	}

	explainSelectSql = fmt.Sprintf("%v where %v", explainSelectSuffix, explainSelectWhere)

	return explainSelectSql, nil
}
